package testhttpgrpcapi

import (
	"context"
	"testing"

	testhttpgrpc "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc"
	testHTTPGrpctest "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc/test_http_grpctest"
)

// TestHTTPTransport tests all HTTP-specific endpoints
func TestHTTPTransport(t *testing.T) {
	svc := NewTestHTTPGrpc()
	harness := testHTTPGrpctest.NewHarness(t, svc)
	defer harness.Close()

	ctx := context.Background()

	t.Run("NonStreaming", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			payload := &testhttpgrpc.HTTPNoStreamPayload{Msg: "test"}
			res, err := harness.Client.HTTPNoStream(ctx, payload)
			if err != nil {
				t.Fatalf("HTTPNoStream failed: %v", err)
			}
			if res == nil {
				t.Error("expected result, got nil")
			}
		})

		t.Run("WithError", func(t *testing.T) {
			// Test successful case
			payload := &testhttpgrpc.HTTPNoStreamErrorPayload{Msg: "success"}
			res, err := harness.Client.HTTPNoStreamError(ctx, payload)
			if err != nil {
				t.Fatalf("HTTPNoStreamError failed: %v", err)
			}
			if res == nil || res.Out != "processed: success" {
				t.Errorf("unexpected result: %v", res)
			}

			// Test error case
			errorPayload := &testhttpgrpc.HTTPNoStreamErrorPayload{Msg: "error"}
			_, err = harness.Client.HTTPNoStreamError(ctx, errorPayload)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	})

	t.Run("ServerSentEvents", func(t *testing.T) {
		stream, err := harness.Client.HTTPServerStreamSse(ctx)
		if err != nil {
			t.Fatalf("HTTPServerStreamSse failed: %v", err)
		}
		if stream != nil {
			// SSE streams may not have data immediately
			_, _ = stream.Recv()
		}
	})

	t.Run("WebSocket", func(t *testing.T) {
		t.Run("ServerStream", func(t *testing.T) {
			stream, err := harness.Client.HTTPServerStreamWs(ctx)
			if err != nil {
				t.Fatalf("HTTPServerStreamWs failed: %v", err)
			}
			if stream != nil {
				res, _ := stream.Recv()
				if res != nil && res.Message != "test message" {
					t.Errorf("unexpected message: %s", res.Message)
				}
			}
		})

		t.Run("ClientStream", func(t *testing.T) {
			stream, err := harness.Client.HTTPClientStreamWs(ctx)
			if err != nil {
				t.Fatalf("HTTPClientStreamWs failed: %v", err)
			}
			if stream != nil {
				err := stream.Send(&testhttpgrpc.HTTPClientStreamWsStreamingPayload{
					Message: "test",
				})
				if err != nil {
					t.Errorf("failed to send: %v", err)
				}
				res, err := stream.CloseAndRecv()
				if err != nil {
					t.Errorf("failed to close and receive: %v", err)
				}
				if res != nil && res.Out != "received" {
					t.Errorf("unexpected result: %s", res.Out)
				}
			}
		})

		t.Run("Bidirectional", func(t *testing.T) {
			stream, err := harness.Client.HTTPBidiStreamWs(ctx)
			if err != nil {
				t.Fatalf("HTTPBidiStreamWs failed: %v", err)
			}
			if stream != nil {
				err := stream.Send(&testhttpgrpc.HTTPBidiStreamWsStreamingPayload{
					Message: "ping",
				})
				if err != nil {
					t.Errorf("failed to send: %v", err)
				}
				res, _ := stream.Recv()
				if res != nil && res.Message != "echo: ping" {
					t.Errorf("unexpected message: %s", res.Message)
				}
				_ = stream.Close()
			}
		})
	})
}
