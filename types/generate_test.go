package types

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"goa.design/goa/v3/codegen"
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
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
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
			golden := filepath.Join("testdata", fmt.Sprintf("%s.go_", c.Name))
			if *update {
				ioutil.WriteFile(golden, buf.Bytes(), 0644)
			}
			expected, _ := ioutil.ReadFile(golden)
			if buf.String() != string(expected) {
				t.Errorf("invalid content for %s: got\n%s\ngot vs. expected:\n%s",
					fs[0].Path, buf.String(), codegen.Diff(t, buf.String(), string(expected)))
			}
		})
	}
}
