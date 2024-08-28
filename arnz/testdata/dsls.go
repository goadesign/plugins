package testdata

import (
	. "goa.design/goa/v3/dsl"
	. "goa.design/plugins/v3/arnz/dsl"
)

var WrongScope = func() {
	Service("WrongScope", func() {
		AllowArnsMatching("admin")
		Method("create", func() {
			Result(Empty)
			HTTP(func() {
				POST("/")
				Response(StatusOK)
			})
		})
	})
}

var BadMatcher = func() {
	Service("BadMatcher", func() {
		AllowArnsMatching(`^arn:aws:iam::123456789012:user/([a-zA-Z0-9_+=,.@-]+$`)
		Method("create", func() {
			Result(Empty)
			HTTP(func() {
				GET("/")
				Response(StatusOK)
			})
		})
	})
}
