package docs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
	"goa.design/goa/v3/http/codegen/openapi"
	plugexpr "goa.design/plugins/v3/docs/expr"
)

// Register registers the docs plugin generator exactly once.
var registerOnce sync.Once

func Register() {
	registerOnce.Do(func() {
		codegen.RegisterPlugin("docs", "gen", nil, Generate)
	})
}

// Backward compatibility: register on package import.
func init() { Register() }

// Generate produces the documentation JSON file.
func Generate(_ string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	// First, build a complete map of all definitions from all roots.
	// This ensures that cross-package type references can be resolved.
	allDefs := make(map[string]*openapi.Schema)
	for _, root := range roots {
		if r, ok := root.(*expr.RootExpr); ok {
			// Create a temporary, isolated context for each root to avoid global state pollution.
			prev := openapi.Definitions
			openapi.Definitions = make(map[string]*openapi.Schema)

			for _, tpe := range r.Types {
				if ut, ok := tpe.(*expr.UserTypeExpr); ok {
					openapi.GenerateTypeDefinition(r.API, ut)
				}
			}
			for _, rt := range r.ResultTypes {
				openapi.GenerateResultTypeDefinition(r.API, rt, expr.DefaultView)
			}

			// Merge the generated definitions into the global map.
			for n, s := range openapi.Definitions {
				if _, exists := allDefs[n]; !exists {
					allDefs[n] = dupSchema(s)
				}
			}

			// Restore the original global definitions to maintain isolation.
			openapi.Definitions = prev
		}
	}

	for _, root := range roots {
		if r, ok := root.(*expr.RootExpr); ok {
			files = append(files, docsFile(r, allDefs))
		}
	}
	return files, nil
}

func docsFile(r *expr.RootExpr, allDefs map[string]*openapi.Schema) *codegen.File {
	docs := &data{
		API:      apiDocs(r.API),
		Services: servicesDocs(r),
	}

	// Default behavior: use global OpenAPI definitions to preserve ordering and
	// compatibility with existing golden tests.
	defs := allDefs

	// If either option is enabled, build a local definition map for this root
	// and apply transforms/inlining as needed, isolating from global state.
	if plugexpr.Root.UseJSONTags {
		// Re-scope the definitions to only those present in the current root,
		// but use the globally-aware `allDefs` for lookups during transforms.
		local := make(map[string]*openapi.Schema)
		prev := openapi.Definitions
		openapi.Definitions = make(map[string]*openapi.Schema)
		for _, tpe := range r.Types {
			if ut, ok := tpe.(*expr.UserTypeExpr); ok {
				openapi.GenerateTypeDefinition(r.API, ut)
			}
		}
		for _, rt := range r.ResultTypes {
			openapi.GenerateResultTypeDefinition(r.API, rt, expr.DefaultView)
		}
		for n := range openapi.Definitions {
			if def, ok := allDefs[n]; ok {
				local[n] = dupSchema(def)
			}
		}
		openapi.Definitions = prev

		// Apply JSON tag transforms if requested
		if plugexpr.Root.UseJSONTags {
			local = transformDefinitionsWithJSONTagsHybrid(r, local, nil)
		}
		defs = local
	} else {
		// When not transforming, avoid leaking the synthetic Empty type produced by
		// Goa internals in some scenarios by filtering it out, but only when present.
		if _, hasEmpty := defs["Empty"]; hasEmpty {
			filtered := make(map[string]*openapi.Schema, len(defs))
			for n, s := range defs {
				if n == "Empty" {
					continue
				}
				filtered[n] = s
			}
			defs = filtered
		}
	}

	docs.Definitions = defs

	// Inline $refs if requested via DSL flag.
	if plugexpr.Root.InlineRefs {
		// When inlining, use the complete set of definitions from all roots
		// to ensure cross-package references can be resolved.
		inliningDefs := allDefs
		if plugexpr.Root.UseJSONTags {
			inliningDefs = transformDefinitionsWithJSONTagsHybrid(r, allDefs, nil)
		}

		// Inline inside service payloads/results/errors.
		for _, svc := range docs.Services {
			for _, m := range svc.Methods {
				if m.Payload != nil && m.Payload.Type != nil {
					inlineRefsInSchema(m.Payload.Type, inliningDefs, make(map[string]bool))
				}
				if m.StreamingPayload != nil && m.StreamingPayload.Type != nil {
					inlineRefsInSchema(m.StreamingPayload.Type, inliningDefs, make(map[string]bool))
				}
				if m.Result != nil && m.Result.Type != nil {
					inlineRefsInSchema(m.Result.Type, inliningDefs, make(map[string]bool))
				}
				if m.StreamingResult != nil && m.StreamingResult.Type != nil {
					inlineRefsInSchema(m.StreamingResult.Type, inliningDefs, make(map[string]bool))
				}
				for _, e := range m.Errors {
					if e != nil && e.Type != nil {
						inlineRefsInSchema(e.Type, inliningDefs, make(map[string]bool))
					}
				}
			}
		}
		// Inline inside definitions themselves (properties that refer to other defs).
		for _, def := range defs {
			inlineRefsInSchema(def, inliningDefs, make(map[string]bool))
		}
	}

	jsonPath := filepath.Join(codegen.Gendir, "docs.json")
	if _, err := os.Stat(jsonPath); !os.IsNotExist(err) {
		// Goa does not delete files in the top-level gen folder.
		// https://github.com/goadesign/goa/pull/2194
		// The plugin must delete docs.json so that the generator does not append
		// to any existing docs.json.
		if err := os.Remove(jsonPath); err != nil {
			panic(err)
		}
	}
	jsonSection := &codegen.SectionTemplate{
		Name:    "docs",
		FuncMap: template.FuncMap{"toJSON": toJSON},
		Source:  "{{ toJSON .}}",
		Data:    docs,
	}
	return &codegen.File{
		Path:             jsonPath,
		SectionTemplates: []*codegen.SectionTemplate{jsonSection},
	}
}

// dupSchema creates a safe deep copy of the given schema, ensuring maps are initialized.
func dupSchema(s *openapi.Schema) *openapi.Schema {
	if s == nil {
		return nil
	}
	js := openapi.Schema{
		ID:                   s.ID,
		Description:          s.Description,
		Schema:               s.Schema,
		Type:                 s.Type,
		DefaultValue:         s.DefaultValue,
		Title:                s.Title,
		Media:                s.Media,
		ReadOnly:             s.ReadOnly,
		PathStart:            s.PathStart,
		Links:                s.Links,
		Ref:                  s.Ref,
		Enum:                 s.Enum,
		Format:               s.Format,
		Pattern:              s.Pattern,
		Minimum:              s.Minimum,
		Maximum:              s.Maximum,
		MinLength:            s.MinLength,
		MaxLength:            s.MaxLength,
		MinItems:             s.MinItems,
		MaxItems:             s.MaxItems,
		Required:             s.Required,
		AdditionalProperties: s.AdditionalProperties,
		Properties:           make(map[string]*openapi.Schema, len(s.Properties)),
		Definitions:          make(map[string]*openapi.Schema, len(s.Definitions)),
		AnyOf:                nil,
		Example:              s.Example,
		Extensions:           s.Extensions,
	}
	for n, p := range s.Properties {
		js.Properties[n] = dupSchema(p)
	}
	if s.Items != nil {
		js.Items = dupSchema(s.Items)
	}
	for n, d := range s.Definitions {
		js.Definitions[n] = dupSchema(d)
	}
	if len(s.AnyOf) > 0 {
		js.AnyOf = make([]*openapi.Schema, len(s.AnyOf))
		for i := range s.AnyOf {
			js.AnyOf[i] = dupSchema(s.AnyOf[i])
		}
	}
	return &js
}

func apiDocs(api *expr.APIExpr) *apiData {
	data := &apiData{
		Name:        api.Name,
		Title:       api.Title,
		Description: api.Description,
		Version:     api.Version,
		Terms:       api.TermsOfService,
	}
	if len(api.Servers) > 0 {
		data.Servers = make(map[string]*serverData, len(api.Servers))
		for _, s := range api.Servers {
			data.Servers[s.Name] = generateServer(s)
		}
	}
	if c := api.Contact; c != nil {
		data.Contact = &contactData{c.Name, c.Email, c.URL}
	}
	if l := api.License; l != nil {
		data.License = &licenseData{l.Name, l.URL}
	}
	if d := api.Docs; d != nil {
		data.Docs = &docsData{d.Description, d.URL}
	}
	data.Requirements = make([]*requirementData, len(api.Requirements))
	for i, req := range api.Requirements {
		data.Requirements[i] = generateRequirement(req)
	}

	return data
}

func servicesDocs(r *expr.RootExpr) map[string]*serviceData {
	svcs := make(map[string]*serviceData, len(r.Services))
	var nameScope = codegen.NewNameScope()

	for _, svc := range r.Services {
		n := svc.Name
		svcs[n] = &serviceData{
			Name:        n,
			Description: svc.Description,
		}

		svcs[n].Methods = make(map[string]*methodData, len(svc.Methods))
		for _, meth := range svc.Methods {
			svcs[n].Methods[meth.Name] = generateMethod(r.API, meth, nameScope)
		}

		svcs[n].Requirements = make([]*requirementData, len(svc.Requirements))
		for i, req := range svc.Requirements {
			svcs[n].Requirements[i] = generateRequirement(req)
		}
	}
	return svcs
}

func generateServer(s *expr.ServerExpr) *serverData {
	data := &serverData{
		Name:        s.Name,
		Description: s.Description,
		Services:    s.Services,
	}
	if len(s.Hosts) > 0 {
		data.Hosts = make(map[string]*hostData)
		for _, h := range s.Hosts {
			data.Hosts[h.Name] = &hostData{
				Name:        h.Name,
				ServerName:  h.ServerName,
				Description: h.Description,
			}
			if len(h.URIs) > 0 {
				data.Hosts[h.Name].URIs = make([]string, len(h.URIs))
				for i, u := range h.URIs {
					data.Hosts[h.Name].URIs[i] = string(u)
				}
			}
			if o := expr.AsObject(h.Variables.Type); o != nil {
				data.Hosts[h.Name].Variables = make([]*variableData, len(*o))
				for i, na := range *o {
					var def string
					if na.Attribute.DefaultValue != nil {
						def = fmt.Sprintf("%v", na.Attribute.DefaultValue)
					}
					var e []string
					if na.Attribute.Validation != nil && len(na.Attribute.Validation.Values) > 0 {
						e = make([]string, len(na.Attribute.Validation.Values))
						for j, v := range na.Attribute.Validation.Values {
							e[j] = fmt.Sprintf("%v", v)
						}
					}
					data.Hosts[h.Name].Variables[i] = &variableData{na.Name, def, e}
				}
			}
		}
	}
	return data
}

func generateRequirement(req *expr.SecurityExpr) *requirementData {
	r := &requirementData{Scopes: req.Scopes}
	if len(req.Schemes) > 0 {
		r.Schemes = make([]*schemeData, len(req.Schemes))
		for i, sch := range req.Schemes {
			r.Schemes[i] = &schemeData{
				Type:        sch.Type(),
				Description: sch.Description,
				Name:        sch.Name,
				In:          sch.In,
				Scheme:      sch.SchemeName,
			}
			if len(sch.Flows) > 0 {
				r.Schemes[i].Flows = make([]*flowData, len(sch.Flows))
				for j, f := range sch.Flows {
					r.Schemes[i].Flows[j] = &flowData{f.Type(), f.AuthorizationURL, f.TokenURL, f.RefreshURL}
				}
			}
		}
	}
	return r
}

func generateMethod(api *expr.APIExpr, meth *expr.MethodExpr, scope *codegen.NameScope) *methodData {
	m := &methodData{
		Name:             meth.Name,
		Description:      meth.Description,
		Payload:          generatePayload(api, meth.Payload, scope),
		StreamingPayload: generatePayload(api, meth.StreamingPayload, scope),
	}
	if meth.Stream == expr.BidirectionalStreamKind || meth.Stream == expr.ServerStreamKind {
		m.StreamingResult = generatePayload(api, meth.Result, scope)
	} else {
		m.Result = generatePayload(api, meth.Result, scope)
	}
	m.Errors = make(map[string]*errorData, len(meth.Errors))
	for _, er := range meth.Errors {
		m.Errors[er.Name] = generateError(api, er)
	}
	m.Requirements = make([]*requirementData, len(meth.Requirements))
	for i, req := range meth.Requirements {
		m.Requirements[i] = generateRequirement(req)
	}
	return m
}

func generatePayload(api *expr.APIExpr, att *expr.AttributeExpr, nameScope *codegen.NameScope) *payloadData {
	// since the definitions section is global to the API, we need to ensure uniqueness of TypeName
	if ut, ok := att.Type.(*expr.UserTypeExpr); ok {
		if ut == expr.Empty {
			return nil
		}
		ut.TypeName = nameScope.Unique(ut.TypeName)
	}

	schema := openapi.AttributeTypeSchema(api, att)
	ex := att.Example(api.ExampleGenerator)
	if plugexpr.Root.UseJSONTags {
		// avoid mutating shared schema nodes
		schema = dupSchema(schema)
		applyJSONTagsToSchema(att, schema)
		ex = transformExampleWithJSONTags(att, ex)
	}
	return &payloadData{
		Type:    schema,
		Example: ex,
	}
}

func generateError(api *expr.APIExpr, er *expr.ErrorExpr) *errorData {
	_, temporary := er.Meta["goa:error:temporary"]
	_, timeout := er.Meta["goa:error:timeout"]
	_, fault := er.Meta["goa:error:fault"]
	sch := openapi.AttributeTypeSchema(api, er.AttributeExpr)
	if plugexpr.Root.UseJSONTags {
		sch = dupSchema(sch)
		applyJSONTagsToSchema(er.AttributeExpr, sch)
	}
	return &errorData{
		Name:        er.Name,
		Description: er.Description,
		Type:        sch,
		Temporary:   temporary,
		Timeout:     timeout,
		Fault:       fault,
	}
}

func toJSON(d interface{}) string {
	b, err := json.Marshal(d)
	if err != nil {
		panic("openapi: " + err.Error()) // bug
	}
	return string(b)
}

// inlineRefsInSchema replaces $ref with a deep copy of the referenced
// definition where possible. It recurses through properties/items/anyOf.
// A small visited set prevents infinite recursion on cycles.
func inlineRefsInSchema(s *openapi.Schema, defs map[string]*openapi.Schema, visiting map[string]bool) {
	if s == nil {
		return
	}
	// Inline reference if it points to local definitions.
	if after, ok := strings.CutPrefix(s.Ref, "#/definitions/"); ok {
		name := after
		if !visiting[name] {
			if def, ok := defs[name]; ok {
				visiting[name] = true
				dup := dupSchema(def)
				// Recurse into the dup first to inline nested refs.
				inlineRefsInSchema(dup, defs, visiting)
				// Replace s with contents of dup (shallow copy fields and maps)
				*s = *dup
				visiting[name] = false
			}
		}
	}
	// Recurse
	if s.Items != nil {
		inlineRefsInSchema(s.Items, defs, visiting)
	}
	if s.Properties != nil {
		for _, p := range s.Properties {
			inlineRefsInSchema(p, defs, visiting)
		}
	}
	if s.AdditionalProperties != nil {
		if asp, ok := s.AdditionalProperties.(*openapi.Schema); ok {
			inlineRefsInSchema(asp, defs, visiting)
		}
	}
	if s.AnyOf != nil {
		for _, a := range s.AnyOf {
			inlineRefsInSchema(a, defs, visiting)
		}
	}
}

// transformDefinitionsWithJSONTagsHybrid tries attrIndex first; if not found, falls back to Root.UserType lookup.
func transformDefinitionsWithJSONTagsHybrid(r *expr.RootExpr, defs map[string]*openapi.Schema, attrIndex map[string]*expr.AttributeExpr) map[string]*openapi.Schema {
	if len(defs) == 0 {
		return defs
	}
	out := make(map[string]*openapi.Schema, len(defs))
	for name, sch := range defs {
		dup := dupSchema(sch)
		if att, ok := attrIndex[name]; ok && att != nil {
			applyJSONTagsToSchema(att, dup)
		} else if ut := r.UserType(name); ut != nil {
			applyJSONTagsToSchema(ut.Attribute(), dup)
		}
		out[name] = dup
	}
	return out
}

// applyJSONTagsToSchema mutates s to use JSON tag names from Meta on the given
// attribute and its descendants. It preserves required field semantics and
// updates examples when present.
func applyJSONTagsToSchema(att *expr.AttributeExpr, s *openapi.Schema) {
	if att == nil || s == nil {
		return
	}
	// If this schema is a ref, we expect the referenced definition to be transformed separately.
	if s.Ref != "" {
		return
	}

	// Unwrap user/result types to their underlying attributes before processing.
	for {
		switch t := att.Type.(type) {
		case *expr.ResultTypeExpr:
			att = t.AttributeExpr
			continue
		case expr.UserType:
			att = t.Attribute()
			continue
		}
		break
	}

	// Recurse into composite types first so nested structures are handled.
	switch t := att.Type.(type) {
	case *expr.Array:
		if s.Items != nil {
			applyJSONTagsToSchema(t.ElemType, s.Items)
		}
	case *expr.Map:
		if as, ok := s.AdditionalProperties.(*openapi.Schema); ok {
			applyJSONTagsToSchema(t.ElemType, as)
		}
	case *expr.Union:
		for i, v := range t.Values {
			if i < len(s.AnyOf) {
				applyJSONTagsToSchema(v.Attribute, s.AnyOf[i])
			}
		}
	}

	// Handle object property renaming and example/required updates.
	if obj := expr.AsObject(att.Type); obj != nil && s.Properties != nil {
		// Build new properties map using JSON tag names. Use walkAttribute so bases/references are included.
		newProps := make(map[string]*openapi.Schema, len(s.Properties))
		nameMap := make(map[string]string, len(*obj))
		_ = walkAttribute(att, func(oldName string, child *expr.AttributeExpr) error {
			jsonName := jsonTagName(child)
			if jsonName == "" || jsonName == "-" {
				jsonName = oldName
			}
			// Find property schema by old key (fallback to jsonName if already renamed).
			prop := s.Properties[oldName]
			if prop == nil {
				prop = s.Properties[jsonName]
			}
			if prop != nil {
				applyJSONTagsToSchema(child, prop)
				if _, exists := newProps[jsonName]; !exists {
					newProps[jsonName] = prop
				}
			}
			nameMap[oldName] = jsonName
			return nil
		})
		s.Properties = newProps

		if len(s.Required) > 0 {
			newReq := make([]string, 0, len(s.Required))
			for _, rn := range s.Required {
				if jn, ok := nameMap[rn]; ok {
					if _, exists := newProps[jn]; exists {
						newReq = append(newReq, jn)
					}
				} else if _, exists := newProps[rn]; exists {
					newReq = append(newReq, rn)
				}
			}
			s.Required = newReq
		}

		if s.Example != nil {
			s.Example = transformExampleWithJSONTags(att, s.Example)
		}
	}
}

// walkAttribute iterates over the given attribute, its bases and references (if any),
// calling the iterator for each field of object types. This mirrors goa's internal
// expr.walkAttribute to ensure bases and references are considered when transforming.
func walkAttribute(att *expr.AttributeExpr, it func(name string, a *expr.AttributeExpr) error) error { //nolint:cyclop
	switch dt := att.Type.(type) {
	case expr.UserType:
		if err := walkAttribute(dt.Attribute(), it); err != nil {
			return err
		}
	case *expr.Object:
		for _, nat := range *dt {
			if err := it(nat.Name, nat.Attribute); err != nil {
				return err
			}
		}
	}
	for _, b := range att.Bases {
		if err := walkAttribute(&expr.AttributeExpr{Type: b}, it); err != nil {
			return err
		}
	}
	for _, r := range att.References {
		if err := walkAttribute(&expr.AttributeExpr{Type: r}, it); err != nil {
			return err
		}
	}
	return nil
}

// jsonTagName extracts the JSON tag field name from the attribute Meta if set.
// It supports values like "name,omitempty" and returns "name".
func jsonTagName(att *expr.AttributeExpr) string {
	if att == nil || att.Meta == nil {
		return ""
	}
	if vals, ok := att.Meta["struct:tag:json"]; ok && len(vals) > 0 {
		tag := vals[len(vals)-1]
		if idx := strings.Index(tag, ","); idx >= 0 {
			tag = tag[:idx]
		}
		return tag
	}
	return ""
}

// transformExampleWithJSONTags rewrites example map keys to match JSON tags from
// Meta. It recurses through objects and arrays. For user and result types it
// recurses into the underlying attribute.
func transformExampleWithJSONTags(att *expr.AttributeExpr, ex any) any {
	if att == nil || ex == nil {
		return ex
	}
	switch t := att.Type.(type) {
	case *expr.ResultTypeExpr:
		return transformExampleWithJSONTags(t.AttributeExpr, ex)
	case expr.UserType:
		return transformExampleWithJSONTags(t.Attribute(), ex)
	case *expr.Object:
		m, ok := ex.(map[string]any)
		if !ok {
			return ex
		}
		res := make(map[string]any, len(m))
		for _, nat := range *t {
			oldName := nat.Name
			jsonName := jsonTagName(nat.Attribute)
			if jsonName == "" || jsonName == "-" {
				jsonName = oldName
			}
			if val, ok := m[oldName]; ok {
				res[jsonName] = transformExampleWithJSONTags(nat.Attribute, val)
			}
		}
		return res
	case *expr.Array:
		if arr, ok := ex.([]any); ok {
			out := make([]any, len(arr))
			for i := range arr {
				out[i] = transformExampleWithJSONTags(t.ElemType, arr[i])
			}
			return out
		}
		return ex
	case *expr.Map:
		// Only transform element values.
		switch m := ex.(type) {
		case map[string]any:
			out := make(map[string]any, len(m))
			for k, v := range m {
				out[k] = transformExampleWithJSONTags(t.ElemType, v)
			}
			return out
		case map[any]any:
			out := make(map[any]any, len(m))
			for k, v := range m {
				out[k] = transformExampleWithJSONTags(t.ElemType, v)
			}
			return out
		default:
			return ex
		}
	case *expr.Union:
		// Attempt best-effort transform by applying first variant.
		if len(t.Values) > 0 {
			return transformExampleWithJSONTags(t.Values[0].Attribute, ex)
		}
		return ex
	default:
		return ex
	}
}
