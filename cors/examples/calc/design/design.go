package design

import (
	. "goa.design/goa/v3/dsl"
	cors "goa.design/plugins/v3/cors/dsl"
)

var _ = API("calc", func() {
	Title("CORS Example Calc API")
	Description("This API demonstrates the use of the goa CORS plugin")
	cors.Origin("http://127.0.0.1", func() {
		cors.Headers("X-Shared-Secret")
		cors.Methods("GET", "POST")
		cors.Expose("X-Time")
		cors.MaxAge(600)
		cors.Credentials()
	})
})

var _ = Service("calc", func() {
	Description("The calc service exposes public endpoints that defines CORS policy.")
	cors.Origin("/.*localhost.*/", func() {
		cors.Methods("GET", "POST")
		cors.Expose("X-Time", "X-Api-Version")
		cors.MaxAge(100)
	})

	Method("add", func() {
		Description("Add adds up the two integer parameters and returns the results.")
		Payload(func() {
			Attribute("a", Int, func() {
				Description("Left operand")
				Example(1)
			})
			Attribute("b", Int, func() {
				Description("Right operand")
				Example(2)
			})
			Required("a", "b")
		})
		Result(Int, func() {
			Description("Result of addition")
			Example(3)
		})
		HTTP(func() {
			GET("/add/{a}/{b}")

			Response(StatusOK)
		})
	})

	Files("/", "/index.html")
})
