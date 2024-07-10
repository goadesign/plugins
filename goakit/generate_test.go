package goakit

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/generator"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
	httpcodegen "goa.design/goa/v3/http/codegen"
	"goa.design/plugins/v3/goakit/testdata"
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
			require.NoError(t, err)
			// Before state: Collect all files with goa endpoint.
			for _, f := range newFiles {
				goaEndpointFiles[f.Path] = containsGoaEndpoint(t, f)
			}
			// After state: files with goa endpoint should be replaced by gokit endpoint
			for _, f := range files {
				if goaEndpointFiles[f.Path] {
					assert.False(t, containsGoaEndpoint(t, f), "file %s still has goa endpoints", f.Path)
					buf := new(bytes.Buffer)
					require.NoError(t, f.SectionTemplates[0].Write(buf))
					code := buf.String()
					assert.Contains(t, code, "github.com/go-kit/kit/endpoint", "go-kit not imported in file %s", f.Path)
				}
			}
		})
	}
}

func TestGoakitifyExample(t *testing.T) {
	cases := map[string]struct {
		DSL     func()
		Code    map[string]string
		Imports []string
	}{
		"mixed": {
			DSL: testdata.MixedDSL,
			Code: map[string]string{
				"service-main-logger":      testdata.MixedMainLoggerCode,
				"service-main-middleware":  testdata.MixedMainMiddlewareCode,
				"service-main-server-init": testdata.MixedMainServerInitCode,
			},
			Imports: []string{"http/mixed_service/kitserver"},
		},
		"multi-services": {
			DSL: testdata.MultiServiceDSL,
			Code: map[string]string{
				"server-http-init": testdata.MultiServicesServerInitCode,
			},
			Imports: []string{"http/service1/kitserver", "http/service2/kitserver"},
		},
		"with-error": {
			DSL: testdata.WithErrorDSL,
			Code: map[string]string{
				"server-http-init": testdata.WithErrorServerInitCode,
			},
			Imports: []string{"http/with_error_service/kitserver"},
		},
		"goifyable": {
			DSL: testdata.GoifyableServiceDSL,
			Code: map[string]string{
				"server-http-init": testdata.GoifyableServerInitCode,
			},
			Imports: []string{"http/goifyable_service/kitserver"},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			httpcodegen.RunHTTPDSL(t, c.DSL)
			roots := []eval.Root{expr.Root}
			files := generateExamples(t, roots)
			files, err := GoakitifyExample("", roots, files)
			require.NoError(t, err)
			for _, f := range files {
				source, ok := containsStdlibLogger(t, f)
				require.False(t, ok, "file %s still has stdlib logger instances:\n%s", f.Path, source)
				for _, s := range f.SectionTemplates {
					expCode, ok := c.Code[s.Name]
					if !ok {
						continue
					}
					buf := new(bytes.Buffer)
					require.NoError(t, s.Write(buf))
					code := buf.String()
					code = codegen.FormatTestCode(t, "package foo\nfunc example(){"+code+"}")
					require.Equal(t, expCode, code, "invalid code for %s", s.Name)
				}
				if strings.HasSuffix(f.Path, "/http.go") && !strings.HasSuffix(f.Path, "cli/http.go") {
					for _, imp := range c.Imports {
						requireImport(t, f, imp)
					}
				}
			}
		})
	}
}

func generateFiles(t *testing.T, roots []eval.Root) []*codegen.File {
	files, err := generator.Service("", roots)
	require.NoError(t, err)
	httpFiles, err := generator.Transport("", roots)
	require.NoError(t, err)
	files = append(files, httpFiles...)
	return files
}

func generateExamples(t *testing.T, roots []eval.Root) []*codegen.File {
	files, err := generator.Example("", roots)
	require.NoError(t, err)
	return files
}

func containsGoaEndpoint(t *testing.T, f *codegen.File) bool {
	t.Helper()
	for _, s := range f.SectionTemplates {
		if goaEndpointRegexp.MatchString(s.Source) {
			return true
		}
	}
	return false
}

func containsStdlibLogger(t *testing.T, f *codegen.File) (string, bool) {
	t.Helper()
	for _, s := range f.SectionTemplates {
		if strings.Contains(s.Source, "*log.Logger") {
			return s.Source, true
		}
	}
	return "", false
}
