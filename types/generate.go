package types

import (
	"fmt"
	"path/filepath"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
)

type (
	typeData struct {
		VarName     string
		Description string
		Def         string
	}

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
func Generate(_ string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, root := range roots {
		if r, ok := root.(*expr.RootExpr); ok {
			files = append(files, typesFile(r))
		}
	}
	return files, nil
}

func typesFile(r *expr.RootExpr) *codegen.File {
	path := filepath.Join(codegen.Gendir, "types", "types.go")
	header := codegen.Header("Data types", "types",
		[]*codegen.ImportSpec{
			codegen.GoaImport(""),
		},
	)
	scope := codegen.NewNameScope()
	sections := []*codegen.SectionTemplate{header}

	data := make([]typeData, len(r.Types))
	for i, t := range r.Types {
		data[i] = typeData{
			VarName:     codegen.Goify(t.Name(), true),
			Description: t.Attribute().Description,
			Def:         scope.GoTypeDef(t.Attribute(), false, false),
		}
	}
	sections = append(sections, &codegen.SectionTemplate{
		Name:   "type-decl",
		Source: typeDeclT,
		Data:   data,
	})

	var vdata []validateData
	attCtx := codegen.AttributeContext{Scope: codegen.NewAttributeScope(scope)}
	for _, t := range r.Types {
		def := codegen.RecursiveValidationCode(t.Attribute(), &attCtx, true, expr.IsAlias(t), "v")
		fmt.Println(def)
		vdata = append(vdata, validateData{
			VarName:     codegen.Goify(t.Name(), true),
			Name:        t.Name(),
			ValidateDef: def,
			Ref:         scope.GoTypeRef(&expr.AttributeExpr{Type: t}),
		})
	}
	sections = append(sections, &codegen.SectionTemplate{
		Name:   "type-validation",
		Source: validateT,
		Data:   vdata,
	})

	return &codegen.File{
		Path:             path,
		SectionTemplates: sections,
	}
}

// input: TypeData
const typeDeclT = `type (
	{{range . }}{{ if .Description }}{{ comment .Description }}
	{{ end }}{{ .VarName }} {{ .Def }}
{{ end }})

`

const validateT = `{{range . }}{{ printf "Validate%s runs the validations defined on %s" .VarName .Name | comment }}
func Validate{{ .VarName }}(v {{ .Ref }}) (err error) {
	{{ .ValidateDef }}
	return
}
{{ end }}

`
