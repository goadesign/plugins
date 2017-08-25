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
	ep := new(Endpoints)

	ep.Pick = func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*Criteria)
		return s.Pick(ctx, p)
	}

	return ep
}
