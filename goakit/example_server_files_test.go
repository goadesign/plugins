package goakit

import (
	"path/filepath"
	"testing"

	"goa.design/goa/codegen"
	httpcodegen "goa.design/goa/http/codegen"
	httpdesign "goa.design/goa/http/design"
	"goa.design/plugins/goakit/testdata"
)

func TestExampleServerFiles(t *testing.T) {
	cases := map[string]struct {
		DSL           func()
		SvcStructCode string
		ExpEndpoints  int
	}{
		"multi-endpoints": {testdata.MultiEndpointDSL, testdata.MultiEndpointServiceStructCode, 2},
		"mixed":           {testdata.MixedDSL, testdata.MixedServiceStructCode, 1},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			httpcodegen.RunHTTPDSL(t, c.DSL)
			fs := ExampleServerFiles("", httpdesign.Root)
			for _, f := range fs {
				if filepath.Base(f.Path) == "main.go" {
					continue
				}
				endpointSections := f.Section("goakit-dummy-endpoint")
				if len(endpointSections) != c.ExpEndpoints {
					t.Errorf("goakit-dummy-endpoint: invalid section code: expected %d, got %d sections", c.ExpEndpoints, len(endpointSections))
				}

				sections := f.Section("goakit-dummy-service-struct")
				if len(sections) != 1 {
					t.Errorf("goakit-dummy-service-struct: invalid section code: expected 1, got %d sections", len(sections))
				}
				code := codegen.SectionCode(t, sections[0])
				if code != c.SvcStructCode {
					t.Errorf("invalid code, got:\n%s\ngot vs. expected:\n%s", code, codegen.Diff(t, code, c.SvcStructCode))
				}
			}
		})
	}
}
