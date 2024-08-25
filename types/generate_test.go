package types

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/eval"

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
		{"multiple", testdata.Multiple},
		{"alias", testdata.Alias},
		{"example", testdata.Exampl},
		{"array", testdata.Array},
		{"recArray", testdata.ArrayArray},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			service.Services = make(service.ServicesData)
			root := codegen.RunDSL(t, c.DSL)
			fs, err := Generate("", []eval.Root{root}, nil)
			require.NoError(t, err)
			require.NotEmpty(t, fs)
			require.NotEmpty(t, fs[0].SectionTemplates)
			var buf bytes.Buffer
			for _, s := range fs[0].SectionTemplates {
				require.NoError(t, s.Write(&buf))
			}
			got := buf.String()[strings.Index(buf.String(), "\n")+1:]
			golden := filepath.Join("testdata", fmt.Sprintf("%s.go_", c.Name))
			if *update {
				os.WriteFile(golden, buf.Bytes(), 0644)
				return
			}
			expected, err := os.ReadFile(golden)
			require.NoError(t, err)
			assert.Equal(t, string(expected), got)
		})
	}
}
