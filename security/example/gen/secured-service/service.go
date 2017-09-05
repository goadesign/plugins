// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// secured_service service
//
// Command:
// $ goa gen goa.design/plugins/security/example/design

package securedService

import "context"

// Service is the secured_service service interface.
type Service interface {
	// Creates a valid JWT
	Signin(context.Context, string) (string, error)
	// This action is secured with the jwt scheme
	Secure(context.Context, *SecurePayload) (string, error)
	// This action is secured with the jwt scheme and also requires an API key query string.
	DoublySecure(context.Context, *DoublySecurePayload) (string, error)
	// This action is secured with the jwt scheme and also requires an API key header.
	AlsoDoublySecure(context.Context, *AlsoDoublySecurePayload) (string, error)
}

// SecurePayload is the payload type of the secured_service service secure
// method.
type SecurePayload struct {
	// Whether to force auth failure even with a valid JWT
	Fail *bool
	// JWT used for authentication
	Token *string
}

// DoublySecurePayload is the payload type of the secured_service service
// doubly_secure method.
type DoublySecurePayload struct {
	// API key
	Key *string
	// JWT used for authentication
	Token *string
}

// AlsoDoublySecurePayload is the payload type of the secured_service service
// also_doubly_secure method.
type AlsoDoublySecurePayload struct {
	// API key
	Key *string
	// JWT used for authentication
	Token *string
}

type Unauthorized struct {
	// Credentials are invalid
	Value string
}

// Error returns "unauthorized".
func (e *Unauthorized) Error() string {
	return "unauthorized"
}
