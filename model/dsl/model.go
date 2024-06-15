package dsl

import (
	"goa.design/goa/v3/eval"
	goaexpr "goa.design/goa/v3/expr"
	"goa.design/plugins/v3/model/expr"

	// Register code generators for the Model plugin
	_ "goa.design/plugins/v3/model"
)

// Model specifies the import path of the Go package that contains the model as
// well as the name of the system to be validated.
//
// Model must appear in the API expression.
//
// Example:
//
//	import . "goa.design/plugins/v3/model"
//
//	var _ = API("calc", func() {
//	   Model("goa.design/model/examples/basic/model", "calc")
//	})
func Model(path, systemName string) {
	expr.Root.ModelPkgPath = path
	expr.Root.SystemName = systemName
}

// ModelContainerFormat specifies a fmt.Sprintf format string that is used
// to compute the name of the model container from the service name.
// The default value is "%s Service".
//
// ModelContainerFormat must appear in the API expression.
//
// Example:
//
//	import . "goa.design/plugins/v3/model"
//
//	var _ = API("calc", func() {
//	   Model("goa.design/model/examples/basic/model", "calc")
//	   ModelContainerFormat("%s Service")
//	})
func ModelContainerFormat(format string) {
	expr.Root.ContainerNameFormat = format
}

// ModelExcludedTags specifies a list of tags. Containers that have any of these
// tags will be excluded from the check.
// By default, the tags "external", "thirdparty" and "database" are excluded.
//
// ModelExcludedTags must appear in the API expression.
//
// Example:
//
//	import . "goa.design/plugins/v3/model"
//
//	var _ = API("calc", func() {
//	   Model("goa.design/model/examples/basic/model", "calc")
//	   ModelExcludedTags("external", "thirdparty", "database")
//	})
func ModelExcludedTags(tags ...string) {
	expr.Root.ExcludedTags = tags
}

// ModelComplete marks the model as complete, meaning that the model plugin
// should validate that all containers have a corresponding service in addition
// to validating that each service has a corresponding container.
//
// ModelComplete must appear in the API expression.
//
// Example:
//
//	import . "goa.design/plugins/v3/model"
//
//	var _ = API("calc", func() {
//	   Model("goa.design/model/examples/basic/model", "calc")
//	   ModelComplete()
//	})
func ModelComplete() {
	expr.Root.Complete = true
}

// ModelContainer sets the name of the container for the enclosing service.
//
// ModelContainer must appear in a Service expression.
//
// Example:
//
//	import . "goa.design/plugins/v3/model"
//
//	var _ = API("calc", func() {
//		Service("adder", func() {
//			ModelContainer("Adder")
//		})
//	})
func ModelContainer(name string) {
	svc, ok := eval.Current().(*goaexpr.ServiceExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	expr.Root.ServiceContainer[svc.Name] = name
}

// ModelNone disables the model validation for the enclosing service.
//
// ModelNone must appear in a Service expression.
//
// Example:
//
//	import . "goa.design/plugins/v3/model"
//
//	var _ = API("calc", func() {
//		Service("adder", func() {
//			ModelNone()
//		})
//	})
func ModelNone() {
	svc, ok := eval.Current().(*goaexpr.ServiceExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	expr.Root.ServiceContainer[svc.Name] = ""
}
