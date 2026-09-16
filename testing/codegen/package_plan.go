// This file records every package, name, and import emitted by the testing
// plugin so Goa can report conflicts before it writes source files.
package codegen

import (
	"path"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/expr"
)

// PlanPackage records the package, imports, and fixed Go names emitted for one
// service test helper package. Public names stay fixed because applications
// refer to them from handwritten tests.
func PlanPackage(
	generation *codegen.Generation,
	servicePlan *service.Plan,
	root *expr.RootExpr,
	svc *expr.ServiceExpr,
) error {
	serviceImport, _, err := servicePlan.ServicePackageImports(svc)
	if err != nil {
		return err
	}
	servicePath := path.Base(serviceImport.Path)
	testPackage, err := generation.ClaimPackage(path.Join(serviceImport.Path, servicePath+"test"))
	if err != nil {
		return err
	}
	if err := planPackageImports(testPackage, generation.GenPkg(), root, svc, serviceImport, servicePath); err != nil {
		return err
	}
	return planPackageNames(testPackage, servicePlan, root, svc, designedMethodTransports(root))
}

// planPackageImports records every import name written directly in testing
// templates for one service.
func planPackageImports(
	pkg *codegen.GeneratedPackage,
	genpkg string,
	root *expr.RootExpr,
	svc *expr.ServiceExpr,
	serviceImport *codegen.ImportSpec,
	servicePath string,
) error {
	imports := []*codegen.ImportSpec{
		{Path: "context"},
		{Path: "encoding/json"},
		{Path: "fmt"},
		{Path: "io"},
		{Path: "net/http"},
		{Path: "net/http/httptest"},
		{Path: "os"},
		{Path: "strings"},
		{Path: "testing"},
		{Path: "time"},
		{Path: "goa.design/goa/v3/pkg"},
		{Path: "gopkg.in/yaml.v3", Name: "yaml"},
		{Path: serviceImport.Path, Name: serviceImport.Name},
	}
	if hasHTTPTransport(root, svc) {
		imports = append(imports,
			&codegen.ImportSpec{Path: "bufio"},
			&codegen.ImportSpec{Path: "bytes"},
			&codegen.ImportSpec{Path: "net/url"},
			&codegen.ImportSpec{Path: "github.com/gorilla/websocket"},
			&codegen.ImportSpec{Path: "goa.design/goa/v3/http", Name: "goahttp"},
			&codegen.ImportSpec{Path: path.Join(genpkg, "http", servicePath, "client"), Name: "httpcli"},
			&codegen.ImportSpec{Path: path.Join(genpkg, "http", servicePath, "server"), Name: "httpsvr"},
		)
	}
	if hasGRPCTransport(root, svc) {
		imports = append(imports,
			&codegen.ImportSpec{Path: "net"},
			&codegen.ImportSpec{Path: "google.golang.org/grpc"},
			&codegen.ImportSpec{Path: "google.golang.org/grpc/test/bufconn"},
			&codegen.ImportSpec{Path: "goa.design/goa/v3/grpc", Name: "goagrpc"},
			&codegen.ImportSpec{Path: path.Join(genpkg, "grpc", servicePath, "client"), Name: "grpccli"},
			&codegen.ImportSpec{Path: path.Join(genpkg, "grpc", servicePath, "server"), Name: "grpcsvr"},
			&codegen.ImportSpec{Path: path.Join(genpkg, "grpc", servicePath, "pb"), Name: serviceImport.Name + "pb"},
		)
	}
	if hasJSONRPCTransport(root, svc) {
		imports = append(imports,
			&codegen.ImportSpec{Path: "net/url"},
			&codegen.ImportSpec{Path: "goa.design/goa/v3/http", Name: "goahttp"},
			&codegen.ImportSpec{Path: "goa.design/goa/v3/jsonrpc", Name: "jsonrpc"},
			&codegen.ImportSpec{Path: path.Join(genpkg, "jsonrpc", servicePath, "client"), Name: "jsonrpccli"},
			&codegen.ImportSpec{Path: path.Join(genpkg, "jsonrpc", servicePath, "server"), Name: "jsonrpcsvr"},
		)
	}
	validator := ExtractValidatorsFromYAML()
	if validator.Path != "" && validator.Package != "" && validator.Path != pkg.ImportPath() {
		imports = append(imports, &codegen.ImportSpec{Path: validator.Path, Name: validator.Package})
	}
	for _, spec := range imports {
		if err := pkg.RequireImport(spec); err != nil {
			return err
		}
	}
	return nil
}

// planPackageNames records every package-level declaration written by the
// testing templates. Receiver methods do not need separate package names.
func planPackageNames(
	pkg *codegen.GeneratedPackage,
	servicePlan *service.Plan,
	root *expr.RootExpr,
	svc *expr.ServiceExpr,
	transports *methodTransports,
) error {
	types := []string{
		"Transport",
		"Client",
		"Harness",
		"Options",
		"Option",
		"Scenario",
		"Step",
		"Expectation",
		"ScenarioConfig",
		"Validators",
		"ScenarioRunner",
		"TestData",
	}
	if hasErrors(svc) {
		types = append(types, "ErrorAsserter")
	}
	examples := expr.NewExampleGenerator(root.API.RandomizerFactory)
	for _, method := range svc.Methods {
		if methodHasPayloadExample(method, examples) {
			types = append(types, codegen.Goify(method.Name, true)+"PayloadBuilder")
		}
		if transport, ok := transports.jsonRPC[method]; ok && transport.serverSentEvents && method.Stream == expr.ServerStreamKind {
			names, err := servicePlan.HTTPMethodNames(method)
			if err != nil {
				return err
			}
			types = append(types, "jsonrpcSSE"+names.Method+"Wrapper")
		}
	}
	for _, name := range types {
		if err := pkg.DeclareName(codegen.NewExactName(codegen.NameType, name)); err != nil {
			return err
		}
	}

	constants := []string{"AutoTransport"}
	if hasHTTPTransport(root, svc) {
		constants = append(constants, "HTTPTransport")
	}
	if hasGRPCTransport(root, svc) {
		constants = append(constants, "GRPCTransport")
	}
	if hasJSONRPCTransport(root, svc) {
		constants = append(constants, "JSONRPCTransport")
	}
	for _, name := range constants {
		if err := pkg.DeclareName(codegen.NewExactName(codegen.NameConstant, name)); err != nil {
			return err
		}
	}

	functions := []string{
		"WithContext",
		"WithTimeout",
		"NewHarness",
		"LoadScenarios",
		"NewScenarioRunner",
		"defaultValidateResult",
		"NewTestData",
	}
	if hasErrors(svc) {
		functions = append(functions, "NewErrorAsserter")
	}
	for _, name := range functions {
		if err := pkg.DeclareName(codegen.NewExactName(codegen.NameFunction, name)); err != nil {
			return err
		}
	}

	for _, name := range []string{"ValidTransports", "TransportAvailability"} {
		if err := pkg.DeclareName(codegen.NewExactName(codegen.NameVariable, name)); err != nil {
			return err
		}
	}
	return nil
}

// methodHasPayloadExample reports whether the test data template writes a
// payload builder type for method.
func methodHasPayloadExample(method *expr.MethodExpr, examples *expr.ExampleGenerator) bool {
	if method.StreamingPayload.Type != expr.Empty {
		return method.StreamingPayload.Example(examples.At(expr.MethodStreamingPayloadExampleIdentity(method))) != nil
	}
	if method.Payload.Type != expr.Empty {
		return method.Payload.Example(examples.At(expr.MethodPayloadExampleIdentity(method))) != nil
	}
	return false
}
