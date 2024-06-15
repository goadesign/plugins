package testdata

import (
	. "goa.design/goa/v3/dsl"
	. "goa.design/plugins/v3/model/dsl"
)

func OneService(system string, dsls ...func()) func() {
	return design(system, "A", dsls)
}
func TwoServices(system string, dsls ...func()) func() {
	return design(system, "A", "B", dsls)
}
func ThreeServices(system string, dsls ...func()) func() {
	return design(system, "A", "B", "C", dsls)
}
func OtherService(system string, dsls ...func()) func() {
	return design(system, "C", dsls)
}
func NoPackage() func() {
	return func() {
		Service("A", func() {
			Method("Method", func() {
				HTTP(func() { GET("/") })
			})
		})
	}
}

func design(system string, args ...any) func() {
	return func() {
		Model("goa.design/plugins/v3/model/testdata/model", system)
		var names []string
		var apidsl, svcdsl func()
		for _, arg := range args {
			switch v := arg.(type) {
			case string:
				names = append(names, v)
			case []func():
				if len(v) > 0 {
					apidsl = v[0]
				}
				if len(v) > 1 {
					svcdsl = v[1]
				}
			}
		}
		if apidsl != nil {
			apidsl()
		}
		for _, name := range names {
			Service(name, func() {
				if svcdsl != nil {
					svcdsl()
				}
				Method("Method", func() {
					HTTP(func() {
						GET("/")
					})
				})
			})
		}
	}
}
