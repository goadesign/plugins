package arnz_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	httpcodegen "goa.design/goa/v3/http/codegen"
	d "goa.design/plugins/v3/arnz/testdata"
)

func TestConflicting(t *testing.T) {
	assert.Panics(t, func() {
		httpcodegen.RunHTTPDSL(t, d.WrongScope)
	})
}
