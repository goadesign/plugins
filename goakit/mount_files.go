package goakit

import (
	"fmt"
	"path/filepath"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/expr"
	httpcodegen "goa.design/goa/v3/http/codegen"
)

// MountFiles produces the files containing the HTTP handler mount functions
// that configure the mux to serve the requests.
func MountFiles(root *expr.RootExpr) []*codegen.File {
	fw := make([]*codegen.File, len(root.API.HTTP.Services))
	for i, svc := range root.API.HTTP.Services {
		fw[i] = mountFile(svc)
	}
	return fw
}

// mountFile returns the file defining the mount handler functions for the given
// service.
func mountFile(svc *expr.HTTPServiceExpr) *codegen.File {
	path := filepath.Join(codegen.Gendir, "http", codegen.SnakeCase(svc.Name()), "kitserver", "mount.go")
	data := httpcodegen.HTTPServices.Get(svc.Name())
	title := fmt.Sprintf("%s go-kit HTTP server encoders and decoders", svc.Name())
	sections := []*codegen.SectionTemplate{
		codegen.Header(title, "server", []*codegen.ImportSpec{
			{Path: "net/http"},
			{Path: "goa.design/goa/v3/http", Name: "goahttp"},
		}),
	}
	for _, e := range data.Endpoints {
		sections = append(sections, &codegen.SectionTemplate{
			Name:   "goakit-mount-handler",
			Source: mountHandlerT,
			Data:   e,
		})
	}
	fm := codegen.TemplateFuncs()
	fm["join"] = strings.Join
	for _, fs := range data.FileServers {
		sections = append(sections, &codegen.SectionTemplate{
			Name:    "goakit-mount-file-server",
			Source:  mountFileServerT,
			Data:    fs,
			FuncMap: fm,
		})
	}

	return &codegen.File{Path: path, SectionTemplates: sections}
}

// input: EndpointData
const mountHandlerT = `{{ printf "%s configures the mux to serve the %q service %q endpoint." .MountHandler .ServiceName .Method.Name | comment }}
func {{ .MountHandler }}(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	{{- range .Routes }}
	mux.Handle("{{ .Verb }}", "{{ .Path }}", f)
	{{- end }}
}
`

// input: FileServerData
const mountFileServerT = `{{ printf "%s configures the mux to serve GET request made to %q." .MountHandler (join .RequestPaths ", ") | comment }}
func {{ .MountHandler }}(mux goahttp.Muxer) {
{{- if .IsDir }}
	{{- range .RequestPaths }}
	mux.Handle("GET", "{{ . }}", http.FileServer(http.Dir({{ printf "%q" $.FilePath }})))
	{{- end }}
{{- else }}
	{{- range .RequestPaths }}
	mux.Handle("GET", "{{ . }}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, {{ printf "%q" $.FilePath }})
		}))
	{{- end }}
{{- end }}
}
`
