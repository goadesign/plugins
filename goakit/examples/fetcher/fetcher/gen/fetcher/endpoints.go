// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// fetcher endpoints
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/fetcher/fetcher/design

package fetchersvc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints wraps the "fetcher" service endpoints.
type Endpoints struct {
	Fetch endpoint.Endpoint
}

// NewEndpoints wraps the methods of the "fetcher" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Fetch: NewFetchEndpoint(s),
	}
}

// NewFetchEndpoint returns an endpoint function that calls the method "fetch"
// of service "fetcher".
func NewFetchEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*FetchPayload)
		return s.Fetch(ctx, p)
	}
}
