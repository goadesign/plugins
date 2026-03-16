package dsl

import (
	"goa.design/goa/v3/eval"
	goaexpr "goa.design/goa/v3/expr"
	"goa.design/plugins/v3/tscap"
	"goa.design/plugins/v3/tscap/auth"
)

// Require declares that the method requires a Tailscale app capability
// with the specified action and resource.
func Require(capability, action, resource string) {
	gate := get()
	if gate == nil {
		return
	}

	if capability == "" {
		eval.ReportError("capability name cannot be empty in Require")
		return
	}
	if action == "" {
		eval.ReportError("action cannot be empty in Require")
		return
	}
	if resource == "" {
		eval.ReportError("resource cannot be empty in Require")
		return
	}

	gate.Requirement = &auth.Requirement{
		Capability: capability,
		Action:     action,
		Resource:   resource,
	}
}

// AllowAnonymous marks the method as not requiring any capability check.
func AllowAnonymous() {
	gate := get()
	if gate == nil {
		return
	}
	gate.AllowAnonymous = true
}

func get() *auth.Gate {
	m, ok := eval.Current().(*goaexpr.MethodExpr)
	if !ok {
		eval.IncompatibleDSL()
		return nil
	}

	if _, exists := tscap.MethodGates[m.Service.Name]; !exists {
		tscap.MethodGates[m.Service.Name] = make(map[string]*auth.Gate)
	}

	if _, exists := tscap.MethodGates[m.Service.Name][m.Name]; !exists {
		tscap.MethodGates[m.Service.Name][m.Name] = &auth.Gate{
			MethodName: m.Name,
		}
	}

	return tscap.MethodGates[m.Service.Name][m.Name]
}
