package design

import (
	. "goa.design/goa/v3/dsl"
	_ "goa.design/plugins/goakit"
)

var _ = API("archiver", func() {
	Title("The goakit example downstream service")
	Description("Archiver is a service that manages the content of HTTP responses")
})

var _ = Service("archiver", func() {
	HTTP(func() {
		Path("/archive")
	})

	Method("archive", func() {
		Description("Archive HTTP response")
		Payload(ArchivePayload)
		Result(ArchiveMedia)
		HTTP(func() {
			POST("/")
			Response(StatusOK)
		})
	})

	Method("read", func() {
		Description("Read HTTP response from archive")
		Payload(func() {
			Attribute("id", Int, "ID of archive", func() {
				Minimum(0)
			})
			Required("id")
		})
		Result(ArchiveMedia)
		Error("not_found")
		Error("bad_request")
		HTTP(func() {
			GET("/{id}")
			Response(StatusOK)
			Response("not_found", StatusNotFound, func() {
				Tag("code", "not_found")
			})
			Response("bad_request", StatusBadRequest, func() {
				Tag("code", "bad_request")
			})
		})
	})
})

var ArchivePayload = Type("ArchivePayload", func() {
	Attribute("status", Int, "HTTP status", func() {
		Minimum(0)
		Example(200)
	})
	Attribute("body", String, "HTTP response body content")
	Required("status", "body")
})

var ArchiveMedia = ResultType("application/vnd.goa.archive", func() {
	Description("Archive is an HTTP response archive")
	TypeName("ArchiveMedia")
	Reference(ArchivePayload)
	Attributes(func() {
		Attribute("href", String, "The archive resouce href", func() {
			Pattern("^/archive/[0-9]+$")
			Example("/archive1/")
		})
		Attribute("status")
		Attribute("body")
		Required("href", "status", "body")
	})
	View("default", func() {
		Attribute("href")
		Attribute("status")
		Attribute("body")
	})
})

var _ = Service("health", func() {
	HTTP(func() {
		Path("/health")
	})

	Method("show", func() {
		Description("Health check endpoint")
		Result(String)
		HTTP(func() {
			GET("/")
			Response(func() {
				Code(StatusOK)
				ContentType("text/plain")
			})
		})
	})
})
