// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// adder HTTP client types
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/adder/design

package client

import (
	addersvc "goa.design/plugins/security/examples/calc/adder/gen/adder"
)

// InvalidScopes is the type of the "adder" service "add" endpoint HTTP
// response body for the "invalid-scopes" error.
type InvalidScopes string

// Unauthorized is the type of the "adder" service "add" endpoint HTTP response
// body for the "unauthorized" error.
type Unauthorized string

// NewAddInvalidScopes builds a adder service add endpoint invalid-scopes error.
func NewAddInvalidScopes(body InvalidScopes) addersvc.InvalidScopes {
	v := addersvc.InvalidScopes(body)
	return v
}

// NewAddUnauthorized builds a adder service add endpoint unauthorized error.
func NewAddUnauthorized(body Unauthorized) addersvc.Unauthorized {
	v := addersvc.Unauthorized(body)
	return v
}
