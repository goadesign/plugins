package testdata

import (
	. "goa.design/goa/v3/dsl"
)

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
