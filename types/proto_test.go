package types

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"goa.design/goa/v3/codegen"
	. "goa.design/goa/v3/dsl"
	"goa.design/goa/v3/eval"

	"goa.design/plugins/v3/types/testdata"
)

func TestProto(t *testing.T) {
	cases := []struct {
		Name string
		DSL  func()
	}{
		{"empty", testdata.Empt},
		{"noval", testdata.NoValidation},
		{"required", testdata.Require},
		{"validation", testdata.Validation},
		{"multiple", testdata.Multiple},
		{"alias", testdata.Alias},
		{"array", testdata.Array},
		{"recArray", testdata.ArrayArray},
		{"protojson", testdata.ProtoJSON},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			root := codegen.RunDSL(t, c.DSL)
			fs, err := GenerateProto("example.com/test", []eval.Root{root}, nil)
			require.NoError(t, err)

			if c.Name == "empty" {
				// Empty designs don't generate proto files
				require.Empty(t, fs)
				return
			}

			require.NotEmpty(t, fs)
			require.NotEmpty(t, fs[0].SectionTemplates)
			var buf bytes.Buffer
			for _, s := range fs[0].SectionTemplates {
				require.NoError(t, s.Write(&buf))
			}
			got := buf.String()
			golden := filepath.Join("testdata", fmt.Sprintf("%s.proto_", c.Name))
			if *update {
				os.WriteFile(golden, buf.Bytes(), 0644)
				return
			}
			expected, err := os.ReadFile(golden)
			if os.IsNotExist(err) {
				t.Logf("Golden file doesn't exist yet: %s", golden)
				t.Logf("Generated proto:\n%s", got)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, string(expected), got)
		})
	}
}

func TestProtoTypeMapping(t *testing.T) {
	tests := []struct {
		name     string
		dsl      func()
		expected string
	}{
		{
			"Boolean",
			func() {
				Type("TestType", func() {
					Attribute("field", Boolean)
				})
			},
			"bool",
		},
		{
			"Int32",
			func() {
				Type("TestType", func() {
					Attribute("field", Int32)
				})
			},
			"int32",
		},
		{
			"Int64",
			func() {
				Type("TestType", func() {
					Attribute("field", Int64)
				})
			},
			"int64",
		},
		{
			"String",
			func() {
				Type("TestType", func() {
					Attribute("field", String)
				})
			},
			"string",
		},
		{
			"Bytes",
			func() {
				Type("TestType", func() {
					Attribute("field", Bytes)
				})
			},
			"bytes",
		},
		{
			"Any",
			func() {
				Type("TestType", func() {
					Attribute("field", Any)
				})
			},
			"google.protobuf.Any",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := codegen.RunDSL(t, tt.dsl)
			fs, err := GenerateProto("test", []eval.Root{root}, nil)
			require.NoError(t, err)
			require.NotEmpty(t, fs)

			var buf bytes.Buffer
			for _, s := range fs[0].SectionTemplates {
				require.NoError(t, s.Write(&buf))
			}

			got := buf.String()
			assert.Contains(t, got, tt.expected, "Proto should contain the correct type mapping")
		})
	}
}
