package testhttpgrpcapi

import (
	"context"
	"testing"

	testhttpgrpc "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc"
	testHTTPGrpctest "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc/test_http_grpctest"
)

// TestErrorAssertions tests the error assertion helpers
func TestErrorAssertions(t *testing.T) {
	// Create a harness with the service implementation
	svc := NewTestHTTPGrpc()
	harness := testHTTPGrpctest.NewHarness(t, svc)
	defer harness.Close()

	// Create error asserter
	asserter := testHTTPGrpctest.NewErrorAsserter(t)

	ctx := context.Background()

	t.Run("HTTPNoStreamError_Success", func(t *testing.T) {
		// Test successful case (no error)
		payload := &testhttpgrpc.HTTPNoStreamErrorPayload{
			Msg: "success",
		}
		res, err := harness.Client.HTTPNoStreamError(ctx, payload)
		if err != nil {
			t.Fatalf("Expected success, got error: %v", err)
		}
		if res == nil || res.Out != "processed: success" {
			t.Errorf("Unexpected result: %v", res)
		}
	})

	t.Run("HTTPNoStreamError_DivisionByZero", func(t *testing.T) {
		// Test error case
		payload := &testhttpgrpc.HTTPNoStreamErrorPayload{
			Msg: "error",
		}
		_, err := harness.Client.HTTPNoStreamError(ctx, payload)
		asserter.AssertDivisionByZero(err)
	})

	t.Run("GrpcNoStreamError_Success", func(t *testing.T) {
		// Test successful division
		payload := &testhttpgrpc.GrpcNoStreamErrorDivByZeroPayload{
			Dividend: 10,
			Divisor:  2,
		}
		res, err := harness.Client.GrpcNoStreamErrorDivByZero(ctx, payload)
		if err != nil {
			t.Fatalf("Expected success, got error: %v", err)
		}
		if res == nil || res.Quotient != 5 {
			t.Errorf("Unexpected result: %v", res)
		}
	})

	t.Run("GrpcNoStreamError_DivisionByZero", func(t *testing.T) {
		// Test division by zero
		payload := &testhttpgrpc.GrpcNoStreamErrorDivByZeroPayload{
			Dividend: 10,
			Divisor:  0,
		}
		_, err := harness.Client.GrpcNoStreamErrorDivByZero(ctx, payload)
		asserter.AssertDivisionByZero(err)
	})

	t.Run("ExpectDivisionByZero", func(t *testing.T) {
		// Test using ExpectDivisionByZero helper with HTTP transport
		asserter.ExpectDivisionByZero(func() error {
			payload := &testhttpgrpc.HTTPNoStreamErrorPayload{
				Msg: "error",
			}
			_, err := harness.Client.HTTPNoStreamError(ctx, payload)
			return err
		})
	})

	t.Run("TransportSpecificErrors_HTTP", func(t *testing.T) {
		// Force HTTP transport and test error
		payload := &testhttpgrpc.HTTPNoStreamErrorPayload{
			Msg: "error",
		}
		_, err := harness.Client.HTTP().HTTPNoStreamError(ctx, payload)
		asserter.AssertDivisionByZero(err)
	})

	t.Run("TransportSpecificErrors_GRPC", func(t *testing.T) {
		// Force gRPC transport and test error
		payload := &testhttpgrpc.GrpcNoStreamErrorDivByZeroPayload{
			Dividend: 42,
			Divisor:  0,
		}
		_, err := harness.Client.GRPC().GrpcNoStreamErrorDivByZero(ctx, payload)
		asserter.AssertDivisionByZero(err)
	})
}
