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
	goaexpr "goa.design/goa/v3/expr"
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
	// Always emit a single aggregated docs.json at the top level.
	jsonPath := filepath.Join(codegen.Gendir, "docs.json")
	if _, err := os.Stat(jsonPath); err == nil {
		if err := os.Remove(jsonPath); err != nil {
			return nil, fmt.Errorf("remove stale docs.json: %w", err)
		}
	}

	var agg *data

	for _, root := range roots {
		r, ok := root.(*goaexpr.RootExpr)
		if !ok {
			continue
		}
		d := docsDataForRoot(r)
		if d == nil {
			continue
		}
		if agg == nil {
			// First root initializes the aggregate.
			agg = &data{
				API:         d.API,
				Services:    make(map[string]*serviceData, len(d.Services)),
				Definitions: make(map[string]*openapi.Schema, len(d.Definitions)),
			}
		}
		// Merge services (merge methods if service reappears).
		for sn, svc := range d.Services {
			if existing, ok := agg.Services[sn]; ok && existing != nil && svc != nil {
				if existing.Methods == nil {
					existing.Methods = make(map[string]*methodData)
				}
				for mn, md := range svc.Methods {
					existing.Methods[mn] = md
				}
				// Preserve description/requirements if present on either.
				if existing.Description == "" {
					existing.Description = svc.Description
				}
				if len(existing.Requirements) == 0 && len(svc.Requirements) > 0 {
					existing.Requirements = svc.Requirements
				}
			} else {
				agg.Services[sn] = svc
			}
		}
		// Merge definitions (last-wins for same key).
		for dn, sch := range d.Definitions {
			agg.Definitions[dn] = sch
		}
	}

	if agg == nil {
		// No valid roots; nothing to do.
		return files, nil
	}

	jsonSection := &codegen.SectionTemplate{
		Name:    "docs",
		FuncMap: template.FuncMap{"mustJSON": mustJSON},
		Source:  "{{ mustJSON .}}",
		Data:    agg,
	}
	files = append(files, &codegen.File{
		Path:             jsonPath,
		SectionTemplates: []*codegen.SectionTemplate{jsonSection},
	})
	return files, nil
}

// docsDataForRoot builds the docs data for a single root while scoping
// openapi.Definitions to that root only. The global definitions map is restored
// before returning.
func docsDataForRoot(r *goaexpr.RootExpr) *data {
	prev := openapi.Definitions
	openapi.Definitions = make(map[string]*openapi.Schema)
	defer func() { openapi.Definitions = prev }()

	f := docsFile(r)
	if len(f.SectionTemplates) == 0 || f.SectionTemplates[0] == nil {
		return nil
	}
	d, _ := f.SectionTemplates[0].Data.(*data)
	return d
}

func docsFile(r *goaexpr.RootExpr) *codegen.File {
	// Per-run generation state
	state := newGenState()
	docs := &data{
		API:      apiDocs(r.API),
		Services: servicesDocs(r, state),
	}

	// Default behavior: use root-local OpenAPI definitions produced during
	// payload/result schema generation to preserve golden stability.
	defs := make(map[string]*openapi.Schema, len(openapi.Definitions))
	for n, d := range openapi.Definitions {
		defs[n] = dupSchema(d)
	}

	// Force-emit definitions for explicitly included user types (and deps).
	if its := plugexpr.Root.IncludedTypes; len(its) > 0 {
		forced := forcedDefinitions(r.API, its, state)
		for n, s := range forced {
			defs[n] = s
		}
	}

	// Apply JSON tag transforms if requested.
	if plugexpr.Root.UseJSONTags {
		defs = transformDefinitionsWithJSONTagsHybrid(r, defs)
	}

	// Filter synthetic Empty type if present.
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

	// Inline $ref occurrences if requested.
	if plugexpr.Root.InlineRefs {
		inlineAllServiceSchemas(docs, defs)
		// Also inline inside definitions so nested refs are expanded there too.
		stack := make(map[string]bool)
		for _, sch := range defs {
			inlineRefs(sch, defs, stack)
		}
	}

	docs.Definitions = defs

	jsonPath := filepath.Join(codegen.Gendir, "docs.json")
	jsonSection := &codegen.SectionTemplate{
		Name:    "docs",
		FuncMap: template.FuncMap{"mustJSON": mustJSON},
		Source:  "{{ mustJSON .}}",
		Data:    docs,
	}
	return &codegen.File{
		Path:             jsonPath,
		SectionTemplates: []*codegen.SectionTemplate{jsonSection},
	}
}

// forcedDefinitions builds JSON Schemas for the given user types and their
// transitive dependencies and returns a name->schema map suitable for inclusion
// in the top-level definitions object.
func forcedDefinitions(api *goaexpr.APIExpr, uts []goaexpr.UserType, state *genState) map[string]*openapi.Schema { //nolint:cyclop
	seen := make(map[string]*goaexpr.UserTypeExpr)

	var addUT func(goaexpr.UserType)
	addUT = func(ut goaexpr.UserType) {
		if ut == nil {
			return
		}
		if ute, ok := ut.(*goaexpr.UserTypeExpr); ok {
			name := ute.TypeName
			if _, ok := seen[name]; ok {
				return
			}
			seen[name] = ute
			walkDeps(ute.AttributeExpr, addUT)
		}
	}

	for _, ut := range uts {
		addUT(ut)
	}

	out := make(map[string]*openapi.Schema, len(seen))
	for name, ute := range seen {
		sch := schemaForAttribute(api, ute.AttributeExpr, state)
		out[name] = dupSchema(sch)
	}
	return out
}

// walkDeps recursively visits attribute data types and invokes add on any
// discovered user types (including result types). It follows bases and
// references to mirror behavior elsewhere in the plugin.
func walkDeps(att *goaexpr.AttributeExpr, add func(goaexpr.UserType)) {
	if att == nil || att.Type == nil {
		return
	}
	switch t := att.Type.(type) {
	case goaexpr.UserType:
		add(t)
		walkDeps(t.Attribute(), add)
	case *goaexpr.Array:
		walkDeps(t.ElemType, add)
	case *goaexpr.Map:
		walkDeps(t.KeyType, add)
		walkDeps(t.ElemType, add)
	case *goaexpr.Object:
		for _, nat := range *t {
			walkDeps(nat.Attribute, add)
		}
	case *goaexpr.ResultTypeExpr:
		add(t)
		walkDeps(t.AttributeExpr, add)
	case *goaexpr.Union:
		for _, v := range t.Values {
			walkDeps(v.Attribute, add)
		}
	}
	for _, b := range att.Bases {
		walkDeps(&goaexpr.AttributeExpr{Type: b}, add)
	}
	for _, r := range att.References {
		walkDeps(&goaexpr.AttributeExpr{Type: r}, add)
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

func apiDocs(api *goaexpr.APIExpr) *apiData {
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

// mustGenerateService returns false if the service (or its HTTP service) is
// marked with Meta("openapi:generate", "false") or legacy
// Meta("swagger:generate", "false").
func mustGenerateService(r *goaexpr.RootExpr, svc *goaexpr.ServiceExpr) bool {
	// Explicit docs DSL can disable docs.json without impacting OpenAPI
	if vals, ok := svc.Meta["docs:generate"]; ok && len(vals) > 0 && strings.EqualFold(vals[len(vals)-1], "false") {
		return false
	}
	if !openapi.MustGenerate(svc.Meta) {
		return false
	}
	if r != nil && r.API != nil && r.API.HTTP != nil {
		for _, hs := range r.API.HTTP.Services {
			if hs != nil && hs.ServiceExpr == svc {
				if !openapi.MustGenerate(hs.Meta) {
					return false
				}
				break
			}
		}
	}
	return true
}

// mustGenerateMethod returns false if the method is marked with
// Meta("openapi:generate", "false") or legacy Meta("swagger:generate",
// "false"). If HTTP endpoints exist for the method, then the method is only
// considered generatable if at least one endpoint is also generatable (mirrors
// Goa behavior where endpoint-level meta gates individual operations).
func mustGenerateMethod(r *goaexpr.RootExpr, meth *goaexpr.MethodExpr) bool { //nolint:cyclop
	if vals, ok := meth.Meta["docs:generate"]; ok && len(vals) > 0 && strings.EqualFold(vals[len(vals)-1], "false") {
		return false
	}
	if !openapi.MustGenerate(meth.Meta) {
		return false
	}
	// If there is HTTP configuration for this method, require at least one
	// endpoint to be marked generatable; otherwise keep default true.
	if r != nil && r.API != nil && r.API.HTTP != nil && meth != nil && meth.Service != nil {
		for _, hs := range r.API.HTTP.Services {
			if hs == nil || hs.ServiceExpr != meth.Service {
				continue
			}
			// Found the HTTP service matching the method's service.
			hasEndpoint := false
			for _, e := range hs.HTTPEndpoints {
				if e == nil || e.MethodExpr != meth {
					continue
				}
				hasEndpoint = true
				if openapi.MustGenerate(e.Meta) {
					// At least one endpoint is generatable - keep the method.
					return true
				}
			}
			// If endpoints exist but none are generatable, skip the method.
			if hasEndpoint {
				return false
			}
			break
		}
	}
	return true
}

func servicesDocs(r *goaexpr.RootExpr, state *genState) map[string]*serviceData {
	svcs := make(map[string]*serviceData, len(r.Services))

	for _, svc := range r.Services {
		if !mustGenerateService(r, svc) {
			continue
		}
		n := svc.Name
		svcs[n] = &serviceData{
			Name:        n,
			Description: svc.Description,
		}

		svcs[n].Methods = make(map[string]*methodData, len(svc.Methods))
		for _, meth := range svc.Methods {
			if !mustGenerateMethod(r, meth) {
				continue
			}
			svcs[n].Methods[meth.Name] = generateMethod(r.API, meth, state)
		}

		svcs[n].Requirements = make([]*requirementData, len(svc.Requirements))
		for i, req := range svc.Requirements {
			svcs[n].Requirements[i] = generateRequirement(req)
		}
	}
	return svcs
}

func generateServer(s *goaexpr.ServerExpr) *serverData {
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
			if o := goaexpr.AsObject(h.Variables.Type); o != nil {
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

func generateRequirement(req *goaexpr.SecurityExpr) *requirementData {
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

func generateMethod(api *goaexpr.APIExpr, meth *goaexpr.MethodExpr, state *genState) *methodData {
	m := &methodData{
		Name:             meth.Name,
		Description:      meth.Description,
		Payload:          generatePayload(api, meth.Payload, state),
		StreamingPayload: generatePayload(api, meth.StreamingPayload, state),
	}
	if meth.Stream == goaexpr.BidirectionalStreamKind || meth.Stream == goaexpr.ServerStreamKind {
		m.StreamingResult = generatePayload(api, meth.Result, state)
	} else {
		m.Result = generatePayload(api, meth.Result, state)
	}
	m.Errors = make(map[string]*errorData, len(meth.Errors))
	for _, er := range meth.Errors {
		m.Errors[er.Name] = generateError(api, er, state)
	}
	m.Requirements = make([]*requirementData, len(meth.Requirements))
	for i, req := range meth.Requirements {
		m.Requirements[i] = generateRequirement(req)
	}
	return m
}

func generatePayload(api *goaexpr.APIExpr, att *goaexpr.AttributeExpr, state *genState) *payloadData {
	// Do not generate payload for Empty
	if ut, ok := att.Type.(*goaexpr.UserTypeExpr); ok && ut == goaexpr.Empty {
		return nil
	}

	schema := schemaForAttribute(api, att, state)
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

func generateError(api *goaexpr.APIExpr, er *goaexpr.ErrorExpr, state *genState) *errorData {
	_, temporary := er.Meta["goa:error:temporary"]
	_, timeout := er.Meta["goa:error:timeout"]
	_, fault := er.Meta["goa:error:fault"]
	sch := schemaForAttribute(api, er.AttributeExpr, state)
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

func mustJSON(d interface{}) string {
	b, err := json.Marshal(d)
	if err != nil {
		panic("docs: " + err.Error()) // bug
	}
	return string(b)
}

// inline $ref logic removed

// transformDefinitionsWithJSONTagsHybrid tries Root.UserType lookup.
func transformDefinitionsWithJSONTagsHybrid(r *goaexpr.RootExpr, defs map[string]*openapi.Schema) map[string]*openapi.Schema {
	if len(defs) == 0 {
		return defs
	}
	out := make(map[string]*openapi.Schema, len(defs))
	for name, sch := range defs {
		dup := dupSchema(sch)
		if ut := r.UserType(name); ut != nil {
			applyJSONTagsToSchema(ut.Attribute(), dup)
		}
		out[name] = dup
	}
	return out
}

// applyJSONTagsToSchema mutates s to use JSON tag names from Meta on the given
// attribute and its descendants. It preserves required field semantics and
// updates examples when present.
func applyJSONTagsToSchema(att *goaexpr.AttributeExpr, s *openapi.Schema) {
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
		case *goaexpr.ResultTypeExpr:
			att = t.AttributeExpr
			continue
		case goaexpr.UserType:
			att = t.Attribute()
			continue
		}
		break
	}

	// Recurse into composite types first so nested structures are handled.
	switch t := att.Type.(type) {
	case *goaexpr.Array:
		if s.Items != nil {
			applyJSONTagsToSchema(t.ElemType, s.Items)
		}
	case *goaexpr.Map:
		if as, ok := s.AdditionalProperties.(*openapi.Schema); ok {
			applyJSONTagsToSchema(t.ElemType, as)
		}
	case *goaexpr.Union:
		for i, v := range t.Values {
			if i < len(s.AnyOf) {
				applyJSONTagsToSchema(v.Attribute, s.AnyOf[i])
			}
		}
	}

	// Handle object property renaming and example/required updates.
	// If multiple fields resolve to the same JSON tag, the first seen wins.
	if obj := goaexpr.AsObject(att.Type); obj != nil && s.Properties != nil {
		// Build new properties map using JSON tag names. Use walkAttribute so bases/references are included.
		newProps := make(map[string]*openapi.Schema, len(s.Properties))
		nameMap := make(map[string]string, len(*obj))
		_ = walkAttribute(att, func(oldName string, child *goaexpr.AttributeExpr) error {
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
func walkAttribute(att *goaexpr.AttributeExpr, it func(name string, a *goaexpr.AttributeExpr) error) error { //nolint:cyclop
	switch dt := att.Type.(type) {
	case goaexpr.UserType:
		if err := walkAttribute(dt.Attribute(), it); err != nil {
			return err
		}
	case *goaexpr.Object:
		for _, nat := range *dt {
			if err := it(nat.Name, nat.Attribute); err != nil {
				return err
			}
		}
	}
	for _, b := range att.Bases {
		if err := walkAttribute(&goaexpr.AttributeExpr{Type: b}, it); err != nil {
			return err
		}
	}
	for _, r := range att.References {
		if err := walkAttribute(&goaexpr.AttributeExpr{Type: r}, it); err != nil {
			return err
		}
	}
	return nil
}

// jsonTagName extracts the JSON tag field name from the attribute Meta if set.
// It supports values like "name,omitempty" and returns "name".
func jsonTagName(att *goaexpr.AttributeExpr) string {
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
func transformExampleWithJSONTags(att *goaexpr.AttributeExpr, ex any) any {
	if att == nil || ex == nil {
		return ex
	}
	switch t := att.Type.(type) {
	case *goaexpr.ResultTypeExpr:
		return transformExampleWithJSONTags(t.AttributeExpr, ex)
	case goaexpr.UserType:
		return transformExampleWithJSONTags(t.Attribute(), ex)
	case *goaexpr.Object:
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
	case *goaexpr.Array:
		if arr, ok := ex.([]any); ok {
			out := make([]any, len(arr))
			for i := range arr {
				out[i] = transformExampleWithJSONTags(t.ElemType, arr[i])
			}
			return out
		}
		return ex
	case *goaexpr.Map:
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
	case *goaexpr.Union:
		// Attempt best-effort transform by applying first variant.
		if len(t.Values) > 0 {
			return transformExampleWithJSONTags(t.Values[0].Attribute, ex)
		}
		return ex
	default:
		return ex
	}
}
