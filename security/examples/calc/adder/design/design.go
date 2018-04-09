package design

import (
	. "goa.design/goa/http/design"
	_ "goa.design/plugins/security"
	. "goa.design/plugins/security/dsl"
)

var _ = API("adder", func() {
	Title("Security Example Calc Adder service API")
	Description("This API demonstrates the use of the goa security plugin")
	Docs(func() { // Documentation links
		Description("Security plugin README")
		URL("https://github.com/goadesign/plugins/tree/master/security/README.md")
	})
})

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
			Attribute("a", Int, func() {
				Description("Left operand")
				Example(2)
			})
			Attribute("b", Int, func() {
				Description("Right operand")
				Example(3)
			})
			Required("key", "a", "b")
		})
		Result(Int, func() {
			Example(5)
		})
		Error("unauthorized", String)
		HTTP(func() {
			GET("/add/{a}/{b}")

			Param("key") // Use query string parameter to auth request

			Response(StatusOK)
			Response("unauthorized", StatusUnauthorized)
		})
	})
})
