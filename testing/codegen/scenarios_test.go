package codegen

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	"goa.design/plugins/v3/testing/codegen/testdata"
)

func TestGenerateScenarios(t *testing.T) {
	cases := map[string]struct {
		DSL  func()
		Code map[string][]string
		Path string
	}{
		"with-payload": {
			DSL: testdata.WithPayloadDSL,
			Code: map[string][]string{
				"scenario-types":  {testdata.ScenarioTypesWithPayloadCode},
				"scenario-runner": {testdata.ScenarioRunnerWithPayloadCode},
			},
			Path: "gen/with_payload_service/with_payload_servicetest/scenarios.go",
		},
		"with-result": {
			DSL: testdata.WithResultDSL,
			Code: map[string][]string{
				"scenario-types":  {testdata.ScenarioTypesWithResultCode},
				"scenario-runner": {testdata.ScenarioRunnerWithResultCode},
			},
			Path: "gen/with_result_service/with_result_servicetest/scenarios.go",
		},
		"without-payload-result": {
			DSL: testdata.WithoutPayloadResultDSL,
			Code: map[string][]string{
				"scenario-types":  {testdata.ScenarioTypesWithoutPayloadResultCode},
				"scenario-runner": {testdata.ScenarioRunnerWithoutPayloadResultCode},
			},
			Path: "gen/without_payload_result_service/without_payload_result_servicetest/scenarios.go",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			root := codegen.RunDSL(t, c.DSL)
			services := service.NewServicesData(root)
			svc := root.Services[0]
			svcData := services.Get(svc.Name)
			fs := generateScenarios("", svcData, root, svc)
			f := fs[0]
			assert.Equal(t, c.Path, f.Path)
			for sec, secCode := range c.Code {
				testCode(t, f, sec, secCode)
			}
		})
	}
}

func TestGenerateScenarios_ArrayResultTypeAssertion(t *testing.T) {
	root := codegen.RunDSL(t, testdata.WithArrayResultDSL)
	services := service.NewServicesData(root)
	svc := root.Services[0]
	svcData := services.Get(svc.Name)
	fs := generateScenarios("", svcData, root, svc)
	f := fs[0]

	sections := f.Section("scenario-runner")
	if len(sections) != 1 {
		t.Fatalf("expected 1 scenario-runner section, got %d", len(sections))
	}
	code := codegen.SectionCode(t, sections[0])

	// This is the canonical fully-qualified Go type reference that should be used
	// in the generated type assertion.
	wantRef := svcData.Scope.GoFullTypeRef(svc.Methods[0].Result, svcData.PkgName)
	assert.Contains(t, code, "typedResult := result.("+wantRef+")")

	// Regression guard for invalid formatting like "pkg.[]T" (issue #234).
	assert.NotContains(t, code, svcData.PkgName+".[]")
	assert.NotContains(t, code, "*"+svcData.PkgName+".[]")
}

func TestGenerateExampleScenarios(t *testing.T) {
	cases := map[string]struct {
		DSL  func()
		Code map[string][]string
	}{
		"with-payload": {
			DSL: testdata.WithPayloadDSL,
			Code: map[string][]string{
				"example-scenarios": {testdata.ExampleScenariosWithPayloadCode},
			},
		},
		"with-result": {
			DSL: testdata.WithResultDSL,
			Code: map[string][]string{
				"example-scenarios": {testdata.ExampleScenariosWithResultCode},
			},
		},
		"without-payload-result": {
			DSL: testdata.WithoutPayloadResultDSL,
			Code: map[string][]string{
				"example-scenarios": {testdata.ExampleScenariosWithoutPayloadResultCode},
			},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			root := codegen.RunDSL(t, c.DSL)
			svc := root.Services[0]
			f := generateExampleScenarios("", root, svc)
			assert.Equal(t, "scenarios.yaml", f.Path)
			for sec, secCode := range c.Code {
				sections := f.Section(sec)
				require.Len(t, sections, len(secCode))
				for i, c := range secCode {
					var buf bytes.Buffer
					assert.NoError(t, sections[i].Write(&buf))
					assert.Equal(t, c, buf.String())
				}
			}
		})
	}
}
