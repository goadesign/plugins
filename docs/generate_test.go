package docs_test

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

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
			if err := fs[0].SectionTemplates[0].Write(&buf); err != nil {
				t.Fatal(err)
			}
			golden := filepath.Join("testdata", fmt.Sprintf("%s.json", c.Name))
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
