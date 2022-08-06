package types

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
			if err != nil {
				t.Fatal(err)
			}
			if len(fs) == 0 {
				t.Fatalf("got 0 files, expected 1")
			}
			if len(fs[0].SectionTemplates) == 0 {
				t.Fatalf("got 0 sections, expected 1")
			}
			var buf bytes.Buffer
			for _, s := range fs[0].SectionTemplates {
				if err := s.Write(&buf); err != nil {
					t.Fatal(err)
				}
			}
			got := buf.String()[strings.Index(buf.String(), "\n")+1:]
			golden := filepath.Join("testdata", fmt.Sprintf("%s.go_", c.Name))
			if *update {
				os.WriteFile(golden, buf.Bytes(), 0644)
			}
			expected, _ := os.ReadFile(golden)
			if got != string(expected) {
				t.Errorf("invalid content for %s compared to %s: got\n%s\ngot vs. expected:\n%s",
					fs[0].Path, golden, got, codegen.Diff(t, got, string(expected)))
			}
		})
	}
}
