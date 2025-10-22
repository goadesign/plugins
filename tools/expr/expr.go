package expr

import (
	"fmt"

	"goa.design/goa/v3/eval"
	goaexpr "goa.design/goa/v3/expr"
)

type (
	// RootExpr aggregates tool sets grouped by service.
	RootExpr struct {
		ToolSets []*ToolSetExpr
	}

	// ToolSetExpr groups tools under a common name (per-service).
	ToolSetExpr struct {
		Name    string
		Service *goaexpr.ServiceExpr
		Tools   []*ToolExpr
	}

	// ToolExpr represents an individual tool definition.
	ToolExpr struct {
		Name    string
		ToolSet *ToolSetExpr
		Method  *goaexpr.MethodExpr
		DSLFunc func()
		Derived bool
	}
)

// Root is the design root expression.
var Root = &RootExpr{}

func init() {
	if err := eval.Register(Root); err != nil {
		panic(err)
	}
}

// EvalName implements eval.Expression.
func (r *RootExpr) EvalName() string {
	return "tools plugin"
}

// WalkSets satisfies eval.Root.
func (r *RootExpr) WalkSets(walk eval.SetWalker) {
	if len(r.ToolSets) == 0 {
		return
	}
	walk(eval.ToExpressionSet(r.ToolSets))
}

// DependsOn ensures the Goa core DSL executes first.
func (r *RootExpr) DependsOn() []eval.Root {
	return []eval.Root{goaexpr.Root}
}

// Packages returns the DSL import path for better error locations.
func (r *RootExpr) Packages() []string {
	return []string{"goa.design/plugins/v3/tools/dsl"}
}

// Reset clears accumulated tool sets between DSL runs.
func (r *RootExpr) Reset() {
	r.ToolSets = nil
}

// / EvalName implements eval.Expression.
func (t *ToolSetExpr) EvalName() string {
	if t.Service == nil {
		return fmt.Sprintf("toolset %s", t.Name)
	}
	return fmt.Sprintf("toolset %s::%s", t.Service.Name, t.Name)
}

// EvalName implements eval.Expression.
func (t *ToolExpr) EvalName() string {
	if t.ToolSet == nil {
		return fmt.Sprintf("tool %s", t.Name)
	}
	return fmt.Sprintf("tool %s::%s", t.ToolSet.EvalName(), t.Name)
}

// Finalize ensures the payload and result types referenced by tools are marked
// with the force-generate meta so they are emitted by downstream generators
// even if they are not explicitly referenced by a service method.
func (t *ToolSetExpr) Finalize() {
	for _, tool := range t.Tools {
		if tool == nil || tool.Method == nil {
			continue
		}
		if tool.Method.Payload != nil {
			tool.Method.Payload.AddMeta("type:generate:force", "true")
		}
		if tool.Method.Result != nil {
			tool.Method.Result.AddMeta("type:generate:force", "true")
		}
	}
}
