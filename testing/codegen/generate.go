package codegen

import (
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/expr"
)

// Generate produces testing files (client, harness, scenarios, errors, testdata) for the given service.
func Generate(genpkg string, svcData *service.Data, root *expr.RootExpr, svc *expr.ServiceExpr) []*codegen.File {
	files := make([]*codegen.File, 0, 5)
	if f := generateClient(genpkg, svcData, root, svc); f != nil {
		files = append(files, f)
	}
	if f := generateHarness(genpkg, svcData, root, svc); f != nil {
		files = append(files, f)
	}
	if fs := generateScenarios(genpkg, svcData, root, svc); fs != nil {
		files = append(files, fs...)
	}
	if f := generateErrors(genpkg, svcData, root, svc); f != nil {
		files = append(files, f)
	}
	if f := generateTestData(genpkg, svcData, root, svc); f != nil {
		files = append(files, f)
	}
	return files
}

// GenerateSuiteTopLevel returns the top-level example suite file.
func GenerateSuiteTopLevel(genpkg, examplePkg string, root *expr.RootExpr, svc *expr.ServiceExpr) *codegen.File {
	return generateSuiteTopLevel(genpkg, examplePkg, root, svc)
}

// GenerateExampleScenarios returns an example scenarios.yaml file.
func GenerateExampleScenarios(genpkg string, root *expr.RootExpr, svc *expr.ServiceExpr) *codegen.File {
	return generateExampleScenarios(genpkg, root, svc)
}
