// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// adder service
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/adder/design

package addersvc

import (
	"context"
)

// The adder service exposes an add method secured via API keys.
type Service interface {
	// This action returns the sum of two integers and is secured with the API key
	// scheme
	Add(context.Context, *AddPayload) (int, error)
}

// AddPayload is the payload type of the adder service add method.
type AddPayload struct {
	// API key
	Key string
	// Left operand
	A int
	// Right operand
	B int
}

type Unauthorized struct {
	Value string
}

type InvalidScopes struct {
	Value string
}

// Error returns "unauthorized".
func (e *Unauthorized) Error() string {
	return "unauthorized"
}

// Error returns "invalid-scopes".
func (e *InvalidScopes) Error() string {
	return "invalid-scopes"
}
