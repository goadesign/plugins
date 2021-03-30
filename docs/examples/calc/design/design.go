package design

import (
	. "goa.design/goa/v3/dsl"
	_ "goa.design/plugins/v3/docs"
)

// API describes the global properties of the API server.
var _ = API("calc", func() {
	Title("Calculator Service")
})

// Service describes a service
var _ = Service("calc", func() {
	Description("The calc service performs additions on numbers")

	// Method describes a service method (endpoint)
	Method("add", func() {
		// Payload describes the method payload.
		// Here the payload is an object that consists of two fields.
		Payload(func() {
			// Attribute describes an object field
			Attribute("left", Int, "Left operand")
			Attribute("right", Int, "Right operand")
			Required("left", "right")
		})

		// StreamingPayload indicates that the attributes can be streamed to the
		// server (via websocket).
		StreamingPayload(func() {
			Attribute("a", Int, "Left operand")
			Attribute("b", Int, "Right operand")
			Required("a", "b")
		})

		// Result describes the method result.
		// Here the result is a simple integer value.
		Result(Int)

		StreamingResult(Int)

		// HTTP describes the HTTP transport mapping.
		HTTP(func() {
			// Requests to the service consist of HTTP GET requests.
			// The payload fields are encoded as path parameters.
			GET("/add/{left}/{right}")
			// Responses use a "200 OK" HTTP status.
			// The result is encoded in the response body.
			Response(StatusOK)
		})
	})

	// Serve the file with relative path ../../gen/http/openapi.json for
	// requests sent to /swagger.json.
	Files("/swagger.json", "../../gen/http/openapi.json")
})
