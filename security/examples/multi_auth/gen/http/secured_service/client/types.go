// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// secured_service HTTP client types
//
// Command:
// $ goa gen goa.design/plugins/security/examples/multi_auth/design

package client

import (
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

// AlsoDoublySecureRequestBody is the type of the "secured_service" service
// "also_doubly_secure" endpoint HTTP request body.
type AlsoDoublySecureRequestBody struct {
	// Username used to perform signin
	Username *string `form:"username,omitempty" json:"username,omitempty" xml:"username,omitempty"`
	// Username used to perform signin
	Password *string `form:"password,omitempty" json:"password,omitempty" xml:"password,omitempty"`
}

// SigninUnauthorizedResponseBody is the type of the "secured_service" service
// "signin" endpoint HTTP response body for the "unauthorized" error.
type SigninUnauthorizedResponseBody string

// SecureUnauthorizedResponseBody is the type of the "secured_service" service
// "secure" endpoint HTTP response body for the "unauthorized" error.
type SecureUnauthorizedResponseBody string

// DoublySecureUnauthorizedResponseBody is the type of the "secured_service"
// service "doubly_secure" endpoint HTTP response body for the "unauthorized"
// error.
type DoublySecureUnauthorizedResponseBody string

// AlsoDoublySecureUnauthorizedResponseBody is the type of the
// "secured_service" service "also_doubly_secure" endpoint HTTP response body
// for the "unauthorized" error.
type AlsoDoublySecureUnauthorizedResponseBody string

// NewSigninRequestBody builds the HTTP request body from the payload of the
// "signin" endpoint of the "secured_service" service.
func NewSigninRequestBody(p *securedservice.SigninPayload) *SigninRequestBody {
	body := &SigninRequestBody{
		Username: p.Username,
		Password: p.Password,
	}
	return body
}

// NewAlsoDoublySecureRequestBody builds the HTTP request body from the payload
// of the "also_doubly_secure" endpoint of the "secured_service" service.
func NewAlsoDoublySecureRequestBody(p *securedservice.AlsoDoublySecurePayload) *AlsoDoublySecureRequestBody {
	body := &AlsoDoublySecureRequestBody{
		Username: p.Username,
		Password: p.Password,
	}
	return body
}

// NewSigninUnauthorized builds a secured_service service signin endpoint
// unauthorized error.
func NewSigninUnauthorized(body SigninUnauthorizedResponseBody) securedservice.Unauthorized {
	v := securedservice.Unauthorized(body)
	return v
}

// NewSecureUnauthorized builds a secured_service service secure endpoint
// unauthorized error.
func NewSecureUnauthorized(body SecureUnauthorizedResponseBody) securedservice.Unauthorized {
	v := securedservice.Unauthorized(body)
	return v
}

// NewDoublySecureUnauthorized builds a secured_service service doubly_secure
// endpoint unauthorized error.
func NewDoublySecureUnauthorized(body DoublySecureUnauthorizedResponseBody) securedservice.Unauthorized {
	v := securedservice.Unauthorized(body)
	return v
}

// NewAlsoDoublySecureUnauthorized builds a secured_service service
// also_doubly_secure endpoint unauthorized error.
func NewAlsoDoublySecureUnauthorized(body AlsoDoublySecureUnauthorizedResponseBody) securedservice.Unauthorized {
	v := securedservice.Unauthorized(body)
	return v
}
