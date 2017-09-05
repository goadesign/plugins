// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// secured_service HTTP client types
//
// Command:
// $ goa gen goa.design/plugins/security/example/design

package client

import (
	"goa.design/plugins/security/example/gen/securedservice"
)

// SecureRequestBody is the type of the secured_service secure HTTP endpoint
// request body.
type SecureRequestBody struct {
	// JWT used for authentication
	Token *string `form:"token,omitempty" json:"token,omitempty" xml:"token,omitempty"`
}

// DoublySecureRequestBody is the type of the secured_service doubly_secure
// HTTP endpoint request body.
type DoublySecureRequestBody struct {
	// JWT used for authentication
	Token *string `form:"token,omitempty" json:"token,omitempty" xml:"token,omitempty"`
}

// AlsoDoublySecureRequestBody is the type of the secured_service
// also_doubly_secure HTTP endpoint request body.
type AlsoDoublySecureRequestBody struct {
	// JWT used for authentication
	Token *string `form:"token,omitempty" json:"token,omitempty" xml:"token,omitempty"`
}

// SigninUnauthorizedResponseBody is the type of the secured_service "signin"
// HTTP endpoint unauthorized error response body.
type SigninUnauthorizedResponseBody struct {
	// Credentials are invalid
	Value *string `form:"value,omitempty" json:"value,omitempty" xml:"value,omitempty"`
}

// NewSecureRequestBody builds the secured_service service secure endpoint
// request body from a payload.
func NewSecureRequestBody(p *securedservice.SecurePayload) *SecureRequestBody {
	body := &SecureRequestBody{
		Token: p.Token,
	}

	return body
}

// NewDoublySecureRequestBody builds the secured_service service doubly_secure
// endpoint request body from a payload.
func NewDoublySecureRequestBody(p *securedservice.DoublySecurePayload) *DoublySecureRequestBody {
	body := &DoublySecureRequestBody{
		Token: p.Token,
	}

	return body
}

// NewAlsoDoublySecureRequestBody builds the secured_service service
// also_doubly_secure endpoint request body from a payload.
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
