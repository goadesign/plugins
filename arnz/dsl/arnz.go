package dsl

import (
	"goa.design/goa/v3/eval"
	goaexpr "goa.design/goa/v3/expr"
	"goa.design/plugins/v3/arnz"
	"goa.design/plugins/v3/arnz/caller"
)

// AllowArnsLike will check to see if the given substring is contained in the caller ARN.
func AllowArnsLike(arns ...string) {
	rule := get()
	rule.AllowArnsLike = arns

	if len(rule.AllowArnsMatching) > 0 {
		eval.ReportError("arnz.AllowArnsMatching and arnz.AllowArnsLike cannot be used together")
	}
}

// AllowArnsMatching will check to see if the given string is an exact match to the caller ARN.
func AllowArnsMatching(arns ...string) {
	rule := get()
	rule.AllowArnsMatching = arns

	if len(rule.AllowArnsLike) > 0 {
		eval.ReportError("arnz.AllowArnsMatching and arnz.AllowArnsLike cannot be used together")
	}
}

// AllowUnsigned will allow callers skipping the gateway to bypass ARN checks.
func AllowUnsigned() {
	rule := get()
	rule.AllowUnsigned = true
}

func get() *caller.Gate {
	if m, ok := eval.Current().(*goaexpr.MethodExpr); ok {
		if _, exists := arnz.MethodRules[m.Service.Name]; !exists {
			arnz.MethodRules[m.Service.Name] = make(map[string]*caller.Gate)
		}

		if _, exists := arnz.MethodRules[m.Service.Name][m.Name]; !exists {
			arnz.MethodRules[m.Service.Name][m.Name] = &caller.Gate{
				MethodName: m.Name,
			}
		}

		return arnz.MethodRules[m.Service.Name][m.Name]
	}

	eval.IncompatibleDSL()
	return nil
}
