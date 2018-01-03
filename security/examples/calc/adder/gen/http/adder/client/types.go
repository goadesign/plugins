// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// adder HTTP client types
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/adder/design

package client

import (
	goa "goa.design/goa"
	addersvc "goa.design/plugins/security/examples/calc/adder/gen/adder"
)

// AddInvalidScopesResponseBody is the type of the "adder" service "add"
// endpoint HTTP response body for the "invalid-scopes" error.
type AddInvalidScopesResponseBody struct {
	Value *string `form:"value,omitempty" json:"value,omitempty" xml:"value,omitempty"`
}

// AddUnauthorizedResponseBody is the type of the "adder" service "add"
// endpoint HTTP response body for the "unauthorized" error.
type AddUnauthorizedResponseBody struct {
	Value *string `form:"value,omitempty" json:"value,omitempty" xml:"value,omitempty"`
}

// NewAddInvalidScopes builds a adder service add endpoint invalid-scopes error.
func NewAddInvalidScopes(body *AddInvalidScopesResponseBody) *addersvc.InvalidScopes {
	v := &addersvc.InvalidScopes{
		Value: *body.Value,
	}
	return v
}

// NewAddUnauthorized builds a adder service add endpoint unauthorized error.
func NewAddUnauthorized(body *AddUnauthorizedResponseBody) *addersvc.Unauthorized {
	v := &addersvc.Unauthorized{
		Value: *body.Value,
	}
	return v
}

// Validate runs the validations defined on AddInvalidScopesResponseBody
func (body *AddInvalidScopesResponseBody) Validate() (err error) {
	if body.Value == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("value", "body"))
	}
	return
}

// Validate runs the validations defined on AddUnauthorizedResponseBody
func (body *AddUnauthorizedResponseBody) Validate() (err error) {
	if body.Value == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("value", "body"))
	}
	return
}
