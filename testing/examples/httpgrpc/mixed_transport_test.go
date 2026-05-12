package testhttpgrpcapi

import (
	"context"
	"io"
	"testing"

	testhttpgrpc "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc"
	testHTTPGrpctest "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc/test_http_grpctest"
)

// TestMixedTransport tests endpoints that support both HTTP and gRPC transports
func TestMixedTransport(t *testing.T) {
	svc := NewTestHTTPGrpc()
	harness := testHTTPGrpctest.NewHarness(t, svc)
	defer harness.Close()

	ctx := context.Background()

	t.Run("NonStreaming", func(t *testing.T) {
		payload := &testhttpgrpc.MixedNoStreamPayload{Msg: "test"}

		t.Run("AutoTransport", func(t *testing.T) {
			// Should use first available transport (HTTP)
			res, err := harness.Client.MixedNoStream(ctx, payload)
			if err != nil {
				t.Fatalf("MixedNoStream (auto) failed: %v", err)
			}
			if res == nil {
				t.Error("expected result, got nil")
			}
		})

		t.Run("ExplicitHTTP", func(t *testing.T) {
			res, err := harness.Client.HTTP().MixedNoStream(ctx, payload)
			if err != nil {
				t.Fatalf("MixedNoStream (HTTP) failed: %v", err)
			}
			if res == nil {
				t.Error("expected result, got nil")
			}
		})

		t.Run("ExplicitGRPC", func(t *testing.T) {
			res, err := harness.Client.GRPC().MixedNoStream(ctx, payload)
			if err != nil {
				t.Fatalf("MixedNoStream (gRPC) failed: %v", err)
			}
			if res == nil {
				t.Error("expected result, got nil")
			}
		})
	})

	t.Run("ServerStreaming", func(t *testing.T) {
		t.Run("AutoTransport", func(t *testing.T) {
			stream, err := harness.Client.MixedServerStream(ctx)
			if err != nil {
				t.Fatalf("MixedServerStream (auto) failed: %v", err)
			}
			if stream != nil {
				_, _ = stream.Recv()
			}
		})

		t.Run("HTTPServerSentEvents", func(t *testing.T) {
			stream, err := harness.Client.HTTP().MixedServerStream(ctx)
			if err != nil {
				t.Fatalf("MixedServerStream (HTTP SSE) failed: %v", err)
			}
			if stream != nil {
				_, _ = stream.Recv()
			}
		})

		t.Run("GRPCServerStream", func(t *testing.T) {
			stream, err := harness.Client.GRPC().MixedServerStream(ctx)
			if err != nil {
				t.Fatalf("MixedServerStream (gRPC) failed: %v", err)
			}
			if stream != nil {
				_, _ = stream.Recv()
			}
		})
	})

	t.Run("ClientStreaming", func(t *testing.T) {
		sendAndReceive := func(stream testhttpgrpc.MixedClientStreamWsGrpcClientStream) error {
			err := stream.Send(&testhttpgrpc.MixedClientStreamWsGrpcStreamingPayload{
				Message: "test",
			})
			if err != nil {
				return err
			}
			res, err := stream.CloseAndRecv()
			if err != nil {
				return err
			}
			if res != nil && res.Out != "received" {
				t.Errorf("unexpected result: %s", res.Out)
			}
			return nil
		}

		t.Run("AutoTransport", func(t *testing.T) {
			stream, err := harness.Client.MixedClientStreamWsGrpc(ctx)
			if err != nil {
				t.Fatalf("MixedClientStreamWsGrpc (auto) failed: %v", err)
			}
			if stream != nil {
				if err := sendAndReceive(stream); err != nil {
					t.Errorf("stream operation failed: %v", err)
				}
			}
		})

		t.Run("HTTPWebSocket", func(t *testing.T) {
			stream, err := harness.Client.HTTP().MixedClientStreamWsGrpc(ctx)
			if err != nil {
				t.Fatalf("MixedClientStreamWsGrpc (HTTP) failed: %v", err)
			}
			if stream != nil {
				if err := sendAndReceive(stream); err != nil {
					t.Errorf("stream operation failed: %v", err)
				}
			}
		})

		t.Run("GRPCClientStream", func(t *testing.T) {
			stream, err := harness.Client.GRPC().MixedClientStreamWsGrpc(ctx)
			if err != nil {
				t.Fatalf("MixedClientStreamWsGrpc (gRPC) failed: %v", err)
			}
			if stream != nil {
				if err := sendAndReceive(stream); err != nil {
					t.Errorf("stream operation failed: %v", err)
				}
			}
		})
	})

	t.Run("BidirectionalStreaming", func(t *testing.T) {
		sendReceiveClose := func(stream testhttpgrpc.MixedBidiStreamWsGrpcClientStream) error {
			err := stream.Send(&testhttpgrpc.MixedBidiStreamWsGrpcStreamingPayload{
				Message: "ping",
			})
			if err != nil {
				return err
			}
			res, err := stream.Recv()
			if err != nil && err != io.EOF {
				return err
			}
			if res != nil && res.Message != "echo: ping" {
				t.Errorf("unexpected message: %s", res.Message)
			}
			return stream.Close()
		}

		t.Run("AutoTransport", func(t *testing.T) {
			stream, err := harness.Client.MixedBidiStreamWsGrpc(ctx)
			if err != nil {
				t.Fatalf("MixedBidiStreamWsGrpc (auto) failed: %v", err)
			}
			if stream != nil {
				if err := sendReceiveClose(stream); err != nil {
					t.Errorf("stream operation failed: %v", err)
				}
			}
		})

		t.Run("HTTPWebSocket", func(t *testing.T) {
			stream, err := harness.Client.HTTP().MixedBidiStreamWsGrpc(ctx)
			if err != nil {
				t.Fatalf("MixedBidiStreamWsGrpc (HTTP) failed: %v", err)
			}
			if stream != nil {
				if err := sendReceiveClose(stream); err != nil {
					t.Errorf("stream operation failed: %v", err)
				}
			}
		})

		t.Run("GRPCBidirectional", func(t *testing.T) {
			stream, err := harness.Client.GRPC().MixedBidiStreamWsGrpc(ctx)
			if err != nil {
				t.Fatalf("MixedBidiStreamWsGrpc (gRPC) failed: %v", err)
			}
			if stream != nil {
				if err := sendReceiveClose(stream); err != nil {
					t.Errorf("stream operation failed: %v", err)
				}
			}
		})
	})
}
