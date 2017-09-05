// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// secured_service HTTP server types
//
// Command:
// $ goa gen goa.design/plugins/security/example/design

package server

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
	Value string `form:"value" json:"value" xml:"value"`
}

// NewSigninUnauthorizedResponseBody builds the secured_service service signin
// endpoint response body from a result.
func NewSigninUnauthorizedResponseBody(res *securedservice.Unauthorized) *SigninUnauthorizedResponseBody {
	body := &SigninUnauthorizedResponseBody{
		Value: res.Value,
	}

	return body
}

// NewSecureSecurePayload builds a secured_service service secure endpoint
// payload.
func NewSecureSecurePayload(body *SecureRequestBody, fail *bool) *securedservice.SecurePayload {
	v := &securedservice.SecurePayload{
		Token: body.Token,
	}
	v.Fail = fail

	return v
}

// NewDoublySecureDoublySecurePayload builds a secured_service service
// doubly_secure endpoint payload.
func NewDoublySecureDoublySecurePayload(body *DoublySecureRequestBody, key *string) *securedservice.DoublySecurePayload {
	v := &securedservice.DoublySecurePayload{
		Token: body.Token,
	}
	v.Key = key

	return v
}

// NewAlsoDoublySecureAlsoDoublySecurePayload builds a secured_service service
// also_doubly_secure endpoint payload.
func NewAlsoDoublySecureAlsoDoublySecurePayload(body *AlsoDoublySecureRequestBody, key *string) *securedservice.AlsoDoublySecurePayload {
	v := &securedservice.AlsoDoublySecurePayload{
		Token: body.Token,
	}
	v.Key = key

	return v
}
