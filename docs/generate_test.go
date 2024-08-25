package docs_test

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	"goa.design/plugins/v3/docs"
	"goa.design/plugins/v3/docs/testdata"
)

var update = flag.Bool("update", false, "update golden files")

func TestDocs(t *testing.T) {
	cases := []struct {
		Name string
		DSL  func()
	}{
		{"api-only", testdata.APIOnly},
		{"no-payload-no-return", testdata.NoPayloadNoReturn},
		{"primitive-payload-no-return", testdata.PrimitivePayloadNoReturn},
		{"array-payload-no-return", testdata.ArrayPayloadNoReturn},
		{"map-payload-no-return", testdata.MapPayloadNoReturn},
		{"user-payload-no-return", testdata.UserPayloadNoReturn},
		{"no-payload-primitive-return", testdata.NoPayloadPrimitiveReturn},
		{"no-payload-array-return", testdata.NoPayloadArrayReturn},
		{"no-payload-map-return", testdata.NoPayloadMapReturn},
		{"no-payload-user-return", testdata.NoPayloadUserReturn},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			root := codegen.RunDSL(t, c.DSL)
			fs, err := docs.Generate("", []eval.Root{root}, nil)
			require.NoError(t, err)
			require.NotEmpty(t, fs)
			require.NotEmpty(t, fs[0].SectionTemplates)
			var buf bytes.Buffer
			require.NoError(t, fs[0].SectionTemplates[0].Write(&buf))
			golden := filepath.Join("testdata", fmt.Sprintf("%s.json", c.Name))
			if *update {
				os.WriteFile(golden, buf.Bytes(), 0644)
				return
			}
			expected, err := os.ReadFile(golden)
			require.NoError(t, err)
			assert.Equal(t, string(expected), buf.String())
		})
	}
}
