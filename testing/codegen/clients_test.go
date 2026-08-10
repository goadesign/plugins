package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	goacodegen "goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	"goa.design/plugins/v3/testing/codegen/testdata"
)

func TestGenerateClient(t *testing.T) {
	cases := map[string]struct {
		DSL  func()
		Code map[string][]string
		Path string
	}{
		"with-payload": {
			DSL: testdata.WithPayloadDSL,
			Code: map[string][]string{
				"client-methods": {testdata.ClientMethodsWithPayloadCode},
			},
			Path: "gen/with_payload_service/with_payload_servicetest/client.go",
		},
		"with-result": {
			DSL: testdata.WithResultDSL,
			Code: map[string][]string{
				"client-methods": {testdata.ClientMethodsWithResultCode},
			},
			Path: "gen/with_result_service/with_result_servicetest/client.go",
		},
		"without-payload-result": {
			DSL: testdata.WithoutPayloadResultDSL,
			Code: map[string][]string{
				"client-methods": {testdata.ClientMethodsWithoutPayloadResultCode},
			},
			Path: "gen/without_payload_result_service/without_payload_result_servicetest/client.go",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			root := goacodegen.RunDSL(t, c.DSL)
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
