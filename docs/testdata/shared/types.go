package shared

import . "goa.design/goa/v3/dsl"

var CrossPackageType = Type("CrossPackageType", func() {
	Field(1, "C", String)
	Required("C")
})
