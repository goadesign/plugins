package dsl

import (
	"goa.design/goa/v3/eval"
	goaexpr "goa.design/goa/v3/expr"
	"goa.design/plugins/v3/arnz"
)

func AllowArnsLike(arns ...string) {
	if m, ok := eval.Current().(*goaexpr.MethodExpr); ok {
		if _, exists := arnz.MethodARNs[m.Service.Name]; !exists {
			arnz.MethodARNs[m.Service.Name] = make(map[string][]string)
		}
		arnz.MethodARNs[m.Service.Name][m.Name] = append(arnz.MethodARNs[m.Service.Name][m.Name], arns...)
	} else {
		eval.IncompatibleDSL()
	}
}
