package design

import (
	. "goa.design/goa/dsl"
	_ "goa.design/plugins/goakit"
)

var _ = API("calc", func() {
	Title("Goakit example Calc API")
	Description("This API demonstrates the use of the goa goakit plugin")
})

var _ = Service("calc", func() {
	Description("The calc service exposes public endpoints that uses go-kit.")
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
