package design

import (
	. "goa.design/goa/v3/dsl"
	. "goa.design/plugins/v3/arnz/dsl"
)

const (
	AdminArn = "arn:aws:iam::123456789012:user/administrator"
	DevArn   = "arn:aws:iam::123456789012:user/developer"
	ReadArn  = "arn:aws:iam::123456789012:user/reader"
)

var CrudResponse = Type("ResponseBody", func() {
	Attribute("action", String)
	Required("action")
})

var _ = API("Arnz", func() {})

var _ = Service("Like", func() {
	HTTP(func() {
		Path("/like")
	})

	Method("create", func() {
		AllowArnsLike("admin")
		Result(CrudResponse)
		HTTP(func() {
			POST("/")
			Response(StatusOK)
		})
	})

	Method("read", func() {
		Result(CrudResponse)
		HTTP(func() {
			GET("/")
			Response(StatusOK)
		})
	})

	Method("update", func() {
		AllowUnsigned()
		AllowArnsLike("admin", "developer")
		Result(CrudResponse)
		HTTP(func() {
			PUT("/")
			Response(StatusOK)
		})
	})

	Method("delete", func() {
		AllowArnsLike("admin")
		Result(CrudResponse)
		HTTP(func() {
			DELETE("/")
			Response(StatusOK)
		})
	})
})

var _ = Service("Match", func() {
	HTTP(func() {
		Path("/match")
	})

	Method("create", func() {
		AllowArnsMatching(AdminArn)
		Result(CrudResponse)
		HTTP(func() {
			POST("/")
			Response(StatusOK)
		})
	})

	Method("read", func() {
		Result(CrudResponse)
		HTTP(func() {
			GET("/")
			Response(StatusOK)
		})
	})

	Method("update", func() {
		AllowUnsigned()
		AllowArnsLike(AdminArn, DevArn)
		Result(CrudResponse)
		HTTP(func() {
			PUT("/")
			Response(StatusOK)
		})
	})

	Method("delete", func() {
		AllowArnsLike(AdminArn)
		Result(CrudResponse)
		HTTP(func() {
			DELETE("/")
			Response(StatusOK)
		})
	})
})
