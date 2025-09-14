package docs

import (
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/expr"
	"goa.design/goa/v3/http/codegen/openapi"
)

// genState stores per-generation state so we can avoid mutating
// global DSL structures while still ensuring uniqueness for schema names.
type genState struct {
	nameScope     *codegen.NameScope
	assignedNames map[*expr.UserTypeExpr]string
}

// newGenState constructs a fresh generator state for one Generate run.
func newGenState() *genState {
	return &genState{
		nameScope:     codegen.NewNameScope(),
		assignedNames: make(map[*expr.UserTypeExpr]string),
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

// schemaForAttribute computes the OpenAPI schema for the given attribute using
// a temporary, per-call unique name for user types. The original AST is restored
// before returning.
func schemaForAttribute(api *expr.APIExpr, att *expr.AttributeExpr, state *genState) *openapi.Schema {
	if att == nil || att.Type == nil {
		return nil
	}
	if ut, ok := att.Type.(*expr.UserTypeExpr); ok {
		orig := ut.TypeName
		ut.TypeName = state.uniqueNameFor(ut)
		defer func() { ut.TypeName = orig }()
	}
	return openapi.AttributeTypeSchema(api, att)
}
