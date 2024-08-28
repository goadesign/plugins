package arnz

import (
	"path/filepath"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	goahttp "goa.design/goa/v3/http/codegen"
	"goa.design/plugins/v3/arnz/caller"
)

var MethodRules = make(map[string]map[string]*caller.Gate)

func init() {
	codegen.RegisterPlugin("arnz", "gen", nil, Generate)
}

func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, file := range files {
		if filepath.Base(file.Path) == "server.go" {
			for _, s := range file.Section("server-handler") {
				if data, ok := s.Data.(*goahttp.EndpointData); ok {
					if _, ok := MethodRules[data.ServiceName]; ok {
						if _, ok := MethodRules[data.ServiceName][data.Method.Name]; ok {
							codegen.AddImport(file.SectionTemplates[0],
								&codegen.ImportSpec{Path: "encoding/json"},
								&codegen.ImportSpec{Path: "strings"},
								&codegen.ImportSpec{Path: "github.com/aws/aws-lambda-go/events"},
								&codegen.ImportSpec{Path: "goa.design/plugins/v3/arnz/caller"},
							)

							file.SectionTemplates = append(file.SectionTemplates, &codegen.SectionTemplate{
								Name:   "arnz-middleware",
								Source: gated,
								Data:   MethodRules[data.ServiceName][data.Method.Name],
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
								&codegen.ImportSpec{Path: "goa.design/plugins/v3/arnz/caller"},
							)

							file.SectionTemplates = append(file.SectionTemplates, &codegen.SectionTemplate{
								Name:   "arnz-middleware",
								Source: ungated,
								Data: caller.Gate{
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

const ungated = `
{{ printf "for authorization based on AWS ARNs" | comment }}
 func {{ .MethodName }}Arnz(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, pass := caller.Extract(w, r); !pass {
			return
		}
		handler(w, r)
	}
}
`

const gated = `
{{ printf "for authorization based on AWS ARNs" | comment }}
func {{ .MethodName }}Arnz(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { 
		{{- if .AllowUnsigned }}
		if caller.IsUnsigned(r) {
			handler(w, r)
			return
		}
		{{- end }}

		{{- if or (gt (len .AllowArnsLike) 0) (gt (len .AllowArnsMatching) 0) }}
		callerArn, pass := caller.Extract(w, r)
		if !pass {
			return
		}

		{{- if gt (len .AllowArnsLike) 0 }}
		allowedArnsLike := []string{ {{- range .AllowArnsLike }} "{{ . }}", {{- end }} }
		if !caller.ArnLike(w, *callerArn, allowedArnsLike) {
			return
		}
		{{- end }}

		{{- if gt (len .AllowArnsMatching) 0 }}
		allowedArnsMatching := []string{ {{- range .AllowArnsMatching }} "{{ . }}", {{- end }} }
		if !caller.ArnMatch(w, *callerArn, allowedArnsMatching) {
			return
		}
		{{- end }}

		{{- else }}

		_, pass := caller.Extract(w, r)
		if !pass {
			return
		}

		{{- end }}

		handler(w, r)
	}
}
`
