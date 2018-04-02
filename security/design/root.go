package design

import (
	goadesign "goa.design/goa/design"
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

// WalkSets iterates over the security requirements.
func (r *RootExpr) WalkSets(walk eval.SetWalker) {
	expressions := []eval.Expression{}
	for _, svc := range goadesign.Root.Services {
		for _, m := range svc.Methods {
			reqs := Requirements(m)
			for _, req := range reqs {
				expressions = append(expressions, req)
			}
		}
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

// Requirements returns the security requirements for the given method.
func Requirements(m *goadesign.MethodExpr) []*EndpointSecurityExpr {
	var (
		sexpr      []*EndpointSecurityExpr
		found      bool
		noSecurity bool
	)
	for _, es := range Root.EndpointSecurity {
		if es.Method == m {
			// Handle special case of no security
			for _, s := range es.Schemes {
				if s.Kind == NoKind {
					noSecurity = true
					break
				}
			}
			if !noSecurity {
				sexpr = append(sexpr, es)
				found = true
			}
		}
	}
	if found || noSecurity {
		return sexpr
	}
	for _, ss := range Root.ServiceSecurity {
		if ss.Service == m.Service {
			sexpr = append(sexpr, &EndpointSecurityExpr{SecurityExpr: ss.SecurityExpr, Method: m})
			found = true
		}
	}
	if found {
		return sexpr
	}
	for _, as := range Root.APISecurity {
		sexpr = append(sexpr, &EndpointSecurityExpr{SecurityExpr: as, Method: m})
	}
	return sexpr
}
