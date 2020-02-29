package cors_test

import (
	"path/filepath"
	"strings"
	"testing"

	"goa.design/goa/codegen"
	"goa.design/goa/eval"
	"goa.design/goa/expr"
	httpcodegen "goa.design/goa/http/codegen"
	"goa.design/plugins/cors"
	"goa.design/plugins/cors/testdata"
)

func TestGenerate(t *testing.T) {
	var corsHandler = `// NewCORSHandler creates a HTTP handler which returns a simple 200 response.
func NewCORSHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
}
`
	cases := []struct {
		Name             string
		DSL              func()
		HandleOriginCode string
		MountCORSCode    string
		ServerInitCode   string
		FileCount        int
	}{
		{"simple-origin", testdata.SimpleOriginDSL, testdata.SimpleOriginHandleCode, testdata.SimpleOriginMountCode, testdata.SimpleOriginServerInitCode, 2},
		{"regexp-origin", testdata.RegexpOriginDSL, testdata.RegexpOriginHandleCode, testdata.RegexpOriginMountCode, testdata.RegexpOriginServerInitCode, 2},
		{"multi-origin", testdata.MultiOriginDSL, testdata.MultiOriginHandleCode, testdata.MultiOriginMountCode, testdata.MultiOriginServerInitCode, 2},
		{"origin-file-server", testdata.OriginFileServerDSL, testdata.OriginFileServerHandleCode, testdata.OriginFileServerMountCode, testdata.OriginFileServerServerInitCode, 1},
		{"origin-multi-endpoint", testdata.OriginMultiEndpointDSL, testdata.OriginMultiEndpointHandleCode, testdata.OriginMultiEndpointMountCode, testdata.OriginMultiEndpointServerInitCode, 2},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			httpcodegen.RunHTTPDSL(t, c.DSL)
			fs := httpcodegen.ServerFiles("", expr.Root)
			if len(fs) != c.FileCount {
				t.Fatalf("got %d files, expected two", len(fs))
			}
			cors.Generate("", []eval.Root{expr.Root}, fs)
			for _, f := range fs {
				if filepath.Base(f.Path) != "server.go" {
					continue
				}
				testCode(t, f, "handle-cors", c.HandleOriginCode)
				testCode(t, f, "mount-cors", c.MountCORSCode)
				testCode(t, f, "cors-handler-init", corsHandler)
				testCode(t, f, "server-init", c.ServerInitCode)
				var originHndlr string
				for _, s := range f.Section("handle-cors") {
					data := s.Data.(*cors.ServiceData)
					originHndlr = data.OriginHandler
				}
				for _, s := range f.Section("server-handler") {
					if !strings.Contains(s.Source, originHndlr) {
						t.Errorf("server-handler: invalid code, expected to contain %s", originHndlr)
					}
				}
				for _, s := range f.Section("server-files") {
					if !strings.Contains(s.Source, originHndlr) {
						t.Errorf("server-handler: invalid code, expected to contain %s", originHndlr)
					}
				}
			}
		})
	}
}

func testCode(t *testing.T, file *codegen.File, section, expCode string) {
	sections := file.Section(section)
	if len(sections) < 1 {
		t.Fatalf("%s: got %d sections, expected at least 1", section, len(sections))
	}
	code := codegen.SectionCode(t, sections[0])
	if code != expCode {
		t.Errorf("invalid code, got:\n%s\ngot vs. expected:\n%s", code, codegen.Diff(t, code, expCode))
	}
}
