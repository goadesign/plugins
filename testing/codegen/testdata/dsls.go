package testdata

import (
	. "goa.design/goa/v3/dsl"
)

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

var WithoutResultDSL = func() {
	Service("WithoutResultService", func() {
		Method("WithoutResultMethod", func() {
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
