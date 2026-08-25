// This file writes go-kit mount functions using the names chosen while Goa
// planned the generated server package.
package goakit

import (
	"fmt"
	"path/filepath"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/expr"
	httpcodegen "goa.design/goa/v3/http/codegen"
)

type (
	// goakitFileServerData contains the chosen mount name rendered by the go-kit
	// template alongside the file server values supplied by Goa.
	goakitFileServerData struct {
		*httpcodegen.FileServerData
		MountHandler string
	}
)

// MountFiles produces the files containing the HTTP handler mount functions
// that configure the mux to serve the requests.
func MountFiles(root *expr.RootExpr) []*codegen.File {
	plan, err := planHTTP("generated.local/gen", root)
	if err != nil {
		panic(err)
	}
	return mountFiles(plan)
}

// mountFiles builds mount functions from the HTTP names chosen for the current
// generation run.
func mountFiles(plan *goakitRootPlan) []*codegen.File {
	fw := make([]*codegen.File, len(plan.root.API.HTTP.Services))
	for i, service := range plan.root.API.HTTP.Services {
		data := httpServiceData(plan.http, service)
		fw[i] = mountFile(service, data, plan.services[service])
	}
	return fw
}

// mountFile returns the file defining the mount handler functions for the given
// service.
func mountFile(
	svc *expr.HTTPServiceExpr,
	data *httpcodegen.ServiceData,
	names *goakitServicePlan,
) *codegen.File {
	path := filepath.Join(codegen.Gendir, "http", data.Service.PathName, "kitserver", "mount.go")
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
			Data:   plannedEndpointData(e, names.endpoints[e.Method.Name]),
		})
	}
	fm := codegen.TemplateFuncs()
	fm["join"] = strings.Join
	for index, fs := range data.FileServers {
		sections = append(sections, &codegen.SectionTemplate{
			Name:   "goakit-mount-file-server",
			Source: mountFileServerT,
			Data: &goakitFileServerData{
				FileServerData: fs,
				MountHandler:   names.fileServers[index].Name(),
			},
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
