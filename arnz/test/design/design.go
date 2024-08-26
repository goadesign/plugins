package design

import (
	. "goa.design/goa/v3/dsl"
	. "goa.design/plugins/v3/arnz/dsl"
)

var _ = API("Arnz", func() {
	Server("test", func() {
		Host("localhost", func() {
			URI("http://localhost:8085")
		})
	})
})

var _ = Service("Arnz", func() {
	Method("Authenticated", func() {
		Result(Empty)
		HTTP(func() {
			GET("/autho")
			Response(StatusOK)
		})
	})

	Method("Authorized", func() {
		AllowArnsLike("Admin")
		Result(Empty)
		HTTP(func() {
			GET("/authz")
			Response(StatusOK)
		})
	})
})
