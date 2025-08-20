package design

import (
	. "goa.design/goa/v3/dsl"
	_ "goa.design/plugins/v3/testing"
)

// Service exercises the testing plugin across HTTP and gRPC transports.
var _ = Service("test-http-grpc", func() {
	Description("Testing plugin matrix across HTTP and gRPC transports")

	// HTTP non-stream
	Method("http_no_stream", func() {
		Payload(func() { Field(1, "msg", String); Required("msg") })
		Result(func() { Field(1, "out", String); Required("out") })
		HTTP(func() { POST("/http/no-stream") })
	})

	// gRPC non-stream
	Method("grpc_no_stream", func() {
		Payload(func() { Field(1, "msg", String); Required("msg") })
		Result(func() { Field(1, "out", String); Required("out") })
		GRPC(func() {})
	})

	// HTTP non-stream with error
	Method("http_no_stream_error", func() {
		Payload(func() { Field(1, "msg", String); Required("msg") })
		Result(func() { Field(1, "out", String); Required("out") })
		Error("division_by_zero", DivisionByZeroError, "Division by zero error")
		HTTP(func() { 
			POST("/http/no-stream-error")
			Response("division_by_zero", StatusBadRequest, func() {
				Description("Division by zero error")
			})
		})
	})

	// gRPC non-stream with error - includes validation for testing edge cases
	Method("grpc_no_stream_error_div_by_zero", func() {
		Payload(func() { 
			Field(1, "dividend", Float64, func() {
				Minimum(-1000)
				Maximum(1000)
			})
			Field(2, "divisor", Float64, func() {
				Minimum(-100)
				Maximum(100)
			})
			Required("dividend", "divisor")
		})
		Result(func() { Field(1, "quotient", Float64); Required("quotient") })
		Error("division_by_zero", DivisionByZeroError, "Division by zero error")
		GRPC(func() {
			Response("division_by_zero", CodeInvalidArgument, func() {
				Description("Division by zero error")
			})
		})
	})

	// HTTP server stream via SSE
	Method("http_server_stream_sse", func() {
		StreamingResult(func() { Field(1, "event", String); Required("event") })
		HTTP(func() {
			POST("/http/sse")
			ServerSentEvents(func() {})
			Response(StatusOK, func() { ContentType("text/event-stream") })
		})
	})

	// HTTP server stream via WebSocket (default)
	Method("http_server_stream_ws", func() {
		StreamingResult(func() { Field(1, "message", String); Required("message") })
		HTTP(func() { GET("/http/ws/server") })
	})

	// gRPC server stream
	Method("grpc_server_stream", func() {
		Payload(func() { Field(1, "from", Int); Required("from") })
		StreamingResult(func() { Field(1, "value", Int); Required("value") })
		GRPC(func() {})
	})

	// HTTP client stream via WebSocket
	Method("http_client_stream_ws", func() {
		StreamingPayload(func() { Field(1, "message", String); Required("message") })
		Result(func() { Field(1, "out", String); Required("out") })
		HTTP(func() { POST("/http/ws/client") })
	})

	// gRPC client stream
	Method("grpc_client_stream", func() {
		StreamingPayload(func() { Field(1, "value", Int); Required("value") })
		Result(func() { Field(1, "sum", Int); Required("sum") })
		GRPC(func() {})
	})

	// HTTP bidirectional stream via WebSocket
	Method("http_bidi_stream_ws", func() {
		StreamingPayload(func() { Field(1, "message", String); Required("message") })
		StreamingResult(func() { Field(1, "message", String); Required("message") })
		HTTP(func() { POST("/http/ws/bidi") })
	})

	// gRPC bidirectional stream
	Method("grpc_bidi_stream", func() {
		StreamingPayload(func() { Field(1, "in", String); Required("in") })
		StreamingResult(func() { Field(1, "out", String); Required("out") })
		GRPC(func() {})
	})

	// Mixed transports: non-stream over both HTTP and gRPC
	Method("mixed_no_stream", func() {
		Payload(func() { Field(1, "msg", String); Required("msg") })
		Result(func() { Field(1, "out", String); Required("out") })
		HTTP(func() { POST("/mixed/no-stream") })
		GRPC(func() {})
	})

	// Mixed transports: server stream over HTTP (SSE) and gRPC
	Method("mixed_server_stream", func() {
		StreamingResult(func() { Field(1, "event", String); Required("event") })
		HTTP(func() {
			POST("/mixed/server-stream")
			ServerSentEvents(func() {})
			Response(StatusOK, func() { ContentType("text/event-stream") })
		})
		GRPC(func() {})
	})

	// Mixed transports: client stream over HTTP (WebSocket) and gRPC
	Method("mixed_client_stream_ws_grpc", func() {
		StreamingPayload(func() { Field(1, "message", String); Required("message") })
		Result(func() { Field(1, "out", String); Required("out") })
		HTTP(func() { POST("/mixed/client-stream/ws") })
		GRPC(func() {})
	})

	// Mixed transports: bidirectional stream over HTTP (WebSocket) and gRPC
	Method("mixed_bidi_stream_ws_grpc", func() {
		StreamingPayload(func() { Field(1, "message", String); Required("message") })
		StreamingResult(func() { Field(1, "message", String); Required("message") })
		HTTP(func() { POST("/mixed/bidi-stream/ws") })
		GRPC(func() {})
	})
})

// DivisionByZeroError is the custom error type for division by zero errors
var DivisionByZeroError = Type("DivisionByZeroError", func() {
	Field(1, "message", String, "division by zero error message")
	Required("message")
})
