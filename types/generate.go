// This file plans and writes standalone Go types for every named type in a
// Goa design. Planning records names, imports, field layouts, and validations
// before Goa freezes the generated packages. Rendering only formats those
// recorded values, so it never changes the evaluated design.
package types

import (
	"cmp"
	"fmt"
	"path"
	"path/filepath"
	"sort"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/generator"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
)

type (
	// generationPlan stores everything needed to write one types.go file after
	// Goa has chosen the final declaration and import names.
	generationPlan struct {
		pkg        *codegen.GeneratedPackage
		imports    map[string]struct{}
		types      []*typePlan
		hasGoaRoot bool
	}

	// typePlan stores one named type and its optional validation function.
	typePlan struct {
		name        string
		description string
		declaration *codegen.TypeDeclaration
		definition  *codegen.GoTypePlan
		reference   *codegen.GoTypePlan
		validator   *codegen.NameDeclaration
		validation  *codegen.ValidationPlan
	}

	// typeData contains the final strings used by the type template.
	typeData struct {
		Name        string
		Description string
		Definition  string
	}

	// validateData contains the final strings used by the validation template.
	validateData struct {
		VarName     string
		Name        string
		ValidateDef string
		Ref         string
	}

	// validatorNameOrder gives each generated validator a stable position when
	// another plugin requests the same function name.
	validatorNameOrder string
)

// Gendir is the directory that contains the generated standalone types.
const Gendir = "types"

// init registers the Go and Protocol Buffer type generators.
func init() {
	generator.RegisterPlugin("types", "gen", newPlugin)
	codegen.RegisterPlugin("types-proto", "gen", nil, GenerateProto)
}

// Generate adds the standalone Go types to files. Goa generation commands use
// the planning-aware plugin registered by init; this function remains for code
// that invoked the released plugin callback directly.
func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	generation, err := codegen.NewGeneration(genpkg, roots)
	if err != nil {
		return nil, err
	}
	planned, err := planTypes(generation)
	if err != nil {
		return nil, err
	}
	if err := generation.Freeze(); err != nil {
		return nil, err
	}
	return generateTypes(planned, files)
}

// newPlugin creates isolated planning state for one Goa generation command.
func newPlugin() generator.Plugin {
	var planned *generationPlan
	return generator.Plugin{
		Plan: func(plan *generator.Plan) error {
			var err error
			planned, err = planTypes(plan.Generation())
			return err
		},
		Generate: func(_ *generator.Plan, files []*codegen.File) ([]*codegen.File, error) {
			return generateTypes(planned, files)
		},
	}
}

// planTypes declares every standalone type, validator, and import before Goa
// chooses their final names. It reads each evaluated design without changing
// its services or named types.
func planTypes(generation *codegen.Generation) (*generationPlan, error) {
	planned := &generationPlan{imports: make(map[string]struct{})}
	var userTypes []expr.UserType
	for _, evaluated := range generation.Roots() {
		root, ok := evaluated.(*expr.RootExpr)
		if !ok {
			continue
		}
		planned.hasGoaRoot = true
		userTypes = append(userTypes, root.Types...)
	}
	if !planned.hasGoaRoot {
		return planned, nil
	}

	generatedPackage, err := generation.ClaimPackage(path.Join(generation.GenPkg(), Gendir))
	if err != nil {
		return nil, err
	}
	planned.pkg = generatedPackage

	declarations := make(map[expr.UserType]*codegen.TypeDeclaration, len(userTypes))
	for _, userType := range userTypes {
		origin := userType.Origin()
		if _, exists := declarations[origin]; exists {
			continue
		}
		declaration, err := generatedPackage.DeclareUserType(userType)
		if err != nil {
			return nil, err
		}
		declarations[origin] = declaration
	}

	policy := codegen.GoLayoutPolicy{UseDefault: true, SumType: true}
	bindType := typeBinder(generatedPackage, declarations)
	validators := make(map[expr.UserType]*codegen.NameDeclaration, len(declarations))
	for _, userType := range userTypes {
		origin := userType.Origin()
		if _, exists := validators[origin]; exists || !codegen.NeedsValidation(userType.Attribute(), policy) {
			continue
		}
		validator, err := generatedPackage.DeclareDependentName(
			codegen.NameFunction,
			declarations[origin].Declaration(),
			"Validate",
			"",
			validatorNameOrder(userType.Name()),
		)
		if err != nil {
			return nil, err
		}
		validators[origin] = validator
	}

	seen := make(map[expr.UserType]struct{}, len(declarations))
	for _, userType := range userTypes {
		origin := userType.Origin()
		if _, exists := seen[origin]; exists {
			continue
		}
		seen[origin] = struct{}{}
		definition, err := codegen.PlanGoType(userType.Attribute(), codegen.GoTypePlanOptions{
			Owner:  generatedPackage.ImportPath(),
			Policy: policy,
			Bind:   bindType,
		})
		if err != nil {
			return nil, err
		}
		reference, err := codegen.PlanGoType(&expr.AttributeExpr{Type: userType}, codegen.GoTypePlanOptions{
			Owner:  generatedPackage.ImportPath(),
			Policy: policy,
			Bind:   bindType,
		})
		if err != nil {
			return nil, err
		}
		current := &typePlan{
			name:        userType.Name(),
			description: userType.Attribute().Description,
			declaration: declarations[origin],
			definition:  definition,
			reference:   reference,
			validator:   validators[origin],
		}
		if current.validator != nil {
			current.validation, err = codegen.NewValidationPlan(
				userType.Attribute(),
				definition,
				codegen.ValidationPlanOptions{
					Required: true,
					Alias:    expr.IsAlias(userType),
					Bind:     validatorBinder(validators),
				},
			)
			if err != nil {
				return nil, err
			}
		}
		if err := planImports(planned, definition, current.validation); err != nil {
			return nil, err
		}
		planned.types = append(planned.types, current)
	}
	return planned, nil
}

// generateTypes appends the planned types.go file after Goa has frozen all
// package names. A generation with no Goa service-design root adds no file.
func generateTypes(planned *generationPlan, files []*codegen.File) ([]*codegen.File, error) {
	if !planned.hasGoaRoot {
		return files, nil
	}
	return append(files, typesFile(planned)), nil
}

// typesFile formats recorded layouts and validations without reading the
// evaluated design again.
func typesFile(planned *generationPlan) *codegen.File {
	qualifier := planned.pkg.ImportName
	owner := planned.pkg.ImportPath()
	types := make([]typeData, len(planned.types))
	var validations []validateData
	for index, current := range planned.types {
		types[index] = typeData{
			Name:        current.declaration.Name(),
			Description: current.description,
			Definition:  current.definition.Link(owner, qualifier).Def(),
		}
		if current.validation == nil {
			continue
		}
		linked, err := current.validation.Link(current.definition.Link(owner, qualifier))
		if err != nil {
			panic(err)
		}
		validations = append(validations, validateData{
			VarName:     current.validator.Name(),
			Name:        current.name,
			ValidateDef: linked.Render("v", "v"),
			Ref:         current.reference.Link(owner, qualifier).Ref(),
		})
	}
	sort.Slice(types, func(left, right int) bool {
		return types[left].Name < types[right].Name
	})

	paths := make([]string, 0, len(planned.imports))
	for importPath := range planned.imports {
		paths = append(paths, importPath)
	}
	sort.Strings(paths)
	imports := make([]*codegen.ImportSpec, len(paths))
	for index, importPath := range paths {
		imports[index] = planned.pkg.Import(importPath)
	}
	return &codegen.File{
		Path: filepath.Join(codegen.Gendir, Gendir, "types.go"),
		SectionTemplates: []*codegen.SectionTemplate{
			codegen.Header("Data types", "types", imports),
			{Name: "type-definitions", Source: typesT, Data: types},
			{Name: "type-validation", Source: validateT, Data: validations},
		},
	}
}

// typeBinder returns the type declaration already assigned to each named
// design type. Standalone OneOf values are rejected because this plugin has no
// public OneOf representation to preserve.
func typeBinder(generatedPackage *codegen.GeneratedPackage, declarations map[expr.UserType]*codegen.TypeDeclaration) codegen.GoTypeBinder {
	return func(request codegen.GoTypeBindingRequest) (codegen.GoTypeBinding, error) {
		if request.Kind == codegen.GoUnion {
			return codegen.GoTypeBinding{}, fmt.Errorf("types plugin does not support Goa OneOf types")
		}
		userType, ok := request.Attribute.Type.(expr.UserType)
		if !ok {
			return codegen.GoTypeBinding{}, fmt.Errorf("types plugin cannot bind %s as a named type", request.Kind)
		}
		declaration := declarations[userType.Origin()]
		if declaration == nil {
			return codegen.GoTypeBinding{}, fmt.Errorf("types plugin has no declaration for type %q", userType.Name())
		}
		return codegen.GoTypeBinding{Owner: generatedPackage.ImportPath(), Type: declaration}, nil
	}
}

// validatorBinder returns the validation function already assigned to each
// nested named type that needs validation.
func validatorBinder(validators map[expr.UserType]*codegen.NameDeclaration) codegen.ValidatorDeclarationBinder {
	return func(request codegen.ValidatorBindingRequest) (*codegen.NameDeclaration, error) {
		userType, ok := request.Attribute.Type.(expr.UserType)
		if !ok {
			return nil, fmt.Errorf("types plugin cannot validate unnamed nested type")
		}
		validator := validators[userType.Origin()]
		if validator == nil {
			return nil, fmt.Errorf("types plugin has no validator for type %q", userType.Name())
		}
		return validator, nil
	}
}

// planImports registers each package referenced by a type definition or its
// validation code and remembers which import lines the file must contain.
func planImports(planned *generationPlan, definition *codegen.GoTypePlan, validation *codegen.ValidationPlan) error {
	imports := definition.ImportPreferences()
	if validation != nil {
		imports = append(imports, validation.ImportPreferences()...)
	}
	for _, goImport := range imports {
		if goImport.Path == planned.pkg.ImportPath() {
			continue
		}
		if err := planned.pkg.DeclareImport(&codegen.ImportSpec{Name: goImport.Name, Path: goImport.Path}); err != nil {
			return err
		}
		planned.imports[goImport.Path] = struct{}{}
	}
	return nil
}

// ComparePackageName orders validators by the authored type name they check.
func (o validatorNameOrder) ComparePackageName(other codegen.PackageNameOrder) int {
	return cmp.Compare(o, other.(validatorNameOrder))
}

const typesT = `{{ range . }}
{{- if .Description }}
{{ .Description | comment }}
{{- end }}
type {{ .Name }} {{ .Definition }}
{{ end }}
`

const validateT = `{{ range . }}
{{ printf "%s runs the validations defined on %s" .VarName .Name | comment }}
func {{ .VarName }}(v {{ .Ref }}) (err error) {
	{{ .ValidateDef }}
	return
}
{{ end }}
`
