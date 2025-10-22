package dsl

import (
	"goa.design/goa/v3/eval"
	goaexpr "goa.design/goa/v3/expr"
	"goa.design/plugins/v3/tools/expr"

	// Register code generators for the tools plugin when the DSL is imported.
	_ "goa.design/plugins/v3/tools"
)

// ToolSet groups tool definitions under a logical name scoped to the current
// service.
//
// ToolSet must appear in a Service expression.
//
// ToolSet takes two arguments: the name of the tool set and the defining DSL.
// The DSL typically calls Tool or ToolFromMethod to register tools.
//
// Example:
//
//	import . "goa.design/plugins/v3/tools/dsl"
//
//	var _ = Service("inventory", func() {
//	    ToolSet("admin", func() {
//	        // Define a pure tool with its own payload/result types
//	        Tool("rebuild_index", func() {
//	            Description("Rebuilds the search index")
//	            Payload(func() {
//	                Attribute("force", Boolean)
//	                Required("force")
//	            })
//	            Result(Empty)
//	        })
//
//	        // Reuse an existing service method as a tool
//	        ToolFromMethod("delete")
//	    })
//	})
func ToolSet(name string, fn func()) {
	svc, ok := eval.Current().(*goaexpr.ServiceExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	ts := &expr.ToolSetExpr{
		Service: svc,
		Name:    name,
	}
	expr.Root.ToolSets = append(expr.Root.ToolSets, ts)
	if fn != nil {
		_ = eval.Execute(fn, ts)
	}
}

// Tool defines a pure tool that does not map to an existing service method.
// The tool DSL may use standard Goa Method DSL helpers such as Payload, Result,
// and Description.
//
// Tool must appear in a ToolSet expression.
//
// Tool accepts the tool name and an optional defining DSL.
//
// Example:
//
//	Tool("rotate_keys", func() {
//	    Description("Rotates API keys for downstream systems")
//	    Payload(func() {
//	        Attribute("dry_run", Boolean)
//	    })
//	    Result(Empty)
//	})
func Tool(name string, fn ...func()) {
	ts, ok := eval.Current().(*expr.ToolSetExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	var dsl func()
	if len(fn) > 0 && fn[len(fn)-1] != nil {
		dsl = fn[len(fn)-1]
	}
	method := &goaexpr.MethodExpr{
		Name:    name,
		Service: ts.Service,
	}
	tool := &expr.ToolExpr{
		Name:    name,
		ToolSet: ts,
		Method:  method,
		DSLFunc: func() {
			eval.Execute(dsl, method)
		},
	}
	ts.Tools = append(ts.Tools, tool)
	if tool.DSLFunc != nil {
		tool.DSLFunc()
		tool.DSLFunc = nil
	}
}

// ToolFromMethod registers an existing service method as a tool within the
// current tool set.  The tool's payload and result types are automatically
// inferred from the referenced method expression.
//
// ToolFromMethod must appear in a ToolSet expression.
//
// ToolFromMethod accepts one or two arguments. The first argument is the
// unqualified name of a method.  An optional second argument may be provided to
// specify a different tool name; if omitted, the method name is used as the
// tool name by default.
//
// Example:
//
//	var _ = Service("inventory", func() {
//	    Method("delete", func() {
//	        Payload(func() { Attribute("id", String) ; Required("id") })
//	        Result(Empty)
//	    })
//
//	    ToolSet("ops", func() {
//	        ToolFromMethod("delete", "delete_item")
//	    })
//	})
func ToolFromMethod(method string, toolName ...string) {
	ts, ok := eval.Current().(*expr.ToolSetExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	m := ts.Service.Method(method)
	if m == nil {
		eval.ReportError("tool references unknown method %s::%s", ts.Service.Name, method)
		return
	}
	tool := method
	if len(toolName) > 0 {
		tool = toolName[0]
	}
	ts.Tools = append(ts.Tools, &expr.ToolExpr{
		Name:    tool,
		ToolSet: ts,
		Method:  m,
		Derived: true,
	})
}
