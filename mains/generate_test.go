package mains

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/generator"
	"goa.design/goa/v3/eval"
	"goa.design/plugins/v3/mains/testdata"
)

func TestSingleServiceMainRelocation(t *testing.T) {
	root := codegen.RunDSL(t, testdata.SingleServiceDSL)
	files, err := generator.Example("generated.local/gen", []eval.Root{root})
	require.NoError(t, err)

	// Sanity: default example mains produce cmd/<dir>/main.go and cmd/<dir>/http.go
	require.GreaterOrEqual(t, len(files), 2)

	out, err := Generate("generated.local/gen", []eval.Root{root}, files)
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
	files, err := generator.Example("generated.local/gen", []eval.Root{root})
	require.NoError(t, err)

	out, err := Generate("generated.local/gen", []eval.Root{root}, files)
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
	files, err := generator.Example("generated.local/gen", []eval.Root{root})
	require.NoError(t, err)

	out, err := Generate("generated.local/gen", []eval.Root{root}, files)
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

func TestMainsAddsFileServerNils(t *testing.T) {
	root := codegen.RunDSL(t, testdata.FileServerServiceDSL)
	files, err := generator.Example("generated.local/gen", []eval.Root{root})
	require.NoError(t, err)

	out, err := Generate("generated.local/gen", []eval.Root{root}, files)
	require.NoError(t, err)

	// Expect relocated main under services/static/cmd/static/main.go
	var mainFile *codegen.File
	for _, f := range out {
		if f.Path == "services/static/cmd/static/main.go" {
			mainFile = f
			break
		}
	}
	require.NotNil(t, mainFile)

	// Render the mains section and look for exactly 2 (errhandler, formatter)
	// + 3 (file servers) nil arguments at the end of the New(...) call.
	sections := mainFile.Section("mains-main")
	require.Greater(t, len(sections), 0)
	var buf bytes.Buffer
	require.NoError(t, sections[0].Write(&buf))
	code := buf.String()

	// Match a New(...) call that ends with five consecutive `, nil` args
	// (2 standard + 3 file servers)
	re := regexp.MustCompile(`New\([\s\S]*,\s*nil(?:,\s*nil){4}\)`) // total 5 nils
	assert.Regexp(t, re, code)
}
