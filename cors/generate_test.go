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
		HandleOriginCode []string
		MountCORSCode    []string
		ServerInitCode   []string
		CodeGenCount     int
	}{
		{"simple-origin", testdata.SimpleOriginDSL, []string{testdata.SimpleOriginHandleCode}, []string{testdata.SimpleOriginMountCode}, []string{testdata.SimpleOriginServerInitCode}, 2},
		{"regexp-origin", testdata.RegexpOriginDSL, []string{testdata.RegexpOriginHandleCode}, []string{testdata.RegexpOriginMountCode}, []string{testdata.RegexpOriginServerInitCode}, 2},
		{"multi-origin", testdata.MultiOriginDSL, []string{testdata.MultiOriginHandleCode}, []string{testdata.MultiOriginMountCode}, []string{testdata.MultiOriginServerInitCode}, 2},
		{"origin-file-server", testdata.OriginFileServerDSL, []string{testdata.OriginFileServerHandleCode}, []string{testdata.OriginFileServerMountCode}, []string{testdata.OriginFileServerServerInitCode}, 1},
		{"origin-multi-endpoint", testdata.OriginMultiEndpointDSL, []string{testdata.OriginMultiEndpointHandleCode}, []string{testdata.OriginMultiEndpointMountCode}, []string{testdata.OriginMultiEndpointServerInitCode}, 2},
		{"multiservice-origin", testdata.MultiServiceSameOriginDSL, []string{testdata.MultiServiceSameOriginFirstServiceHandleCode, testdata.MultiServiceSameOriginSecondServiceHandleCode}, []string{testdata.MultiServiceSameOriginFirstServiceMountCode, testdata.MultiServiceSameOriginSecondServiceMountCode}, []string{testdata.MultiServiceSameOriginFirstServiceInitCode, testdata.MultiServiceSameOriginSecondServiceInitCode}, 4},
		{"files", testdata.FilesDSL, []string{testdata.FilesHandleCode}, []string{testdata.FilesMountCode}, []string{testdata.FilesServerInitCode}, 1},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			httpcodegen.RunHTTPDSL(t, c.DSL)
			fs := httpcodegen.ServerFiles("", expr.Root)
			if len(fs) != c.CodeGenCount {
				t.Fatalf("got %d files, expected %d", len(fs), c.CodeGenCount)
			}
			cors.Generate("", []eval.Root{expr.Root}, fs)
			expectedCodeIndex := -1
			for _, f := range fs {
				if filepath.Base(f.Path) != "server.go" {
					continue
				}
				expectedCodeIndex += 1
				testCode(t, f, "handle-cors", c.HandleOriginCode[expectedCodeIndex])
				testCode(t, f, "mount-cors", c.MountCORSCode[expectedCodeIndex])
				testCode(t, f, "cors-handler-init", corsHandler)
				testCode(t, f, "server-init", c.ServerInitCode[expectedCodeIndex])
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
