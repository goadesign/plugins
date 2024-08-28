package arnz

import (
	"path/filepath"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	goahttp "goa.design/goa/v3/http/codegen"
	"goa.design/plugins/v3/arnz/auth"
)

var MethodGates = make(map[string]map[string]*auth.Gate)

func init() {
	codegen.RegisterPlugin("arnz", "gen", nil, Generate)
}

func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, file := range files {
		if filepath.Base(file.Path) == "server.go" {
			for _, s := range file.Section("server-handler") {
				if data, ok := s.Data.(*goahttp.EndpointData); ok {
					if _, ok := MethodGates[data.ServiceName]; ok {
						if _, ok := MethodGates[data.ServiceName][data.Method.Name]; ok {
							codegen.AddImport(file.SectionTemplates[0],
								&codegen.ImportSpec{Path: "encoding/json"},
								&codegen.ImportSpec{Path: "strings"},
								&codegen.ImportSpec{Path: "github.com/aws/aws-lambda-go/events"},
								&codegen.ImportSpec{Path: "goa.design/plugins/v3/arnz/auth"},
							)

							file.SectionTemplates = append(file.SectionTemplates, &codegen.SectionTemplate{
								Name:   "arnz-middleware",
								Source: definedGate,
								Data:   MethodGates[data.ServiceName][data.Method.Name],
							})

							s.Source = strings.Replace(
								s.Source,
								`mux.Handle("{{ .Verb }}", "{{ .Path }}", f)`,
								`mux.Handle("{{ .Verb }}", "{{ .Path }}", `+data.Method.Name+`Arnz(f))`,
								1,
							)
						} else {
							codegen.AddImport(file.SectionTemplates[0],
								&codegen.ImportSpec{Path: "encoding/json"},
								&codegen.ImportSpec{Path: "github.com/aws/aws-lambda-go/events"},
								&codegen.ImportSpec{Path: "goa.design/plugins/v3/arnz/auth"},
							)

							file.SectionTemplates = append(file.SectionTemplates, &codegen.SectionTemplate{
								Name:   "arnz-middleware",
								Source: defaultGate,
								Data: auth.Gate{
									MethodName: data.Method.Name,
								},
							})

							s.Source = strings.Replace(
								s.Source,
								`mux.Handle("{{ .Verb }}", "{{ .Path }}", f)`,
								`mux.Handle("{{ .Verb }}", "{{ .Path }}", `+data.Method.Name+`Arnz(f))`,
								1,
							)
						}
					}
				}
			}
		}
	}
	return files, nil
}

const defaultGate = `
{{ printf "for authorization based on AWS ARNs" | comment }}
 func {{ .MethodName }}Arnz(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, pass := auth.Authenticate(w, r); !pass {
			return
		}
		handler(w, r)
	}
}
`

const definedGate = `
{{ printf "for authorization based on AWS ARNs" | comment }}
func {{ .MethodName }}Arnz(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { 
		{{- if .AllowUnsigned }}
		if auth.IsUnsigned(r) {
			handler(w, r)
			return
		}
		{{- end }}
		{{- if (gt (len .AllowArnsMatching) 0) }}
		callerArn, pass := auth.Authenticate(w, r)
		if !pass {
			return
		}
		allowArnsMatching := []string{
			{{- range .AllowArnsMatching }}
			` + "`{{ . }}`" + `,
			{{- end }}
		}
		if !auth.Authorize(w, *callerArn, allowArnsMatching) {
			return
		}
		{{- else }}
		if _, pass := auth.Authenticate(w, r); !pass {
			return
		}
		{{- end }}
		handler(w, r)
	}
}
`
