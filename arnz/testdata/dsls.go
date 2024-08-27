package testdata

import (
	. "goa.design/goa/v3/dsl"
	. "goa.design/plugins/v3/arnz/dsl"
)

var Conflicting = func() {
	Service("Conflicting", func() {
		Method("create", func() {
			AllowArnsLike("admin")
			AllowArnsMatching("admin")
			Result(Empty)
			HTTP(func() {
				POST("/")
				Response(StatusOK)
			})
		})
	})
}

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
