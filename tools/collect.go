package tools

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	expr "goa.design/goa/v3/expr"
	"goa.design/goa/v3/http/codegen/openapi"

	toolsexpr "goa.design/plugins/v3/tools/expr"
)

type generator struct {
	genpkg      string
	services    *service.ServicesData
	typeScope   *codegen.NameScope
	types       map[string]*typeInfo
	ordered     []*typeInfo
	tools       []*toolEntry
	imports     map[string]*codegen.ImportSpec
	needsGoa    bool
	typeImports map[string]*codegen.ImportSpec
}

type toolEntry struct {
	Name    string
	Service string
	Set     string
	Method  *string
	Payload *typeInfo
	Result  *typeInfo
}

type typeInfo struct {
	Key           string
	TypeName      string
	Doc           string
	Def           string
	SchemaVar     string
	SchemaLiteral string
	ExportedCodec string
	GenericCodec  string
	MarshalFunc   string
	UnmarshalFunc string
	ValidateFunc  string
	Validation    string
	HasValidation bool
	FullRef       string
	ElemRef       string
	Pointer       bool
	CheckNil      bool
	MarshalArg    string
	UnmarshalArg  string
	ValidationSrc []string
	NeedType      bool
	Import        *codegen.ImportSpec
	NilError      string
	DecodeError   string
	ValidateError string
	EmptyError    string
	Usage         typeUsage
}

type typeUsage string

const (
	usagePayload typeUsage = "payload"
	usageResult  typeUsage = "result"
)

func (u typeUsage) String() string {
	return string(u)
}

func newGenerator(genpkg string) (*generator, error) {
	return &generator{
		genpkg:      genpkg,
		services:    service.NewServicesData(expr.Root),
		typeScope:   codegen.NewNameScope(),
		types:       make(map[string]*typeInfo),
		imports:     make(map[string]*codegen.ImportSpec),
		typeImports: make(map[string]*codegen.ImportSpec),
	}, nil
}

func (g *generator) collect() error {
	if toolsexpr.Root == nil {
		return nil
	}
	for _, ts := range toolsexpr.Root.ToolSets {
		var (
			svcData *service.Data
			svcName string
		)
		if ts.Service != nil {
			svcName = ts.Service.Name
			svcData = g.services.Get(svcName)
			if svcData == nil {
				return fmt.Errorf("tools: service %q not found in design", svcName)
			}
		}
		for _, tool := range ts.Tools {
			executeToolDSL(tool)
			if tool.Method == nil {
				return fmt.Errorf("tools: tool %q missing method expression", tool.Name)
			}
			payload, err := g.typeFor(tool, tool.Method.Payload, usagePayload, svcData)
			if err != nil {
				return err
			}
			result, err := g.typeFor(tool, tool.Method.Result, usageResult, svcData)
			if err != nil {
				return err
			}
			entry := &toolEntry{
				Name:    tool.Name,
				Service: svcName,
				Set:     ts.Name,
				Payload: payload,
				Result:  result,
			}
			if isDerived(tool) && tool.Method != nil {
				method := tool.Method.Name
				entry.Method = &method
			}
			g.tools = append(g.tools, entry)
		}
	}
	sort.Slice(g.tools, func(i, j int) bool {
		if g.tools[i].Service == g.tools[j].Service {
			return g.tools[i].Name < g.tools[j].Name
		}
		return g.tools[i].Service < g.tools[j].Service
	})
	return nil
}

func (g *generator) typeFor(tool *toolsexpr.ToolExpr, att *expr.AttributeExpr, usage typeUsage, svcData *service.Data) (*typeInfo, error) {
	if att == nil || att.Type == nil || att.Type == expr.Empty {
		return nil, nil
	}
	if isDerived(tool) {
		return g.ensureDerivedType(tool, att, usage, svcData)
	}
	return g.ensurePureType(tool, att, usage)
}

func (g *generator) ensurePureType(tool *toolsexpr.ToolExpr, att *expr.AttributeExpr, usage typeUsage) (*typeInfo, error) {
	base := codegen.Goify(tool.Name, true)
	switch usage {
	case usagePayload:
		base += "Payload"
	case usageResult:
		base += "Result"
	}
	typeName := g.typeScope.Unique(base, "")
	key := "*" + typeName
	if existing := g.types[key]; existing != nil {
		return existing, nil
	}
	doc := fmt.Sprintf("%s defines the JSON %s for the %s tool.", typeName, usage.String(), tool.Name)
	def := fmt.Sprintf("%s %s", typeName, g.typeScope.GoTypeDef(att, false, true))
	schemaBytes, err := schemaForAttribute(tool.Method, att)
	if err != nil {
		return nil, err
	}
	schemaLiteral := formatSchema(schemaBytes)
	schemaVar := ""
	if schemaLiteral != "" {
		schemaVar = lowerCamel(typeName) + "Schema"
	}
	validation := buildValidation(att, "", g.typeScope)
	info := &typeInfo{
		Key:           key,
		TypeName:      typeName,
		Doc:           doc,
		Def:           def,
		SchemaVar:     schemaVar,
		SchemaLiteral: schemaLiteral,
		ExportedCodec: typeName + "Codec",
		GenericCodec:  lowerCamel(typeName) + "Codec",
		MarshalFunc:   "Marshal" + typeName,
		UnmarshalFunc: "Unmarshal" + typeName,
		ValidateFunc:  "Validate" + typeName,
		Validation:    validation,
		HasValidation: validation != "",
		FullRef:       "*" + typeName,
		ElemRef:       typeName,
		NeedType:      true,
		NilError:      fmt.Sprintf("%s is nil", lowerCamel(typeName)),
		DecodeError:   fmt.Sprintf("decode %s", lowerCamel(typeName)),
		ValidateError: fmt.Sprintf("validate %s", lowerCamel(typeName)),
		EmptyError:    fmt.Sprintf("%s JSON is empty", lowerCamel(typeName)),
		Usage:         usage,
	}
	if info.HasValidation {
		g.needsGoa = true
	}
	g.addTypeImports(att)
	finalizeTypeInfo(info)
	g.types[key] = info
	g.ordered = append(g.ordered, info)
	return info, nil
}

func (g *generator) ensureDerivedType(tool *toolsexpr.ToolExpr, att *expr.AttributeExpr, usage typeUsage, svcData *service.Data) (*typeInfo, error) {
	if svcData == nil {
		return nil, fmt.Errorf("tools: derived tool %q missing service data", tool.Name)
	}
	method := svcData.Method(tool.Method.Name)
	if method == nil {
		return nil, fmt.Errorf("tools: method %q for tool %q not found", tool.Method.Name, tool.Name)
	}
	var (
		typeName string
		typeRef  string
		loc      *codegen.Location
	)
	switch usage {
	case usagePayload:
		typeName = method.Payload
		typeRef = method.PayloadRef
		loc = method.PayloadLoc
	case usageResult:
		typeName = method.Result
		typeRef = method.ResultRef
		loc = method.ResultLoc
	default:
		typeName = method.Payload
		typeRef = method.PayloadRef
		loc = method.PayloadLoc
	}
	if typeName == "" {
		typeName = svcData.Scope.GoTypeName(att)
	}
	pkgName := packageName(loc, svcData)
	if typeRef == "" {
		typeRef = svcData.Scope.GoFullTypeRef(att, pkgName)
	}
	if typeRef == "" {
		return nil, fmt.Errorf("tools: unable to compute type reference for tool %q %s", tool.Name, usage)
	}
	key := typeRef
	if existing := g.types[key]; existing != nil {
		return existing, nil
	}
	schemaBytes, err := schemaForAttribute(tool.Method, att)
	if err != nil {
		return nil, err
	}
	schemaLiteral := formatSchema(schemaBytes)
	schemaVar := ""
	if schemaLiteral != "" {
		schemaVar = lowerCamel(typeName) + "Schema"
	}
	validation := buildValidation(att, pkgName, svcData.Scope)
	info := &typeInfo{
		Key:           key,
		TypeName:      typeName,
		SchemaVar:     schemaVar,
		SchemaLiteral: schemaLiteral,
		ExportedCodec: typeName + "Codec",
		GenericCodec:  lowerCamel(typeName) + "Codec",
		MarshalFunc:   "Marshal" + typeName,
		UnmarshalFunc: "Unmarshal" + typeName,
		ValidateFunc:  "Validate" + typeName,
		Validation:    validation,
		HasValidation: validation != "",
		FullRef:       typeRef,
		ElemRef:       strings.TrimPrefix(typeRef, "*"),
		NilError:      fmt.Sprintf("%s is nil", lowerCamel(typeName)),
		DecodeError:   fmt.Sprintf("decode %s", lowerCamel(typeName)),
		ValidateError: fmt.Sprintf("validate %s", lowerCamel(typeName)),
		EmptyError:    fmt.Sprintf("%s JSON is empty", lowerCamel(typeName)),
		Usage:         usage,
	}
	if info.HasValidation {
		g.needsGoa = true
	}
	finalizeTypeInfo(info)
	if importPath := importPath(loc, svcData, g.genpkg); importPath != "" {
		spec := &codegen.ImportSpec{Path: importPath}
		info.Import = spec
		g.imports[importPath] = spec
	}
	g.types[key] = info
	g.ordered = append(g.ordered, info)
	return info, nil
}

func packageName(loc *codegen.Location, svc *service.Data) string {
	if loc != nil {
		return loc.PackageName()
	}
	if svc != nil {
		return svc.PkgName
	}
	return ""
}

func importPath(loc *codegen.Location, svc *service.Data, genpkg string) string {
	if loc != nil {
		if loc.RelImportPath != "" {
			return path.Join(genpkg, loc.RelImportPath)
		}
		return ""
	}
	if svc == nil {
		return ""
	}
	return path.Join(genpkg, svc.PathName)
}

func (g *generator) addTypeImports(att *expr.AttributeExpr) {
	for _, im := range gatherAttributeImports(g.genpkg, att) {
		if im == nil || im.Path == "" {
			continue
		}
		g.typeImports[im.Path] = im
	}
}

func gatherAttributeImports(genpkg string, att *expr.AttributeExpr) []*codegen.ImportSpec {
	if att == nil {
		return nil
	}
	var specs []*codegen.ImportSpec
	switch dt := att.Type.(type) {
	case expr.UserType:
		if loc := codegen.UserTypeLocation(dt); loc != nil && loc.RelImportPath != "" {
			specs = append(specs, &codegen.ImportSpec{Path: path.Join(genpkg, loc.RelImportPath)})
		}
		specs = append(specs, gatherAttributeImports(genpkg, dt.Attribute())...)
	case *expr.Array:
		specs = append(specs, gatherAttributeImports(genpkg, dt.ElemType)...)
	case *expr.Map:
		specs = append(specs, gatherAttributeImports(genpkg, dt.KeyType)...)
		specs = append(specs, gatherAttributeImports(genpkg, dt.ElemType)...)
	case expr.CompositeExpr:
		specs = append(specs, gatherAttributeImports(genpkg, dt.Attribute())...)
	}
	specs = append(specs, codegen.GetMetaTypeImports(att)...)
	return specs
}

func schemaForAttribute(method *expr.MethodExpr, att *expr.AttributeExpr) ([]byte, error) {
	if method == nil || att == nil || att.Type == nil || att.Type == expr.Empty {
		return nil, nil
	}
	prev := openapi.Definitions
	openapi.Definitions = make(map[string]*openapi.Schema)
	defer func() { openapi.Definitions = prev }()
	schema := openapi.AttributeTypeSchema(expr.Root.API, att)
	if schema == nil {
		return nil, nil
	}
	if len(openapi.Definitions) > 0 {
		schema.Definitions = openapi.Definitions
	}
	return schema.JSON()
}

func buildValidation(att *expr.AttributeExpr, pkg string, scope *codegen.NameScope) string {
	if att == nil || att.Type == nil || att.Type == expr.Empty {
		return ""
	}
	ctx := codegen.NewAttributeContext(false, false, true, pkg, scope)
	return strings.TrimSpace(codegen.ValidationCode(att, nil, ctx, true, expr.IsAlias(att.Type), false, "body"))
}

func formatSchema(schema []byte) string {
	if len(schema) == 0 {
		return ""
	}
	content := string(schema)
	return "[]byte(`\n" + content + "\n`)"
}

func lowerCamel(s string) string {
	return codegen.Goify(s, false)
}

func finalizeTypeInfo(info *typeInfo) {
	info.Pointer = strings.HasPrefix(info.FullRef, "*")
	if info.Pointer {
		info.CheckNil = true
		info.MarshalArg = "v"
		info.UnmarshalArg = "&v"
	} else {
		info.MarshalArg = "v"
		info.UnmarshalArg = "v"
	}
	if info.Validation != "" {
		info.ValidationSrc = strings.Split(info.Validation, "\n")
	}
}

func isDerived(tool *toolsexpr.ToolExpr) bool {
	return tool != nil && tool.Derived
}

func executeToolDSL(tool *toolsexpr.ToolExpr) {
	if tool == nil || tool.DSLFunc == nil {
		return
	}
	tool.DSLFunc()
	tool.DSLFunc = nil
}
