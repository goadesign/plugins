package testhttpgrpcapi

import (
	"context"
	"io"
	"testing"
	"time"

	testhttpgrpc "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc"
	testHTTPGrpctest "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc/test_http_grpctest"
)

// Runtest-http-grpcHarness exercises the generated harness against your
// service implementation.// Call this helper from your test, passing your service implementation.
func RunTestHTTPGrpcHarness(t *testing.T, svc testhttpgrpc.Service) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	h := testHTTPGrpctest.NewHarness(t, svc)
	defer h.Close()

	td := testHTTPGrpctest.NewTestData()
	t.Run("http_no_stream", func(t *testing.T) {
		result, err := h.Client.HTTPNoStream(ctx, td.ValidHTTPNoStreamPayload())
		if err != nil {
			t.Errorf("http_no_stream failed: %v", err)
		}
		if result == nil {
			t.Error("http_no_stream returned nil result")
		}
	})
	t.Run("grpc_no_stream", func(t *testing.T) {
		result, err := h.Client.GrpcNoStream(ctx, td.ValidGrpcNoStreamPayload())
		if err != nil {
			t.Errorf("grpc_no_stream failed: %v", err)
		}
		if result == nil {
			t.Error("grpc_no_stream returned nil result")
		}
	})
	t.Run("http_no_stream_error", func(t *testing.T) {
		result, err := h.Client.HTTPNoStreamError(ctx, td.ValidHTTPNoStreamErrorPayload())
		if err != nil {
			t.Errorf("http_no_stream_error failed: %v", err)
		}
		if result == nil {
			t.Error("http_no_stream_error returned nil result")
		}
	})
	t.Run("grpc_no_stream_error_div_by_zero", func(t *testing.T) {
		result, err := h.Client.GrpcNoStreamErrorDivByZero(ctx, td.ValidGrpcNoStreamErrorDivByZeroPayload())
		if err != nil {
			t.Errorf("grpc_no_stream_error_div_by_zero failed: %v", err)
		}
		if result == nil {
			t.Error("grpc_no_stream_error_div_by_zero returned nil result")
		}
	})
	t.Run("mixed_no_stream", func(t *testing.T) {
		result, err := h.Client.MixedNoStream(ctx, td.ValidMixedNoStreamPayload())
		if err != nil {
			t.Errorf("mixed_no_stream failed: %v", err)
		}
		if result == nil {
			t.Error("mixed_no_stream returned nil result")
		}
	})
	t.Run("http_server_stream_sse_Stream", func(t *testing.T) {
		stream, err := h.Client.HTTPServerStreamSse(ctx)
		if err != nil {
			t.Errorf("Failed to create http_server_stream_sse stream: %v", err)
		}
		if stream == nil {
			t.Fatal("http_server_stream_sse returned nil stream")
		}
		// Server stream - receive at least one message
		_, err = stream.Recv()
		if err != nil && err != io.EOF {
			t.Errorf("http_server_stream_sse recv failed: %v", err)
		}
	})
	t.Run("http_server_stream_ws_Stream", func(t *testing.T) {
		stream, err := h.Client.HTTPServerStreamWs(ctx)
		if err != nil {
			t.Errorf("Failed to create http_server_stream_ws stream: %v", err)
		}
		if stream == nil {
			t.Fatal("http_server_stream_ws returned nil stream")
		}
		// Server stream - receive at least one message
		_, err = stream.Recv()
		if err != nil && err != io.EOF {
			t.Errorf("http_server_stream_ws recv failed: %v", err)
		}
	})
	t.Run("grpc_server_stream_Stream", func(t *testing.T) {
		stream, err := h.Client.GrpcServerStream(ctx, td.ValidGrpcServerStreamPayload())
		if err != nil {
			t.Errorf("Failed to create grpc_server_stream stream: %v", err)
		}
		if stream == nil {
			t.Fatal("grpc_server_stream returned nil stream")
		}
		// Server stream - receive at least one message
		_, err = stream.Recv()
		if err != nil && err != io.EOF {
			t.Errorf("grpc_server_stream recv failed: %v", err)
		}
	})
	t.Run("http_client_stream_ws_Stream", func(t *testing.T) {
		stream, err := h.Client.HTTPClientStreamWs(ctx)
		if err != nil {
			t.Errorf("Failed to create http_client_stream_ws stream: %v", err)
		}
		if stream == nil {
			t.Fatal("http_client_stream_ws returned nil stream")
		}
		// Client stream - send test data and close
		// Stream has typed payloads, send multiple
		for i := 0; i < 3; i++ {
			payload := td.ValidHTTPClientStreamWsPayload()
			if err := stream.Send(payload); err != nil {
				t.Errorf("http_client_stream_ws send failed: %v", err)
				break
			}
		}
		result, err := stream.CloseAndRecv()
		if err != nil {
			t.Errorf("http_client_stream_ws close and recv failed: %v", err)
		}
		if result == nil {
			t.Error("http_client_stream_ws returned nil result")
		}
	})
	t.Run("grpc_client_stream_Stream", func(t *testing.T) {
		stream, err := h.Client.GrpcClientStream(ctx)
		if err != nil {
			t.Errorf("Failed to create grpc_client_stream stream: %v", err)
		}
		if stream == nil {
			t.Fatal("grpc_client_stream returned nil stream")
		}
		// Client stream - send test data and close
		// Stream has typed payloads, send multiple
		for i := 0; i < 3; i++ {
			payload := td.ValidGrpcClientStreamPayload()
			if err := stream.Send(payload); err != nil {
				t.Errorf("grpc_client_stream send failed: %v", err)
				break
			}
		}
		result, err := stream.CloseAndRecv()
		if err != nil {
			t.Errorf("grpc_client_stream close and recv failed: %v", err)
		}
		if result == nil {
			t.Error("grpc_client_stream returned nil result")
		}
	})
	t.Run("http_bidi_stream_ws_Stream", func(t *testing.T) {
		stream, err := h.Client.HTTPBidiStreamWs(ctx)
		if err != nil {
			t.Errorf("Failed to create http_bidi_stream_ws stream: %v", err)
		}
		if stream == nil {
			t.Fatal("http_bidi_stream_ws returned nil stream")
		}
		// Bidirectional stream - send and receive
		// Send a test message
		payload := td.ValidHTTPBidiStreamWsPayload()
		if err := stream.Send(payload); err != nil {
			t.Errorf("http_bidi_stream_ws send failed: %v", err)
		}

		// Try to receive a response
		_, err = stream.Recv()
		if err != nil && err != io.EOF {
			t.Errorf("http_bidi_stream_ws recv failed: %v", err)
		}

		// Close the stream
		if err := stream.Close(); err != nil {
			t.Errorf("http_bidi_stream_ws close failed: %v", err)
		}
	})
	t.Run("grpc_bidi_stream_Stream", func(t *testing.T) {
		stream, err := h.Client.GrpcBidiStream(ctx)
		if err != nil {
			t.Errorf("Failed to create grpc_bidi_stream stream: %v", err)
		}
		if stream == nil {
			t.Fatal("grpc_bidi_stream returned nil stream")
		}
		// Bidirectional stream - send and receive
		// Send a test message
		payload := td.ValidGrpcBidiStreamPayload()
		if err := stream.Send(payload); err != nil {
			t.Errorf("grpc_bidi_stream send failed: %v", err)
		}

		// Try to receive a response
		_, err = stream.Recv()
		if err != nil && err != io.EOF {
			t.Errorf("grpc_bidi_stream recv failed: %v", err)
		}

		// Close the stream
		if err := stream.Close(); err != nil {
			t.Errorf("grpc_bidi_stream close failed: %v", err)
		}
	})
	t.Run("mixed_server_stream_Stream", func(t *testing.T) {
		stream, err := h.Client.MixedServerStream(ctx)
		if err != nil {
			t.Errorf("Failed to create mixed_server_stream stream: %v", err)
		}
		if stream == nil {
			t.Fatal("mixed_server_stream returned nil stream")
		}
		// Server stream - receive at least one message
		_, err = stream.Recv()
		if err != nil && err != io.EOF {
			t.Errorf("mixed_server_stream recv failed: %v", err)
		}
	})
	t.Run("mixed_client_stream_ws_grpc_Stream", func(t *testing.T) {
		stream, err := h.Client.MixedClientStreamWsGrpc(ctx)
		if err != nil {
			t.Errorf("Failed to create mixed_client_stream_ws_grpc stream: %v", err)
		}
		if stream == nil {
			t.Fatal("mixed_client_stream_ws_grpc returned nil stream")
		}
		// Client stream - send test data and close
		// Stream has typed payloads, send multiple
		for i := 0; i < 3; i++ {
			payload := td.ValidMixedClientStreamWsGrpcPayload()
			if err := stream.Send(payload); err != nil {
				t.Errorf("mixed_client_stream_ws_grpc send failed: %v", err)
				break
			}
		}
		result, err := stream.CloseAndRecv()
		if err != nil {
			t.Errorf("mixed_client_stream_ws_grpc close and recv failed: %v", err)
		}
		if result == nil {
			t.Error("mixed_client_stream_ws_grpc returned nil result")
		}
	})
	t.Run("mixed_bidi_stream_ws_grpc_Stream", func(t *testing.T) {
		stream, err := h.Client.MixedBidiStreamWsGrpc(ctx)
		if err != nil {
			t.Errorf("Failed to create mixed_bidi_stream_ws_grpc stream: %v", err)
		}
		if stream == nil {
			t.Fatal("mixed_bidi_stream_ws_grpc returned nil stream")
		}
		// Bidirectional stream - send and receive
		// Send a test message
		payload := td.ValidMixedBidiStreamWsGrpcPayload()
		if err := stream.Send(payload); err != nil {
			t.Errorf("mixed_bidi_stream_ws_grpc send failed: %v", err)
		}

		// Try to receive a response
		_, err = stream.Recv()
		if err != nil && err != io.EOF {
			t.Errorf("mixed_bidi_stream_ws_grpc recv failed: %v", err)
		}

		// Close the stream
		if err := stream.Close(); err != nil {
			t.Errorf("mixed_bidi_stream_ws_grpc close failed: %v", err)
		}
	})
}
