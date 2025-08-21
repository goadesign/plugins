package testjsonrpcapi

import (
	"context"

	"goa.design/clue/log"
	testjsonrpc "goa.design/plugins/v3/testing/examples/jsonrpc/gen/test_jsonrpc"
)

// test-jsonrpc service example implementation.
// The example methods log the requests and return zero values.
type testJsonrpcsrvc struct{}

// NewTestJsonrpc returns the test-jsonrpc service implementation.
func NewTestJsonrpc() testjsonrpc.Service {
	return &testJsonrpcsrvc{}
}

// JsonrpcNoStream implements jsonrpc_no_stream.
func (s *testJsonrpcsrvc) JsonrpcNoStream(ctx context.Context, p *testjsonrpc.JsonrpcNoStreamPayload) (res *testjsonrpc.JsonrpcNoStreamResult, err error) {
	res = &testjsonrpc.JsonrpcNoStreamResult{}
	log.Printf(ctx, "testJsonrpc.jsonrpc_no_stream")
	return
}

// JsonrpcNoStreamError implements jsonrpc_no_stream_error.
func (s *testJsonrpcsrvc) JsonrpcNoStreamError(ctx context.Context, p *testjsonrpc.JsonrpcNoStreamErrorPayload) (res *testjsonrpc.JsonrpcNoStreamErrorResult, err error) {
	res = &testjsonrpc.JsonrpcNoStreamErrorResult{}
	log.Printf(ctx, "testJsonrpc.jsonrpc_no_stream_error")
	return
}
