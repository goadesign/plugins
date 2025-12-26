package mains

import (
	"testing"

	"bytes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/example"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/eval"
	httpcodegen "goa.design/goa/v3/http/codegen"
	"goa.design/plugins/v3/mains/testdata"
)

func TestSingleServiceMainRelocation(t *testing.T) {
	root := codegen.RunDSL(t, testdata.SingleServiceDSL)
	svcs := service.NewServicesData(root)
	httpSvcs := httpcodegen.NewServicesData(svcs, root.API.HTTP)
	files := append(example.ServerFiles("gen", root, svcs), httpcodegen.ExampleServerFiles("gen", httpSvcs)...)

	// Sanity: default example mains produce cmd/<dir>/main.go and cmd/<dir>/http.go
	require.GreaterOrEqual(t, len(files), 2)

	out, err := Generate("gen", []eval.Root{root}, files)
	require.NoError(t, err)

	// Expect relocated main under services/calc/cmd/calc/main.go
	var hasRelocated bool
	var hasOldMain, hasOldHTTP bool
	for _, f := range out {
		switch f.Path {
		case "services/calc/cmd/calc/main.go":
			hasRelocated = true
		case "cmd/edge/main.go":
			hasOldMain = true
		case "cmd/edge/http.go":
			hasOldHTTP = true
		}
	}
	assert.True(t, hasRelocated, "expected relocated main not found")
	assert.False(t, hasOldMain, "default example main should be removed")
	assert.False(t, hasOldHTTP, "default http.go should be removed")
}

func TestMultiServiceMainStaysUnderCmd(t *testing.T) {
	root := codegen.RunDSL(t, testdata.MultiServiceDSL)
	svcs := service.NewServicesData(root)
	httpSvcs := httpcodegen.NewServicesData(svcs, root.API.HTTP)
	files := append(example.ServerFiles("gen", root, svcs), httpcodegen.ExampleServerFiles("gen", httpSvcs)...)

	out, err := Generate("gen", []eval.Root{root}, files)
	require.NoError(t, err)

	var hasCmdMain, hasHTTP bool
	for _, f := range out {
		if f.Path == "cmd/api/main.go" {
			hasCmdMain = true
		}
		if f.Path == "cmd/api/http.go" {
			hasHTTP = true
		}
	}
	assert.True(t, hasCmdMain, "expected cmd/api/main.go for multi-service server")
	assert.False(t, hasHTTP, "default http.go should be removed")
}

func TestWebSocketMainIncludesUpgrader(t *testing.T) {
	root := codegen.RunDSL(t, testdata.WebSocketServiceDSL)
	svcs := service.NewServicesData(root)
	httpSvcs := httpcodegen.NewServicesData(svcs, root.API.HTTP)
	files := append(example.ServerFiles("gen", root, svcs), httpcodegen.ExampleServerFiles("gen", httpSvcs)...)

	out, err := Generate("gen", []eval.Root{root}, files)
	require.NoError(t, err)

	// Find relocated main
	var mainFile *codegen.File
	for _, f := range out {
		if f.Path == "services/chat/cmd/chat/main.go" {
			mainFile = f
			break
		}
	}
	require.NotNil(t, mainFile)

	// Assert WS import is present in the generated header.
	header := mainFile.Section("source-header")
	require.Greater(t, len(header), 0)
	var hbuf bytes.Buffer
	require.NoError(t, header[0].Write(&hbuf))
	assert.Contains(t, hbuf.String(), "github.com/gorilla/websocket")

	// Assert upgrader usage is present in the main body.
	body := mainFile.Section("mains-main")
	require.Greater(t, len(body), 0)
	var bbuf bytes.Buffer
	require.NoError(t, body[0].Write(&bbuf))
	assert.Contains(t, bbuf.String(), "websocket.Upgrader")
}
