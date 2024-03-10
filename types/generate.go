package types

import (
	"path/filepath"
	"regexp"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
)

type (
	validateData struct {
		VarName     string
		Name        string
		ValidateDef string
		Ref         string
	}
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
			files = append(files, typesFile(genpkg, r))
		}
	}
	return files, nil
}

func typesFile(genpkg string, r *expr.RootExpr) *codegen.File {
	types := r.Types
	path := filepath.Join(codegen.Gendir, "types", "types.go")
	header := codegen.Header("Data types", "types",
		[]*codegen.ImportSpec{
			codegen.GoaImport(""),
			{Path: "unicode/utf8"},
		},
	)
	sections := []*codegen.SectionTemplate{header}

	// Create dummy service so we can leverage Goa's code generation.
	svc := &expr.ServiceExpr{Name: "dummy"}
	expr.Root.Services = []*expr.ServiceExpr{svc}

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

	var vdata []validateData
	scope := codegen.NewNameScope()
	attCtx := codegen.AttributeContext{Scope: codegen.NewAttributeScope(scope)}
	addValidation := func(t expr.UserType) {
		def := codegen.ValidationCode(t.Attribute(), t, &attCtx, true, expr.IsAlias(t), false, "v")
		if def == "" {
			return
		}
		vdata = append(vdata, validateData{
			VarName:     codegen.Goify(t.Name(), true),
			Name:        t.Name(),
			ValidateDef: def,
			Ref:         scope.GoTypeRef(&expr.AttributeExpr{Type: t}),
		})
	}
	seen := make(map[string]struct{})
	for _, t := range r.Types {
		collectUserTypes(t, addValidation, seen)
	}
	sections = append(sections, &codegen.SectionTemplate{
		Name:   "type-validation",
		Source: validateT,
		Data:   vdata,
	})

	return &codegen.File{Path: path, SectionTemplates: sections}
}

// collectUserTypes traverses the given data type recursively and calls back the
// given function for each attribute using a user type.
func collectUserTypes(dt expr.DataType, cb func(expr.UserType), seen map[string]struct{}) {
	if dt == expr.Empty {
		return
	}
	switch actual := dt.(type) {
	case *expr.Object:
		for _, nat := range *actual {
			collectUserTypes(nat.Attribute.Type, cb, seen)
		}
	case *expr.Array:
		collectUserTypes(actual.ElemType.Type, cb, seen)
	case *expr.Map:
		collectUserTypes(actual.KeyType.Type, cb, seen)
		collectUserTypes(actual.ElemType.Type, cb, seen)
	case expr.UserType:
		if _, ok := seen[actual.ID()]; ok {
			return
		}
		seen[actual.ID()] = struct{}{}
		cb(actual)
		collectUserTypes(actual.Attribute().Type, cb, seen)
	}
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

const validateT = `{{range . }}{{ printf "Validate%s runs the validations defined on %s" .VarName .Name | comment }}
func Validate{{ .VarName }}(v {{ .Ref }}) (err error) {
	{{ .ValidateDef }}
	return
}
{{ end }}

`
