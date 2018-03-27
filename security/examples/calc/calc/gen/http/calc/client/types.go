// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// calc HTTP client types
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/calc/design

package client

import (
	calcsvc "goa.design/plugins/security/examples/calc/calc/gen/calc"
)

// LoginRequestBody is the type of the "calc" service "login" endpoint HTTP
// request body.
type LoginRequestBody struct {
	User     string `form:"user" json:"user" xml:"user"`
	Password string `form:"password" json:"password" xml:"password"`
}

// AddRequestBody is the type of the "calc" service "add" endpoint HTTP request
// body.
type AddRequestBody struct {
	// JWT used for authentication
	Token string `form:"token" json:"token" xml:"token"`
}

// LoginUnauthorizedResponseBody is the type of the "calc" service "login"
// endpoint HTTP response body for the "unauthorized" error.
type LoginUnauthorizedResponseBody string

// NewLoginRequestBody builds the HTTP request body from the payload of the
// "login" endpoint of the "calc" service.
func NewLoginRequestBody(p *calcsvc.LoginPayload) *LoginRequestBody {
	body := &LoginRequestBody{
		User:     p.User,
		Password: p.Password,
	}
	return body
}

// NewAddRequestBody builds the HTTP request body from the payload of the "add"
// endpoint of the "calc" service.
func NewAddRequestBody(p *calcsvc.AddPayload) *AddRequestBody {
	body := &AddRequestBody{
		Token: p.Token,
	}
	return body
}

// NewLoginUnauthorized builds a calc service login endpoint unauthorized error.
func NewLoginUnauthorized(body LoginUnauthorizedResponseBody) calcsvc.Unauthorized {
	v := calcsvc.Unauthorized(body)
	return v
}
