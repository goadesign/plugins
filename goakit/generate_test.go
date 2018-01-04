package goakit

import (
	"bytes"
	"strings"
	"testing"

	"goa.design/goa/codegen"
	"goa.design/goa/codegen/generator"
	goadesign "goa.design/goa/design"
	"goa.design/goa/eval"
	httpcodegen "goa.design/goa/http/codegen"
	httpdesign "goa.design/goa/http/design"
	"goa.design/plugins/goakit/testdata"
)

func TestGenerate(t *testing.T) {
	// Map all files containing goaEndpoint indexed by file path
	goaEndpointFiles := map[string]bool{}
	cases := map[string]struct {
		DSL      func()
		ExpFiles int
	}{
		"multi-endpoints": {testdata.MultiEndpointDSL, 3},
		"multi-services":  {testdata.MultiServiceDSL, 6},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			httpcodegen.RunHTTPDSL(t, c.DSL)
			roots := []eval.Root{goadesign.Root, httpdesign.Root}
			files := generateFiles(t, roots)
			// Before state: Collect all files with goa endpoint.
			for _, f := range files {
				goaEndpointFiles[f.Path] = containsGoaEndpoint(f)
			}
			newFiles, err := Generate("", roots, files)
			if err != nil {
				t.Fatalf("generate error: %v", err)
			}
			newFilesCount := len(newFiles) - len(files)
			if newFilesCount != c.ExpFiles {
				t.Errorf("invalid code: number of new files expected %d, got %d", c.ExpFiles, newFilesCount)
			}
			// After state: files with goa endpoint should be replaced by gokit endpoint
			for _, f := range files {
				if goaEndpointFiles[f.Path] {
					if containsGoaEndpoint(f) {
						t.Errorf("file %s still has goa endpoints", f.Path)
					}
					buf := new(bytes.Buffer)
					if err := f.SectionTemplates[0].Write(buf); err != nil {
						t.Fatalf("error writing section in file %s", f.Path)
					}
					code := buf.String()
					if !strings.Contains(code, "github.com/go-kit/kit/endpoint") {
						t.Errorf("go-kit not imported in file %s:\n%s", f.Path, code)
					}
				}
			}

		})
	}
}

func TestExample(t *testing.T) {
	cases := map[string]struct {
		DSL      func()
		ExpFiles int
	}{
		"mixed":          {testdata.MixedDSL, 2},
		"multi-services": {testdata.MultiServiceDSL, 3},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			httpcodegen.RunHTTPDSL(t, c.DSL)
			roots := []eval.Root{goadesign.Root, httpdesign.Root}
			files := generateFiles(t, roots)
			newFiles, err := Example("", roots, files)
			if err != nil {
				t.Fatalf("examples generate error: %v", err)
			}
			newFileCount := len(newFiles) - len(files)
			if newFileCount != c.ExpFiles {
				t.Errorf("invalid code, number of new files expected %d, got %d", c.ExpFiles, newFileCount)
			}
		})
	}
}

func generateFiles(t *testing.T, roots []eval.Root) []*codegen.File {
	files, err := generator.Service("", roots)
	if err != nil {
		t.Fatalf("error in code generation: %v", err)
	}
	httpFiles, err := generator.Transport("", roots)
	if err != nil {
		t.Fatalf("error in HTTP code generation: %v", err)
	}
	files = append(files, httpFiles...)
	return files
}

func containsGoaEndpoint(f *codegen.File) bool {
	for _, s := range f.SectionTemplates {
		if goaEndpointRegexp.MatchString(s.Source) {
			return true
		}
	}
	return false
}
