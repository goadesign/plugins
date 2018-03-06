// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// calc service
//
// Command:
// $ goa gen goa.design/plugins/cors/examples/calc/design

package calcsvc

import (
	"context"
)

// The calc service exposes public endpoints that defines CORS policy.
type Service interface {
	// Add adds up the two integer parameters and returns the results.
	Add(context.Context, *AddPayload) (int, error)
}

// AddPayload is the payload type of the calc service add method.
type AddPayload struct {
	// Left operand
	A int
	// Right operand
	B int
}
