package docs

import (
	"strings"

	openapi "goa.design/goa/v3/http/codegen/openapi"
)

// inlineAllServiceSchemas walks all service method schemas and inlines any
// $ref occurrences that point to local definitions.
func inlineAllServiceSchemas(d *data, defs map[string]*openapi.Schema) {
	if d == nil || len(d.Services) == 0 {
		return
	}
	stack := make(map[string]bool)
	for _, svc := range d.Services {
		if svc == nil || len(svc.Methods) == 0 {
			continue
		}
		for _, m := range svc.Methods {
			if m == nil {
				continue
			}
			if m.Payload != nil && m.Payload.Type != nil {
				inlineRefs(m.Payload.Type, defs, stack)
			}
			if m.StreamingPayload != nil && m.StreamingPayload.Type != nil {
				inlineRefs(m.StreamingPayload.Type, defs, stack)
			}
			if m.Result != nil && m.Result.Type != nil {
				inlineRefs(m.Result.Type, defs, stack)
			}
			if m.StreamingResult != nil && m.StreamingResult.Type != nil {
				inlineRefs(m.StreamingResult.Type, defs, stack)
			}
			if len(m.Errors) > 0 {
				for _, e := range m.Errors {
					if e != nil && e.Type != nil {
						inlineRefs(e.Type, defs, stack)
					}
				}
			}
		}
	}
}

// inlineRefs replaces local definition references with their fully expanded
// schema copy. Cycles are broken by retaining a single $ref at the cycle edge.
func inlineRefs(s *openapi.Schema, defs map[string]*openapi.Schema, stack map[string]bool) {
	if s == nil {
		return
	}

	if s.Ref != "" {
		// JSON Schema 2020-12 uses $defs for local schema definitions.
		const prefix = "#/$defs/"
		if !strings.HasPrefix(s.Ref, prefix) {
			// Unexpected external ref; leave intact.
			return
		}
		name := strings.TrimPrefix(s.Ref, prefix)
		if name == "" {
			return
		}
		if stack[name] {
			// Cycle detected; keep the $ref here to break infinite expansion.
			return
		}
		def, ok := defs[name]
		if !ok || def == nil {
			// Defensive: unresolved ref, leave as-is.
			return
		}
		stack[name] = true
		copy := dupSchema(def)
		inlineRefs(copy, defs, stack)
		*s = *copy
		s.Ref = ""
		delete(stack, name)
		return
	}

	// Recurse into container positions.
	if len(s.Properties) > 0 {
		for _, p := range s.Properties {
			inlineRefs(p, defs, stack)
		}
	}
	if s.Items != nil {
		inlineRefs(s.Items, defs, stack)
	}
	if ap, ok := s.AdditionalProperties.(*openapi.Schema); ok && ap != nil {
		inlineRefs(ap, defs, stack)
	}
	if len(s.AnyOf) > 0 {
		for _, a := range s.AnyOf {
			inlineRefs(a, defs, stack)
		}
	}
	if len(s.Defs) > 0 {
		for _, d := range s.Defs {
			inlineRefs(d, defs, stack)
		}
	}
}
