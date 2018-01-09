package design

import (
	. "goa.design/goa/http/design"
	_ "goa.design/plugins/cors"
	. "goa.design/plugins/cors/dsl"
)

var _ = API("calc", func() {
	Title("CORS Example Calc API")
	Description("This API demonstrates the use of the goa CORS plugin")
	Origin("http://127.0.0.1", func() {
		Headers("X-Shared-Secret")
		Methods("GET", "POST")
		Expose("X-Time")
		MaxAge(600)
		Credentials()
	})
})

var _ = Service("calc", func() {
	Description("The calc service exposes public endpoints that defines CORS policy.")
	Origin("/.*localhost.*/", func() {
		Methods("GET", "POST")
		Expose("X-Time", "X-Api-Version")
		MaxAge(100)
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
})
