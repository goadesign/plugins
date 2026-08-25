// This file defines small Goa services used to verify generated testing code.
package testdata

import (
	. "goa.design/goa/v3/dsl"
)

var WithPayloadDSL = func() {
	Service("WithPayloadService", func() {
		Method("WithPayloadMethod", func() {
			Payload(func() {
				Field(1, "Attribute", String)
			})
			HTTP(func() {
				GET("/")
			})
			GRPC(func() {})
			JSONRPC(func() {})
		})
	})
}

var WithResultDSL = func() {
	Service("WithResultService", func() {
		Method("WithResultMethod", func() {
			Result(func() {
				Field(1, "Attribute", String)
			})
			HTTP(func() {
				GET("/")
			})
			GRPC(func() {})
			JSONRPC(func() {})
		})
	})
}

var WithoutPayloadResultDSL = func() {
	Service("WithoutPayloadResultService", func() {
		Method("WithoutPayloadResultMethod", func() {
			HTTP(func() {
				GET("/")
			})
			GRPC(func() {})
			JSONRPC(func() {})
		})
	})
}

var WithArrayResultDSL = func() {
	var AccessControl = ResultType("AccessControl", func() {
		Attribute("id", String)
		Required("id")
	})
	Service("WithArrayResultService", func() {
		Method("ListAccessControl", func() {
			Result(ArrayOf(AccessControl))
			HTTP(func() {
				GET("/")
			})
		})
	})
}

var WithStreamDSL = func() {
	Service("WithStreamService", func() {
		Method("WithStreamMethod", func() {
			StreamingPayload(String)
			HTTP(func() {
				GET("/")
			})
		})
	})
}

var WithoutStreamDSL = func() {
	Service("WithoutStreamService", func() {
		Method("WithoutStreamMethod", func() {
			HTTP(func() {
				GET("/")
			})
		})
	})
}

var JSONRPCTransportsDSL = func() {
	API("testing-jsonrpc", func() {
		JSONRPC(func() {})
	})
	Service("JSONRPCService", func() {
		JSONRPC(func() {
			POST("/rpc")
		})
		Method("Plain", func() {
			JSONRPC(func() {})
		})
		Method("Events", func() {
			StreamingResult(String)
			JSONRPC(func() {
				ServerSentEvents()
			})
		})
	})
}
