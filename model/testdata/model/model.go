package testdata

import (
	"fmt"

	. "goa.design/model/dsl"
	"goa.design/plugins/v3/model/expr"
)

var _ = Design("Test", func() {
	defineSystem("OneContainer", "A")
	defineSystem("TwoContainers", "A", "B")
	defineSystem("ThreeContainers", "A", "B", "C")
	defineSystem("OtherContainer", "C")
	defineSystem("OtherContainers", "C", "D")

	SoftwareSystem("FormattedContainer", func() {
		Container("C A", func() {
			Tag("A")
		})
	})
})

func defineSystem(name string, services ...string) {
	SoftwareSystem(name, func() {
		for _, svc := range services {
			Container(fmt.Sprintf(expr.Root.ContainerNameFormat, svc), func() {
				Tag(svc)
			})
		}
	})
}
