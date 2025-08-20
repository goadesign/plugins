package design

import (
	. "goa.design/goa/v3/dsl"
	_ "goa.design/plugins/v3/testing"
)

// Dedicated JSON-RPC service (non-streaming only) to keep generation clean.
var _ = Service("test-jsonrpc", func() {
	Description("Testing plugin matrix across JSON-RPC transports (non-streaming)")

	// JSON-RPC non-stream
	Method("jsonrpc_no_stream", func() {
		Payload(func() { Field(1, "msg", String); Required("msg") })
		Result(func() { Field(1, "out", String); Required("out") })
		JSONRPC(func() {})
	})

	// JSON-RPC non-stream with error
	Method("jsonrpc_no_stream_error", func() {
		Payload(func() { Field(1, "msg", String); Required("msg") })
		Result(func() { Field(1, "out", String); Required("out") })
		Error("division_by_zero", ErrorResult, "Division by zero error")
		JSONRPC(func() {})
	})
})
