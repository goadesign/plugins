// This file checks compatibility entry points that existing applications call
// without access to Goa's internal generation plan.
package testing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	"goa.design/plugins/v3/testing/codegen/testdata"
)

func TestReleasedCallbacksPlanServiceDataWithoutInputFiles(t *testing.T) {
	root := codegen.RunDSL(t, testdata.WithPayloadDSL)

	generated, err := Generate("generated.local/gen", []eval.Root{root}, nil)
	require.NoError(t, err)
	assert.Len(t, generated, 4)

	examples, err := GenerateExample("generated.local/gen", []eval.Root{root}, nil)
	require.NoError(t, err)
	assert.Len(t, examples, 2)
}

func TestReleasedCallbacksRejectEmptyGeneratedPackage(t *testing.T) {
	root := codegen.RunDSL(t, testdata.WithPayloadDSL)

	_, err := Generate("", []eval.Root{root}, nil)
	require.Error(t, err)
	_, err = GenerateExample("", []eval.Root{root}, nil)
	require.Error(t, err)
}
