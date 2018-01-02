// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// adder HTTP server types
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/adder/design

package server

import (
	addersvc "goa.design/plugins/security/examples/calc/adder/gen/adder"
)

// AddInvalidScopesResponseBody is the type of the "adder" service "add"
// endpoint HTTP response body for the "invalid-scopes" error.
type AddInvalidScopesResponseBody struct {
	Value string `form:"value" json:"value" xml:"value"`
}

// AddUnauthorizedResponseBody is the type of the "adder" service "add"
// endpoint HTTP response body for the "unauthorized" error.
type AddUnauthorizedResponseBody struct {
	Value string `form:"value" json:"value" xml:"value"`
}

// NewAddInvalidScopesResponseBody builds the HTTP response body from the
// result of the "add" endpoint of the "adder" service.
func NewAddInvalidScopesResponseBody(res *addersvc.InvalidScopes) *AddInvalidScopesResponseBody {
	body := &AddInvalidScopesResponseBody{
		Value: res.Value,
	}
	return body
}

// NewAddUnauthorizedResponseBody builds the HTTP response body from the result
// of the "add" endpoint of the "adder" service.
func NewAddUnauthorizedResponseBody(res *addersvc.Unauthorized) *AddUnauthorizedResponseBody {
	body := &AddUnauthorizedResponseBody{
		Value: res.Value,
	}
	return body
}

// NewAddAddPayload builds a adder service add endpoint payload.
func NewAddAddPayload(a int, b int, key string) *addersvc.AddPayload {
	return &addersvc.AddPayload{
		A:   a,
		B:   b,
		Key: key,
	}
}
