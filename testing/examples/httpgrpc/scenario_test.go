package testhttpgrpcapi

import (
	"testing"

	test "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc/test_http_grpctest"
)

// TestScenarios runs all scenarios via the generated runner
func TestScenarios(t *testing.T) {
	runner, err := test.LoadScenarios("scenarios.yaml")
	if err != nil {
		t.Fatalf("failed to load scenarios: %v", err)
	}
	harness := test.NewHarness(t, NewTestHTTPGrpc())
	runner.Run(t, harness.Client)
}
