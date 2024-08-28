package dsl

import (
	"regexp"

	"goa.design/goa/v3/eval"
	goaexpr "goa.design/goa/v3/expr"
	"goa.design/plugins/v3/arnz"
	"goa.design/plugins/v3/arnz/caller"
)

// AllowArnsMatching accepts regex patterns to match caller ARNs.
func AllowArnsMatching(regex ...string) {
	rule := get()
	for _, given := range regex {
		regexp.MustCompile(given)
		rule.AllowArnsMatching = append(rule.AllowArnsMatching, given)
	}
}

// AllowUnsigned will allow callers skipping the gateway to bypass ARN checks.
func AllowUnsigned() {
	rule := get()
	rule.AllowUnsigned = true
}

func get() *caller.Gate {
	if m, ok := eval.Current().(*goaexpr.MethodExpr); ok {
		if _, exists := arnz.MethodGates[m.Service.Name]; !exists {
			arnz.MethodGates[m.Service.Name] = make(map[string]*caller.Gate)
		}

		if _, exists := arnz.MethodGates[m.Service.Name][m.Name]; !exists {
			arnz.MethodGates[m.Service.Name][m.Name] = &caller.Gate{
				MethodName: m.Name,
			}
		}

		return arnz.MethodGates[m.Service.Name][m.Name]
	}

	eval.IncompatibleDSL()
	return nil
}
