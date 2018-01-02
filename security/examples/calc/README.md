# goa v2 Security Plugin Example

This example illustrates how to secure microservice endpoints using the goa v2
security plugin. It consists of two services `calc` and `adder`, the `calc`
service is a client of the `adder`. The example thus illustrates how to
implement both a secure service and a secure client.

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
The design first defines both the `BasicAuth` and `JWTAuth` security schemes
then uses these definitions in the `login` and `add` methods respectively.

The `adder` service only makes use of a API key scheme to secure requests made
to the `add` endpoint:

```go
// APIKeyAuth defines a security scheme that uses API keys.
var APIKeyAuth = APIKeySecurity("api_key", func() {
	Description("Secures endpoint by requiring an API key.")
})

var _ = Service("adder", func() {
	Description("The adder service exposes an add method secured via API keys.")
	Method("add", func() {
		Description("This action returns the sum of two integers and is secured with the API key scheme")
		Security(APIKeyAuth) // Use API keys to auth requests to this endpoint.
		Payload(func() {
			APIKey("api_key", "key", String, func() {
				Description("API key")
				Example("abcdef12345")
			})
			...
```

The payloads for all secure endpoints list the fields holding credentials using
specific DSL that identifies the corresponding security attributes e.g.
`Username`, `Password` or `APIKey`. The service method implementations may use
the fields to perform validations while clients use it to set the credentials.