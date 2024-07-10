package dsl

import (
	"goa.design/goa/v3/eval"
	goaexpr "goa.design/goa/v3/expr"
	"goa.design/plugins/v3/arnz"
)

func AllowArnsLike(arns ...string) {
	if m, ok := eval.Current().(*goaexpr.MethodExpr); ok {
		arnz.MethodARNs[m.Name] = append(arnz.MethodARNs[m.Name], arns...)
	} else {
		eval.IncompatibleDSL()
	}
}
