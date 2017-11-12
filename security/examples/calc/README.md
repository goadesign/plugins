# goa v2 Security Plugin Example

This example illustrates how to secure microservice endpoints using the goa v2
security plugin. It consists of two services `calc` and `adder`.

The `calc` service exposes an `add` endpoint to public clients. It uses JWT
tokens to perform auth. The tokens are created using the `login` endpoint which
accepts a username and password using basic authentication.

The `adder` service exposes an `add` endpoint to private clients. It requires an
API key to validate that that requests come from authorized clients. the `calc`
service uses the `adder` service for the implementation of its `add` endpoint.

The flow of requests starting with login is thus:

```
client --- [ username/password ] --> calc.login
       <--     [ JWT token ]     --- 

client --- [ JWT token ] --> calc.add --- [ API key ] --> adder
```

## Design

The key design sections for the `calc` service define the `login` and `add`
security requirements:

```go
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
...

	Method("add", func() {
        Description("Add adds up the two integer parameters and returns the results. This action is secured with the jwt scheme")

        // Use JWT to auth requests to this endpoint.
		Security(JWTAuth, func() {
			Scope("calc:add") // Enforce presence of "calc:read" in the JWT "scope" claim.
		})
...
```
