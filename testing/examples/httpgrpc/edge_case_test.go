package testhttpgrpcapi

import (
	"context"
	"testing"

	testhttpgrpc "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc"
	testHTTPGrpctest "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc/test_http_grpctest"
)

// TestEdgeCaseGenerators tests the edge case generation functionality
func TestEdgeCaseGenerators(t *testing.T) {
	svc := NewTestHTTPGrpc()
	harness := testHTTPGrpctest.NewHarness(t, svc)
	defer harness.Close()

	ctx := context.Background()
	td := testHTTPGrpctest.NewTestData()

	t.Run("ValidPayloadWorks", func(t *testing.T) {
		// First verify that the valid payload generator works
		payload := td.ValidGrpcNoStreamErrorDivByZeroPayload()
		if payload == nil {
			t.Fatal("expected valid payload, got nil")
		}
		
		// The valid payload should have non-zero divisor
		if payload.Divisor == 0 {
			t.Error("valid payload should not have zero divisor")
		}
		
		// Should succeed with valid payload
		res, err := harness.Client.GrpcNoStreamErrorDivByZero(ctx, payload)
		if err != nil {
			t.Fatalf("unexpected error with valid payload: %v", err)
		}
		if res == nil {
			t.Error("expected result with valid payload")
		}
	})

	t.Run("DivisionByZeroEdgeCase", func(t *testing.T) {
		// NOTE: The testing framework doesn't generate error-triggering payloads
		// because the design doesn't specify what input values cause errors.
		// We manually create a zero divisor payload for testing.
		payload := td.ValidGrpcNoStreamErrorDivByZeroPayload()
		payload.Divisor = 0 // Manually set to trigger the error
		
		_, err := harness.Client.GrpcNoStreamErrorDivByZero(ctx, payload)
		if err == nil {
			t.Error("expected division by zero error")
		}
		
		// Verify it's the right error type
		if _, ok := err.(*testhttpgrpc.DivisionByZeroError); !ok {
			t.Errorf("expected DivisionByZeroError, got %T", err)
		}
	})

	t.Run("MinValuesEdgeCase", func(t *testing.T) {
		// Use the generated edge case method
		payload := td.GrpcNoStreamErrorDivByZeroPayloadWithMinValues()
		
		// Verify minimum values were applied
		if payload.Dividend != -1000.0 {
			t.Errorf("expected dividend -1000, got %f", payload.Dividend)
		}
		if payload.Divisor != -100.0 {
			t.Errorf("expected divisor -100, got %f", payload.Divisor)
		}
		
		res, err := harness.Client.GrpcNoStreamErrorDivByZero(ctx, payload)
		if err != nil {
			t.Fatalf("unexpected error with min values: %v", err)
		}
		if res == nil {
			t.Error("expected result with min values")
		}
		
		// Result should be 10.0 (-1000 / -100)
		if res.Quotient != 10.0 {
			t.Errorf("expected quotient 10.0, got %f", res.Quotient)
		}
	})

	t.Run("MaxValuesEdgeCase", func(t *testing.T) {
		// Use the generated edge case method
		payload := td.GrpcNoStreamErrorDivByZeroPayloadWithMaxValues()
		
		// Verify maximum values were applied
		if payload.Dividend != 1000.0 {
			t.Errorf("expected dividend 1000, got %f", payload.Dividend)
		}
		if payload.Divisor != 100.0 {
			t.Errorf("expected divisor 100, got %f", payload.Divisor)
		}
		
		res, err := harness.Client.GrpcNoStreamErrorDivByZero(ctx, payload)
		if err != nil {
			t.Fatalf("unexpected error with max values: %v", err)
		}
		if res == nil {
			t.Error("expected result with max values")
		}
		
		// Result should be 10.0 (1000 / 100)
		if res.Quotient != 10.0 {
			t.Errorf("expected quotient 10.0, got %f", res.Quotient)
		}
	})
}