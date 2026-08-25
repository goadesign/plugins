package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"goa.design/goa/v3/codegen"
	"goa.design/plugins/v3/testing/codegen/testdata"
)

func TestGenerateHarness(t *testing.T) {
	cases := map[string]struct {
		DSL  func()
		Code map[string][]string
		Path string
	}{
		"with-stream": {
			DSL: testdata.WithStreamDSL,
			Code: map[string][]string{
				"http-harness": {testdata.WithStreamCode},
			},
			Path: "gen/with_stream_service/with_stream_servicetest/harness.go",
		},
		"without-stream": {
			DSL: testdata.WithoutStreamDSL,
			Code: map[string][]string{
				"http-harness": {testdata.WithoutStreamCode},
			},
			Path: "gen/without_stream_service/without_stream_servicetest/harness.go",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			root := codegen.RunDSL(t, c.DSL)
			services := plannedServiceData(t, root)
			svc := root.Services[0]
			svcData := services.Get(svc.Name)
			f := generateHarness("", svcData, root, svc, designedMethodTransports(root))
			assert.Equal(t, c.Path, f.Path)
			for sec, secCode := range c.Code {
				testCode(t, f, sec, secCode)
			}
		})
	}
}

func testCode(t *testing.T, file *codegen.File, section string, expCode []string) {
	sections := file.Section(section)
	require.Len(t, sections, len(expCode))
	for i, c := range expCode {
		code := codegen.SectionCode(t, sections[i])
		assert.Equal(t, c, code)
	}
}
