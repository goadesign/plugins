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

var _ = API("Arnz", func() {})

var _ = Service("Like", func() {
	HTTP(func() {
		Path("/like")
	})

	Method("create", func() {
		AllowArnsLike("admin")
		Result(Empty)
		HTTP(func() {
			POST("/")
			Response(StatusOK)
		})
	})

	Method("read", func() {
		Result(Empty)
		HTTP(func() {
			GET("/")
			Response(StatusOK)
		})
	})

	Method("update", func() {
		AllowUnsigned()
		AllowArnsLike("admin", "developer")
		Result(Empty)
		HTTP(func() {
			PUT("/")
			Response(StatusOK)
		})
	})

	Method("delete", func() {
		AllowArnsLike("admin")
		Result(Empty)
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
		Result(Empty)
		HTTP(func() {
			POST("/")
			Response(StatusOK)
		})
	})

	Method("read", func() {
		Result(Empty)
		HTTP(func() {
			GET("/")
			Response(StatusOK)
		})
	})

	Method("update", func() {
		AllowUnsigned()
		AllowArnsLike(AdminArn, DevArn)
		Result(Empty)
		HTTP(func() {
			PUT("/")
			Response(StatusOK)
		})
	})

	Method("delete", func() {
		AllowArnsLike(AdminArn)
		Result(Empty)
		HTTP(func() {
			DELETE("/")
			Response(StatusOK)
		})
	})
})
