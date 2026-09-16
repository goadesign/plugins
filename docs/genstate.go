package docs

import (
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/expr"
	"goa.design/goa/v3/http/codegen/openapi"
	openapiv2 "goa.design/goa/v3/http/codegen/openapi/v2"
)

// genState stores per-generation state so we can avoid mutating
// global DSL structures while still ensuring uniqueness for schema names.
type genState struct {
	nameScope     *codegen.NameScope
	assignedNames map[*expr.UserTypeExpr]string
	examples      *expr.ExampleGenerator
	definitions   map[string]*openapi.Schema
}

// newGenState constructs the state used to document one API.
func newGenState(api *expr.APIExpr) *genState {
	return &genState{
		nameScope:     codegen.NewNameScope(),
		assignedNames: make(map[*expr.UserTypeExpr]string),
		examples:      expr.NewExampleGenerator(api.RandomizerFactory),
		definitions:   make(map[string]*openapi.Schema),
	}
}

// uniqueNameFor returns a stable, unique name for the given user type within
// this generator run. It never mutates the passed user type.
func (s *genState) uniqueNameFor(ut *expr.UserTypeExpr) string {
	if n, ok := s.assignedNames[ut]; ok {
		return n
	}
	n := s.nameScope.Unique(ut.TypeName)
	s.assignedNames[ut] = n
	return n
}

// schemaForAttribute builds a schema and keeps its named definitions with the
// API currently being documented. It restores temporary type names before it
// returns so documentation generation does not change the evaluated design.
func schemaForAttribute(api *expr.APIExpr, att *expr.AttributeExpr, examples *expr.ExampleGenerator, state *genState) *openapi.Schema {
	if att == nil || att.Type == nil {
		return nil
	}
	if ut, ok := att.Type.(*expr.UserTypeExpr); ok {
		orig := ut.TypeName
		ut.TypeName = state.uniqueNameFor(ut)
		defer func() { ut.TypeName = orig }()
	}
	schema := openapiv2.BuildAttributeSchema(api, att, examples)
	for name, definition := range schema.Defs {
		state.definitions[name] = dupSchema(definition)
	}
	schema.Defs = nil
	return schema
}
