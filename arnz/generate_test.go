package arnz_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
	d "goa.design/plugins/v3/arnz/testdata"
)

func TestWrongScope(t *testing.T) {
	assert.Panics(t, func() {
		codegen.RunDSL(t, d.WrongScope)
	})
}

func TestBadMatcher(t *testing.T) {
	eval.Context = &eval.DSLContext{}
	serviceExpr := &expr.ServiceExpr{}
	eval.Execute(d.BadMatcher, serviceExpr)
	assert.NotNil(t, eval.Context.Errors)
}
