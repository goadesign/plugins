package testhttpgrpcapi

import (
	"context"
	"testing"

	testhttpgrpc "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc"
	testHTTPGrpctest "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc/test_http_grpctest"
)

// TestGRPCTransport tests all gRPC-specific endpoints
func TestGRPCTransport(t *testing.T) {
	svc := NewTestHTTPGrpc()
	harness := testHTTPGrpctest.NewHarness(t, svc)
	defer harness.Close()

	ctx := context.Background()
	td := testHTTPGrpctest.NewTestData()

	t.Run("Unary", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			payload := &testhttpgrpc.GrpcNoStreamPayload{Msg: "test"}
			res, err := harness.Client.GrpcNoStream(ctx, payload)
			if err != nil {
				t.Fatalf("GrpcNoStream failed: %v", err)
			}
			if res == nil {
				t.Error("expected result, got nil")
			}
		})

		t.Run("WithCustomError", func(t *testing.T) {
			// Test successful division
			payload := &testhttpgrpc.GrpcNoStreamErrorDivByZeroPayload{
				Dividend: 10,
				Divisor:  2,
			}
			res, err := harness.Client.GrpcNoStreamErrorDivByZero(ctx, payload)
			if err != nil {
				t.Fatalf("GrpcNoStreamErrorDivByZero failed: %v", err)
			}
			if res == nil || res.Quotient != 5 {
				t.Errorf("unexpected result: %v", res)
			}

			// Test division by zero error
			errorPayload := &testhttpgrpc.GrpcNoStreamErrorDivByZeroPayload{
				Dividend: 10,
				Divisor:  0,
			}
			_, err = harness.Client.GrpcNoStreamErrorDivByZero(ctx, errorPayload)
			if err == nil {
				t.Error("expected DivisionByZero error, got nil")
			}
			// Verify it's the correct error type
			if _, ok := err.(*testhttpgrpc.DivisionByZeroError); !ok {
				t.Errorf("expected DivisionByZeroError, got %T", err)
			}
		})
	})

	t.Run("ServerStreaming", func(t *testing.T) {
		stream, err := harness.Client.GrpcServerStream(ctx, td.ValidGrpcServerStreamPayload())
		if err != nil {
			t.Fatalf("GrpcServerStream failed: %v", err)
		}
		if stream != nil {
			// Server may not send data immediately
			_, _ = stream.Recv()
		}
	})

	t.Run("ClientStreaming", func(t *testing.T) {
		stream, err := harness.Client.GrpcClientStream(ctx)
		if err != nil {
			t.Fatalf("GrpcClientStream failed: %v", err)
		}
		if stream != nil {
			// Send multiple values
			values := []int{5, 10, 15}
			expectedSum := 30

			for _, v := range values {
				err := stream.Send(&testhttpgrpc.GrpcClientStreamStreamingPayload{
					Value: v,
				})
				if err != nil {
					t.Errorf("failed to send value %d: %v", v, err)
				}
			}

			res, err := stream.CloseAndRecv()
			if err != nil {
				t.Errorf("failed to close and receive: %v", err)
			}
			if res != nil && res.Sum != expectedSum {
				t.Errorf("expected sum %d, got %d", expectedSum, res.Sum)
			}
		}
	})

	t.Run("BidirectionalStreaming", func(t *testing.T) {
		stream, err := harness.Client.GrpcBidiStream(ctx)
		if err != nil {
			t.Fatalf("GrpcBidiStream failed: %v", err)
		}
		if stream != nil {
			// Send a message
			err := stream.Send(&testhttpgrpc.GrpcBidiStreamStreamingPayload{
				In: "test message",
			})
			if err != nil {
				t.Errorf("failed to send: %v", err)
			}

			// Try to receive (may get EOF if server doesn't echo)
			_, _ = stream.Recv()

			// Close the stream
			err = stream.Close()
			if err != nil {
				t.Errorf("failed to close: %v", err)
			}
		}
	})
}
