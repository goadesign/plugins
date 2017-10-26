package design

import (
	. "goa.design/goa/http/design"
	. "goa.design/goa/http/dsl"
	_ "goa.design/plugins/goakit"
)

var _ = API("fetcher", func() {
	Title("The goakit example upstream service")
	Description("Fetcher is a service that makes GET requests to arbitrary URLs and stores the results in the downstream 'archiver' service.")
	Server("http://localhost:8081")
})

var _ = Service("fetcher", func() {
	Method("fetch", func() {
		Description("Fetch makes a GET request to the given URL and stores the results in the archiver service which must be running or the request fails")
		Payload(func() {
			Attribute("url", String, "URL to be fetched", func() {
				Format("uri")
			})
			Required("url")
		})
		Result(FetchMedia)
		Error("bad_request")
		Error("internal_error")
		HTTP(func() {
			GET("fetch/{*url}")
			Response(StatusOK)
			Response("bad_request", StatusBadRequest, func() {
				Tag("code", "bad_request")
			})
			Response("internal_error", StatusInternalServerError, func() {
				Tag("code", "internal_error")
			})
		})
	})
})

var FetchMedia = ResultType("application/vnd.goa.fetch", func() {
	Description("FetchResponse contains the HTTP status code returned by the fetched service and the href to the archived results in the archiver service")
	TypeName("FetchMedia")
	Attributes(func() {
		Attribute("status", Int, "HTTP status code returned by fetched service", func() {
			Minimum(0)
			Example(200)
		})
		Attribute("archive_href", String, "The href to the corresponding archive in the archiver service", func() {
			Pattern("^/archive/[0-9]+$")
			Example("/archive/1")
		})
		Required("status", "archive_href")
	})
	View("default", func() {
		Attribute("status")
		Attribute("archive_href")
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
