package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	httpcodegen "goa.design/goa/v3/http/codegen"
	"goa.design/plugins/v3/testing/codegen/testdata"
)

func TestGenerateSuiteTopLevel(t *testing.T) {
	cases := map[string]struct {
		DSL  func()
		Code map[string][]string
		Path string
	}{
		"with-payload": {
			DSL: testdata.WithPayloadDSL,
			Code: map[string][]string{
				"suite-test": {testdata.SuiteTestWithPayloadCode},
			},
			Path: "with_payload_service_suite_test.go",
		},
		"with-result": {
			DSL: testdata.WithResultDSL,
			Code: map[string][]string{
				"suite-test": {testdata.SuiteTestWithResultCode},
			},
			Path: "with_result_service_suite_test.go",
		},
		"without-payload-result": {
			DSL: testdata.WithoutPayloadResultDSL,
			Code: map[string][]string{
				"suite-test": {testdata.SuiteTestWithoutPayloadResultCode},
			},
			Path: "without_payload_result_service_suite_test.go",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			root := httpcodegen.RunHTTPDSL(t, c.DSL)
			svc := root.Services[0]
			f := generateSuiteTopLevel("", "", root, svc)
			assert.Equal(t, c.Path, f.Path)
			for sec, secCode := range c.Code {
				testCode(t, f, sec, secCode)
			}
		})
	}
}
