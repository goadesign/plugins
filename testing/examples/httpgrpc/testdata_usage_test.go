package testhttpgrpcapi

import (
	"context"
	"testing"

	testhttpgrpc "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc"
	testHTTPGrpctest "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc/test_http_grpctest"
)

// TestDataGenerators demonstrates the use of generated test data
func TestDataGenerators(t *testing.T) {
	svc := NewTestHTTPGrpc()
	harness := testHTTPGrpctest.NewHarness(t, svc)
	defer harness.Close()

	ctx := context.Background()
	td := testHTTPGrpctest.NewTestData()

	t.Run("ValidPayloadGenerators", func(t *testing.T) {
		// Test that valid payload generators work
		t.Run("HTTPNoStream", func(t *testing.T) {
			payload := td.ValidHTTPNoStreamPayload()
			if payload == nil {
				t.Fatal("expected valid payload, got nil")
			}
			if payload.Msg == "" {
				t.Error("expected non-empty message in payload")
			}

			// Use the valid payload in a call
			res, err := harness.Client.HTTPNoStream(ctx, payload)
			if err != nil {
				t.Fatalf("call failed with valid payload: %v", err)
			}
			if res == nil {
				t.Error("expected result, got nil")
			}
		})

		t.Run("GrpcNoStreamErrorDivByZero", func(t *testing.T) {
			payload := td.ValidGrpcNoStreamErrorDivByZeroPayload()
			if payload == nil {
				t.Fatal("expected valid payload, got nil")
			}
			// The generated valid payload should have non-zero divisor
			if payload.Divisor == 0 {
				t.Error("valid payload should not have zero divisor")
			}

			// Use the valid payload - should succeed
			res, err := harness.Client.GrpcNoStreamErrorDivByZero(ctx, payload)
			if err != nil {
				t.Fatalf("call failed with valid payload: %v", err)
			}
			if res == nil {
				t.Error("expected result, got nil")
			}
		})
	})

	t.Run("PayloadBuilders", func(t *testing.T) {
		// Test the builder pattern for customizing payloads
		t.Run("CustomizeWithBuilder", func(t *testing.T) {
			// Start with a valid payload and customize it
			builder := td.NewHTTPNoStreamPayloadBuilder()
			payload := builder.Build()

			// Manually customize after building
			payload.Msg = "custom message"

			res, err := harness.Client.HTTPNoStream(ctx, payload)
			if err != nil {
				t.Fatalf("call failed: %v", err)
			}
			if res == nil {
				t.Error("expected result, got nil")
			}
		})
	})

	t.Run("WithAllFieldsGenerators", func(t *testing.T) {
		// Test generators that populate all fields
		payload := td.HTTPNoStreamPayloadWithAllFields()
		if payload == nil {
			t.Fatal("expected payload with all fields, got nil")
		}
		if payload.Msg == "" {
			t.Error("expected all fields to be populated")
		}
	})

	// NOTE: Result generators are not needed conceptually
	// The testing framework tests service implementations which produce results.
	// Tests send payloads (inputs) and verify results (outputs) from the service.

	// NOTE: Edge case generators should also be available
	t.Run("MissingEdgeCaseGenerators", func(t *testing.T) {
		t.Skip("Edge case generators are not implemented yet")

		// This is what SHOULD be available:
		// payload := td.GrpcNoStreamErrorDivByZeroPayloadWithZeroDivisor()
		// payload := td.HTTPNoStreamPayloadWithEmptyString()
		// payload := td.HTTPNoStreamPayloadWithMaxLengthString()
		// payload := td.GrpcServerStreamPayloadWithMinValue()
		// payload := td.GrpcServerStreamPayloadWithMaxValue()
	})
}

// TestDataForStreamingMethods demonstrates test data for streaming methods
func TestDataForStreamingMethods(t *testing.T) {
	svc := NewTestHTTPGrpc()
	harness := testHTTPGrpctest.NewHarness(t, svc)
	defer harness.Close()

	ctx := context.Background()
	td := testHTTPGrpctest.NewTestData()

	t.Run("StreamingPayloads", func(t *testing.T) {
		t.Run("ClientStream", func(t *testing.T) {
			stream, err := harness.Client.GrpcClientStream(ctx)
			if err != nil {
				t.Fatalf("failed to create stream: %v", err)
			}

			// Use valid streaming payload generator
			payload := td.ValidGrpcClientStreamPayload()
			if payload == nil {
				t.Fatal("expected valid streaming payload, got nil")
			}

			// Send multiple payloads
			for i := 0; i < 3; i++ {
				p := td.ValidGrpcClientStreamPayload()
				if err := stream.Send(p); err != nil {
					t.Errorf("failed to send payload %d: %v", i, err)
				}
			}

			res, err := stream.CloseAndRecv()
			if err != nil {
				t.Errorf("failed to close and receive: %v", err)
			}
			if res == nil {
				t.Error("expected result, got nil")
			}
		})

		t.Run("BidirectionalStream", func(t *testing.T) {
			stream, err := harness.Client.GrpcBidiStream(ctx)
			if err != nil {
				t.Fatalf("failed to create stream: %v", err)
			}

			// Use valid streaming payload
			payload := td.ValidGrpcBidiStreamPayload()
			if payload == nil {
				t.Fatal("expected valid streaming payload, got nil")
			}

			if err := stream.Send(payload); err != nil {
				t.Errorf("failed to send: %v", err)
			}

			// Try to receive
			_, _ = stream.Recv()

			if err := stream.Close(); err != nil {
				t.Errorf("failed to close: %v", err)
			}
		})
	})

	// NOTE: Streaming result generators are not needed either
	// Server streams produce results continuously - tests verify the stream behavior
}

// TestInterestingDataPoints demonstrates what interesting data points should be generated
func TestInterestingDataPoints(t *testing.T) {
	t.Run("NumericEdgeCases", func(t *testing.T) {
		t.Skip("Not implemented yet - should generate min/max/zero values")

		// Should have generators for numeric edge cases:
		// - Zero values
		// - Minimum values (negative for signed types)
		// - Maximum values
		// - Common edge cases (1, -1, etc.)

		// Example of what should be available:
		// payload := td.GrpcServerStreamPayloadWithZeroFrom()
		// payload := td.GrpcServerStreamPayloadWithMaxFrom()
		// payload := td.GrpcServerStreamPayloadWithMinFrom()
	})

	t.Run("StringEdgeCases", func(t *testing.T) {
		t.Skip("Not implemented yet - should generate empty/long/special strings")

		// Should have generators for string edge cases:
		// - Empty strings
		// - Single character
		// - Very long strings (if max length defined)
		// - Strings with special characters
		// - Unicode strings

		// Example of what should be available:
		// payload := td.HTTPNoStreamPayloadWithEmptyMsg()
		// payload := td.HTTPNoStreamPayloadWithLongMsg()
		// payload := td.HTTPNoStreamPayloadWithUnicodeMsg()
	})

	t.Run("OptionalFields", func(t *testing.T) {
		t.Skip("Not implemented yet - should handle optional fields")

		// Should have generators for optional field combinations:
		// - All optional fields nil
		// - All optional fields populated
		// - Various combinations

		// Example of what should be available:
		// payload := td.SomePayloadWithNoOptionalFields()
		// payload := td.SomePayloadWithAllOptionalFields()
	})

	t.Run("ValidationBoundaries", func(t *testing.T) {
		t.Skip("Not implemented yet - should test validation boundaries")

		// Should generate data at validation boundaries:
		// - Minimum length strings
		// - Maximum length strings
		// - Values at min/max constraints
		// - Pattern-matching edge cases

		// Example of what should be available:
		// payload := td.SomePayloadAtMinValidation()
		// payload := td.SomePayloadAtMaxValidation()
	})
}

// TestCustomTestData shows how users can extend the generated test data
func TestCustomTestData(t *testing.T) {
	td := testHTTPGrpctest.NewTestData()

	// Users can create custom test data by starting with valid data
	// and modifying it for their specific test cases

	t.Run("CreateInvalidData", func(t *testing.T) {
		// Start with valid data
		payload := td.ValidGrpcNoStreamErrorDivByZeroPayload()

		// Modify to create invalid data for error testing
		payload.Divisor = 0 // Create division by zero scenario

		// Now we have a payload that will trigger an error
		svc := NewTestHTTPGrpc()
		harness := testHTTPGrpctest.NewHarness(t, svc)
		defer harness.Close()

		_, err := harness.Client.GrpcNoStreamErrorDivByZero(context.Background(), payload)
		if err == nil {
			t.Error("expected division by zero error")
		}
	})

	t.Run("CreateTestSequence", func(t *testing.T) {
		// Create a sequence of related test data
		payloads := []*testhttpgrpc.GrpcClientStreamStreamingPayload{
			td.ValidGrpcClientStreamPayload(),
			td.ValidGrpcClientStreamPayload(),
			td.ValidGrpcClientStreamPayload(),
		}

		// Modify to create a specific sequence
		for i := range payloads {
			payloads[i].Value = i + 1 // Sequential values: 1, 2, 3
		}

		// Use in test...
		if len(payloads) != 3 {
			t.Error("expected 3 payloads")
		}
	})
}
