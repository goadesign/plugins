package arnz_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/http/codegen"
	d "goa.design/plugins/v3/arnz/testdata"
)

func TestWrongScope(t *testing.T) {
	assert.Panics(t, func() {
		codegen.RunHTTPDSL(t, d.WrongScope)
	})
}

func TestConflicting(t *testing.T) {
	assert.True(t, eval.Execute(d.Conflicting, nil), eval.Context.Error())
}
