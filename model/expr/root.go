package expr

import (
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
)

// DefaultExcluded specifies a list of default exclusion tags. Containers that
// have any of these tags will be excluded from the check.
var DefaultExcluded []string = []string{"database", "external", "thirdparty"}

// DefaultFormat is the default format string used to compute the name of the
// model container from the service name.
const DefaultFormat string = "%s Service"

// Root is the design root expression.
var Root *RootExpr

type (
	// RootExpr keeps track of the Model package path defined in the design.
	RootExpr struct {
		// ModelPkgPath is the path to the directory containing the model files.
		ModelPkgPath string
		// SystemName is the name of the system to be validated.
		SystemName string
		// ContainerNameFormat is the format string used to compute the name of
		// the model container from the service name.
		ContainerNameFormat string
		// ExcludedTags specifies a list of tags. Containers that have any of these
		// tags will be excluded from the check.
		ExcludedTags []string
		// Complete is true if all containers have a corresponding service.
		Complete bool
		// ServiceContainer maps service names to container names.
		ServiceContainer map[string]string
	}
)

// Register design root with eval engine.
func init() {
	if err := eval.Register(Root); err != nil {
		panic(err)
	}
	InitRoot()
}

// InitRoot initializes the root expression.
func InitRoot() {
	Root = &RootExpr{
		ContainerNameFormat: DefaultFormat,
		ExcludedTags:        DefaultExcluded,
		ServiceContainer:    make(map[string]string),
	}
}

// EvalName returns the name used in error messages.
func (r *RootExpr) EvalName() string {
	return "Model plugin"
}

// WalkSets is a no-op as the Model plugin does not define any embedded DSL.
func (r *RootExpr) WalkSets(walk eval.SetWalker) {
}

// DependsOn tells the eval engine to run the Goa DSL first.
func (r *RootExpr) DependsOn() []eval.Root {
	return []eval.Root{expr.Root}
}

// Packages returns the import path to the Go packages that make
// up the DSL. This is used to skip frames that point to files
// in these packages when computing the location of errors.
func (r *RootExpr) Packages() []string {
	return []string{"goa.design/plugins/v3/model/dsl"}
}
