// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// health endpoints
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/fetcher/archiver/design

package health

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints wraps the "health" service endpoints.
type Endpoints struct {
	Show endpoint.Endpoint
}

// NewEndpoints wraps the methods of the "health" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Show: NewShowEndpoint(s),
	}
}

// NewShowEndpoint returns an endpoint function that calls the method "show" of
// service "health".
func NewShowEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.Show(ctx)
	}
}
