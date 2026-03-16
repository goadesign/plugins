package tscap

import (
	"path/filepath"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	goahttp "goa.design/goa/v3/http/codegen"
	"goa.design/plugins/v3/tscap/auth"
)

// MethodGates stores the authorization gates for each method.
// Keys are service name -> method name -> gate configuration.
var MethodGates = make(map[string]map[string]*auth.Gate)

func init() {
	codegen.RegisterPlugin("tscap", "gen", nil, Generate)
}

// Generate modifies the generated server code to add Tailscale capability
// authorization middleware.
func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, file := range files {
		if filepath.Base(file.Path) != "server.go" {
			continue
		}

		for _, s := range file.Section("server-handler") {
			data, ok := s.Data.(*goahttp.EndpointData)
			if !ok {
				continue
			}

			var gateDefined bool
			var gate *auth.Gate
			if serviceGates, ok := MethodGates[data.ServiceName]; ok {
				if g, ok := serviceGates[data.Method.Name]; ok {
					gateDefined = true
					gate = g
				}
			}

			codegen.AddImport(file.SectionTemplates[0],
				&codegen.ImportSpec{Path: "goa.design/plugins/v3/tscap/auth"},
			)

			if gateDefined {
				file.SectionTemplates = append(file.SectionTemplates, &codegen.SectionTemplate{
					Name:   "tscap-middleware",
					Source: definedGate,
					Data:   gate,
				})
			} else {
				file.SectionTemplates = append(file.SectionTemplates, &codegen.SectionTemplate{
					Name:   "tscap-middleware",
					Source: defaultGate,
					Data: auth.Gate{
						MethodName: data.Method.Name,
					},
				})
			}

			s.Source = strings.Replace(
				s.Source,
				`mux.Handle("{{ .Verb }}", "{{ .Path }}", f)`,
				`mux.Handle("{{ .Verb }}", "{{ .Path }}", `+data.Method.Name+`Tscap(f))`,
				1,
			)
		}
	}
	return files, nil
}

const defaultGate = `
{{ printf "for authorization based on Tailscale app capabilities" | comment }}
func {{ .MethodName }}Tscap(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := auth.ParseCapabilities(w, r); !ok {
			return
		}
		handler(w, r)
	}
}
`

const definedGate = `
{{ printf "for authorization based on Tailscale app capabilities" | comment }}
func {{ .MethodName }}Tscap(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{{- if .AllowAnonymous }}
		if r.Header.Get(auth.Header) == "" {
			handler(w, r)
			return
		}
		{{- end }}
		{{- if .Requirement }}
		caps, ok := auth.ParseCapabilities(w, r)
		if !ok {
			return
		}
		req := auth.Requirement{
			Capability: ` + "`{{ .Requirement.Capability }}`" + `,
			Action:     ` + "`{{ .Requirement.Action }}`" + `,
			Resource:   ` + "`{{ .Requirement.Resource }}`" + `,
		}
		if !auth.Check(w, caps, req) {
			return
		}
		{{- else if not .AllowAnonymous }}
		if _, ok := auth.ParseCapabilities(w, r); !ok {
			return
		}
		{{- end }}
		handler(w, r)
	}
}
`
