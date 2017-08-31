package security

import (
	"strings"

	"goa.design/goa/codegen"
	"goa.design/goa/eval"
	httpdesign "goa.design/goa/http/design"
)

// Register the plugin Generator functions.
func init() {
	codegen.RegisterPlugin("gen", Generate)
}

// Generate produces server code that enforce the security requirements defined
// in the design. Generate also produces client code that makes it possible to
// provide the required security artifacts. Finally Generate also generate code
// that initializes the context given to the service methods with security
// information.
func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	var output []*codegen.File
	for _, f := range files {
		for _, s := range f.SectionTemplates {
			if s.Name == "server-handler-init" {
				s.Source = strings.Replace(s.Source, "{{ .HandlerInit }}", "{{ .HandlerInit }}Unsecured")
			}
		}
	}
	for _, root := range roots {
		if r, ok := root.(*httpdesign.RootExpr); ok {
			output = append(output, SecuredHandlerFiles(genpkg, r)...)
		}
	}
	return output, nil
}
