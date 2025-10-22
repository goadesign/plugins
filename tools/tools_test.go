package tools_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"goa.design/goa/v3/codegen"
	. "goa.design/goa/v3/dsl"
	"goa.design/goa/v3/eval"
	goaexpr "goa.design/goa/v3/expr"
	"goa.design/plugins/v3/tools"
	toolsdsl "goa.design/plugins/v3/tools/dsl"
	toolsexpr "goa.design/plugins/v3/tools/expr"
)

func buildTestDSL(t *testing.T) *goaexpr.RootExpr {
	t.Helper()
	toolsexpr.Root.Reset()
	return codegen.RunDSL(t, func() {
		API("tool-api", func() {})

		var GetKeyEventsParams = Type("GetKeyEventsParams", func() {
			Attribute("org_id", String)
			Required("org_id")
		})

		var GetKeyEventsResult = Type("GetKeyEventsResult", func() {
			Attribute("status", String)
			Required("status")
		})

		var ListDevicesPayload = Type("ListDevicesPayload", func() {
			Attribute("org_id", String)
			Required("org_id")
		})

		Service("toolsvc", func() {
			Method("ListDevices", func() {
				Payload(ListDevicesPayload)
				Result(String)
			})

			toolsdsl.ToolSet("ada_tools", func() {
				toolsdsl.Tool("get_key_events", func() {
					Payload(GetKeyEventsParams)
					Result(GetKeyEventsResult)
				})
				toolsdsl.ToolFromMethod("ListDevices")
			})
		})
	})
}

func TestToolDSL(t *testing.T) {
	t.Cleanup(func() { toolsexpr.Root.Reset() })

	_ = buildTestDSL(t)

	if !assert.Len(t, toolsexpr.Root.ToolSets, 1) {
		return
	}
	ts := toolsexpr.Root.ToolSets[0]
	assert.Equal(t, "toolsvc", ts.Service.Name)
	assert.Equal(t, "ada_tools", ts.Name)
	if !assert.Len(t, ts.Tools, 2) {
		return
	}

	pure := ts.Tools[0]
	if assert.NotNil(t, pure) {
		assert.Equal(t, "get_key_events", pure.Name)
		if assert.NotNil(t, pure.Method.Payload) {
			assert.Equal(t, "GetKeyEventsParams", pure.Method.Payload.Type.Name())
		}
		if assert.NotNil(t, pure.Method.Result) {
			assert.Equal(t, "GetKeyEventsResult", pure.Method.Result.Type.Name())
		}
	}

	derived := ts.Tools[1]
	if assert.NotNil(t, derived) {
		assert.Equal(t, "ListDevices", derived.Name)
		if assert.NotNil(t, derived.Method.Payload) {
			assert.Equal(t, "ListDevicesPayload", derived.Method.Payload.Type.Name())
		}
	}
}

func TestGenerateRegistry(t *testing.T) {
	t.Cleanup(func() { toolsexpr.Root.Reset() })

	root := buildTestDSL(t)

	files, err := tools.Generate("example.com/gen", []eval.Root{root}, nil)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Len(t, files, 4) {
		return
	}

	typesPath := filepath.Join(codegen.Gendir, "tools", "types.go")
	schemasPath := filepath.Join(codegen.Gendir, "tools", "schemas.go")
	codecsPath := filepath.Join(codegen.Gendir, "tools", "codecs.go")
	registryPath := filepath.Join(codegen.Gendir, "tools", "registry.go")

	filesByPath := map[string]*codegen.File{}
	for _, f := range files {
		filesByPath[f.Path] = f
	}

	typesFile, ok := filesByPath[typesPath]
	if !assert.True(t, ok, "types.go not emitted") {
		return
	}
	schemasFile, ok := filesByPath[schemasPath]
	if !assert.True(t, ok, "schemas.go not emitted") {
		return
	}
	codecFile, ok := filesByPath[codecsPath]
	if !assert.True(t, ok, "codecs.go not emitted") {
		return
	}
	registryFile, ok := filesByPath[registryPath]
	if !assert.True(t, ok, "registry.go not emitted") {
		return
	}

	var typeBuf bytes.Buffer
	for _, sec := range typesFile.SectionTemplates {
		if !assert.NoError(t, sec.Write(&typeBuf)) {
			return
		}
	}
	typesOut := typeBuf.String()
	assert.Contains(t, typesOut, "type (")
	assert.Contains(t, typesOut, "GetKeyEventsPayload defines the JSON payload")
	assert.NotContains(t, typesOut, "ListDevicesPayload defines the JSON payload")

	var schemaBuf bytes.Buffer
	for _, sec := range schemasFile.SectionTemplates {
		if !assert.NoError(t, sec.Write(&schemaBuf)) {
			return
		}
	}
	schemasOut := schemaBuf.String()
	assert.Contains(t, schemasOut, "getKeyEventsPayloadSchema")
	assert.Contains(t, schemasOut, "$schema")

	var codecBuf bytes.Buffer
	for _, sec := range codecFile.SectionTemplates {
		if !assert.NoError(t, sec.Write(&codecBuf)) {
			return
		}
	}
	codecOut := codecBuf.String()
	assert.Contains(t, codecOut, "import (\n\t\"encoding/json\"")
	assert.Contains(t, codecOut, "GetKeyEventsPayloadCodec = tools.JSONCodec[*GetKeyEventsPayload]")
	assert.Contains(t, codecOut, "func MarshalGetKeyEventsPayload")

	var registryBuf bytes.Buffer
	for _, sec := range registryFile.SectionTemplates {
		if !assert.NoError(t, sec.Write(&registryBuf)) {
			return
		}
	}
	registryOut := registryBuf.String()
	assert.Contains(t, registryOut, "var ToolRegistry = []tools.ToolSpec")
	assert.Contains(t, registryOut, "Payload: tools.TypeSpec{")
}
