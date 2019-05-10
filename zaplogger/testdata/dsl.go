package testdata

import (
	. "goa.design/goa/v3/dsl"
)

var SimpleServiceDSL = func() {
	Service("SimpleService", func() {
		Method("SimpleMethod", func() {
			HTTP(func() {
				GET("/")
			})
		})
	})
}
