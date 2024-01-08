package otel

import (
	"path/filepath"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
)

// Register the plugin Generator functions.
func init() {
	codegen.RegisterPluginLast("otel", "gen", nil, Generate)
}

// Generate generates the call to otelhttp.WithRouteTag
func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, f := range files {
		if filepath.Base(f.Path) == "server.go" {
			for _, s := range f.SectionTemplates {
				if s.Name == "server-handler" {
					s.Source = strings.Replace(
						s.Source,
						"{{- range Routes }}",
						`f = otelhttp.WithRouteTag("{{ .Path }}", f).ServeHTTP
						 {{- range Routes }}`,
						1,
					)
				}
			}
			imports := f.SectionTemplates[0].Data.(map[string]any)["Imports"].([]*codegen.ImportSpec)
			imports = append(imports, &codegen.ImportSpec{
				Path: "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
				Name: "otelhttp",
			})
			f.SectionTemplates[0].Data.(map[string]any)["Imports"] = imports
		}
	}
	return files, nil
}
