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

// AddRequestBody is the type of the "calc" service "add" endpoint HTTP request
// body.
type AddRequestBody struct {
	// JWT used for authentication
	Token *string `form:"token,omitempty" json:"token,omitempty" xml:"token,omitempty"`
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
func NewAddAddPayload(body *AddRequestBody, a int, b int) *calcsvc.AddPayload {
	v := &calcsvc.AddPayload{
		Token: *body.Token,
	}
	v.A = a
	v.B = b
	return v
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

// Validate runs the validations defined on AddRequestBody
func (body *AddRequestBody) Validate() (err error) {
	if body.Token == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("token", "body"))
	}
	return
}
