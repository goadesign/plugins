// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// calc HTTP server types
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/calc/design

package server

import (
	goa "goa.design/goa"
	calcsvc "goa.design/plugins/security/examples/calc/calc/gen/calc"
)

// LoginRequestBody is the type of the "calc" service "login" endpoint HTTP
// request body.
type LoginRequestBody struct {
	User     *string `form:"user,omitempty" json:"user,omitempty" xml:"user,omitempty"`
	Password *string `form:"password,omitempty" json:"password,omitempty" xml:"password,omitempty"`
}

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
func NewLoginLoginPayload(body *LoginRequestBody) *calcsvc.LoginPayload {
	v := &calcsvc.LoginPayload{
		User:     *body.User,
		Password: *body.Password,
	}
	return v
}

// NewAddAddPayload builds a calc service add endpoint payload.
func NewAddAddPayload(a int, b int, token string) *calcsvc.AddPayload {
	return &calcsvc.AddPayload{
		A:     a,
		B:     b,
		Token: token,
	}
}

// Validate runs the validations defined on LoginRequestBody
func (body *LoginRequestBody) Validate() (err error) {
	if body.User == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("user", "body"))
	}
	if body.Password == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("password", "body"))
	}
	return
}
