package testdata

import (
	. "goa.design/goa/v3/dsl"
	. "goa.design/plugins/v3/arnz/dsl"
)

var Authenticated = func() {
	Service("Arnz", func() {
		Method("authenticated", func() {
			HTTP(func() {
				GET("/authenticated")
			})
		})
	})
}

var Authorized = func() {
	Service("Arz", func() {
		Method("authorized", func() {
			AllowArnsLike("Administrator")
			HTTP(func() {
				GET("/authorized")
			})
		})
	})
}
