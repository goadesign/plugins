package testdata

import (
	. "goa.design/goa/v3/dsl"
	i18n "goa.design/plugins/v3/i18n/dsl"
)

var SimpleI18nDSL = func() {
	API("calc", func() {
		i18n.Title(T.M("title"))
	})
	Service("SimpleOrigin", func() {
		i18n.Description(T.M("title"))
		i18n.Example(T.M("title"))

		Method("SimpleOriginMethod", func() {
			HTTP(func() {
				GET("/")
			})
		})
	})

}
