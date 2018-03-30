package design

import (
	"goa.design/goa/eval"
	httpdesign "goa.design/goa/http/design"
)

// Root is the design root expression.
var Root = new(RootExpr)

// PluginName is the name of the security plugin
const PluginName = "security"

type (
	// RootExpr keeps track of the security schemes defined in the design
	RootExpr struct {
		// Schemes list the registered security schemes.
		Schemes []*SchemeExpr
		// APISecurity list the API level security requirements.
		APISecurity []*SecurityExpr
		// ServiceSecurity list service level security requirements.
		ServiceSecurity []*ServiceSecurityExpr
		// EndpointSecurity list method level security requirements.
		EndpointSecurity []*EndpointSecurityExpr
	}
)

// Register design root with eval engine.
func init() {
	eval.Register(Root)
}

// EvalName returns the name used in error messages.
func (r *RootExpr) EvalName() string {
	return "security plugin"
}

// WalkSets iterates over the schemes.
func (r *RootExpr) WalkSets(walk eval.SetWalker) {
	expressions := make([]eval.Expression, len(r.Schemes))
	for i, s := range r.Schemes {
		expressions[i] = s
	}
	walk(expressions)
}

// DependsOn tells the eval engine to run the goa HTTP DSL first.
func (r *RootExpr) DependsOn() []eval.Root {
	return []eval.Root{httpdesign.Root}
}

// Packages returns the import path to the Go packages that make
// up the DSL. This is used to skip frames that point to files
// in these packages when computing the location of errors.
func (r *RootExpr) Packages() []string {
	return []string{"goa.design/plugins/security/http/dsl"}
}
