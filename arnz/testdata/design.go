package design

import (
	. "goa.design/goa/v3/dsl"
	. "goa.design/plugins/v3/arnz/dsl"
)

var conflicting = Service("ConflictingDSLs", func() {
	Method("conflicting DSLs", func() {
		AllowArnsLike("admin")
		AllowArnsMatching("dev")
		Result(Empty)
		HTTP(func() {
			POST("/")
			Response(StatusOK)
		})
	})
})

var wronglocation = Service("WrongLocation", func() {
	AllowArnsLike("admin")

	Method("wrong location", func() {
		Result(Empty)
		HTTP(func() {
			POST("/")
			Response(StatusOK)
		})
	})
})
