// Code generated by goa v3.1.2, DO NOT EDIT.
//
// fetcher endpoints
//
// Command:
// $ goa gen goa.design/plugins/v3/goakit/examples/fetcher/fetcher/design -o
// $(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/fetcher

package fetcher

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

// Use applies the given middleware to all the "fetcher" service endpoints.
func (e *Endpoints) Use(m func(endpoint.Endpoint) endpoint.Endpoint) {
	e.Fetch = m(e.Fetch)
}

// NewFetchEndpoint returns an endpoint function that calls the method "fetch"
// of service "fetcher".
func NewFetchEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*FetchPayload)
		res, err := s.Fetch(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedFetchMedia(res, "default")
		return vres, nil
	}
}
