package design

import (
	. "goa.design/goa/v3/dsl"
	i18n "goa.design/plugins/v3/i18n/dsl"
)

var _ = API("calc", func() {
	i18n.Title(M("ApiCalcTitle"))
	i18n.Description(M("ApiCalcDescription"))
})

var _ = Service("calc", func() {
	i18n.Description(M("ApiCalcServiceCalcDescription"))

	Method("add", func() {
		i18n.Description(M("ApiCalcServiceCalcMethodAddDescription"))
		Payload(func() {
			Attribute("a", Int, func() {
				// Most basic usage
				i18n.Description(M("ApiCalcServiceCalcMethodAddPayloadAttributeADescription"))
				Example(1)
			})
			Attribute("b", Int, func() {
				i18n.Description(M("ApiCalcServiceCalcMethodAddPayloadAttributeBDescription"))
				Example(2)
			})
			Required("a", "b")
		})
		Result(Int, func() {
			i18n.Description(M("ApiCalcServiceCalcMethodAddResultDescription"))
			Example(3)
		})
		HTTP(func() {
			GET("/add/{a}/{b}")

			Response(StatusOK)
		})
	})
})
