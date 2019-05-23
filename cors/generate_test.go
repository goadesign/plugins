package cors_test

import (
	"path/filepath"
	"strings"
	"testing"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
	httpcodegen "goa.design/goa/v3/http/codegen"
	"goa.design/plugins/v3/cors"
	"goa.design/plugins/v3/cors/testdata"
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
	}{
		{"simple-origin", testdata.SimpleOriginDSL, testdata.SimpleOriginHandleCode, testdata.SimpleOriginMountCode, testdata.SimpleOriginServerInitCode},
		{"regexp-origin", testdata.RegexpOriginDSL, testdata.RegexpOriginHandleCode, testdata.RegexpOriginMountCode, testdata.RegexpOriginServerInitCode},
		{"multi-origin", testdata.MultiOriginDSL, testdata.MultiOriginHandleCode, testdata.MultiOriginMountCode, testdata.MultiOriginServerInitCode},
		{"origin-file-server", testdata.OriginFileServerDSL, testdata.OriginFileServerHandleCode, testdata.OriginFileServerMountCode, testdata.OriginFileServerServerInitCode},
		{"origin-multi-endpoint", testdata.OriginMultiEndpointDSL, testdata.OriginMultiEndpointHandleCode, testdata.OriginMultiEndpointMountCode, testdata.OriginMultiEndpointServerInitCode},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			httpcodegen.RunHTTPDSL(t, c.DSL)
			fs := httpcodegen.ServerFiles("", expr.Root)
			if len(fs) != 2 {
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
