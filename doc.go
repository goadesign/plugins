/*
Packages in the Plugin directory contain plugins for [goa
v2](https://godoc.org/goa.design/goa/v3). Plugins can extend the goa DSL, generate
new artifacts and modify the output of existing generators.

There are currently two plugins in the directory:

  * The [goakit](https://godoc.org/goa.design/plugins/goakit) plugin generates
    code that integrates with the [go-kit](https://github.com/go-kit/kit)
    library.

  * The [cors](https://godoc.org/goa.design/plugins/cors) plugin adds new DSL to
    define CORS policies. The plugin generates HTTP server code that respond to
    CORS requests in compliance with the policies defined in the design.

Writing a Plugin

Writing a plugin consists of two steps:

1. Writing the functions that gets called by the goa tool during code generation.

2. Registering the functions with the goa tool.

A plugin may implement the "gen" goa command, the "example" goa command or both.
In each case a plugin may register one or two functions: the first function
called "Prepare" gets called prior to any code generation actually happening.
The function can modify the design data structures before goa uses them to
generate code. The second function called "Generate" does the actual code
generation.

The signature of the Generate function is:

    func (genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error)

where:

    "genpkg" is the Go import path to the top level generated package ("gen")
    "roots" is the set of design roots created by the DSL.
    "files" is the current set of generated files.

The function must return the entire set of generated files (even the files that
the plugin does not modify).

The functions must then be registered with the goa code generator tool using the
"RegisterPlugin" function of the "codegen" package. This is typically done in a
package "init" function, for example:

	// Register the plugin.
	func init() {
		codegen.RegisterPlugin("gen", Prepare, Generate)
		codegen.RegisterPlugin("example", Prepare, Example)
	}

The first argument of RegisterPlugin must be one of "gen" or "example" and
specifies the "goa" command that triggers the call to the plugin functions. The
second argument is the Prepare function if any, nil otherwise. The last argument
is the code generator function if any, nil otherwise.

Extending The DSL

A plugin may introduce new DSL "keywords" (typically Go package functions). The
DSL functions initialize the content of a design root object (an object that
implements the "eval.Root" interface), for example:

	func MyDSLMethod(arg string) {
		// validate where DSL function is called, in this example it must be called
		// outside of all other expression
		if _, ok := eval.Current().(eval.TopExpr); !ok {
			eval.IncompatibleDSL()
			return nil
		}

		// create corresponding design object
		expr := &design.MyExpr{ Name: arg }

		// store design object in root
		design.Root.Exprs = append(design.Root.Exprs, expr)
	}

The [eval](https://godoc.org/goa.design/goa/v3/codegen/eval) package contains a
number of functions that can be leveraged to implement the DSL such as error
reporting functions.

The plugin design root object is then given to the plugin Generate function (if
there's one) which may take advantage of it to generate the code.
*/
package Plugin
