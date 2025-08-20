package testhttpgrpcapi

import (
	"context"
	"io"
	"testing"
	"time"

	testhttpgrpc "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc"
	testHTTPGrpctest "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc/test_http_grpctest"
)

// mockService implements a minimal test service for stream testing
type mockService struct{}

func (s *mockService) HTTPNoStream(ctx context.Context, p *testhttpgrpc.HTTPNoStreamPayload) (*testhttpgrpc.HTTPNoStreamResult, error) {
	return &testhttpgrpc.HTTPNoStreamResult{Out: "test"}, nil
}

func (s *mockService) HTTPNoStreamError(ctx context.Context, p *testhttpgrpc.HTTPNoStreamErrorPayload) (*testhttpgrpc.HTTPNoStreamErrorResult, error) {
	return &testhttpgrpc.HTTPNoStreamErrorResult{Out: "test"}, nil
}

func (s *mockService) GrpcNoStreamErrorDivByZero(ctx context.Context, p *testhttpgrpc.GrpcNoStreamErrorDivByZeroPayload) (*testhttpgrpc.GrpcNoStreamErrorDivByZeroResult, error) {
	if p.Divisor == 0 {
		return nil, &testhttpgrpc.DivisionByZeroError{Message: "division by zero"}
	}
	return &testhttpgrpc.GrpcNoStreamErrorDivByZeroResult{Quotient: float64(p.Dividend) / float64(p.Divisor)}, nil
}

func (s *mockService) HTTPServerStreamSse(ctx context.Context, stream testhttpgrpc.HTTPServerStreamSseServerStream) error {
	// Send one message and return (returning closes the SSE stream)
	return stream.Send(&testhttpgrpc.HTTPServerStreamSseResult{Event: "test"})
}

func (s *mockService) HTTPServerStreamWs(ctx context.Context, stream testhttpgrpc.HTTPServerStreamWsServerStream) error {
	// Send one message and return
	return stream.Send(&testhttpgrpc.HTTPServerStreamWsResult{Message: "test"})
}

func (s *mockService) HTTPClientStreamWs(ctx context.Context, stream testhttpgrpc.HTTPClientStreamWsServerStream) error {
	// Read until client closes, then send response
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&testhttpgrpc.HTTPClientStreamWsResult{Out: "test"})
		}
		if err != nil {
			return err
		}
	}
}

func (s *mockService) HTTPBidiStreamWs(ctx context.Context, stream testhttpgrpc.HTTPBidiStreamWsServerStream) error {
	// Echo back one message
	msg, err := stream.Recv()
	if err != nil && err != io.EOF {
		return err
	}
	if msg != nil {
		return stream.Send(&testhttpgrpc.HTTPBidiStreamWsResult{Message: msg.Message})
	}
	return stream.Close()
}

func (s *mockService) GrpcNoStream(ctx context.Context, p *testhttpgrpc.GrpcNoStreamPayload) (*testhttpgrpc.GrpcNoStreamResult, error) {
	return &testhttpgrpc.GrpcNoStreamResult{Out: "test"}, nil
}

func (s *mockService) GrpcServerStream(ctx context.Context, p *testhttpgrpc.GrpcServerStreamPayload, stream testhttpgrpc.GrpcServerStreamServerStream) error {
	// Send one message and return
	return stream.Send(&testhttpgrpc.GrpcServerStreamResult{Value: 42})
}

func (s *mockService) GrpcClientStream(ctx context.Context, stream testhttpgrpc.GrpcClientStreamServerStream) error {
	// Read until client closes, then send response
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&testhttpgrpc.GrpcClientStreamResult{Sum: 100})
		}
		if err != nil {
			return err
		}
	}
}

func (s *mockService) GrpcBidiStream(ctx context.Context, stream testhttpgrpc.GrpcBidiStreamServerStream) error {
	// Echo back one message
	msg, err := stream.Recv()
	if err != nil && err != io.EOF {
		return err
	}
	if msg != nil {
		return stream.Send(&testhttpgrpc.GrpcBidiStreamResult{Out: "test"})
	}
	return nil
}

func (s *mockService) MixedNoStream(ctx context.Context, p *testhttpgrpc.MixedNoStreamPayload) (*testhttpgrpc.MixedNoStreamResult, error) {
	return &testhttpgrpc.MixedNoStreamResult{Out: "test"}, nil
}

func (s *mockService) MixedServerStream(ctx context.Context, stream testhttpgrpc.MixedServerStreamServerStream) error {
	// Send one message and return
	return stream.Send(&testhttpgrpc.MixedServerStreamResult{Event: "test"})
}

func (s *mockService) MixedClientStreamWsGrpc(ctx context.Context, stream testhttpgrpc.MixedClientStreamWsGrpcServerStream) error {
	// Read until client closes, then send response
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&testhttpgrpc.MixedClientStreamWsGrpcResult{Out: "test"})
		}
		if err != nil {
			return err
		}
	}
}

func (s *mockService) MixedBidiStreamWsGrpc(ctx context.Context, stream testhttpgrpc.MixedBidiStreamWsGrpcServerStream) error {
	// Echo back one message
	msg, err := stream.Recv()
	if err != nil && err != io.EOF {
		return err
	}
	if msg != nil {
		return stream.Send(&testhttpgrpc.MixedBidiStreamWsGrpcResult{Message: msg.Message})
	}
	return stream.Close()
}

func TestStreamHarnessBasics(t *testing.T) {
	svc := NewTestHTTPGrpc()
	h := testHTTPGrpctest.NewHarness(t, svc)
	defer h.Close()
	_ = context.Background()
}

// TestStreamMethodsCompile verifies that the generated test suite correctly
// handles different stream types. This test will fail to compile if the
// generated suite calls the wrong methods on streams.
func TestStreamMethodsCompile(t *testing.T) {
	// Setup the mock service
	service := &mockService{}

	// Create test harness
	h := testHTTPGrpctest.NewHarness(t, service)
	defer h.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	td := testHTTPGrpctest.NewTestData()

	// Test server streams - should only have Recv() method
	t.Run("ServerStreams", func(t *testing.T) {
		// HTTP SSE server stream
		stream1, err := h.Client.HTTPServerStreamSse(ctx)
		if err != nil {
			t.Fatalf("Failed to create HTTPServerStreamSse: %v", err)
		}
		// This should compile - server streams have Recv()
		_, err = stream1.Recv()
		if err != nil && err != io.EOF {
			t.Errorf("Recv failed: %v", err)
		}

		// HTTP WebSocket server stream
		stream2, err := h.Client.HTTPServerStreamWs(ctx)
		if err != nil {
			t.Fatalf("Failed to create HTTPServerStreamWs: %v", err)
		}
		_, err = stream2.Recv()
		if err != nil && err != io.EOF {
			t.Errorf("Recv failed: %v", err)
		}

		// gRPC server stream
		stream3, err := h.Client.GrpcServerStream(ctx, td.ValidGrpcServerStreamPayload())
		if err != nil {
			t.Fatalf("Failed to create GrpcServerStream: %v", err)
		}
		_, err = stream3.Recv()
		if err != nil && err != io.EOF {
			t.Errorf("Recv failed: %v", err)
		}

		// Mixed transport server stream
		stream4, err := h.Client.MixedServerStream(ctx)
		if err != nil {
			t.Fatalf("Failed to create MixedServerStream: %v", err)
		}
		_, err = stream4.Recv()
		if err != nil && err != io.EOF {
			t.Errorf("Recv failed: %v", err)
		}
	})

	// Test client streams - should have Send() and CloseAndRecv() methods
	t.Run("ClientStreams", func(t *testing.T) {
		// HTTP WebSocket client stream
		stream1, err := h.Client.HTTPClientStreamWs(ctx)
		if err != nil {
			t.Fatalf("Failed to create HTTPClientStreamWs: %v", err)
		}
		// Send a message
		err = stream1.Send(&testhttpgrpc.HTTPClientStreamWsStreamingPayload{Message: "test"})
		if err != nil {
			t.Errorf("Send failed: %v", err)
		}
		// This should compile - client streams have CloseAndRecv()
		_, err = stream1.CloseAndRecv()
		if err != nil {
			t.Errorf("CloseAndRecv failed: %v", err)
		}

		// gRPC client stream
		stream2, err := h.Client.GrpcClientStream(ctx)
		if err != nil {
			t.Fatalf("Failed to create GrpcClientStream: %v", err)
		}
		err = stream2.Send(&testhttpgrpc.GrpcClientStreamStreamingPayload{Value: 10})
		if err != nil {
			t.Errorf("Send failed: %v", err)
		}
		_, err = stream2.CloseAndRecv()
		if err != nil {
			t.Errorf("CloseAndRecv failed: %v", err)
		}

		// Mixed transport client stream
		stream3, err := h.Client.MixedClientStreamWsGrpc(ctx)
		if err != nil {
			t.Fatalf("Failed to create MixedClientStreamWsGrpc: %v", err)
		}
		err = stream3.Send(&testhttpgrpc.MixedClientStreamWsGrpcStreamingPayload{Message: "test"})
		if err != nil {
			t.Errorf("Send failed: %v", err)
		}
		_, err = stream3.CloseAndRecv()
		if err != nil {
			t.Errorf("CloseAndRecv failed: %v", err)
		}
	})

	// Test bidirectional streams - should have Send(), Recv(), and Close() methods
	t.Run("BidirectionalStreams", func(t *testing.T) {
		// HTTP WebSocket bidirectional stream
		stream1, err := h.Client.HTTPBidiStreamWs(ctx)
		if err != nil {
			t.Fatalf("Failed to create HTTPBidiStreamWs: %v", err)
		}
		// Send a message
		err = stream1.Send(&testhttpgrpc.HTTPBidiStreamWsStreamingPayload{Message: "test"})
		if err != nil {
			t.Errorf("Send failed: %v", err)
		}
		// Receive a message
		_, err = stream1.Recv()
		if err != nil && err != io.EOF {
			t.Errorf("Recv failed: %v", err)
		}
		// This should compile - bidirectional streams have Close()
		err = stream1.Close()
		if err != nil {
			t.Errorf("Close failed: %v", err)
		}

		// gRPC bidirectional stream
		stream2, err := h.Client.GrpcBidiStream(ctx)
		if err != nil {
			t.Fatalf("Failed to create GrpcBidiStream: %v", err)
		}
		err = stream2.Send(&testhttpgrpc.GrpcBidiStreamStreamingPayload{In: "test"})
		if err != nil {
			t.Errorf("Send failed: %v", err)
		}
		_, err = stream2.Recv()
		if err != nil && err != io.EOF {
			t.Errorf("Recv failed: %v", err)
		}
		err = stream2.Close()
		if err != nil {
			t.Errorf("Close failed: %v", err)
		}

		// Mixed transport bidirectional stream
		stream3, err := h.Client.MixedBidiStreamWsGrpc(ctx)
		if err != nil {
			t.Fatalf("Failed to create MixedBidiStreamWsGrpc: %v", err)
		}
		err = stream3.Send(&testhttpgrpc.MixedBidiStreamWsGrpcStreamingPayload{Message: "test"})
		if err != nil {
			t.Errorf("Send failed: %v", err)
		}
		_, err = stream3.Recv()
		if err != nil && err != io.EOF {
			t.Errorf("Recv failed: %v", err)
		}
		err = stream3.Close()
		if err != nil {
			t.Errorf("Close failed: %v", err)
		}
	})
}

func TestStreamMethods(t *testing.T) {
	svc := NewTestHTTPGrpc()
	h := testHTTPGrpctest.NewHarness(t, svc)
	defer h.Close()
	ctx := context.Background()

	// Server stream (HTTP SSE)
	s1, err := h.Client.HTTPServerStreamSse(ctx)
	if err == nil && s1 != nil {
		_, _ = s1.Recv()
	}

	// Server stream (gRPC)
	s2, err := h.Client.GrpcServerStream(ctx, &testhttpgrpc.GrpcServerStreamPayload{From: 0})
	if err == nil && s2 != nil {
		_, _ = s2.Recv()
	}

	// Bidi stream (HTTP WS)
	s3, err := h.Client.HTTPBidiStreamWs(ctx)
	if err == nil && s3 != nil {
		_ = s3.Send(&testhttpgrpc.HTTPBidiStreamWsStreamingPayload{Message: "x"})
		_, _ = s3.Recv()
	}

	// Client stream (gRPC)
	s4, err := h.Client.GrpcClientStream(ctx)
	if err == nil && s4 != nil {
		_ = s4.Send(&testhttpgrpc.GrpcClientStreamStreamingPayload{Value: 1})
		_, _ = s4.CloseAndRecv()
	}

	// Mixed streams
	ms1, err := h.Client.MixedServerStream(ctx)
	if err == nil && ms1 != nil {
		_, _ = ms1.Recv()
	}
	mws, err := h.Client.MixedClientStreamWsGrpc(ctx)
	if err == nil && mws != nil {
		_ = mws.Send(&testhttpgrpc.MixedClientStreamWsGrpcStreamingPayload{Message: "m"})
		_, _ = mws.CloseAndRecv()
	}
	mb, err := h.Client.MixedBidiStreamWsGrpc(ctx)
	if err == nil && mb != nil {
		_ = mb.Send(&testhttpgrpc.MixedBidiStreamWsGrpcStreamingPayload{Message: "m"})
		_, _ = mb.Recv()
		_ = mb.Close()
	}
}
