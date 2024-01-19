package otel

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
	httpcodegen "goa.design/goa/v3/http/codegen"
	"goa.design/plugins/v3/otel/testdata"
)

var update = flag.Bool("update", false, "update golden files")

func TestOtel(t *testing.T) {
	cases := []struct {
		Name string
		DSL  func()
	}{
		{"one route", testdata.OneRoute},
		{"multiple routes", testdata.MultipleRoutes},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			root := codegen.RunDSL(t, c.DSL)
			serverFiles := httpcodegen.ServerFiles("gen", root)
			require.Len(t, serverFiles, 2)
			fs, err := Generate("", []eval.Root{root}, serverFiles)
			assert.NoError(t, err)
			require.Len(t, fs, 2)
			sections := fs[0].Section("server-handler")
			require.Len(t, sections, 1)
			section := sections[0]
			var buf bytes.Buffer
			assert.NoError(t, section.Write(&buf))
			golden := filepath.Join("testdata", fmt.Sprintf("%s.golden", c.Name))
			if *update {
				assert.NoError(t, os.WriteFile(golden, buf.Bytes(), 0644))
			}
			expected, _ := os.ReadFile(golden)
			assert.Equal(t, buf.String(), string(expected))
		})
	}
}
