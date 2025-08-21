package validators

import (
	"testing"

	testhttpgrpc "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc"
)

// ValidateHTTPNoStream validates the http_no_stream method result
func ValidateHTTPNoStream(t *testing.T, result *testhttpgrpc.HTTPNoStreamResult, expected map[string]any) {
	if result.Out == "" {
		t.Errorf("result should not be empty")
	}
	if expectedResult, ok := expected["out"].(string); ok {
		if result.Out != expectedResult {
			t.Errorf("expected result %q, got %q", expectedResult, result.Out)
		}
	}
}

// ValidateMixedResult validates the mixed_no_stream method result
func ValidateMixedResult(t *testing.T, result *testhttpgrpc.MixedNoStreamResult, expected map[string]any) {
	if result.Out == "" {
		t.Errorf("result should not be empty")
	}
}
