// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// archiver endpoints
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/fetcher/archiver/design

package archiversvc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints wraps the "archiver" service endpoints.
type Endpoints struct {
	Archive endpoint.Endpoint
	Read    endpoint.Endpoint
}

// NewEndpoints wraps the methods of the "archiver" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Archive: NewArchiveEndpoint(s),
		Read:    NewReadEndpoint(s),
	}
}

// NewArchiveEndpoint returns an endpoint function that calls the method
// "archive" of service "archiver".
func NewArchiveEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ArchivePayload)
		return s.Archive(ctx, p)
	}
}

// NewReadEndpoint returns an endpoint function that calls the method "read" of
// service "archiver".
func NewReadEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ReadPayload)
		return s.Read(ctx, p)
	}
}