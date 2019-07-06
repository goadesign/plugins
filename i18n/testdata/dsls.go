package testdata

import (
	. "goa.design/goa/v3/dsl"
	i18n "goa.design/plugins/v3/i18n/dsl"
)

var SimpleI18nDSL = func() {
	API("calc", func() {
		i18n.Title(M("title"))
	})
	Service("SimpleOrigin", func() {
		i18n.Description(M("title"))
		i18n.Example(M("title"))

		Method("SimpleOriginMethod", func() {
			HTTP(func() {
				GET("/")
			})
		})
	})

}
