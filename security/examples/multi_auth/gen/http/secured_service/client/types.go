// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// secured_service HTTP client types
//
// Command:
// $ goa gen goa.design/plugins/security/examples/multi_auth/design

package client

import (
	goa "goa.design/goa"
	securedservice "goa.design/plugins/security/examples/multi_auth/gen/secured_service"
)

// SigninRequestBody is the type of the "secured_service" service "signin"
// endpoint HTTP request body.
type SigninRequestBody struct {
	// Username used to perform signin
	Username *string `form:"username,omitempty" json:"username,omitempty" xml:"username,omitempty"`
	// Username used to perform signin
	Password *string `form:"password,omitempty" json:"password,omitempty" xml:"password,omitempty"`
}

// SecureRequestBody is the type of the "secured_service" service "secure"
// endpoint HTTP request body.
type SecureRequestBody struct {
	// JWT used for authentication
	Token *string `form:"token,omitempty" json:"token,omitempty" xml:"token,omitempty"`
}

// DoublySecureRequestBody is the type of the "secured_service" service
// "doubly_secure" endpoint HTTP request body.
type DoublySecureRequestBody struct {
	// JWT used for authentication
	Token *string `form:"token,omitempty" json:"token,omitempty" xml:"token,omitempty"`
}

// AlsoDoublySecureRequestBody is the type of the "secured_service" service
// "also_doubly_secure" endpoint HTTP request body.
type AlsoDoublySecureRequestBody struct {
	// JWT used for authentication
	Token *string `form:"token,omitempty" json:"token,omitempty" xml:"token,omitempty"`
}

// SigninUnauthorizedResponseBody is the type of the "secured_service" service
// "signin" endpoint HTTP response body for the "unauthorized" error.
type SigninUnauthorizedResponseBody struct {
	// Credentials are invalid
	Value *string `form:"value,omitempty" json:"value,omitempty" xml:"value,omitempty"`
}

// NewSigninRequestBody builds the HTTP request body from the payload of the
// "signin" endpoint of the "secured_service" service.
func NewSigninRequestBody(p *securedservice.SigninPayload) *SigninRequestBody {
	body := &SigninRequestBody{
		Username: p.Username,
		Password: p.Password,
	}
	return body
}

// NewSecureRequestBody builds the HTTP request body from the payload of the
// "secure" endpoint of the "secured_service" service.
func NewSecureRequestBody(p *securedservice.SecurePayload) *SecureRequestBody {
	body := &SecureRequestBody{
		Token: p.Token,
	}
	return body
}

// NewDoublySecureRequestBody builds the HTTP request body from the payload of
// the "doubly_secure" endpoint of the "secured_service" service.
func NewDoublySecureRequestBody(p *securedservice.DoublySecurePayload) *DoublySecureRequestBody {
	body := &DoublySecureRequestBody{
		Token: p.Token,
	}
	return body
}

// NewAlsoDoublySecureRequestBody builds the HTTP request body from the payload
// of the "also_doubly_secure" endpoint of the "secured_service" service.
func NewAlsoDoublySecureRequestBody(p *securedservice.AlsoDoublySecurePayload) *AlsoDoublySecureRequestBody {
	body := &AlsoDoublySecureRequestBody{
		Token: p.Token,
	}
	return body
}

// NewSigninUnauthorized builds a secured_service service signin endpoint
// unauthorized error.
func NewSigninUnauthorized(body *SigninUnauthorizedResponseBody) *securedservice.Unauthorized {
	v := &securedservice.Unauthorized{
		Value: *body.Value,
	}
	return v
}

// Validate runs the validations defined on SigninUnauthorizedResponseBody
func (body *SigninUnauthorizedResponseBody) Validate() (err error) {
	if body.Value == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("value", "body"))
	}
	return
}
