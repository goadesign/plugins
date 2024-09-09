package design

import (
	. "goa.design/goa/v3/dsl"
	arnz "goa.design/plugins/v3/arnz/dsl"
)

var Admin = []string{"^arn:aws:iam::123456789012:user/administrator$"}
var Dev = append([]string{"^arn:aws:iam::123456789012:user/developer$"}, Admin...)
var ReadOnly = append([]string{"arn:aws:iam::123456789012:user/read-only"}, Dev...)

var CrudResponse = Type("ResponseBody", func() {
	Attribute("action", String)
	Required("action")
})

var _ = API("Arnz", func() {})

var _ = Service("Arnz", func() {
	Method("create", func() {
		arnz.AllowArnsMatching(Admin...)
		Result(CrudResponse)
		HTTP(func() {
			POST("/")
			Response(StatusOK)
		})
	})

	Method("read", func() {
		arnz.AllowArnsMatching(ReadOnly...)
		Result(CrudResponse)
		HTTP(func() {
			GET("/")
			Response(StatusOK)
		})
	})

	Method("update", func() {
		arnz.AllowArnsMatching(Dev...)
		Result(CrudResponse)
		HTTP(func() {
			PUT("/")
			Response(StatusOK)
		})
	})

	Method("delete", func() {
		arnz.AllowArnsMatching(Admin...)
		Result(CrudResponse)
		HTTP(func() {
			DELETE("/")
			Response(StatusOK)
		})
	})

	Method("health", func() {
		arnz.AllowUnsignedCallers()
		Result(CrudResponse)
		HTTP(func() {
			GET("/health")
			Response(StatusOK)
		})
	})
})
