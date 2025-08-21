package testhttpgrpcapi

import (
	"context"
	"io"

	"goa.design/clue/log"
	testhttpgrpc "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc"
)

// test-http-grpc service example implementation.
// The example methods log the requests and return zero values.
type testHTTPGrpcsrvc struct{}

// NewTestHTTPGrpc returns the test-http-grpc service implementation.
func NewTestHTTPGrpc() testhttpgrpc.Service {
	return &testHTTPGrpcsrvc{}
}

// HTTPNoStream implements http_no_stream.
func (s *testHTTPGrpcsrvc) HTTPNoStream(ctx context.Context, p *testhttpgrpc.HTTPNoStreamPayload) (res *testhttpgrpc.HTTPNoStreamResult, err error) {
	log.Printf(ctx, "testHTTPGrpc.http_no_stream")
	return &testhttpgrpc.HTTPNoStreamResult{Out: p.Msg}, nil
}

// GrpcNoStream implements grpc_no_stream.
func (s *testHTTPGrpcsrvc) GrpcNoStream(ctx context.Context, p *testhttpgrpc.GrpcNoStreamPayload) (res *testhttpgrpc.GrpcNoStreamResult, err error) {
	res = &testhttpgrpc.GrpcNoStreamResult{}
	log.Printf(ctx, "testHTTPGrpc.grpc_no_stream")
	return
}

// HTTPNoStreamError implements http_no_stream_error.
func (s *testHTTPGrpcsrvc) HTTPNoStreamError(ctx context.Context, p *testhttpgrpc.HTTPNoStreamErrorPayload) (res *testhttpgrpc.HTTPNoStreamErrorResult, err error) {
	log.Printf(ctx, "testHTTPGrpc.http_no_stream_error")
	// Return division by zero error for testing
	if p.Msg == "error" {
		return nil, &testhttpgrpc.DivisionByZeroError{
			Message: "cannot divide by zero",
		}
	}
	res = &testhttpgrpc.HTTPNoStreamErrorResult{Out: "processed: " + p.Msg}
	return
}

// GrpcNoStreamErrorDivByZero implements grpc_no_stream_error_div_by_zero.
func (s *testHTTPGrpcsrvc) GrpcNoStreamErrorDivByZero(ctx context.Context, p *testhttpgrpc.GrpcNoStreamErrorDivByZeroPayload) (res *testhttpgrpc.GrpcNoStreamErrorDivByZeroResult, err error) {
	log.Printf(ctx, "testHTTPGrpc.grpc_no_stream_error_div_by_zero")
	// Return division by zero error when divisor is 0
	if p.Divisor == 0 {
		return nil, &testhttpgrpc.DivisionByZeroError{
			Message: "cannot divide by zero",
		}
	}
	res = &testhttpgrpc.GrpcNoStreamErrorDivByZeroResult{Quotient: p.Dividend / p.Divisor}
	return
}

// HTTPServerStreamSse implements http_server_stream_sse.
func (s *testHTTPGrpcsrvc) HTTPServerStreamSse(ctx context.Context, stream testhttpgrpc.HTTPServerStreamSseServerStream) (err error) {
	log.Printf(ctx, "testHTTPGrpc.http_server_stream_sse")
	return
}

// HTTPServerStreamWs implements http_server_stream_ws.
func (s *testHTTPGrpcsrvc) HTTPServerStreamWs(ctx context.Context, stream testhttpgrpc.HTTPServerStreamWsServerStream) (err error) {
	log.Printf(ctx, "testHTTPGrpc.http_server_stream_ws")
	// Send at least one message for testing
	return stream.Send(&testhttpgrpc.HTTPServerStreamWsResult{
		Message: "test message",
	})
}

// GrpcServerStream implements grpc_server_stream.
func (s *testHTTPGrpcsrvc) GrpcServerStream(ctx context.Context, p *testhttpgrpc.GrpcServerStreamPayload, stream testhttpgrpc.GrpcServerStreamServerStream) (err error) {
	log.Printf(ctx, "testHTTPGrpc.grpc_server_stream")
	return
}

// HTTPClientStreamWs implements http_client_stream_ws.
func (s *testHTTPGrpcsrvc) HTTPClientStreamWs(ctx context.Context, stream testhttpgrpc.HTTPClientStreamWsServerStream) (err error) {
	log.Printf(ctx, "testHTTPGrpc.http_client_stream_ws")
	// Read from the stream until EOF
	for {
		_, err := stream.Recv()
		if err != nil {
			break
		}
	}
	// Send result back
	return stream.SendAndClose(&testhttpgrpc.HTTPClientStreamWsResult{Out: "received"})
}

// GrpcClientStream implements grpc_client_stream.
func (s *testHTTPGrpcsrvc) GrpcClientStream(ctx context.Context, stream testhttpgrpc.GrpcClientStreamServerStream) (err error) {
	log.Printf(ctx, "testHTTPGrpc.grpc_client_stream")
	// Read from the stream until EOF
	sum := 0
	for {
		msg, err := stream.Recv()
		if err != nil {
			break
		}
		if msg != nil {
			sum += msg.Value
		}
	}
	// Send result back
	return stream.SendAndClose(&testhttpgrpc.GrpcClientStreamResult{Sum: sum})
}

// HTTPBidiStreamWs implements http_bidi_stream_ws.
func (s *testHTTPGrpcsrvc) HTTPBidiStreamWs(ctx context.Context, stream testhttpgrpc.HTTPBidiStreamWsServerStream) (err error) {
	log.Printf(ctx, "testHTTPGrpc.http_bidi_stream_ws")
	// Read one message and send a response
	msg, err := stream.Recv()
	if err != nil && err != io.EOF {
		return err
	}
	if msg != nil {
		// Echo back
		return stream.Send(&testhttpgrpc.HTTPBidiStreamWsResult{
			Message: "echo: " + msg.Message,
		})
	}
	return nil
}

// GrpcBidiStream implements grpc_bidi_stream.
func (s *testHTTPGrpcsrvc) GrpcBidiStream(ctx context.Context, stream testhttpgrpc.GrpcBidiStreamServerStream) (err error) {
	log.Printf(ctx, "testHTTPGrpc.grpc_bidi_stream")
	return
}

// MixedNoStream implements mixed_no_stream.
func (s *testHTTPGrpcsrvc) MixedNoStream(ctx context.Context, p *testhttpgrpc.MixedNoStreamPayload) (res *testhttpgrpc.MixedNoStreamResult, err error) {
	log.Printf(ctx, "testHTTPGrpc.mixed_no_stream")
	return &testhttpgrpc.MixedNoStreamResult{Out: p.Msg}, nil
}

// MixedServerStream implements mixed_server_stream.
func (s *testHTTPGrpcsrvc) MixedServerStream(ctx context.Context, stream testhttpgrpc.MixedServerStreamServerStream) (err error) {
	log.Printf(ctx, "testHTTPGrpc.mixed_server_stream")
	return
}

// MixedClientStreamWsGrpc implements mixed_client_stream_ws_grpc.
func (s *testHTTPGrpcsrvc) MixedClientStreamWsGrpc(ctx context.Context, stream testhttpgrpc.MixedClientStreamWsGrpcServerStream) (err error) {
	log.Printf(ctx, "testHTTPGrpc.mixed_client_stream_ws_grpc")
	// Read from the stream until EOF
	for {
		_, err := stream.Recv()
		if err != nil {
			break
		}
	}
	// Send result back
	return stream.SendAndClose(&testhttpgrpc.MixedClientStreamWsGrpcResult{Out: "received"})
}

// MixedBidiStreamWsGrpc implements mixed_bidi_stream_ws_grpc.
func (s *testHTTPGrpcsrvc) MixedBidiStreamWsGrpc(ctx context.Context, stream testhttpgrpc.MixedBidiStreamWsGrpcServerStream) (err error) {
	log.Printf(ctx, "testHTTPGrpc.mixed_bidi_stream_ws_grpc")
	// Read one message and send a response
	msg, err := stream.Recv()
	if err != nil && err != io.EOF {
		return err
	}
	if msg != nil {
		// Echo back
		return stream.Send(&testhttpgrpc.MixedBidiStreamWsGrpcResult{
			Message: "echo: " + msg.Message,
		})
	}
	return nil
}
