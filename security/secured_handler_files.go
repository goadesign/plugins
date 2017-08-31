package security

import (
	"fmt"
	"path/filepath"

	"goa.design/goa/codegen"
	httpdesign "goa.design/goa/http/design"
)

// SecuredHandlerFiles produces the code for the constructors of server handlers
// that enforce the security requirements defined in the design.
func SecuredHandlerFiles(genpkg string, r *httpdesign.RootExpr) []*codegen.File {
	fw := make([]*codegen.File, 2*len(root.HTTPServices))
	for i, svc := range root.HTTPServices {
		fw[i] = securedHhandlerInit(genpkg, svc)
	}
}

// securedHandlerInit returns the file containing the service secured handler
// constructors.
func securedHandlerInit(genpkg string, svc *httpdesign.ServiceExpr) *codegen.File {
	path := filepath.Join(codegen.Gendir, "http", codegen.SnakeCase(svc.Name()), "server", "secured_server.go")
	data := HTTPServices.Get(svc.Name())
	title := fmt.Sprintf("%s HTTP server secured handler constructors", svc.Name())
	sections := []*codegen.SectionTemplate{
		codegen.Header(title, "server", []*codegen.ImportSpec{
			{Path: "context"},
			{Path: "fmt"},
			{Path: "io"},
			{Path: "net/http"},
			{Path: "goa.design/goa", Name: "goa"},
			{Path: "goa.design/goa/http", Name: "goahttp"},
			{Path: genpkg + "/" + data.Service.PkgName},
		}),
	}

	for _, e := range data.Endpoints {
		sections = append(sections, &codegen.SectionTemplate{Name: "secured-server-handler-init", Source: serverHandlerInitT, Data: e})
	}

	return &codegen.File{Path: path, SectionTemplates: sections}
}

// input: EndpointData
const serverHandlerInitT = `{{ printf "%s creates a HTTP handler which loads the HTTP request and calls the \"%s\" service \"%s\" endpoint." .HandlerInit .ServiceName .Method.Name | comment }}
func {{ .HandlerInit }}(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
) http.Handler {
	unsecured := {{ .HandlerInit }}Unsecured(endpoint, mux, dec, enc)
	var (
		{{- if .Payload.Ref }}
		decodeRequest  = {{ .RequestDecoder }}(mux, dec)
		{{- end }}
		encodeResponse = {{ .ResponseEncoder }}(enc)
		encodeError    = {{ if .Errors }}{{ .ErrorEncoder }}{{ else }}goahttp.ErrorEncoder{{ end }}(enc)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		ctx := context.WithValue(r.Context(), goahttp.ContextKeyAcceptType, accept)

		{{- if .Payload.Ref }}
		payload, err := decodeRequest(r)
		if err != nil {
			encodeError(ctx, w, err)
			return
		}

		res, err := endpoint(ctx, payload)
		{{- else }}
		res, err := endpoint(ctx, nil)
		{{- end }}

		if err != nil {
			encodeError(ctx, w, err)
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			encodeError(ctx, w, err)
		}
	})
}
`
