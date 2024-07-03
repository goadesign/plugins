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
// imports and replacing the following instances
// * "goa.Endpoint" with "github.com/go-kit/kit/endpoint".Endpoint
// * "log.Logger" with "github.com/go-kit/log".Logger
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

// logPrintfRegexp matches occurrences of "log.Printf".
var logPrintfRegexp = regexp.MustCompile(`log\.Printf\((ctx|logCtx), (.*)\)`)

// logFatalRegexp matches occurrences of "log.Fatal"
var logFatalRegexp = regexp.MustCompile(`log\.Fatal\(ctx, (.*)\)`)

// logFatalfRegexp matches occurrences of "log.Fatalf".
var logFatalfRegexp = regexp.MustCompile(`log\.Fatalf\(ctx, err, (.*)\)`)

// deletedImports contains the list of imports that should be removed from the
// generated files.
var deletedImports = []string{"log", "goa.design/clue/log"}

// gokitifyExampleServer imports gokit endpoint, logger, and transport
// packages in the example server implementation. It also replaces every stdlib
// logger with gokit logger.
func gokitifyExampleServer(genpkg string, file *codegen.File) {
	goakitify(file)
	for _, section := range file.SectionTemplates {
		deleteImports(section)
		codegen.AddImport(file.SectionTemplates[0], &codegen.ImportSpec{Path: "github.com/go-kit/log"})
		codegen.AddImport(file.SectionTemplates[0], &codegen.ImportSpec{Path: "fmt"})
		section.Source = strings.Replace(section.Source, "*log.Logger", "log.Logger", -1)
		section.Source = logPrintfRegexp.ReplaceAllString(section.Source, "logger.Log(\"info\", ${2})")
		section.Source = logFatalRegexp.ReplaceAllString(section.Source, "logger.Log(\"fatal\", ${1})\nos.Exit(1)")
		section.Source = logFatalfRegexp.ReplaceAllString(section.Source, "logger.Log(\"fatal\", fmt.Sprintf(${1}))\nos.Exit(1)")
		switch section.Name {
		case "basic-service-struct":
		case "server-main-logger":
			codegen.AddImport(file.SectionTemplates[0], &codegen.ImportSpec{Path: "github.com/go-kit/log"})
			section.Source = gokitLoggerT
		case "server-main-services":
			oldinit := "{{ .VarName }}Svc = {{ $.APIPkg }}.New{{ .StructName }}()"
			newinit := "{{ .VarName }}Svc = {{ $.APIPkg }}.New{{ .StructName }}(logger)"
			section.Source = strings.Replace(section.Source, oldinit, newinit, -1)
		case "server-main-endpoints":
			rm := `			{{ .VarName }}Endpoints.Use(debug.LogPayloads())
			{{ .VarName }}Endpoints.Use(log.Endpoint)`
			section.Source = strings.Replace(section.Source, rm, "", -1)
		case "server-main-handler":
			section.Source = strings.Replace(section.Source, "handle{{ toUpper $u.Transport.Name }}Server(ctx", "handle{{ toUpper $u.Transport.Name }}Server(ctx, logger", 1)
		case "server-main-interrupts":
			section.Source = strings.Replace(section.Source, "ctx, cancel := context.WithCancel(ctx)", "ctx, cancel := context.WithCancel(context.Background())", 1)
		case "server-http-configure":
			section.Source = strings.Replace(section.Source, "eh := errorHandler(ctx)", "eh := errorHandler(logger)", 1)
		case "server-http-errorhandler":
			section.Source = strings.Replace(section.Source, "func errorHandler(logCtx context.Context", "func errorHandler(logger log.Logger", 1)
		case "server-http-start":
			section.Source = strings.Replace(section.Source, "ctx context.Context", "ctx context.Context, logger log.Logger", 1)
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

const gokitLoggerT = `
  // Setup gokit logger.
  var (
    logger log.Logger
  )
  {
    logger = log.NewLogfmtLogger(os.Stderr)
    logger = log.With(logger, "ts", log.DefaultTimestampUTC)
    logger = log.With(logger, "caller", log.DefaultCaller)
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
    eh := errorHandler(logger)
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
