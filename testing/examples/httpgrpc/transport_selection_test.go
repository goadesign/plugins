package testhttpgrpcapi

import (
	"context"
	"io"
	"testing"

	testhttpgrpc "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc"
	testHTTPGrpctest "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc/test_http_grpctest"
)

// TestTransportSelection tests the fluent API for explicit transport selection
// using the .HTTP() and .GRPC() methods to force specific transports
func TestTransportSelection(t *testing.T) {
	// Create a harness with the service implementation
	svc := NewTestHTTPGrpc()
	harness := testHTTPGrpctest.NewHarness(t, svc)
	defer harness.Close()

	ctx := context.Background()

	t.Run("ForceHTTP", func(t *testing.T) {
		// Test forcing HTTP transport for mixed endpoint
		payload := &testhttpgrpc.MixedNoStreamPayload{
			Msg: "test HTTP",
		}
		
		res, err := harness.Client.HTTP().MixedNoStream(ctx, payload)
		if err != nil {
			t.Fatalf("HTTP forced call failed: %v", err)
		}
		
		if res == nil {
			t.Error("expected result, got nil")
		}
	})

	t.Run("ForceGRPC", func(t *testing.T) {
		// Test forcing gRPC transport for mixed endpoint
		payload := &testhttpgrpc.MixedNoStreamPayload{
			Msg: "test gRPC",
		}
		
		res, err := harness.Client.GRPC().MixedNoStream(ctx, payload)
		if err != nil {
			t.Fatalf("gRPC forced call failed: %v", err)
		}
		
		if res == nil {
			t.Error("expected result, got nil")
		}
	})

	t.Run("ForceHTTPStreaming", func(t *testing.T) {
		// Test forcing HTTP SSE for mixed server stream
		stream, err := harness.Client.HTTP().MixedServerStream(ctx)
		if err != nil {
			t.Fatalf("HTTP SSE forced call failed: %v", err)
		}
		
		if stream != nil {
			// Try to receive from the stream
			_, err := stream.Recv()
			if err != nil && err != io.EOF {
				t.Errorf("unexpected error from stream: %v", err)
			}
		}
	})

	t.Run("ForceGRPCStreaming", func(t *testing.T) {
		// Test forcing gRPC for mixed server stream
		stream, err := harness.Client.GRPC().MixedServerStream(ctx)
		if err != nil {
			t.Fatalf("gRPC stream forced call failed: %v", err)
		}
		
		if stream != nil {
			// Try to receive from the stream
			_, err := stream.Recv()
			if err != nil && err != io.EOF {
				t.Errorf("unexpected error from stream: %v", err)
			}
		}
	})
	
	t.Run("MixedClientStreamHTTP", func(t *testing.T) {
		// Test forcing HTTP WebSocket for mixed client stream
		stream, err := harness.Client.HTTP().MixedClientStreamWsGrpc(ctx)
		if err != nil {
			t.Fatalf("HTTP WebSocket client stream failed: %v", err)
		}
		
		if stream != nil {
			// Send some data and close
			err := stream.Send(&testhttpgrpc.MixedClientStreamWsGrpcStreamingPayload{
				Message: "test message",
			})
			if err != nil {
				t.Errorf("failed to send: %v", err)
			}
			
			_, err = stream.CloseAndRecv()
			if err != nil {
				t.Errorf("failed to close and receive: %v", err)
			}
		}
	})
	
	t.Run("MixedClientStreamGRPC", func(t *testing.T) {
		// Test forcing gRPC for mixed client stream
		stream, err := harness.Client.GRPC().MixedClientStreamWsGrpc(ctx)
		if err != nil {
			t.Fatalf("gRPC client stream failed: %v", err)
		}
		
		if stream != nil {
			// Send some data and close
			err := stream.Send(&testhttpgrpc.MixedClientStreamWsGrpcStreamingPayload{
				Message: "test message",
			})
			if err != nil {
				t.Errorf("failed to send: %v", err)
			}
			
			_, err = stream.CloseAndRecv()
			if err != nil {
				t.Errorf("failed to close and receive: %v", err)
			}
		}
	})
	
	t.Run("MixedBidiStreamHTTP", func(t *testing.T) {
		// Test forcing HTTP WebSocket for mixed bidirectional stream
		stream, err := harness.Client.HTTP().MixedBidiStreamWsGrpc(ctx)
		if err != nil {
			t.Fatalf("HTTP WebSocket bidi stream failed: %v", err)
		}
		
		if stream != nil {
			// Send and receive
			err := stream.Send(&testhttpgrpc.MixedBidiStreamWsGrpcStreamingPayload{
				Message: "ping",
			})
			if err != nil {
				t.Errorf("failed to send: %v", err)
			}
			
			_, err = stream.Recv()
			if err != nil && err != io.EOF {
				t.Errorf("failed to receive: %v", err)
			}
			
			err = stream.Close()
			if err != nil {
				t.Errorf("failed to close: %v", err)
			}
		}
	})
	
	t.Run("MixedBidiStreamGRPC", func(t *testing.T) {
		// Test forcing gRPC for mixed bidirectional stream
		stream, err := harness.Client.GRPC().MixedBidiStreamWsGrpc(ctx)
		if err != nil {
			t.Fatalf("gRPC bidi stream failed: %v", err)
		}
		
		if stream != nil {
			// Send and receive
			err := stream.Send(&testhttpgrpc.MixedBidiStreamWsGrpcStreamingPayload{
				Message: "ping",
			})
			if err != nil {
				t.Errorf("failed to send: %v", err)
			}
			
			_, err = stream.Recv()
			if err != nil && err != io.EOF {
				t.Errorf("failed to receive: %v", err)
			}
			
			err = stream.Close()
			if err != nil {
				t.Errorf("failed to close: %v", err)
			}
		}
	})
}