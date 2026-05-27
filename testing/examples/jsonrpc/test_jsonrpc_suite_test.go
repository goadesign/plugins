package testjsonrpcapi

import (
	"context"
	"testing"
	"time"

	testjsonrpc "goa.design/plugins/v3/testing/examples/jsonrpc/gen/test_jsonrpc"
	testJsonrpctest "goa.design/plugins/v3/testing/examples/jsonrpc/gen/test_jsonrpc/test_jsonrpctest"
)

// RunTestJsonrpcHarness exercises the generated harness against your service
// implementation.
// Call this helper from your test, passing your service implementation.
func RunTestJsonrpcHarness(t *testing.T, svc testjsonrpc.Service) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	h := testJsonrpctest.NewHarness(t, svc)
	defer h.Close()

	td := testJsonrpctest.NewTestData()
	t.Run("jsonrpc_no_stream", func(t *testing.T) {
		result, err := h.Client.JsonrpcNoStream(ctx, td.ValidJsonrpcNoStreamPayload())
		if err != nil {
			t.Errorf("jsonrpc_no_stream failed: %v", err)
		}
		if result == nil {
			t.Error("jsonrpc_no_stream returned nil result")
		}
	})
	t.Run("jsonrpc_no_stream_error", func(t *testing.T) {
		result, err := h.Client.JsonrpcNoStreamError(ctx, td.ValidJsonrpcNoStreamErrorPayload())
		if err != nil {
			t.Errorf("jsonrpc_no_stream_error failed: %v", err)
		}
		if result == nil {
			t.Error("jsonrpc_no_stream_error returned nil result")
		}
	})
}
