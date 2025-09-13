package expr

import (
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
)

// Root is the design root expression for the docs plugin.
var Root = &RootExpr{}

type (
	// RootExpr keeps track of docs plugin configuration toggles.
	RootExpr struct {
		// UseJSONTags instructs the docs generator to use JSON struct tags
		// specified via Meta ("struct:tag:json") as field names.
		UseJSONTags bool
		// InlineRefs instructs the docs generator to inline $ref schemas by
		// replacing them with copies of their referenced definitions where
		// possible. Cycles are preserved by leaving $ref in place when needed.
		InlineRefs bool
	}
)

// Register design root with eval engine.
func init() {
	_ = eval.Register(Root)
}

// EvalName returns the name used in error messages.
func (r *RootExpr) EvalName() string { return "Docs plugin" }

// WalkSets implements eval.Root. No-op; configuration is global.
func (*RootExpr) WalkSets(eval.SetWalker) {}

// DependsOn tells the eval engine to run the goa DSL first.
func (*RootExpr) DependsOn() []eval.Root { return []eval.Root{expr.Root} }

// Packages returns the import path to the Go packages that make up the DSL.
// This is used to skip frames that point to files in these packages when
// computing the location of errors.
func (*RootExpr) Packages() []string { return []string{"goa.design/plugins/v3/docs/dsl"} }
