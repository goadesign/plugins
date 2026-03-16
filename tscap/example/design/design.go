package design

import (
	. "goa.design/goa/v3/dsl"
	tscap "goa.design/plugins/v3/tscap/dsl"
)

var _ = API("tscap_example", func() {
	Title("Tailscale Capabilities Example API")
	Description("An example API demonstrating Tailscale app capabilities authorization")
})

var _ = Service("tscap", func() {
	Description("A service demonstrating Tailscale app capabilities")

	Method("list", func() {
		Description("List items - requires read capability")
		tscap.Require("example.com/cap/tscap", "read", "*")
		Result(ArrayOf(String))
		HTTP(func() {
			GET("/items")
			Response(StatusOK)
		})
	})

	Method("create", func() {
		Description("Create an item - requires write capability")
		tscap.Require("example.com/cap/tscap", "write", "items/*")
		Payload(func() {
			Attribute("name", String, "Item name")
			Required("name")
		})
		Result(String)
		HTTP(func() {
			POST("/items")
			Response(StatusCreated)
		})
	})

	Method("admin", func() {
		Description("Admin action - requires admin capability")
		tscap.Require("example.com/cap/tscap", "admin", "*")
		Payload(func() {
			Attribute("id", String, "Item ID")
			Required("id")
		})
		HTTP(func() {
			DELETE("/items/{id}")
			Response(StatusNoContent)
		})
	})

	Method("health", func() {
		Description("Health check - no capability required")
		tscap.AllowAnonymous()
		Result(String)
		HTTP(func() {
			GET("/health")
			Response(StatusOK)
		})
	})
})

var _ = Service("openapi", func() {
	Files("/openapi.json", "gen/http/openapi3.json")
})
