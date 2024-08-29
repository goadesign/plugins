package dsl

import (
	"regexp"

	"goa.design/goa/v3/eval"
	goaexpr "goa.design/goa/v3/expr"
	"goa.design/plugins/v3/arnz"
	"goa.design/plugins/v3/arnz/auth"
)

// AllowArnsMatching accepts regex patterns to match caller ARNs.
func AllowArnsMatching(regex ...string) {
	rule := get()
	for _, given := range regex {
		_, err := regexp.Compile(given)
		if err != nil {
			eval.ReportError("invalid regex pattern in AllowArnsMatching: %s", given)
		}

		rule.AllowArnsMatching = append(rule.AllowArnsMatching, given)
	}
}

// AllowUnsigned will allow callers skipping the gateway to bypass ARN checks.
func AllowUnsignedCallers() {
	rule := get()
	rule.AllowUnsigned = true
}

func get() *auth.Gate {
	if m, ok := eval.Current().(*goaexpr.MethodExpr); ok {
		if _, exists := arnz.MethodGates[m.Service.Name]; !exists {
			arnz.MethodGates[m.Service.Name] = make(map[string]*auth.Gate)
		}

		if _, exists := arnz.MethodGates[m.Service.Name][m.Name]; !exists {
			arnz.MethodGates[m.Service.Name][m.Name] = &auth.Gate{
				MethodName: m.Name,
			}
		}

		return arnz.MethodGates[m.Service.Name][m.Name]
	}

	eval.IncompatibleDSL()
	return nil
}
