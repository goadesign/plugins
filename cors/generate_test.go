package cors_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/eval"
	goaexpr "goa.design/goa/v3/expr"
	httpcodegen "goa.design/goa/v3/http/codegen"

	"goa.design/plugins/v3/cors"
	"goa.design/plugins/v3/cors/testdata"
)

func TestGenerate(t *testing.T) {
	var corsHandler = `// NewCORSHandler creates a HTTP handler which returns a simple 204 response.
func NewCORSHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
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
		{"simple-env-var-origin", testdata.SimpleEnvVarOriginDSL, []string{testdata.SimpleEnvVarOriginHandleCode}, []string{testdata.SimpleEnvVarOriginMountCode}, []string{testdata.SimpleEnvVarOriginServerInitCode}, 2},
		{"multi-origin", testdata.MultiOriginDSL, []string{testdata.MultiOriginHandleCode}, []string{testdata.MultiOriginMountCode}, []string{testdata.MultiOriginServerInitCode}, 2},
		{"origin-file-server", testdata.OriginFileServerDSL, []string{testdata.OriginFileServerHandleCode}, []string{testdata.OriginFileServerMountCode}, []string{testdata.OriginFileServerServerInitCode}, 1},
		{"origin-multi-endpoint", testdata.OriginMultiEndpointDSL, []string{testdata.OriginMultiEndpointHandleCode}, []string{testdata.OriginMultiEndpointMountCode}, []string{testdata.OriginMultiEndpointServerInitCode}, 2},
		{"multiservice-origin", testdata.MultiServiceSameOriginDSL, []string{testdata.MultiServiceSameOriginFirstServiceHandleCode, testdata.MultiServiceSameOriginSecondServiceHandleCode}, []string{testdata.MultiServiceSameOriginFirstServiceMountCode, testdata.MultiServiceSameOriginSecondServiceMountCode}, []string{testdata.MultiServiceSameOriginFirstServiceInitCode, testdata.MultiServiceSameOriginSecondServiceInitCode}, 4},
		{"files", testdata.FilesDSL, []string{testdata.FilesHandleCode}, []string{testdata.FilesMountCode}, []string{testdata.FilesServerInitCode}, 1},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			root := codegen.RunDSL(t, c.DSL)
			fs := serverFiles(t, root)
			require.Len(t, fs, c.CodeGenCount)
			cors.Generate("", []eval.Root{root}, fs)
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
					assert.Contains(t, s.Source, originHndlr)
				}
				for _, s := range f.Section("server-files") {
					assert.Contains(t, s.Source, originHndlr)
				}
			}
		})
	}
}

// serverFiles builds the HTTP server files through the same planning steps as
// a Goa generation run.
func serverFiles(t *testing.T, root *goaexpr.RootExpr) []*codegen.File {
	t.Helper()
	generation, err := codegen.NewGeneration("generated.local/gen", []eval.Root{root})
	require.NoError(t, err)
	servicePlan, err := service.NewPlan(root, generation, goaexpr.NewExampleGenerator(root.API.RandomizerFactory))
	require.NoError(t, err)
	plans, err := httpcodegen.NewPlans(generation, httpcodegen.PlanInput{Root: root, Service: servicePlan})
	require.NoError(t, err)
	require.NoError(t, generation.Freeze())
	require.NoError(t, servicePlan.Link())
	require.NoError(t, plans[0].Link())
	return plans[0].ServerFiles()
}

func testCode(t *testing.T, file *codegen.File, section, expCode string) {
	sections := file.Section(section)
	require.Greater(t, len(sections), 0)
	code := codegen.SectionCode(t, sections[0])
	assert.Equal(t, expCode, code)
}
