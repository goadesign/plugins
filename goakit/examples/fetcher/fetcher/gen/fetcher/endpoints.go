// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// fetcher endpoints
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/client/fetcher/design

package fetcher

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type (
	// Endpoints wraps the fetcher service endpoints.
	Endpoints struct {
		Fetch endpoint.Endpoint
	}
)

// NewEndpoints wraps the methods of a fetcher service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Fetch: NewFetchEndpoint(s),
	}
}

// NewFetchEndpoint returns an endpoint function that calls method "fetch" of
// service "fetcher".
func NewFetchEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*FetchPayload)
		return s.Fetch(ctx, p)
	}
}
