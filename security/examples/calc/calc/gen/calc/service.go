// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// calc service
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/calc/design

package calcsvc

import (
	"context"
)

// The calc service exposes public endpoints that require valid authorization
// credentials.
type Service interface {
	// Creates a valid JWT
	Login(context.Context, *LoginPayload) (string, error)
	// Add adds up the two integer parameters and returns the results. This
	// endpoint is secured with the JWT scheme
	Add(context.Context, *AddPayload) (int, error)
}

// Credentials used to authenticate to retrieve JWT token
type LoginPayload struct {
	User     string
	Password string
}

// AddPayload is the payload type of the calc service add method.
type AddPayload struct {
	// Left operand
	A int
	// Right operand
	B int
	// JWT used for authentication
	Token string
}

type Unauthorized struct {
	// Credentials are invalid
	Value string
}

// Error returns "unauthorized".
func (e *Unauthorized) Error() string {
	return "unauthorized"
}
