package http

import (
	"testing"

	goadesign "goa.design/goa/design"
	"goa.design/goa/eval"
	httpdesign "goa.design/goa/http/design"
	"goa.design/plugins/security/design"
)

// RunHTTPDSL returns the http DSL root resulting from running the given DSL.
func RunHTTPDSL(t *testing.T, dsl func(), schemes []*design.SchemeExpr) *httpdesign.RootExpr {
	// reset all roots and codegen data structures
	eval.Reset()
	goadesign.Root = new(goadesign.RootExpr)
	httpdesign.Root = &httpdesign.RootExpr{Design: goadesign.Root}
	design.Root = &design.RootExpr{}
	eval.Register(goadesign.Root)
	eval.Register(httpdesign.Root)
	eval.Register(design.Root)

	// Add the top-level security schemes to the newly created security RootExpr
	for _, s := range schemes {
		design.Root.Schemes = append(design.Root.Schemes, s)
	}

	goadesign.Root.API = &goadesign.APIExpr{
		Name:    "test api",
		Servers: []*goadesign.ServerExpr{{URL: "http://localhost"}},
	}

	// run DSL (first pass)
	if !eval.Execute(dsl, nil) {
		t.Fatal(eval.Context.Error())
	}

	// run DSL (second pass)
	if err := eval.RunDSL(); err != nil {
		t.Fatal(err)
	}

	// return generated root
	return httpdesign.Root
}
