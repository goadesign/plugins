package expr

import (
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
)

// Root is the design root expression.
var Root = &RootExpr{
	Description: map[eval.Expression]*I18nExpr{},
	Example:     map[eval.Expression]*I18nExpr{},
	Title:       map[eval.Expression]*I18nExpr{},
}

type (
	// RootExpr keeps track of the CORS origins defined in the design.
	RootExpr struct {
		// APIOrigins lists all the CORS definitions indexed by origin string
		// at the API level.
		Description map[eval.Expression]*I18nExpr
		Example     map[eval.Expression]*I18nExpr
		Title       map[eval.Expression]*I18nExpr
	}
)

// Register design root with eval engine.
func init() {
	eval.Register(Root)
}

// EvalName returns the name used in error messages.
func (r *RootExpr) EvalName() string {
	return "I18n plugin"
}

// WalkSets iterates over the API-level and service-level I18n definitions.
func (r *RootExpr) WalkSets(walk eval.SetWalker) {
	oexps := make(eval.ExpressionSet, 0, len(r.Description))
	for _, o := range r.Description {
		oexps = append(oexps, o)
	}
	walk(oexps)

}

// DependsOn tells the eval engine to run the goa DSL first.
func (r *RootExpr) DependsOn() []eval.Root {
	return []eval.Root{expr.Root}
}

// Packages returns the import path to the Go packages that make
// up the DSL. This is used to skip frames that point to files
// in these packages when computing the location of errors.
func (r *RootExpr) Packages() []string {
	return []string{"goa.design/plugins/v3/i18n/dsl"}
}
