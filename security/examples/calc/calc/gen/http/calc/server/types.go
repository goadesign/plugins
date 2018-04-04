// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// calc HTTP server types
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/calc/design

package server

import (
	calcsvc "goa.design/plugins/security/examples/calc/calc/gen/calc"
)

// LoginUnauthorizedResponseBody is the type of the "calc" service "login"
// endpoint HTTP response body for the "unauthorized" error.
type LoginUnauthorizedResponseBody string

// NewLoginUnauthorizedResponseBody builds the HTTP response body from the
// result of the "login" endpoint of the "calc" service.
func NewLoginUnauthorizedResponseBody(res calcsvc.Unauthorized) LoginUnauthorizedResponseBody {
	body := LoginUnauthorizedResponseBody(res)
	return body
}

// NewLoginLoginPayload builds a calc service login endpoint payload.
func NewLoginLoginPayload(user string, password string) *calcsvc.LoginPayload {
	return &calcsvc.LoginPayload{
		User:     user,
		Password: password,
	}
}

// NewAddAddPayload builds a calc service add endpoint payload.
func NewAddAddPayload(a int, b int, token string) *calcsvc.AddPayload {
	return &calcsvc.AddPayload{
		A:     a,
		B:     b,
		Token: token,
	}
}
