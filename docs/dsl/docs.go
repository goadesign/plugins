package dsl

import (
	"goa.design/plugins/v3/docs"
	"goa.design/plugins/v3/docs/expr"
)

// init registers the docs plugin when importing the DSL.
func init() { docs.Register() }

// UseJSONTags configures the docs plugin to use JSON struct tags declared via
// Meta("struct:tag:json", ...) as field names in generated docs. This setting
// affects definitions, payloads, results, and error schemas and examples.
func UseJSONTags() { expr.Root.UseJSONTags = true }

// InlineRefs configures the docs plugin to inline referenced schemas in JSON
// Schema output where possible. Cycles are preserved by keeping $ref.
func InlineRefs() { expr.Root.InlineRefs = true }
