// This file assembles the generated testing package and preserves the released
// helper functions for applications that call the plugin directly.
package codegen

import (
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
	jsonrpccodegen "goa.design/goa/v3/jsonrpc/codegen"
)

// Generate produces testing files (client, harness, scenarios, errors, testdata) for the given service.
func Generate(genpkg string, svcData *service.Data, root *expr.RootExpr, svc *expr.ServiceExpr) []*codegen.File {
	return generate(genpkg, svcData, root, svc, designedMethodTransports(root))
}

// GeneratePlanned produces testing files from the service and JSON-RPC values
// finalized by the current Goa generation run.
func GeneratePlanned(genpkg string, svcData *service.Data, root *expr.RootExpr, svc *expr.ServiceExpr, jsonPlan *jsonrpccodegen.Plan) []*codegen.File {
	return generate(genpkg, svcData, root, svc, plannedMethodTransports(root, jsonPlan))
}

// GenerateSuiteTopLevel returns the top-level example suite file.
func GenerateSuiteTopLevel(genpkg, examplePkg string, root *expr.RootExpr, svc *expr.ServiceExpr) *codegen.File {
	svcData := isolatedServiceData(genpkg, root, svc)
	return GenerateSuiteTopLevelFromData(genpkg, examplePkg, svcData, root, svc)
}

// GenerateSuiteTopLevelFromData returns the top-level example suite from
// service values already finalized by a Goa service plan.
func GenerateSuiteTopLevelFromData(genpkg, examplePkg string, svcData *service.Data, root *expr.RootExpr, svc *expr.ServiceExpr) *codegen.File {
	return generateSuiteTopLevel(genpkg, examplePkg, svcData, root, svc, designedMethodTransports(root))
}

// GeneratePlannedSuiteTopLevel returns an example suite using the transport
// choices finalized by the current Goa generation run.
func GeneratePlannedSuiteTopLevel(genpkg, examplePkg string, svcData *service.Data, root *expr.RootExpr, svc *expr.ServiceExpr, jsonPlan *jsonrpccodegen.Plan) *codegen.File {
	return generateSuiteTopLevel(genpkg, examplePkg, svcData, root, svc, plannedMethodTransports(root, jsonPlan))
}

// GenerateExampleScenarios returns an example scenarios.yaml file.
func GenerateExampleScenarios(genpkg string, root *expr.RootExpr, svc *expr.ServiceExpr) *codegen.File {
	svcData := isolatedServiceData(genpkg, root, svc)
	return GenerateExampleScenariosFromData(genpkg, svcData, root, svc)
}

// GenerateExampleScenariosFromData returns an example scenarios.yaml file
// from service values already finalized by a Goa service plan.
func GenerateExampleScenariosFromData(genpkg string, svcData *service.Data, root *expr.RootExpr, svc *expr.ServiceExpr) *codegen.File {
	return generateExampleScenarios(genpkg, svcData, root, svc, designedMethodTransports(root))
}

// GeneratePlannedExampleScenarios returns example scenarios using the
// transport choices finalized by the current Goa generation run.
func GeneratePlannedExampleScenarios(genpkg string, svcData *service.Data, root *expr.RootExpr, svc *expr.ServiceExpr, jsonPlan *jsonrpccodegen.Plan) *codegen.File {
	return generateExampleScenarios(genpkg, svcData, root, svc, plannedMethodTransports(root, jsonPlan))
}

// generate produces the test files after transport choices have been reduced
// to the branches that this service can use.
func generate(genpkg string, svcData *service.Data, root *expr.RootExpr, svc *expr.ServiceExpr, transports *methodTransports) []*codegen.File {
	files := make([]*codegen.File, 0, 5)
	if f := generateClient(genpkg, svcData, root, svc, transports); f != nil {
		files = append(files, f)
	}
	if f := generateHarness(genpkg, svcData, root, svc, transports); f != nil {
		files = append(files, f)
	}
	if fs := generateScenarios(genpkg, svcData, root, svc, transports); fs != nil {
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

// isolatedServiceData runs Goa's current service planner for callers of the
// older helper API, which supplies a design but not a generation plan.
func isolatedServiceData(genpkg string, root *expr.RootExpr, svc *expr.ServiceExpr) *service.Data {
	generation, err := codegen.NewGeneration(genpkg, []eval.Root{root})
	if err != nil {
		panic(err)
	}
	plan, err := service.NewPlan(root, generation, expr.NewExampleGenerator(root.API.RandomizerFactory))
	if err != nil {
		panic(err)
	}
	if err := generation.Freeze(); err != nil {
		panic(err)
	}
	if err := plan.Link(); err != nil {
		panic(err)
	}
	return plan.Services().Get(svc.Name)
}
