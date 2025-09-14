package docs_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"goa.design/goa/v3/codegen"
	. "goa.design/goa/v3/dsl"
	"goa.design/goa/v3/eval"
	openapi "goa.design/goa/v3/http/codegen/openapi"
	"goa.design/plugins/v3/docs"
	. "goa.design/plugins/v3/docs/dsl"
	plugexpr "goa.design/plugins/v3/docs/expr"
	"goa.design/plugins/v3/docs/testdata"
)

var update = flag.Bool("update", false, "update golden files")

// genDocs unmarshals the generated docs JSON into a generic map for assertions.
func genDocs(t *testing.T, dsl func()) map[string]any {
	t.Helper()
	root := codegen.RunDSL(t, dsl)
	prev := openapi.Definitions
	openapi.Definitions = make(map[string]*openapi.Schema)
	fs, err := docs.Generate("", []eval.Root{root}, nil)
	openapi.Definitions = prev
	require.NoError(t, err)
	require.NotEmpty(t, fs)
	require.NotEmpty(t, fs[0].SectionTemplates)
	var w bytes.Buffer
	require.NoError(t, fs[0].SectionTemplates[0].Write(&w))
	var m map[string]any
	require.NoError(t, json.Unmarshal(w.Bytes(), &m))
	return m
}

func TestDocs(t *testing.T) {
	cases := []struct {
		Name string
		DSL  func()
	}{
		{"api-only", testdata.APIOnly},
		{"no-payload-no-return", testdata.NoPayloadNoReturn},
		{"primitive-payload-no-return", testdata.PrimitivePayloadNoReturn},
		{"array-payload-no-return", testdata.ArrayPayloadNoReturn},
		{"map-payload-no-return", testdata.MapPayloadNoReturn},
		{"user-payload-no-return", testdata.UserPayloadNoReturn},
		{"no-payload-primitive-return", testdata.NoPayloadPrimitiveReturn},
		{"no-payload-array-return", testdata.NoPayloadArrayReturn},
		{"no-payload-map-return", testdata.NoPayloadMapReturn},
		{"no-payload-user-return", testdata.NoPayloadUserReturn},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			root := codegen.RunDSL(t, c.DSL)
			fs, err := docs.Generate("", []eval.Root{root}, nil)
			require.NoError(t, err)
			require.NotEmpty(t, fs)
			require.NotEmpty(t, fs[0].SectionTemplates)
			var buf bytes.Buffer
			require.NoError(t, fs[0].SectionTemplates[0].Write(&buf))
			golden := filepath.Join("testdata", fmt.Sprintf("%s.json", c.Name))
			if *update {
				os.WriteFile(golden, buf.Bytes(), 0644)
				return
			}
			expected, err := os.ReadFile(golden)
			require.NoError(t, err)
			assert.Equal(t, string(expected), buf.String())
		})
	}
}

func TestUseJSONTags(t *testing.T) {
	t.Cleanup(func() { plugexpr.Root.UseJSONTags = false; plugexpr.Root.InlineRefs = false })

	docsMap := genDocs(t, func() {
		UseJSONTags()
		API("Test", func() {})
		var Address = Type("Address", func() {
			Field(1, "street", String, func() { Meta("struct:tag:json", "street_name,omitempty") })
			Field(2, "zip", Int)
			Required("street")
		})
		var User = Type("User", func() {
			Field(1, "name", String, func() { Meta("struct:tag:json", "full_name") })
			Field(2, "address", Address, func() { Meta("struct:tag:json", "addr") })
			Field(3, "skip", String, func() { Meta("struct:tag:json", "-") })
			Required("name", "address")
		})
		Service("S", func() { Method("M1", func() { Payload(User); HTTP(func() { GET("/") }); GRPC(func() {}) }) })
	})

	methods := docsMap["services"].(map[string]any)["S"].(map[string]any)["methods"].(map[string]any)
	m1 := methods["M1"].(map[string]any)
	payload := m1["payload"].(map[string]any)
	ex := payload["example"].(map[string]any)
	if _, ok := ex["full_name"]; !ok {
		t.Fatalf("expected example to include key 'full_name', got keys: %#v", ex)
	}
	if _, ok := ex["addr"]; !ok {
		t.Fatalf("expected example to include key 'addr', got keys: %#v", ex)
	}

	defs := docsMap["definitions"].(map[string]any)
	userDef := defs["User"].(map[string]any)
	props := userDef["properties"].(map[string]any)
	if _, ok := props["full_name"]; !ok {
		t.Fatalf("expected 'full_name' property in User definition, got: %#v", props)
	}
	if _, ok := props["addr"]; !ok {
		t.Fatalf("expected 'addr' property in User definition, got: %#v", props)
	}
	if _, ok := props["skip"]; !ok {
		t.Fatalf("expected 'skip' property to remain unchanged when tag is '-', got: %#v", props)
	}
	req := userDef["required"].([]any)
	assert.Contains(t, req, any("full_name"))
	assert.Contains(t, req, any("addr"))

	addrDef := defs["Address"].(map[string]any)
	addrReq := addrDef["required"].([]any)
	assert.Contains(t, addrReq, any("street_name"))
}

func TestInlineRefs(t *testing.T) {
	t.Cleanup(func() { plugexpr.Root.UseJSONTags = false; plugexpr.Root.InlineRefs = false })

	docsMap := genDocs(t, func() {
		InlineRefs()
		API("Test", func() {})
		var Inner = Type("Inner", func() { Field(1, "value", String) })
		var Outer = Type("Outer", func() { Field(1, "inner", Inner) })
		var MapHolder = Type("MapHolder", func() { Field(1, "m", MapOf(String, Inner)) })
		Service("S", func() {
			Method("M1", func() { Payload(Outer); HTTP(func() { GET("/") }); GRPC(func() {}) })
			Method("M2", func() { Payload(ArrayOf(Inner)); HTTP(func() { GET("/array") }); GRPC(func() {}) })
			Method("M3", func() { Payload(MapHolder); HTTP(func() { GET("/map") }); GRPC(func() {}) })
		})
	})

	svc := docsMap["services"].(map[string]any)["S"].(map[string]any)
	methods := svc["methods"].(map[string]any)
	m1Type := methods["M1"].(map[string]any)["payload"].(map[string]any)["type"].(map[string]any)
	if _, hasRef := m1Type["$ref"]; hasRef {
		t.Fatalf("expected inlined schema for M1 payload, found $ref: %#v", m1Type)
	}
	if _, hasProps := m1Type["properties"]; !hasProps {
		t.Fatalf("expected properties in inlined schema for M1 payload, got: %#v", m1Type)
	}
	m2Type := methods["M2"].(map[string]any)["payload"].(map[string]any)["type"].(map[string]any)
	items := m2Type["items"].(map[string]any)
	if _, hasRef := items["$ref"]; hasRef {
		t.Fatalf("expected inlined schema for M2 payload items, found $ref: %#v", items)
	}
	if _, hasProps := items["properties"]; !hasProps {
		t.Fatalf("expected properties in inlined schema for M2 payload items, got: %#v", items)
	}
	defs := docsMap["definitions"].(map[string]any)
	innerProp := defs["Outer"].(map[string]any)["properties"].(map[string]any)["inner"].(map[string]any)
	if _, hasRef := innerProp["$ref"]; hasRef {
		t.Fatalf("expected inlined nested property in definitions, found $ref: %#v", innerProp)
	}
	mprop := defs["MapHolder"].(map[string]any)["properties"].(map[string]any)["m"].(map[string]any)
	addl := mprop["additionalProperties"].(map[string]any)
	if _, hasRef := addl["$ref"]; hasRef {
		t.Fatalf("expected inlined schema in map additionalProperties, found $ref: %#v", addl)
	}
}

func TestBothOptions(t *testing.T) {
	t.Cleanup(func() { plugexpr.Root.UseJSONTags = false; plugexpr.Root.InlineRefs = false })
	docsMap := genDocs(t, func() {
		UseJSONTags()
		InlineRefs()
		API("Test", func() {})
		var X = Type("X", func() {
			Field(1, "A", String, func() { Meta("struct:tag:json", "aa") })
			Field(2, "b", Int)
			Required("A")
		})
		Service("S", func() { Method("M1", func() { Payload(X); HTTP(func() { GET("/") }); GRPC(func() {}) }) })
	})
	pt := docsMap["services"].(map[string]any)["S"].(map[string]any)["methods"].(map[string]any)["M1"].(map[string]any)["payload"].(map[string]any)["type"].(map[string]any)
	if _, hasRef := pt["$ref"]; hasRef {
		t.Fatalf("expected inlined schema for payload, found $ref: %#v", pt)
	}
	props := pt["properties"].(map[string]any)
	if _, ok := props["aa"]; !ok {
		t.Fatalf("expected JSON tag property 'aa' in inlined payload schema, got: %#v", props)
	}
	if req, ok := pt["required"].([]any); ok {
		assert.Equal(t, 1, len(req))
		assert.Contains(t, req, any("aa"))
	} else {
		t.Fatalf("expected required array in inlined payload schema, got: %#v", pt)
	}
	defs := docsMap["definitions"].(map[string]any)
	if xdef, ok := defs["X"].(map[string]any); ok {
		if xreq, ok := xdef["required"].([]any); ok {
			assert.Contains(t, xreq, any("aa"))
		}
	}
}

func TestJSONTagsAndInlineRefs_Complex(t *testing.T) {
	t.Cleanup(func() { plugexpr.Root.UseJSONTags = false; plugexpr.Root.InlineRefs = false })
	docsMap := genDocs(t, func() {
		UseJSONTags()
		InlineRefs()
		API("Test", func() {})
		var TimeWindow = Type("TimeWindow", func() {
			Field(1, "Start", String, func() { Meta("struct:tag:json", "start") })
			Field(2, "End", String, func() { Meta("struct:tag:json", "end") })
			Required("Start", "End")
		})
		var SeriesSource = Type("SeriesSource", func() {
			Field(1, "DeviceAlias", String, func() { Meta("struct:tag:json", "device_alias") })
			Field(2, "SignalAlias", String, func() { Meta("struct:tag:json", "signal_alias") })
			Required("DeviceAlias", "SignalAlias")
		})
		var TimeSeriesInput = Type("TimeSeriesInput", func() {
			Field(1, "SessionID", String, func() { Meta("struct:tag:json", "session_id") })
			Field(2, "Sources", ArrayOf(SeriesSource), func() { Meta("struct:tag:json", "sources") })
			Field(3, "Window", TimeWindow, func() { Meta("struct:tag:json", "window") })
			Required("SessionID", "Sources", "Window")
		})
		var TSResult = Type("TSResult", func() {
			Field(1, "Alarms", ArrayOf(String), func() { Meta("struct:tag:json", "alarms") })
			Field(2, "EvidenceRefs", ArrayOf(String), func() { Meta("struct:tag:json", "evidence_refs") })
			Required("Alarms", "EvidenceRefs")
		})
		Service("S", func() {
			Method("Get", func() { Payload(TimeSeriesInput); Result(TSResult); HTTP(func() { GET("/") }); GRPC(func() {}) })
		})
	})
	defs := docsMap["definitions"].(map[string]any)
	tsi := defs["TimeSeriesInput"].(map[string]any)
	props := tsi["properties"].(map[string]any)
	_, hasSID := props["session_id"]
	_, hasSrc := props["sources"]
	_, hasWin := props["window"]
	assert.True(t, hasSID)
	assert.True(t, hasSrc)
	assert.True(t, hasWin)
	req := tsi["required"].([]any)
	assert.Contains(t, req, any("session_id"))
	assert.Contains(t, req, any("sources"))
	assert.Contains(t, req, any("window"))
	svc := docsMap["services"].(map[string]any)["S"].(map[string]any)
	methods := svc["methods"].(map[string]any)
	get := methods["Get"].(map[string]any)
	pt := get["payload"].(map[string]any)["type"].(map[string]any)
	if _, hasRef := pt["$ref"]; hasRef {
		t.Fatalf("expected inlined payload schema, found $ref: %#v", pt)
	}
	ppt := pt["properties"].(map[string]any)
	_, ok1 := ppt["session_id"]
	_, ok2 := ppt["sources"]
	_, ok3 := ppt["window"]
	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.True(t, ok3)
	if reqP, ok := pt["required"].([]any); ok {
		assert.Contains(t, reqP, any("session_id"))
		assert.Contains(t, reqP, any("sources"))
		assert.Contains(t, reqP, any("window"))
	} else {
		t.Fatalf("expected required array in payload type, got: %#v", pt)
	}
	rt := get["result"].(map[string]any)["type"].(map[string]any)
	if _, hasRef := rt["$ref"]; hasRef {
		t.Fatalf("expected inlined result schema, found $ref: %#v", rt)
	}
	rprops := rt["properties"].(map[string]any)
	_, okA := rprops["alarms"]
	_, okE := rprops["evidence_refs"]
	assert.True(t, okA)
	assert.True(t, okE)
	if reqR, ok := rt["required"].([]any); ok {
		assert.Contains(t, reqR, any("alarms"))
		assert.Contains(t, reqR, any("evidence_refs"))
	} else {
		t.Fatalf("expected required array in result type, got: %#v", rt)
	}
}

func TestInlineRefs_MethodRootPayloadResult(t *testing.T) {
	t.Cleanup(func() { plugexpr.Root.UseJSONTags = false; plugexpr.Root.InlineRefs = false })

	docsMap := genDocs(t, func() {
		InlineRefs()
		API("Test", func() {})
		var PayloadType = Type("PayloadType", func() {
			Field(1, "Foo", String)
			Required("Foo")
		})
		var ResultType = Type("ResultType", func() {
			Field(1, "Bar", String)
			Required("Bar")
		})
		Service("S", func() {
			Method("M1", func() {
				Payload(PayloadType)
				Result(ResultType)
				HTTP(func() { GET("/") })
				GRPC(func() {})
			})
		})
	})

	svc := docsMap["services"].(map[string]any)["S"].(map[string]any)
	methods := svc["methods"].(map[string]any)
	m1 := methods["M1"].(map[string]any)
	pt := m1["payload"].(map[string]any)["type"].(map[string]any)
	if _, hasRef := pt["$ref"]; hasRef {
		t.Fatalf("expected inlined payload schema at method root, found $ref: %#v", pt)
	}
	rt := m1["result"].(map[string]any)["type"].(map[string]any)
	if _, hasRef := rt["$ref"]; hasRef {
		t.Fatalf("expected inlined result schema at method root, found $ref: %#v", rt)
	}
}

func TestInlineRefs_CrossPackage(t *testing.T) {
	t.Cleanup(func() { plugexpr.Root.UseJSONTags = false; plugexpr.Root.InlineRefs = false })
	docsMap := genDocs(t, func() {
		InlineRefs()
		API("Test", func() {})
		// Service "S1" uses a type from a different package.
		Service("S1", func() {
			Method("M1", func() {
				Payload(testdata.CrossPackageType)
				HTTP(func() { GET("/") })
				GRPC(func() {})
			})
		})
	})

	svc := docsMap["services"].(map[string]any)["S1"].(map[string]any)
	methods := svc["methods"].(map[string]any)
	m1 := methods["M1"].(map[string]any)
	pt := m1["payload"].(map[string]any)["type"].(map[string]any)
	if _, hasRef := pt["$ref"]; hasRef {
		t.Fatalf("expected inlined payload schema for cross-package type, found $ref: %#v", pt)
	}
	assert.Equal(t, "object", pt["type"])
	props := pt["properties"].(map[string]any)
	_, hasA := props["A"]
	assert.True(t, hasA)
}

func TestInlineRefs_CrossService(t *testing.T) {
	t.Cleanup(func() { plugexpr.Root.UseJSONTags = false; plugexpr.Root.InlineRefs = false })
	docsMap := genDocs(t, func() {
		InlineRefs()
		API("Test", func() {})
		var SharedType = Type("SharedType", func() {
			Field(1, "A", String)
			Required("A")
		})
		Service("S1", func() {
			Method("M1", func() {
				Payload(SharedType)
				HTTP(func() { GET("/") })
				GRPC(func() {})
			})
		})
		Service("S2", func() {
			Method("M2", func() {
				Payload(SharedType)
				HTTP(func() { GET("/s2") })
				GRPC(func() {})
			})
		})
	})

	// Check S1
	s1 := docsMap["services"].(map[string]any)["S1"].(map[string]any)
	m1 := s1["methods"].(map[string]any)["M1"].(map[string]any)
	pt1 := m1["payload"].(map[string]any)["type"].(map[string]any)
	if _, hasRef := pt1["$ref"]; hasRef {
		t.Fatalf("S1: expected inlined payload schema for shared type, found $ref: %#v", pt1)
	}

	// Check S2
	s2 := docsMap["services"].(map[string]any)["S2"].(map[string]any)
	m2 := s2["methods"].(map[string]any)["M2"].(map[string]any)
	pt2 := m2["payload"].(map[string]any)["type"].(map[string]any)
	if _, hasRef := pt2["$ref"]; hasRef {
		t.Fatalf("S2: expected inlined payload schema for shared type, found $ref: %#v", pt2)
	}
}
