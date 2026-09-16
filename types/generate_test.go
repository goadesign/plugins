package types

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"

	"goa.design/plugins/v3/types/testdata"
)

var update = flag.Bool("update", false, "update golden files")

func TestTypes(t *testing.T) {
	cases := []struct {
		Name string
		DSL  func()
	}{
		{"empty", testdata.Empt},
		{"noval", testdata.NoValidation},
		{"required", testdata.Require},
		{"validation", testdata.Validation},
		{"validator_collision", testdata.ValidatorCollision},
		{"multiple", testdata.Multiple},
		{"alias", testdata.Alias},
		{"example", testdata.Exampl},
		{"array", testdata.Array},
		{"recArray", testdata.ArrayArray},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			root := codegen.RunDSL(t, c.DSL)
			fs, err := Generate("generated.local/gen", []eval.Root{root}, nil)
			require.NoError(t, err)
			require.NotEmpty(t, fs)
			require.NotEmpty(t, fs[0].SectionTemplates)
			filename, err := fs[0].Render(t.TempDir())
			require.NoError(t, err)
			got, err := os.ReadFile(filename)
			require.NoError(t, err)
			golden := filepath.Join("testdata", fmt.Sprintf("%s.go_", c.Name))
			if *update {
				require.NoError(t, os.WriteFile(golden, got, 0644))
				return
			}
			expected, err := os.ReadFile(golden)
			require.NoError(t, err)
			assert.Equal(t, string(expected), string(got))
		})
	}
}

// TestGeneratePreservesEvaluatedServices catches generators that add a fake
// service to the shared design while borrowing Goa's service templates.
func TestGeneratePreservesEvaluatedServices(t *testing.T) {
	root := codegen.RunDSL(t, testdata.Empt)
	services := append([]*expr.ServiceExpr(nil), root.Services...)

	_, err := Generate("generated.local/gen", []eval.Root{root}, nil)

	require.NoError(t, err)
	assert.Equal(t, services, root.Services)
}
