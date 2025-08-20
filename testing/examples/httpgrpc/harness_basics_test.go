package testhttpgrpcapi

import (
	"testing"

	testhttpgrpc "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc"
	testHTTPGrpctest "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc/test_http_grpctest"
)

// basic harness lifecycle test
func TestTestHTTPGrpcHarness(t *testing.T) {
	svc := NewTestHTTPGrpc()
	h := testHTTPGrpctest.NewHarness(t, svc)
	h.Close()
}

// ensure client methods are available across transports
func TestClientMethodsAvailability(t *testing.T) {
	svc := NewTestHTTPGrpc()
	h := testHTTPGrpctest.NewHarness(t, svc)
	defer h.Close()

	if h.Client == nil {
		t.Fatal("client is nil")
	}
	_ = testhttpgrpc.Service(nil)
}
