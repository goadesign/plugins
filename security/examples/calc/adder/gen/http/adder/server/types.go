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

// InvalidScopes is the type of the "adder" service "add" endpoint HTTP
// response body for the "invalid-scopes" error.
type InvalidScopes string

// Unauthorized is the type of the "adder" service "add" endpoint HTTP response
// body for the "unauthorized" error.
type Unauthorized string

// NewInvalidScopes builds the HTTP response body from the result of the "add"
// endpoint of the "adder" service.
func NewInvalidScopes(res addersvc.InvalidScopes) InvalidScopes {
	body := InvalidScopes(res)
	return body
}

// NewUnauthorized builds the HTTP response body from the result of the "add"
// endpoint of the "adder" service.
func NewUnauthorized(res addersvc.Unauthorized) Unauthorized {
	body := Unauthorized(res)
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
