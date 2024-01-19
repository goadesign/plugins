package testdata

import (
	. "goa.design/goa/v3/dsl"
	_ "goa.design/plugins/v3/docs"
)

var OneRoute = func() {
	Service("Service", func() {
		Method("Method", func() {
			HTTP(func() {
				GET("/")
			})
		})
	})
}

var MultipleRoutes = func() {
	Service("Service2", func() {
		Method("Method", func() {
			HTTP(func() {
				GET("/")
				GET("/other")
			})
			GRPC(func() {})
		})
	})
}
