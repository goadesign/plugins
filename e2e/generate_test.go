package e2e_test

import (
	"bytes"
	"flag"
	"log"
	"testing"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	"goa.design/plugins/v3/e2e"
	"goa.design/plugins/v3/e2e/testdata"
)

var update = flag.Bool("update", false, "update golden files")

func TestE2E(t *testing.T) {
	cases := []struct {
		Name string
		DSL  func()
	}{
		{"api-only", testdata.NoPayloadNoReturn},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			root := codegen.RunDSL(t, c.DSL)
			fs, err := e2e.Generate("", []eval.Root{root}, nil)

			if err != nil {
				t.Fatal(err)
			}
			if len(fs) == 0 {
				t.Fatalf("got 0 files, expected 1")
			}
			for _, f := range fs {
				for _, st := range f.SectionTemplates {
					log.Printf("::: %+v :::", st.Name)
				}
			}
			if len(fs[0].SectionTemplates) < 1 {
				t.Fatalf("got 0 sections, expected 1")
			}
			var buf bytes.Buffer
			if err := fs[0].SectionTemplates[0].Write(&buf); err != nil {
				t.Fatal(err)
			}
		})
	}
}
