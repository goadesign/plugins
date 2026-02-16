// Package otel is a Goa plugin that sets the http.route OpenTelemetry span
// attribute on generated HTTP handlers. It does this by wrapping each handler
// to call trace.SpanFromContext(r.Context()).SetAttributes(semconv.HTTPRoute(...))
// before the handler executes.
//
// Import the package in the service design with a blank identifier:
//
//	import _ "goa.design/plugins/v3/otel"
//
// When using otelhttp as a mux middleware (recommended), the Goa muxer sets
// r.Pattern before middlewares run, so otelhttp picks up the route
// automatically. In that case this plugin is not necessary:
//
//	mux := goahttp.NewMuxer()
//	mux.Use(otelhttp.NewMiddleware("service"))
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

// Generate modifies the generated HTTP server code to set the http.route
// OpenTelemetry span attribute on each handler. This ensures the attribute
// is present regardless of how otelhttp is wired (external handler wrapping
// or mux middleware).
func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, f := range files {
		if filepath.Base(f.Path) != "server.go" {
			continue
		}
		for _, s := range f.SectionTemplates {
			if s.Name == "server-handler" {
				s.Source = strings.Replace(
					s.Source,
					`mux.Handle("{{ .Verb }}", "{{ .Path }}", f)`,
					`mux.Handle("{{ .Verb }}", "{{ .Path }}", func(w http.ResponseWriter, r *http.Request) {
		trace.SpanFromContext(r.Context()).SetAttributes(semconv.HTTPRoute("{{ .Path }}"))
		f(w, r)
	})`,
					1,
				)
			}
		}
		imports := f.SectionTemplates[0].Data.(map[string]any)["Imports"].([]*codegen.ImportSpec)
		imports = append(imports,
			&codegen.ImportSpec{Path: "go.opentelemetry.io/otel/trace"},
			&codegen.ImportSpec{Path: "go.opentelemetry.io/otel/semconv/v1.38.0", Name: "semconv"},
		)
		f.SectionTemplates[0].Data.(map[string]any)["Imports"] = imports
	}
	return files, nil
}
