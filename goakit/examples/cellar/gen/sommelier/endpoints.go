// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// sommelier endpoints
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/cellar/design

package sommelier

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type (
	// Endpoints wraps the sommelier service endpoints.
	Endpoints struct {
		Pick endpoint.Endpoint
	}
)

// NewEndpoints wraps the methods of a sommelier service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Pick: NewPickEndpoint(s),
	}
}

// NewPickEndpoint returns an endpoint function that calls method "pick" of
// service "sommelier".
func NewPickEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*Criteria)
		return s.Pick(ctx, p)
	}
}
