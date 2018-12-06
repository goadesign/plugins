package cors

import (
	"path/filepath"
	"strings"

	"goa.design/goa/codegen"
	"goa.design/goa/codegen/service"
	"goa.design/goa/eval"
	httpcodegen "goa.design/goa/http/codegen"
	"goa.design/plugins/cors/expr"
)

// ServicesData holds the all the ServiceData indexed by service name.
var ServicesData = make(map[string]*ServiceData)

type (
	// ServiceData contains the data necessary to generate origin handlers
	ServiceData struct {
		// Name is the name of the service.
		Name string
		// Origins is a list of origin expressions defined in API and service levels.
		Origins []*expr.OriginExpr
		// OriginHandler is the name of the handler function that sets CORS headers.
		OriginHandler string
		// PreflightPaths is the list of paths that should handle OPTIONS requests.
		PreflightPaths []string
		// Endpoint is the CORS endpoint data.
		Endpoint *httpcodegen.EndpointData
	}
)

// Register the plugin Generator functions.
func init() {
	codegen.RegisterPlugin("cors", "gen", nil, Generate)
}

// Generate produces server code that handle preflight requests and updates
// the HTTP responses with the appropriate CORS headers.
func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, f := range files {
		serverCORS(f)
	}
	return files, nil
}

// buildServiceData builds the data needed to render the CORS handlers.
func buildServiceData(svc string) *ServiceData {
	preflights := expr.PreflightPaths(svc)
	routes := make([]*httpcodegen.RouteData, len(preflights))
	for i, p := range preflights {
		routes[i] = &httpcodegen.RouteData{Verb: "OPTIONS", Path: p}
	}

	return &ServiceData{
		Name:           svc,
		Origins:        expr.Origins(svc),
		PreflightPaths: preflights,
		OriginHandler:  "handle" + codegen.Goify(svc, true) + "Origin",
		Endpoint: &httpcodegen.EndpointData{
			Method: &service.MethodData{
				VarName: "CORS",
			},
			MountHandler: "MountCORSHandler",
			HandlerInit:  "NewCORSHandler",
			Routes:       routes,
		},
	}
}

// serverCORS updates the HTTP server file to handle preflight paths and
// adds the required CORS headers to the response.
func serverCORS(f *codegen.File) {
	if filepath.Base(f.Path) != "server.go" {
		return
	}

	var svcData *ServiceData
	for _, s := range f.Section("server-struct") {
		codegen.AddImport(f.SectionTemplates[0],
			&codegen.ImportSpec{Path: "goa.design/plugins/cors"})

		data := s.Data.(*httpcodegen.ServiceData)
		if d, ok := ServicesData[data.Service.Name]; !ok {
			svcData = buildServiceData(data.Service.Name)
			ServicesData[data.Service.Name] = svcData
		} else {
			svcData = d
		}
		for _, o := range svcData.Origins {
			if o.Regexp {
				codegen.AddImport(f.SectionTemplates[0],
					&codegen.ImportSpec{Path: "regexp"})
				break
			}
		}
		data.Endpoints = append(data.Endpoints, svcData.Endpoint)
		fm := codegen.TemplateFuncs()
		f.SectionTemplates = append(f.SectionTemplates, &codegen.SectionTemplate{
			Name:    "mount-cors",
			Source:  mountCORST,
			Data:    svcData,
			FuncMap: fm,
		})
		f.SectionTemplates = append(f.SectionTemplates, &codegen.SectionTemplate{
			Name:    "cors-handler-init",
			Source:  corsHandlerInitT,
			Data:    svcData,
			FuncMap: fm,
		})
		fm["join"] = strings.Join
		f.SectionTemplates = append(f.SectionTemplates, &codegen.SectionTemplate{
			Name:    "handle-cors",
			Source:  handleCORST,
			Data:    svcData,
			FuncMap: fm,
		})
	}
	for _, s := range f.Section("server-init") {
		s.Source = strings.Replace(s.Source,
			"e.{{ .Method.VarName }}, mux, {{ if .MultipartRequestDecoder }}{{ .MultipartRequestDecoder.InitName }}(mux, {{ .MultipartRequestDecoder.VarName }}){{ else }}dec{{ end }}, enc, eh",
			`{{ if ne .Method.VarName "CORS" }}e.{{ .Method.VarName }}, mux, {{ if .MultipartRequestDecoder }}{{ .MultipartRequestDecoder.InitName }}(mux, {{ .MultipartRequestDecoder.VarName }}){{ else }}dec{{ end }}, enc, eh{{ end }}`,
			-1)
	}
	for _, s := range f.Section("server-handler") {
		s.Source = strings.Replace(s.Source, "h.(http.HandlerFunc)", svcData.OriginHandler+"(h).(http.HandlerFunc)", -1)
	}
	for _, s := range f.Section("server-files") {
		s.Source = strings.Replace(s.Source, "h.ServeHTTP", svcData.OriginHandler+"(h).ServeHTTP", -1)
	}
}

// Data: ServiceData
var corsHandlerInitT = `{{ printf "%s creates a HTTP handler which returns a simple 200 response." .Endpoint.HandlerInit | comment }}
func {{ .Endpoint.HandlerInit }}() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
}
`

// Data: ServiceData
var mountCORST = `{{ printf "%s configures the mux to serve the CORS endpoints for the service %s." .Endpoint.MountHandler .Name | comment }}
func {{ .Endpoint.MountHandler }}(mux goahttp.Muxer, h http.Handler) {
	h = {{ .OriginHandler }}(h)
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	{{- range $p := .PreflightPaths }}
	mux.Handle("OPTIONS", "{{ $p }}", f)
	{{- end }}
}
`

// Data: ServiceData
var handleCORST = `{{ printf "%s applies the CORS response headers corresponding to the origin for the service %s." .OriginHandler .Name | comment }}
func {{ .OriginHandler }}(h http.Handler) http.Handler {
{{- range $i, $policy := .Origins }}
	{{- if $policy.Regexp }}
	spec{{$i}} := regexp.MustCompile({{ printf "%q" $policy.Origin }})
	{{- end }}
{{- end }}
	origHndlr := h.(http.HandlerFunc)
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    origin := r.Header.Get("Origin")
    if origin == "" {
      // Not a CORS request
			origHndlr(w, r)
			return
    }
	{{- range $i, $policy := .Origins }}
		{{- if $policy.Regexp }}
		if cors.MatchOriginRegexp(origin, spec{{$i}}) {
		{{- else }}
		if cors.MatchOrigin(origin, {{ printf "%q" $policy.Origin }}) {
		{{- end }}
      w.Header().Set("Access-Control-Allow-Origin", origin)
			{{- if not (eq $policy.Origin "*") }}
			w.Header().Set("Vary", "Origin")
			{{- end }}
			{{- if $policy.Exposed }}
			w.Header().Set("Access-Control-Expose-Headers", "{{ join $policy.Exposed ", " }}")
			{{- end }}
			{{- if gt $policy.MaxAge 0 }}
			w.Header().Set("Access-Control-Max-Age", "{{ $policy.MaxAge }}")
			{{- end }}
			w.Header().Set("Access-Control-Allow-Credentials", "{{ $policy.Credentials }}")
      if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
        // We are handling a preflight request
				{{- if $policy.Methods }}
				w.Header().Set("Access-Control-Allow-Methods", "{{ join $policy.Methods ", " }}")
				{{- end }}
				{{- if $policy.Headers }}
				w.Header().Set("Access-Control-Allow-Headers", "{{ join $policy.Headers ", " }}")
				{{- end }}
			}
			origHndlr(w, r)
			return
    }
	{{- end }}
		origHndlr(w, r)
		return
  })
}
`
