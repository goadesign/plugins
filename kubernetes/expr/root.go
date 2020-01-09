package expr

import "goa.design/goa/v3/eval"

// Root is root kubernetes object built by the DSL.
var Root = &RootExpr{}

type (
	// RootExpr is the struct built by the DSL on process start.
	RootExpr struct {
		// Services contains the list of kubernetes services..
		Services []*ServiceExpr
		// Deployments contains the list of kubernetes deployments.
		Deployments []*DeploymentExpr
	}
)

// Register DSL roots.
func init() {
	if err := eval.Register(Root); err != nil {
		panic(err)
	}
}

// EvalName is the name of the DSL.
func (r *RootExpr) EvalName() string {
	return "kubernetes design"
}

// WalkSets returns the expressions in order of evaluation.
func (r *RootExpr) WalkSets(walk eval.SetWalker) {
	svcs := make(eval.ExpressionSet, len(r.Services))
	for i, s := range r.Services {
		svcs[i] = s
	}
	walk(svcs)

	depls := make(eval.ExpressionSet, len(r.Deployments))
	for i, d := range r.Deployments {
		depls[i] = d
	}
	walk(depls)
}

// DependsOn returns nil, the kubernetes DSL has no dependency.
func (r *RootExpr) DependsOn() []eval.Root { return nil }

// Packages returns the Go import path to this and the dsl packages.
func (r *RootExpr) Packages() []string {
	return []string{
		"goa.design/plugins/v3/kubernetes/expr",
		"goa.design/plugins/v3/kubernetes/dsl",
	}
}
