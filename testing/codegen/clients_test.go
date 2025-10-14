package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"goa.design/goa/v3/codegen/service"
	httpcodegen "goa.design/goa/v3/http/codegen"
	"goa.design/plugins/v3/testing/codegen/testdata"
)

func TestGenerateClient(t *testing.T) {
	cases := map[string]struct {
		DSL  func()
		Code map[string][]string
		Path string
	}{
		"with-result": {
			DSL: testdata.WithResultDSL,
			Code: map[string][]string{
				"client-methods": {testdata.ClientMethodsWithResultCode},
			},
			Path: "gen/with_result_service/with_result_servicetest/client.go",
		},
		"without-result": {
			DSL: testdata.WithoutResultDSL,
			Code: map[string][]string{
				"client-methods": {testdata.ClientMethodsWithoutResultCode},
			},
			Path: "gen/without_result_service/without_result_servicetest/client.go",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			root := httpcodegen.RunHTTPDSL(t, c.DSL)
			services := service.NewServicesData(root)
			svc := root.Services[0]
			svcData := services.Get(svc.Name)
			f := generateClient("", svcData, root, svc)
			assert.Equal(t, c.Path, f.Path)
			for sec, secCode := range c.Code {
				testCode(t, f, sec, secCode)
			}
		})
	}
}
