package design

import (
	. "goa.design/goa/http/design"
	_ "goa.design/plugins/security"
	. "goa.design/plugins/security/dsl"
)

var _ = API("calc", func() {
	Title("Security Example Calc API")
	Description("This API demonstrates the use of the goa security plugin")
	Docs(func() { // Documentation links
		Description("Security example README")
		URL("https://github.com/goadesign/plugins/security/tree/master/example/README.md")
	})
})

// BasicAuth defines a security scheme that uses basic authentication.
var BasicAuth = BasicAuthSecurity("basic", func() {
	Description("Secures the login endpoint.")
})

// JWTAuth defines a security scheme that uses JWT tokens.
var JWTAuth = JWTSecurity("jwt", func() {
	Description(`Secures endpoint by requiring a valid JWT token retrieved via the login endpoint. Supports scope "calc:add".`)
	Scope("calc:add", "Allows performing additions")
})

var _ = Service("calc", func() {
	Description("The calc service exposes public endpoints that require valid authorization credentials.")

	Method("login", func() {
		Description("Creates a valid JWT")

		// The signin endpoint is secured via basic auth
		Security(BasicAuth)

		// Payload contains the username and password used to perform
		// basic authentication.
		Payload(func() {
			Description("Credentials used to authenticate to retrieve JWT token")
			Username("user", String, func() {
				Example("username")
			})
			Password("password", String, func() {
				Example("password")
			})
			Required("user", "password")
		})

		// Result is a JWT token
		Result(String, func() {
			Description("New JWT token")
		})

		// Error returned in case login failed because of invalid
		// username/password combination.
		Error("unauthorized", String, "Credentials are invalid")

		HTTP(func() {
			POST("/login")

			Response(StatusOK)                           // Body contains JWT token
			Response("unauthorized", StatusUnauthorized) // Use HTTP status 401 when credentials are invalid
		})
	})

	Method("add", func() {
		Description("Add adds up the two integer parameters and returns the results. This endpoint is secured with the JWT scheme")

		// Use JWT to auth requests to this endpoint.
		Security(JWTAuth, func() {
			Scope("calc:add") // Enforce presence of "calc:add" scope in JWT claims.
		})

		Payload(func() {
			Attribute("a", Int, func() {
				Description("Left operand")
				Example(1)
			})
			Attribute("b", Int, func() {
				Description("Right operand")
				Example(2)
			})
			Token("token", String, func() {
				Description("JWT used for authentication")
			})
			Required("a", "b", "token")
		})
		Result(Int, func() {
			Description("Result of addition")
			Example(3)
		})
		Error("forbidden")
		HTTP(func() {
			GET("/add/{a}/{b}")

			Response(StatusOK)
			Response("forbidden", StatusForbidden)
		})
	})
})
