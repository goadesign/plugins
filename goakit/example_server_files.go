package goakit

import (
	"os"
	"path/filepath"
	"strings"

	"goa.design/goa/codegen"
	"goa.design/goa/design"
	httpcodegen "goa.design/goa/http/codegen"
	httpdesign "goa.design/goa/http/design"
)

// ExampleServerFiles returns and example main and dummy service
// implementations.
func ExampleServerFiles(genpkg string, root *httpdesign.RootExpr) []*codegen.File {
	fw := make([]*codegen.File, len(root.HTTPServices)+1)
	for i, svc := range root.HTTPServices {
		fw[i] = dummyServiceFile(genpkg, svc)
	}
	fw[len(root.HTTPServices)] = exampleMain(genpkg, root)
	return fw
}

// dummyServiceFile returns a dummy implementation of the given service.
func dummyServiceFile(genpkg string, svc *httpdesign.ServiceExpr) *codegen.File {
	path := codegen.SnakeCase(svc.Name()) + ".go"
	data := httpcodegen.HTTPServices.Get(svc.Name())
	pkgName := httpcodegen.HTTPServices.Get(svc.Name()).Service.PkgName
	sections := []*codegen.SectionTemplate{
		codegen.Header("", codegen.KebabCase(design.Root.API.Name), []*codegen.ImportSpec{
			{Path: "context"},
			{Path: "github.com/go-kit/kit/log"},
			{Path: filepath.Join(genpkg, svc.Name()), Name: pkgName},
		}),
		{Name: "goakit-dummy-service-struct", Source: dummyServiceStructT, Data: data},
	}
	for _, e := range data.Endpoints {
		sections = append(sections, &codegen.SectionTemplate{
			Name:   "goakit-dummy-endpoint",
			Source: dummyEndpointImplT,
			Data:   e,
		})
	}

	return &codegen.File{Path: path, SectionTemplates: sections}
}

func exampleMain(genpkg string, root *httpdesign.RootExpr) *codegen.File {
	path := filepath.Join("cmd", codegen.SnakeCase(root.Design.API.Name)+"svc", "main.go")
	idx := strings.LastIndex(genpkg, string(os.PathSeparator))
	rootPath := "."
	if idx > 0 {
		rootPath = genpkg[:idx]
	}
	specs := []*codegen.ImportSpec{
		{Path: "context"},
		{Path: "flag"},
		{Path: "fmt"},
		{Path: "net/http"},
		{Path: "os"},
		{Path: "os/signal"},
		{Path: "time"},
		{Path: "github.com/go-kit/kit/endpoint"},
		{Path: "github.com/go-kit/kit/log"},
		{Path: "github.com/go-kit/kit/transport/http", Name: "kithttp"},
		{Path: "goa.design/goa", Name: "goa"},
		{Path: "goa.design/goa/http", Name: "goahttp"},
		{Path: rootPath, Name: codegen.KebabCase(root.Design.API.Name)},
		{Path: "goa.design/goa/http/middleware"},
	}
	for _, svc := range root.HTTPServices {
		pkgName := httpcodegen.HTTPServices.Get(svc.Name()).Service.PkgName
		specs = append(specs, &codegen.ImportSpec{
			Path: filepath.Join(genpkg, "http", svc.Name(), "kitserver"),
			Name: pkgName + "kitsvr",
		})
		specs = append(specs, &codegen.ImportSpec{
			Path: filepath.Join(genpkg, svc.Name()),
			Name: pkgName,
		})
		specs = append(specs, &codegen.ImportSpec{
			Path: filepath.Join(genpkg, "http", codegen.SnakeCase(svc.Name()), "server"),
			Name: pkgName + "svr",
		})
	}
	sections := []*codegen.SectionTemplate{
		codegen.Header("", "main", specs),
	}
	var svcdata []*httpcodegen.ServiceData
	for _, svc := range root.HTTPServices {
		svcdata = append(svcdata, httpcodegen.HTTPServices.Get(svc.Name()))
	}
	data := map[string]interface{}{
		"Services": svcdata,
		"APIPkg":   codegen.KebabCase(root.Design.API.Name),
	}
	sections = append(sections, &codegen.SectionTemplate{
		Name:   "goakit-main",
		Source: mainT,
		Data:   data,
	})

	return &codegen.File{Path: path, SectionTemplates: sections}
}

// input: ServiceData
const dummyServiceStructT = `{{ printf "%s service example implementation.\nThe example methods log the requests and return zero values." .Service.Name | comment }}
type {{ .Service.PkgName }}svc struct {
	logger log.Logger
}

{{ printf "New%s returns the %s service implementation." .Service.VarName .Service.Name | comment }}
func New{{ .Service.VarName }}(logger log.Logger) {{ .Service.PkgName }}.Service {
	return &{{ .Service.PkgName }}svc{logger}
}
`

// input: EndpointData
const dummyEndpointImplT = `{{ comment .Method.Description }}
func (s *{{ .ServicePkgName }}svc) {{ .Method.VarName }}(ctx context.Context{{ if .Payload.Ref }}, p {{ .Payload.Ref }}{{ end }}) ({{ if .Result.Ref }}{{ .Result.Ref }}, {{ end }}error) {
{{- if .Result.Ref }}
	var res {{ .Result.Ref }}
{{- end }}
	s.logger.Log("msg", "{{ .ServiceName }}.{{ .Method.Name }}")
	return {{ if .Result.Ref }}res, {{ end }}nil
}
`

// input: map[string]interface{}{"Services":[]ServiceData, "APIPkg": string}
const mainT = `func main() {
	// Define command line flags, add any other flag required to configure
	// the service.
	var (
		addr = flag.String("listen", ":8080", "HTTP listen ` + "`" + `address` + "`" + `")
	)
	flag.Parse()

	// Setup logger.
	var (
		logger log.Logger
	)
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Create the structs that implement the services.
	var (
	{{- range .Services }}
		{{-  if .Endpoints }}
		{{ .Service.PkgName }}s {{.Service.PkgName}}.Service
		{{- end }}
	{{- end }}
	)
	{
	{{- range .Services }}
		{{-  if .Endpoints }}
		{{ .Service.PkgName }}s = {{ $.APIPkg }}.New{{ .Service.VarName }}(logger)
		{{- end }}
	{{- end }}
	}

	// Wrap the services in endpoints that can be invoked from other
	// services potentially running in different processes.
	var (
	{{- range .Services }}
		{{-  if .Endpoints }}
		{{ .Service.PkgName }}e *{{.Service.PkgName}}.Endpoints
		{{- end }}
	{{- end }}
	)
	{
	{{- range .Services }}
		{{-  if .Endpoints }}
		{{ .Service.PkgName }}e = {{ .Service.PkgName }}.NewEndpoints({{ .Service.PkgName }}s)
		{{- end }}
	{{- end }}
	}

	// Provide the transport specific request decoder and response encoder.
	// The goa http package has built-in support for JSON, XML and gob.
	// Other encodings can be used by providing the corresponding functions,
	// see goa.design/encoding.
	var (
		dec = goahttp.RequestDecoder
		enc = goahttp.ResponseEncoder
	)

	// Build the service HTTP request router (a.k.a. mux).
	var mux goahttp.Muxer
	{
		mux = goahttp.NewMuxer()
	}

	// Wrap the endpoints with the transport specific layer.
	var (
	{{- range .Services }}
		{{- range .Endpoints }}
		{{ .ServicePkgName }}{{ .Method.VarName }}Handler *kithttp.Server
		{{- end }}
		{{ .Service.PkgName }}Server *{{.Service.PkgName}}svr.Server
	{{- end }}
	)
	{
	{{- range .Services }}
		eh := ErrorHandler(logger)
		{{- range .Endpoints }}
		{{ .ServicePkgName }}{{ .Method.VarName }}Handler = kithttp.NewServer(
			endpoint.Endpoint({{ .ServicePkgName }}e.{{ .Method.VarName }}),
			{{- if .Payload.Ref }}
			{{ .ServicePkgName}}kitsvr.{{ .RequestDecoder }}(mux, dec),
			{{- else }}
			func(context.Context, *http.Request) (request interface{}, err error) { return nil, nil },
			{{- end }}
			{{ .ServicePkgName}}kitsvr.{{ .ResponseEncoder }}(enc),
		)
		{{- end }}
		{{-  if .Endpoints }}
		{{ .Service.PkgName }}Server = {{ .Service.PkgName }}svr.New({{ .Service.PkgName }}e, mux, dec, enc, eh)
		{{-  else }}
		{{ .Service.PkgName }}Server = {{ .Service.PkgName }}svr.New(nil, mux, dec, enc, eh)
		{{-  end }}
	{{- end }}
	}

	// Configure the mux.
	{{- range .Services }}{{ $service := . }}
		{{- range .Endpoints }}
	{{ .ServicePkgName}}kitsvr.{{ .MountHandler }}(mux, {{ .ServicePkgName }}{{ .Method.VarName }}Handler)
		{{- end }}
		{{- range .FileServers }}
	{{ $service.Service.PkgName}}kitsvr.{{ .MountHandler }}(mux)
		{{- end }}
	{{- end }}

	// Create channel used by both the signal handler and server goroutines
	// to notify the main goroutine when to stop the server.
	errc := make(chan error)

	// Setup interrupt handler. This optional step configures the process so
	// that SIGINT and SIGTERM signals cause the service to stop gracefully.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// Start HTTP server using default configuration, change the code to
	// configure the server as required by your service.
	srv := &http.Server{Addr: *addr, Handler: mux}
	go func() {
		{{- range .Services }}
		for _, m := range {{ .Service.PkgName }}Server.Mounts {
			{{- if .FileServers }}
			logger.Log("info", fmt.Sprintf("service %s file %s mounted on %s %s", {{ .Service.PkgName }}Server.Service(), m.Method, m.Verb, m.Pattern))
			{{- else }}
			logger.Log("info", fmt.Sprintf("service %s method %s mounted on %s %s", {{ .Service.PkgName }}Server.Service(), m.Method, m.Verb, m.Pattern))
			{{- end }}
		}
		{{- end }}
		logger.Log("listening", *addr)
		errc <- srv.ListenAndServe()
	}()

	// Wait for signal.
	logger.Log("exiting", <-errc)

	// Shutdown gracefully with a 30s timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

	logger.Log("server", "exited")
}

// ErrorHandler returns a function that writes and logs the given error.
// The function also writes and logs the error unique ID so that it's possible
// to correlate.
func ErrorHandler(logger log.Logger) func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		id := ctx.Value(middleware.RequestIDKey).(string)
		w.Write([]byte("[" + id + "] encoding: " + err.Error()))
		logger.Log("error", fmt.Sprintf("[%s] ERROR: %s", id, err.Error()))
	}
}
`
