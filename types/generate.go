package types

import (
	"path/filepath"
	"regexp"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
)

// Name of directory that contains generated types.
const Gendir = "types"

// init registers the plugin generator function.
func init() {
	codegen.RegisterPlugin("types", "gen", nil, Generate)
}

// Generate produces the documentation JSON file.
func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, root := range roots {
		if r, ok := root.(*expr.RootExpr); ok {
			files = append(files, typesFile(genpkg, r.Types))
		}
	}
	return files, nil
}

func typesFile(genpkg string, types []expr.UserType) *codegen.File {
	path := filepath.Join(codegen.Gendir, "types", "types.go")
	header := codegen.Header("Data types", "types",
		[]*codegen.ImportSpec{
			codegen.GoaImport(""),
		},
	)
	sections := []*codegen.SectionTemplate{header}

	// Create dummy service so we can leverage Goa's code generation.
	svc := &expr.ServiceExpr{Name: "dummy"}
	expr.Root.Services = append(expr.Root.Services, svc)

	// Create a dummy method for each type.
	for _, t := range types {
		svc.Methods = append(svc.Methods, &expr.MethodExpr{
			Name:             t.Name() + "M",
			Payload:          &expr.AttributeExpr{Type: t},
			StreamingPayload: &expr.AttributeExpr{Type: expr.Empty},
			Result:           &expr.AttributeExpr{Type: expr.Empty},
			Service:          svc,
		})
	}

	// Generate the code and retrieve the relevant sections.
	files := service.Files(genpkg, svc, make(map[string][]string))
	for _, f := range files {
		for _, section := range f.SectionTemplates {
			sn := section.Name
			if sn == "service-payload" ||
				sn == "service-union-value-method" ||
				sn == "service-user-type" {
				if sn == "service-payload" {
					// Override the payload comment with the original type description.
					section.FuncMap = map[string]interface{}{
						"comment": func(s string) string { return getDescription(s, types) },
					}
				}
				sections = append(sections, section)
			}
		}
	}

	return &codegen.File{Path: path, SectionTemplates: sections}
}

var nameRegex = regexp.MustCompile(`([^\s]*)`)

func getDescription(comment string, types []expr.UserType) string {
	name := nameRegex.FindStringSubmatch(comment)[1]
	for _, t := range types {
		if t.Name() == name {
			return codegen.Comment(t.Attribute().Description)
		}
	}
	return ""
}
