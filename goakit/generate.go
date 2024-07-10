package goakit

import (
	"path"
	"regexp"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
	httpcodegen "goa.design/goa/v3/http/codegen"
)

// Register the plugin Generator functions.
func init() {
	codegen.RegisterPluginFirst("goakit", "gen", nil, Generate)
	codegen.RegisterPluginLast("goakit-goakitify", "gen", nil, Goakitify)
	codegen.RegisterPluginLast("goakit-goakitify-example", "example", nil, GoakitifyExample)
}

// Generate generates go-kit specific decoders and encoders.
func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, root := range roots {
		if r, ok := root.(*expr.RootExpr); ok {
			files = append(files, EncodeDecodeFiles(genpkg, r)...)
			files = append(files, MountFiles(r)...)
		}
	}
	return files, nil
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
			oldinit := "{{ .VarName }}Svc = {{ $.APIPkg }}.New{{ .StructName }}()"
			section.Source = strings.Replace(section.Source, oldinit, initT, 1)
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
				`{{ .VarName }}Endpoints.Use(debug.LogPayloads())`,
				`{{ .VarName }}Endpoints.Use(wrapMiddleware(debug.LogPayloads()))`,
				1,
			)
			section.Source = strings.Replace(
				section.Source,
				`{{ .VarName }}Endpoints.Use(log.Endpoint)`,
				`{{ .VarName }}Endpoints.Use(wrapMiddleware(log.Endpoint))`,
				1,
			)
		case "server-http-init":
			codegen.AddImport(file.SectionTemplates[0], &codegen.ImportSpec{Path: "github.com/go-kit/kit/transport/http", Name: "kithttp"})
			codegen.AddImport(file.SectionTemplates[0], &codegen.ImportSpec{Path: "github.com/go-kit/kit/endpoint"})
			data := section.Data.(map[string]interface{})
			svcs := data["Services"].([]*httpcodegen.ServiceData)
			for _, svc := range svcs {
				svcData := httpcodegen.HTTPServices.Get(svc.Service.Name).Service
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
	{{ .VarName }}Svc = {{ $.APIPkg }}.New{{ .StructName }}(logger)
}
`

const basicServiceStructT = `
{{ printf "%s service example implementation.\nThe example methods log the requests and return zero values." .Name | comment }}
type {{ .VarName }}srvc struct {
	logger log.Logger
}
`

const basicServiceInitT = `
{{ printf "New%s returns the %s service implementation." .StructName .Name | comment }}
func New{{ .StructName }}(logger log.Logger) {{ .PkgName }}.Service {
	return &{{ .VarName }}srvc{
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
    {{- if needStream .Services }}
      upgrader := &websocket.Upgrader{}
    {{- end }}
  {{- range .Services }}
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
      {{ .Service.VarName }}Server = {{ .Service.PkgName }}svr.New({{ .Service.VarName }}Endpoints, mux, dec, enc, eh, nil{{ if needStream $.Services }}, upgrader, nil{{ end }}{{ range .Endpoints }}{{ if .MultipartRequestDecoder }}, {{ $.APIPkg }}.{{ .MultipartRequestDecoder.FuncName }}{{ end }}{{ end }})
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
