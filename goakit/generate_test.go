package goakit

import (
	"bytes"
	"strings"
	"testing"

	"goa.design/goa/codegen"
	"goa.design/goa/codegen/generator"
	"goa.design/goa/eval"
	"goa.design/goa/expr"
	httpcodegen "goa.design/goa/http/codegen"
	"goa.design/plugins/goakit/testdata"
)

func TestGenerate(t *testing.T) {
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
			roots := []eval.Root{expr.Root}
			files := generateFiles(t, roots)
			newFiles, err := Generate("", roots, files)
			if err != nil {
				t.Fatalf("generate error: %v", err)
			}
			newFilesCount := len(newFiles) - len(files)
			if newFilesCount != c.ExpFiles {
				t.Errorf("invalid code: number of new files expected %d, got %d", c.ExpFiles, newFilesCount)
			}
		})
	}
}

func TestGoakitify(t *testing.T) {
	// Map all files containing goaEndpoint indexed by file path
	goaEndpointFiles := map[string]bool{}
	cases := map[string]func(){
		"multi-endpoints": testdata.MultiEndpointDSL,
		"multi-services":  testdata.MultiServiceDSL,
	}
	for name, dsl := range cases {
		t.Run(name, func(t *testing.T) {
			httpcodegen.RunHTTPDSL(t, dsl)
			roots := []eval.Root{expr.Root}
			files := generateFiles(t, roots)
			newFiles, err := Goakitify("", roots, files)
			if err != nil {
				t.Fatalf("generate error: %v", err)
			}
			// Before state: Collect all files with goa endpoint.
			for _, f := range newFiles {
				goaEndpointFiles[f.Path] = containsGoaEndpoint(f)
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

func TestGoakitifyExample(t *testing.T) {
	cases := map[string]struct {
		DSL  func()
		Code map[string]string
	}{
		"mixed": {
			DSL: testdata.MixedDSL,
			Code: map[string]string{
				"service-main-logger":      testdata.MixedMainLoggerCode,
				"service-main-middleware":  testdata.MixedMainMiddlewareCode,
				"service-main-server-init": testdata.MixedMainServerInitCode,
			},
		},
		"multi-services": {
			DSL: testdata.MultiServiceDSL,
			Code: map[string]string{
				"server-http-init": testdata.MultiServicesServerInitCode,
			},
		},
		"with-error": {
			DSL: testdata.WithErrorDSL,
			Code: map[string]string{
				"server-http-init": testdata.WithErrorServerInitCode,
			},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			httpcodegen.RunHTTPDSL(t, c.DSL)
			roots := []eval.Root{expr.Root}
			files := generateExamples(t, roots)
			files, err := GoakitifyExample("", roots, files)
			if err != nil {
				t.Fatalf("examples generate error: %v", err)
			}
			for _, f := range files {
				if containsStdlibLogger(f) {
					t.Errorf("file %s still has stdlib logger instances", f.Path)
				}
				for _, s := range f.SectionTemplates {
					if expCode, ok := c.Code[s.Name]; ok {
						buf := new(bytes.Buffer)
						if err := s.Write(buf); err != nil {
							t.Fatalf("error writing section in file %s", f.Path)
						}
						code := buf.String()
						code = codegen.FormatTestCode(t, "package foo\nfunc example(){"+code+"}")
						if code != expCode {
							t.Errorf("invalid code for %s, got:\n%s\ngot vs. expected:\n%s", s.Name, code, codegen.Diff(t, code, expCode))
						}
					}
				}
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

func generateExamples(t *testing.T, roots []eval.Root) []*codegen.File {
	files, err := generator.Example("", roots)
	if err != nil {
		t.Fatalf("error in code generation: %v", err)
	}
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

func containsStdlibLogger(f *codegen.File) bool {
	for _, s := range f.SectionTemplates {
		if strings.Contains(s.Source, "*log.Logger") {
			return true
		}
	}
	return false
}
