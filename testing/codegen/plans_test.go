// This file builds real Goa plans used by focused testing-plugin unit tests.
package codegen

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	goacodegen "goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/dsl"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
	httpcodegen "goa.design/goa/v3/http/codegen"
	jsonrpccodegen "goa.design/goa/v3/jsonrpc/codegen"
	"goa.design/plugins/v3/testing/codegen/testdata"
)

// TestJSONRPCSSEWrappersUseFinalMethodNames verifies that the testing plugin
// uses the distinct Go method names chosen by Goa when it writes stream
// wrappers for methods whose design names normalize to the same Go name.
func TestJSONRPCSSEWrappersUseFinalMethodNames(t *testing.T) {
	root := goacodegen.RunDSL(t, func() {
		dsl.API("stream-collision", func() {
			dsl.JSONRPC(func() {})
		})
		event := dsl.Type("Event", func() {
			dsl.Field(1, "message", dsl.String)
			dsl.Required("message")
		})
		dsl.Service("Events", func() {
			dsl.JSONRPC(func() {
				dsl.POST("/rpc")
			})
			for _, name := range []string{"watch-events", "watch_events"} {
				dsl.Method(name, func() {
					dsl.StreamingResult(event)
					dsl.JSONRPC(func() {
						dsl.ServerSentEvents()
					})
				})
			}
		})
	})
	directory, err := os.MkdirTemp(".", "generated_test_")
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(directory))
	})
	genpkg := path.Join("goa.design/plugins/v3/testing/codegen", filepath.Base(directory), "gen")
	generation, err := goacodegen.NewGeneration(genpkg, []eval.Root{root})
	require.NoError(t, err)
	servicePlan, err := service.NewPlan(root, generation, expr.NewExampleGenerator(root.API.RandomizerFactory))
	require.NoError(t, err)
	httpPlans, err := httpcodegen.NewJSONRPCPlans(generation, httpcodegen.PlanInput{Root: root, Service: servicePlan})
	require.NoError(t, err)
	jsonPlans, err := jsonrpccodegen.NewPlans(generation, jsonrpccodegen.PlanInput{
		Root:    root,
		Service: servicePlan,
		HTTP:    httpPlans[0],
	})
	require.NoError(t, err)
	require.NoError(t, PlanPackage(generation, servicePlan, root, root.Service("Events")))
	require.NoError(t, generation.Freeze())
	require.NoError(t, servicePlan.Link())
	require.NoError(t, httpPlans[0].Link())
	require.NoError(t, jsonPlans[0].Link())

	files, err := service.Files(servicePlan)
	require.NoError(t, err)
	files = append(files, jsonPlans[0].ServerFiles()...)
	files = append(files, jsonPlans[0].ClientFiles()...)
	files = append(files, jsonPlans[0].ServerTypeFiles()...)
	files = append(files, jsonPlans[0].ClientTypeFiles()...)
	files = append(files, jsonPlans[0].PathFiles()...)
	files = append(files, GeneratePlanned(
		generation.GenPkg(),
		servicePlan.Services().Get("Events"),
		root,
		root.Service("Events"),
		jsonPlans[0],
	)...)

	for _, file := range files {
		_, err := file.Render(directory)
		require.NoError(t, err)
	}
	command := exec.Command("go", "test", "./...")
	command.Dir = directory
	output, err := command.CombinedOutput()
	require.NoError(t, err, string(output))
}

// TestPlanPackageRejectsPublicNameCollisions verifies that testing helpers
// report a fixed API name conflict during planning instead of emitting invalid
// Go source after names have been chosen.
func TestPlanPackageRejectsPublicNameCollisions(t *testing.T) {
	root := goacodegen.RunDSL(t, testdata.WithArrayResultDSL)
	generation, err := goacodegen.NewGeneration("generated.local/gen", []eval.Root{root})
	require.NoError(t, err)
	servicePlan, err := service.NewPlan(root, generation, expr.NewExampleGenerator(root.API.RandomizerFactory))
	require.NoError(t, err)
	testPackage, err := generation.ClaimPackage("generated.local/gen/with_array_result_service/with_array_result_servicetest")
	require.NoError(t, err)
	require.NoError(t, testPackage.DeclareName(goacodegen.NewExactName(goacodegen.NameType, "Client")))

	err = PlanPackage(generation, servicePlan, root, root.Service("WithArrayResultService"))
	require.ErrorContains(t, err, `cannot declare exact type "Client"`)
}

// plannedServiceData runs Goa's service planning steps and returns the values
// that production plugins receive after Goa has chosen every generated name.
func plannedServiceData(t *testing.T, root *expr.RootExpr) *service.ServicesData {
	t.Helper()

	generation, err := goacodegen.NewGeneration("generated.local/gen", []eval.Root{root})
	require.NoError(t, err)
	plan, err := service.NewPlan(root, generation, expr.NewExampleGenerator(root.API.RandomizerFactory))
	require.NoError(t, err)
	require.NoError(t, generation.Freeze())
	require.NoError(t, plan.Link())
	return plan.Services()
}

// plannedJSONRPCData runs Goa's service and JSON-RPC planning steps and
// returns the finalized values that the testing plugin reads during generation.
func plannedJSONRPCData(t *testing.T, root *expr.RootExpr) (*service.ServicesData, *jsonrpccodegen.Plan) {
	t.Helper()

	generation, err := goacodegen.NewGeneration("generated.local/gen", []eval.Root{root})
	require.NoError(t, err)
	servicePlan, err := service.NewPlan(root, generation, expr.NewExampleGenerator(root.API.RandomizerFactory))
	require.NoError(t, err)
	jsonHTTPPlans, err := httpcodegen.NewJSONRPCPlans(generation, httpcodegen.PlanInput{Root: root, Service: servicePlan})
	require.NoError(t, err)
	jsonPlans, err := jsonrpccodegen.NewPlans(generation, jsonrpccodegen.PlanInput{
		Root:    root,
		Service: servicePlan,
		HTTP:    jsonHTTPPlans[0],
	})
	require.NoError(t, err)
	require.NoError(t, generation.Freeze())
	require.NoError(t, servicePlan.Link())
	require.NoError(t, jsonHTTPPlans[0].Link())
	require.NoError(t, jsonPlans[0].Link())
	return servicePlan.Services(), jsonPlans[0]
}
