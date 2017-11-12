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
var BasicAuth = BasicAuth("basic", func() {
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

		Payload(String, func() {
			Description("Credentials used to authenticate to retrieve JWT token")
			Example("user:password")
		})
		Result(String, func() {
			Description("New JWT")
		})
		Error("unauthorized", String, "Credentials are invalid")

		HTTP(func() {
			POST("/login")
			// Use Authorization header to provide basic auth value.
			Response(StatusNoContent, func() {
				Headers(func() {
					Header("Authorization", String, "Generated JWT")
				})
			})
			Response("unauthorized", StatusUnauthorized)
		})
	})

	Method("add", func() {
		Description("Add adds up the two integer parameters and returns the results. This action is secured with the jwt scheme")

		// Use JWT to auth requests to this endpoint.
		Security(JWTAuth, func() {
			Scope("calc:add") // Enforce presence of "api:read" scope in JWT claims.
		})

		Payload(func() {
			Attribute("left", Int, func() {
				Description("Left operand")
				Example(1)
			})
			Attribute("right", Int, func() {
				Description("Right operand")
				Example(2)
			})
			Attribute("token", String, func() {
				Description("JWT used for authentication")
			})
		})
		Result(Int, func() {
			Description("Result of addition")
			Example(3)
		})
		HTTP(func() {
			GET("/add")

			Param("left")
			Param("right")

			Response(StatusOK)
			Response(StatusUnauthorized)
		})
	})

	Method("doubly_secure", func() {
		Description("This action is secured with the jwt scheme and also requires an API key query string.")
		Security(JWTAuth, APIKeyAuth, func() { // Use JWT and an API key to secure this endpoint.
			Scope("api:read")  // Enforce presence of both "api:read"
			Scope("api:write") // and "api:write" scopes in JWT claims.
		})
		Payload(func() {
			APIKey("api_key", "key", String, func() {
				Description("API key")
				Example("abcdef12345")
			})
			Attribute("token", String, func() {
				Description("JWT used for authentication")
				Example("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
			})
		})
		Result(String, func() {
			Example("JWT secured data")
		})
		Error("unauthorized", String)
		HTTP(func() {
			GET("/secure")

			Param("key:k")

			Response(StatusOK)
			Response(StatusUnauthorized)
		})
	})

	Method("also_doubly_secure", func() {
		Description("This action is secured with the jwt scheme and also requires an API key header.")
		Security(JWTAuth, APIKeyAuth, func() { // Use JWT and an API key to secure this endpoint.
			Scope("api:read")  // Enforce presence of both "api:read"
			Scope("api:write") // and "api:write" scopes in JWT claims.
		})
		Payload(func() {
			APIKey("api_key", "key", String, func() {
				Description("API key")
				Example("abcdef12345")
			})
			Attribute("token", String, func() {
				Description("JWT used for authentication")
				Example("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
			})
		})
		Result(String, func() {
			Example("JWT secured data")
		})
		Error("unauthorized", String)
		HTTP(func() {
			POST("/secure")

			Header("key:Authorization")

			Response(StatusOK)
			Response(StatusUnauthorized)
		})
	})
})
