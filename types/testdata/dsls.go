package testdata

import (
	. "goa.design/goa/v3/dsl"
)

var Empt = func() {
	var _ = Service("foo", func() {})
}

var NoValidation = func() {
	var _ = Type("NoVal", func() {
		Attribute("attr", String)
	})
}

var Require = func() {
	var _ = Type("Require", func() {
		Attribute("attr", String)
		Required("attr")
	})
}

var Validation = func() {
	var _ = Type("Validation", func() {
		Attribute("attr", String, func() {
			Pattern("^[a-zA-Z0-9]*$")
		})
	})
}

var Multiple = func() {
	var AType = Type("AType", func() {
		Attribute("attr", String, func() {
			Pattern("^[a-zA-Z0-9]*$")
		})
	})

	var OtherType = Type("OtherType", func() {
		Attribute("attr", String, func() {
			Pattern("^[a-zA-Z0-9]*$")
		})
		Required("attr")
	})

	var _ = Type("Composite", func() {
		Attribute("attr", AType)
		Attribute("other", OtherType)
		Required("attr", "other")
	})
}

var Alias = func() {
	var _ = Type("Alias", String, func() {
		MinLength(10)
	})
}

var Exampl = func() {
	var _ = Type("MyType", func() {
		Description("My type")
		Attribute("age", Int, "Age", func() {
			Minimum(0)
		})
		Attribute("name", String, "Name")
		Required("age", "name")
	})
}
