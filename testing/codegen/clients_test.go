package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	goacodegen "goa.design/goa/v3/codegen"
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
			services := plannedServiceData(t, root)
			svc := root.Services[0]
			svcData := services.Get(svc.Name)
			f := generateClient("", svcData, root, svc, designedMethodTransports(root))
			assert.Equal(t, c.Path, f.Path)
			for sec, secCode := range c.Code {
				testCode(t, f, sec, secCode)
			}
		})
	}
}

func TestBuildMethodTargetsUsesPlannedJSONRPCMode(t *testing.T) {
	root := goacodegen.RunDSL(t, testdata.JSONRPCTransportsDSL)
	_, jsonPlan := plannedJSONRPCData(t, root)
	svc := root.Services[0]
	transports := plannedMethodTransports(root, jsonPlan)

	plain := buildMethodTargets(root, svc, svc.Method("Plain"), transports)
	assert.Len(t, plain, 1)
	assert.True(t, plain[0].IsJSONRPCPlain)
	assert.True(t, plain[0].IsNoStream)

	events := buildMethodTargets(root, svc, svc.Method("Events"), transports)
	assert.Len(t, events, 1)
	assert.True(t, events[0].IsJSONRPCSSE)
	assert.True(t, events[0].IsServerStream)
}
