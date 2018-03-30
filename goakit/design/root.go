package design

import (
	"goa.design/goa/eval"
	httpdesign "goa.design/goa/http/design"
	security "goa.design/plugins/security/design"
)

// Root is the design root expression.
var Root = new(RootExpr)

type (
	// RootExpr is the Root expression for goakit plugin
	RootExpr struct{}
)

// Register design root with eval engine.
func init() {
	eval.Register(Root)
}

// EvalName returns the name used in error messages.
func (r *RootExpr) EvalName() string {
	return "goakit plugin"
}

// WalkSets iterates over the schemes.
func (r *RootExpr) WalkSets(walk eval.SetWalker) {}

// DependsOn tells the eval engine to run the goa HTTP and security plugin first.
func (r *RootExpr) DependsOn() []eval.Root {
	return []eval.Root{httpdesign.Root, security.Root}
}

// Packages returns the import path to the Go packages that make
// up the DSL. This is used to skip frames that point to files
// in these packages when computing the location of errors.
func (r *RootExpr) Packages() []string {
	return []string{}
}
