// Package otel was a Goa plugin that instrumented HTTP handlers with
// otelhttp.WithRouteTag to set the http.route attribute on spans and metrics.
//
// Deprecated: As of Goa v3.x.x the default HTTP muxer sets r.Pattern on every
// matched request, which otelhttp (v0.65.0+) reads automatically to tag spans
// and metrics with the matched route. This plugin is no longer necessary and
// will be removed in a future release. Remove the blank import from your
// design package:
//
//	import _ "goa.design/plugins/v3/otel" // ← delete this line
package otel

import (
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
)

// Register the plugin Generator functions.
func init() {
	codegen.RegisterPluginLast("otel", "gen", nil, Generate)
}

// Generate is a no-op kept for backward compatibility. The Goa HTTP muxer now
// sets r.Pattern on every request, making explicit route tagging unnecessary.
//
// Deprecated: Remove the otel plugin import from your design package.
func Generate(_ string, _ []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	return files, nil
}
