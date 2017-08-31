package design

import (
	. "goa.design/goa/http/design"
	. "goa.design/goa/http/dsl"
	. "goa.design/plugins/security/dsl"
)

// JWT defines a security scheme that uses JWT tokens.
var JWT = JWTSecurity("jwt", func() {
	Description(`Secures endpoint by requiring a valid JWT token retrieved via the signin endpoint. Supports scopes "api:read" and "api:write".`)
	Scope("api:read", "Read-only access")
	Scope("api:write", "Read and write access")
})

// APIKey defines a security scheme that uses API keys.
var APIKey = APIKeySecurity("api_key", func() {
	Description("Secures endpoint by requiring an API key.")
})

// Basic defines a security scheme using basic authentication. The scheme
// protects the "signin" action used to create JWTs.
var Basic = BasicAuthSecurity("basic", func() {
	Description("Basic authentication used to authenticate security principal during signin")
})

var _ = Service("secured_service", func() {

	Method("signin", func() {
		Description("Creates a valid JWT")

		// The signin endpoint is secured via basic auth
		Security(Basic)

		Payload(String, func() {
			Description("Credentials used to authenticate to retrieve JWT token")
			Example("user:password")
		})
		Result(String, func() {
			Description("New JWT")
		})
		Error("unauthorized", String, "Credentials are invalid")

		HTTP(func() {
			POST("/signin")
			// Use Authorization header to provide basic auth value.
			SecurityHeader("basic", "Authorization")
			Response(StatusNoContent, func() {
				Headers(func() {
					Header("Authorization", String, "Generated JWT")
				})
			})
			Response("unauthorized", StatusUnauthorized)
		})
	})

	Method("secure", func() {
		Description("This action is secured with the jwt scheme")
		Security(JWT, func() { // Use JWT to auth requests to this endpoint.
			Scope("api:read") // Enforce presence of "api:read" scope in JWT claims.
		})
		Payload(func() {
			Attribute("fail", Bool, func() {
				Description("Whether to force auth failure even with a valid JWT")
			})
			Attribute("token", String, func() {
				Description("JWT used for authentication")
			})
		})
		Result(String, func() {
			Example("JWT secured data")
		})
		HTTP(func() {
			GET("/secure")

			Param("fail")
			SecurityHeader("jwt", "token:Authorization")

			Response(StatusOK)
			Response(StatusUnauthorized)
		})
	})

	Method("doubly_secure", func() {
		Description("This action is secured with the jwt scheme and also requires an API key query string.")
		Security(JWT, APIKey, func() { // Use JWT and an API key to secure this endpoint.
			Scope("api:read")  // Enforce presence of both "api:read"
			Scope("api:write") // and "api:write" scopes in JWT claims.
		})
		Payload(func() {
			Attribute("key", String, func() {
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

			SecurityParam("api_key", "key")
			SecurityHeader("jwt", "token:Authorization")

			Response(StatusOK)
			Response(StatusUnauthorized)
		})
	})
})
