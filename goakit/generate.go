// This file registers the go-kit generator and records every package name it
// adds before Goa chooses final names for generated code.
package goakit

import (
	"cmp"
	"path"
	"regexp"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/generator"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
	httpcodegen "goa.design/goa/v3/http/codegen"
)

type (
	// goakitRootPlan keeps the Goa HTTP plan and the go-kit names declared for
	// one design root during the same generation run.
	goakitRootPlan struct {
		root     *expr.RootExpr
		http     *httpcodegen.Plan
		services map[*expr.HTTPServiceExpr]*goakitServicePlan
	}

	// goakitServicePlan keeps the function names emitted in the go-kit client
	// and server packages for one service.
	goakitServicePlan struct {
		endpoints   map[string]*goakitEndpointNames
		fileServers []*codegen.NameDeclaration
	}

	// goakitEndpointNames keeps the wrapper names emitted for one method.
	goakitEndpointNames struct {
		mountHandler    *codegen.NameDeclaration
		requestDecoder  *codegen.NameDeclaration
		responseEncoder *codegen.NameDeclaration
		errorEncoder    *codegen.NameDeclaration
		requestEncoder  *codegen.NameDeclaration
		responseDecoder *codegen.NameDeclaration
	}

	// goakitNameOrder gives every generated go-kit function a stable position
	// when two plugins ask for the same name.
	goakitNameOrder struct {
		api, service, method, subject string
		role                          goakitNameRole
	}

	// goakitNameRole identifies the function written by one go-kit template.
	goakitNameRole uint8
)

const (
	goakitMountHandlerRole goakitNameRole = iota + 1
	goakitRequestDecoderRole
	goakitResponseEncoderRole
	goakitErrorEncoderRole
	goakitRequestEncoderRole
	goakitResponseDecoderRole
	goakitFileServerRole
)

// Register the plugin Generator functions.
func init() {
	generator.RegisterPluginFirst("goakit", "gen", newPlugin)
	codegen.RegisterPluginLast("goakit-goakitify", "gen", nil, Goakitify)
	codegen.RegisterPluginLast("goakit-goakitify-example", "example", nil, GoakitifyExample)
}

// Generate generates go-kit specific decoders and encoders.
func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, root := range roots {
		if r, ok := root.(*expr.RootExpr); ok {
			plan, err := planHTTP(genpkg, r)
			if err != nil {
				return nil, err
			}
			files = append(files, encodeDecodeFiles(genpkg, plan)...)
			files = append(files, mountFiles(plan)...)
		}
	}
	return files, nil
}

// newPlugin creates the HTTP package additions used by one Goa generation.
func newPlugin() generator.Plugin {
	var plans []*goakitRootPlan
	return generator.Plugin{
		Plan: func(plan *generator.Plan) error {
			for _, candidate := range plan.Generation().Roots() {
				root, ok := candidate.(*expr.RootExpr)
				if !ok {
					continue
				}
				httpPlan, ok := plan.HTTP(root)
				if !ok {
					continue
				}
				services, err := planPackages(plan.Generation(), plan.Service(root), root)
				if err != nil {
					return err
				}
				plans = append(plans, &goakitRootPlan{root: root, http: httpPlan, services: services})
			}
			return nil
		},
		Generate: func(plan *generator.Plan, files []*codegen.File) ([]*codegen.File, error) {
			for _, planned := range plans {
				files = append(files, encodeDecodeFiles(plan.Generation().GenPkg(), planned)...)
				files = append(files, mountFiles(planned)...)
			}
			return files, nil
		},
	}
}

// planPackages records every package, import, and function emitted by the
// go-kit files before Goa chooses names for the generation run.
func planPackages(generation *codegen.Generation, servicePlan *service.Plan, root *expr.RootExpr) (map[*expr.HTTPServiceExpr]*goakitServicePlan, error) {
	plans := make(map[*expr.HTTPServiceExpr]*goakitServicePlan, len(root.API.HTTP.Services))
	for _, transportService := range root.API.HTTP.Services {
		serviceImport, _, err := servicePlan.ServicePackageImports(transportService.ServiceExpr)
		if err != nil {
			return nil, err
		}
		servicePath := path.Base(serviceImport.Path)
		serverPath := path.Join(generation.GenPkg(), "http", servicePath, "kitserver")
		serverPackage, err := generation.ClaimPackage(serverPath)
		if err != nil {
			return nil, err
		}
		if err := requireImports(serverPackage, []*codegen.ImportSpec{
			{Path: "context"},
			{Path: "net/http"},
			{Path: "strings"},
			{Path: "github.com/go-kit/kit/transport/http", Name: "kithttp"},
			{Path: "goa.design/goa/v3", Name: "goa"},
			{Path: "goa.design/goa/v3/http", Name: "goahttp"},
			{Path: path.Join(generation.GenPkg(), "http", servicePath, "server")},
		}); err != nil {
			return nil, err
		}
		clientPath := path.Join(generation.GenPkg(), "http", servicePath, "kitclient")
		clientPackage, err := generation.ClaimPackage(clientPath)
		if err != nil {
			return nil, err
		}
		if err := requireImports(clientPackage, []*codegen.ImportSpec{
			{Path: "context"},
			{Path: "net/http"},
			{Path: "strings"},
			{Path: "github.com/go-kit/kit/transport/http", Name: "kithttp"},
			{Path: "goa.design/goa/v3", Name: "goa"},
			{Path: "goa.design/goa/v3/http", Name: "goahttp"},
			{Path: path.Join(generation.GenPkg(), "http", servicePath, "client")},
		}); err != nil {
			return nil, err
		}
		planned, err := planServiceNames(root, servicePlan, transportService, serverPackage, clientPackage)
		if err != nil {
			return nil, err
		}
		plans[transportService] = planned
	}
	return plans, nil
}

// planServiceNames declares every wrapper and mount function emitted for one
// HTTP service.
func planServiceNames(root *expr.RootExpr, servicePlan *service.Plan, transportService *expr.HTTPServiceExpr, serverPackage, clientPackage *codegen.GeneratedPackage) (*goakitServicePlan, error) {
	planned := &goakitServicePlan{endpoints: make(map[string]*goakitEndpointNames, len(transportService.HTTPEndpoints))}
	for _, endpoint := range transportService.HTTPEndpoints {
		methodNames, err := servicePlan.HTTPMethodNames(endpoint.MethodExpr)
		if err != nil {
			return nil, err
		}
		order := goakitNameOrder{api: root.API.Name, service: transportService.Name(), method: endpoint.MethodExpr.Name}
		names := new(goakitEndpointNames)
		if names.mountHandler, err = declareGoakitName(serverPackage, "Mount"+methodNames.Method+"Handler", order.withRole(goakitMountHandlerRole)); err != nil {
			return nil, err
		}
		if names.responseEncoder, err = declareGoakitName(serverPackage, "Encode"+methodNames.Method+"Response", order.withRole(goakitResponseEncoderRole)); err != nil {
			return nil, err
		}
		if endpoint.MethodExpr.Payload.Type != expr.Empty {
			if names.requestDecoder, err = declareGoakitName(serverPackage, "Decode"+methodNames.Method+"Request", order.withRole(goakitRequestDecoderRole)); err != nil {
				return nil, err
			}
		}
		if len(endpoint.HTTPErrors) > 0 {
			if names.errorEncoder, err = declareGoakitName(serverPackage, "Encode"+methodNames.Method+"Error", order.withRole(goakitErrorEncoderRole)); err != nil {
				return nil, err
			}
		}
		if goakitClientRequestEncoderSelected(endpoint) {
			if names.requestEncoder, err = declareGoakitName(clientPackage, "Encode"+methodNames.Method+"Request", order.withRole(goakitRequestEncoderRole)); err != nil {
				return nil, err
			}
		}
		if names.responseDecoder, err = declareGoakitName(clientPackage, "Decode"+methodNames.Method+"Response", order.withRole(goakitResponseDecoderRole)); err != nil {
			return nil, err
		}
		planned.endpoints[endpoint.MethodExpr.Name] = names
	}
	for _, fileServer := range transportService.FileServers {
		order := goakitNameOrder{
			api:     root.API.Name,
			service: transportService.Name(),
			subject: fileServer.FilePath,
			role:    goakitFileServerRole,
		}
		declaration, err := declareGoakitName(serverPackage, "Mount"+codegen.Goify(fileServer.FilePath, true), order)
		if err != nil {
			return nil, err
		}
		planned.fileServers = append(planned.fileServers, declaration)
	}
	return planned, nil
}

// declareGoakitName records one public function that a go-kit template writes.
func declareGoakitName(pkg *codegen.GeneratedPackage, preferred string, order goakitNameOrder) (*codegen.NameDeclaration, error) {
	declaration := codegen.NewPreferredName(codegen.NameFunction, preferred, codegen.ExportedName, order)
	if err := pkg.DeclareName(declaration); err != nil {
		return nil, err
	}
	return declaration, nil
}

// goakitClientRequestEncoderSelected reports whether the HTTP client writes a
// request encoder that the go-kit client must wrap.
func goakitClientRequestEncoderSelected(endpoint *expr.HTTPEndpointExpr) bool {
	if endpoint.IsJSONRPC() {
		return true
	}
	if (!endpoint.SkipRequestBodyEncodeDecode && endpoint.Body.Type != expr.Empty) ||
		endpoint.MapQueryParams != nil ||
		len(*expr.AsObject(endpoint.QueryParams().Type)) > 0 ||
		len(*expr.AsObject(endpoint.Headers.Type)) > 0 ||
		len(*expr.AsObject(endpoint.Cookies.Type)) > 0 {
		return true
	}
	for _, requirement := range endpoint.Requirements {
		for _, scheme := range requirement.Schemes {
			if scheme.Kind == expr.BasicAuthKind {
				return true
			}
		}
	}
	return false
}

// requireImports records import names already written in go-kit templates.
func requireImports(pkg *codegen.GeneratedPackage, imports []*codegen.ImportSpec) error {
	for _, spec := range imports {
		if err := pkg.RequireImport(spec); err != nil {
			return err
		}
	}
	return nil
}

// planHTTP creates linked HTTP data for callers that invoke Generate directly.
func planHTTP(genpkg string, root *expr.RootExpr) (*goakitRootPlan, error) {
	generation, err := codegen.NewGeneration(genpkg, []eval.Root{root})
	if err != nil {
		return nil, err
	}
	servicePlan, err := service.NewPlan(root, generation, expr.NewExampleGenerator(root.API.RandomizerFactory))
	if err != nil {
		return nil, err
	}
	plans, err := httpcodegen.NewPlans(generation, httpcodegen.PlanInput{Root: root, Service: servicePlan})
	if err != nil {
		return nil, err
	}
	services, err := planPackages(generation, servicePlan, root)
	if err != nil {
		return nil, err
	}
	if err := generation.Freeze(); err != nil {
		return nil, err
	}
	if err := servicePlan.Link(); err != nil {
		return nil, err
	}
	if err := plans[0].Link(); err != nil {
		return nil, err
	}
	return &goakitRootPlan{root: root, http: plans[0], services: services}, nil
}

// ComparePackageName orders go-kit names by their complete design identity.
func (o goakitNameOrder) ComparePackageName(other codegen.PackageNameOrder) int {
	right := other.(goakitNameOrder)
	for _, compared := range []int{
		cmp.Compare(o.api, right.api),
		cmp.Compare(o.service, right.service),
		cmp.Compare(o.method, right.method),
		cmp.Compare(o.subject, right.subject),
		cmp.Compare(o.role, right.role),
	} {
		if compared != 0 {
			return compared
		}
	}
	return 0
}

// withRole returns the same name identity for another generated function.
func (o goakitNameOrder) withRole(role goakitNameRole) goakitNameOrder {
	o.role = role
	return o
}

// Goakitify modifies all the previously generated files by adding go-kit
// imports and replacing the following instances "goa.Endpoint" with
// "github.com/go-kit/kit/endpoint".Endpoint
//
// Goakitify also wraps instances of endpoint.Endpoint into instances of
// goa.Endpoint when used as argument of either goagrpc.NewStreamHandler or
// goagrpc.NewUnaryHandler.
func Goakitify(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, f := range files {
		goakitify(f)
	}
	return files, nil
}
func goakitify(f *codegen.File) {
	var hasEndpoint bool
	for _, s := range f.SectionTemplates {
		if !hasEndpoint {
			hasEndpoint = goaEndpointRegexp.MatchString(s.Source)
		}
		s.Source = goaEndpointRegexp.ReplaceAllString(s.Source, "${1}endpoint.Endpoint${2}")
		if s.Name == "grpc-handler-init" {
			s.Source = strings.Replace(s.Source, "Handler(endpoint, ", "Handler(goa.Endpoint(endpoint), ", 1)
		}
	}
	if hasEndpoint {
		codegen.AddImport(
			f.SectionTemplates[0],
			&codegen.ImportSpec{Path: "github.com/go-kit/kit/endpoint"},
		)
	}
}

// GoakitifyExample  modifies all the previously generated example files by
// adding go-kit imports.
func GoakitifyExample(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, f := range files {
		gokitifyExampleServer(genpkg, f)
	}
	return files, nil
}

// goaEndpointRegexp matches occurrences of the "goa.Endpoint" type in Go code.
var goaEndpointRegexp = regexp.MustCompile(`([^\p{L}_])goa\.Endpoint([^\p{L}_])`)

// deletedImports contains the list of imports that should be removed from the
// generated files.
var deletedImports = []string{"log", "goa.design/clue/log"}

// gokitifyExampleServer imports gokit endpoint, logger, and transport
// packages in the example server implementation. It also replaces every stdlib
// logger with gokit logger.
func gokitifyExampleServer(genpkg string, file *codegen.File) {
	goakitify(file)
	hasGoaMiddleware := false
	for _, section := range file.SectionTemplates {
		switch section.Name {
		case "server-main-services":
			deleteImports(file.SectionTemplates[0])
			codegen.AddImport(file.SectionTemplates[0], &codegen.ImportSpec{Name: "kitlog", Path: "github.com/go-kit/log"})
			codegen.AddImport(file.SectionTemplates[0], &codegen.ImportSpec{Path: "goa.design/clue/log"})
			section.Source = strings.Replace(
				section.Source,
				"{{ .ServiceVar }} = {{ $.APIPkg }}.{{ .ExampleConstructorDeclaration.Name }}()",
				initT,
				1,
			)
		case "basic-service-struct":
			deleteImports(file.SectionTemplates[0])
			codegen.AddImport(file.SectionTemplates[0], &codegen.ImportSpec{Path: "github.com/go-kit/log"})
			section.Source = basicServiceStructT
		case "basic-service-init":
			section.Source = basicServiceInitT
		case "basic-endpoint":
			section.Source = strings.Replace(
				section.Source,
				`log.Printf(ctx, "{{ .ServiceVarName }}.{{ .Name }}")`,
				`s.logger.Log("service", "{{ .ServiceVarName}}", "method", "{{ .Name }}")`,
				1,
			)
		case "server-main-endpoints":
			hasGoaMiddleware = true
			section.Source = strings.Replace(
				section.Source,
				`{{ .EndpointsVar }}.Use(debug.LogPayloads())`,
				`{{ .EndpointsVar }}.Use(wrapMiddleware(debug.LogPayloads()))`,
				1,
			)
			section.Source = strings.Replace(
				section.Source,
				`{{ .EndpointsVar }}.Use(log.Endpoint)`,
				`{{ .EndpointsVar }}.Use(wrapMiddleware(log.Endpoint))`,
				1,
			)
		case "server-http-init":
			codegen.AddImport(file.SectionTemplates[0], &codegen.ImportSpec{Path: "github.com/go-kit/kit/transport/http", Name: "kithttp"})
			codegen.AddImport(file.SectionTemplates[0], &codegen.ImportSpec{Path: "github.com/go-kit/kit/endpoint"})
			data := section.Data.(map[string]interface{})
			svcs := data["Services"].([]*httpcodegen.ServiceData)
			for _, svc := range svcs {
				svcData := svc.Service
				codegen.AddImport(file.SectionTemplates[0], &codegen.ImportSpec{
					Path: path.Join(genpkg, "http", svcData.PathName, "kitserver"),
					Name: svcData.PkgName + "kitsvr",
				})
			}
			section.Source = gokitServerInitT
		}
	}
	if hasGoaMiddleware {
		codegen.AddImport(file.SectionTemplates[0], &codegen.ImportSpec{Name: "goa", Path: "goa.design/goa/v3/pkg"})
		codegen.AddImport(file.SectionTemplates[0], &codegen.ImportSpec{Path: "github.com/go-kit/kit/endpoint"})
		file.SectionTemplates = append(file.SectionTemplates, &codegen.SectionTemplate{
			Name:   "middleware-wrapper",
			Source: middlewareWrapperT,
		})
	}
}

// deleteImports removes specified import paths from a section's import specifications.
func deleteImports(section *codegen.SectionTemplate) {
	if data, ok := section.Data.(map[string]interface{}); ok {
		if imports, ok := data["Imports"].([]*codegen.ImportSpec); ok {
			var newimports []*codegen.ImportSpec
		outer:
			for _, imp := range imports {
				for _, del := range deletedImports {
					if imp.Path == del {
						continue outer
					}
				}
				newimports = append(newimports, imp)
			}
			data["Imports"] = newimports
		}
	}
}

const middlewareWrapperT = `
// Wrap goa middleware into go-kit middleware.
func wrapMiddleware(mw func(goa.Endpoint) goa.Endpoint) (func (endpoint.Endpoint) endpoint.Endpoint) {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return endpoint.Endpoint(mw(goa.Endpoint(e)))
	}
}
`

const initT = `{
	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(os.Stderr)
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
	logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	logger = kitlog.With(logger, "service", {{ printf "%q" .Name }})
	{{ .ServiceVar }} = {{ $.APIPkg }}.{{ .ExampleConstructorDeclaration.Name }}(logger)
}
`

const basicServiceStructT = `
{{ printf "%s service example implementation.\nThe example methods log the requests and return zero values." .Name | comment }}
type {{ .ExampleStructDeclaration.Name }} struct {
	logger log.Logger
}
`

const basicServiceInitT = `
{{ printf "New%s returns the %s service implementation." .StructName .Name | comment }}
func {{ .ExampleConstructorDeclaration.Name }}(logger log.Logger) {{ .ServicePkg }}.{{ .ServiceDeclaration.Name }} {
	return &{{ .ExampleStructDeclaration.Name }}{
		logger: logger,
	}
}
`

const gokitServerInitT = `
  // Wrap the endpoints with the transport specific layers. The generated
  // server packages contains code generated from the design which maps
  // the service input and output data structures to HTTP requests and
  // responses.
  var (
  {{- range .Services }}
    {{- range .Endpoints }}
      {{ .ServiceVarName }}{{ .Method.VarName }}Handler *kithttp.Server
    {{- end }}
    {{ .Service.VarName }}Server *{{.Service.PkgName}}svr.Server
  {{- end }}
  )
  {
    eh := errorHandler(ctx)
    {{- if needDialer .Services }}
      upgrader := &websocket.Upgrader{}
    {{- end }}
  {{- range $svc := .Services }}
    {{- if .Endpoints }}
      {{- range .Endpoints }}
        {{ .ServiceVarName }}{{ .Method.VarName }}Handler = kithttp.NewServer(
          endpoint.Endpoint({{ .ServiceVarName }}Endpoints.{{ .Method.VarName }}),
          {{- if .Payload.Ref }}
            {{ .ServicePkgName}}kitsvr.{{ .RequestDecoder }}(mux, dec),
          {{- else }}
            func(context.Context, *http.Request) (request interface{}, err error) { return nil, nil },
          {{- end }}
          {{ .ServicePkgName}}kitsvr.{{ .ResponseEncoder }}(enc),
          {{- if .Errors }}
            kithttp.ServerErrorEncoder({{ .ServicePkgName}}kitsvr.{{ .ErrorEncoder }}(enc, nil)),
          {{- end }}
        )
      {{- end }}
      {{ .Service.VarName }}Server = {{ .Service.PkgName }}svr.New({{ .Service.VarName }}Endpoints, mux, dec, enc, eh, nil{{ if hasWebSocket $svc }}, upgrader, nil{{ end }}{{ range .Endpoints }}{{ if .MultipartRequestDecoder }}, {{ $.APIPkg }}.{{ .MultipartRequestDecoder.FuncName }}{{ end }}{{ end }})
    {{-  else }}
      {{ .Service.VarName }}Server = {{ .Service.PkgName }}svr.New(nil, mux, dec, enc, eh, nil)
    {{-  end }}
  {{- end }}
  }

  // Configure the mux.
  {{- range .Services }}{{ $service := . }}
    {{- range .Endpoints }}
  {{ .ServicePkgName}}kitsvr.{{ .MountHandler }}(mux, {{ .ServiceVarName }}{{ .Method.VarName }}Handler)
    {{- end }}
    {{- range .FileServers }}
  {{ $service.Service.PkgName}}kitsvr.{{ .MountHandler }}(mux)
    {{- end }}
  {{- end }}
`
