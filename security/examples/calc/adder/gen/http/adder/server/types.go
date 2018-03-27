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
type AddInvalidScopesResponseBody string

// AddUnauthorizedResponseBody is the type of the "adder" service "add"
// endpoint HTTP response body for the "unauthorized" error.
type AddUnauthorizedResponseBody string

// NewAddInvalidScopesResponseBody builds the HTTP response body from the
// result of the "add" endpoint of the "adder" service.
func NewAddInvalidScopesResponseBody(res addersvc.InvalidScopes) AddInvalidScopesResponseBody {
	body := AddInvalidScopesResponseBody(res)
	return body
}

// NewAddUnauthorizedResponseBody builds the HTTP response body from the result
// of the "add" endpoint of the "adder" service.
func NewAddUnauthorizedResponseBody(res addersvc.Unauthorized) AddUnauthorizedResponseBody {
	body := AddUnauthorizedResponseBody(res)
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
