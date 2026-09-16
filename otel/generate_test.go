package otel

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
	"goa.design/goa/v3/expr"
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
			generation, err := codegen.NewGeneration("generated.local/gen", []eval.Root{root})
			require.NoError(t, err)
			servicePlan, err := service.NewPlan(root, generation, expr.NewExampleGenerator(root.API.RandomizerFactory))
			require.NoError(t, err)
			plans, err := httpcodegen.NewPlans(generation, httpcodegen.PlanInput{Root: root, Service: servicePlan})
			require.NoError(t, err)
			require.NoError(t, generation.Freeze())
			require.NoError(t, servicePlan.Link())
			require.NoError(t, plans[0].Link())
			serverFiles := plans[0].ServerFiles()
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
			assert.Equal(t, strings.TrimLeft(buf.String(), "\n"), strings.TrimLeft(string(expected), "\n"))
		})
	}
}
