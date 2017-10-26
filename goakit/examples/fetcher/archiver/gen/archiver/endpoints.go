// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// archiver endpoints
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/client/archiver/design

package archiver

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type (
	// Endpoints wraps the archiver service endpoints.
	Endpoints struct {
		Archive endpoint.Endpoint
		Read    endpoint.Endpoint
	}
)

// NewEndpoints wraps the methods of a archiver service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Archive: NewArchiveEndpoint(s),
		Read:    NewReadEndpoint(s),
	}
}

// NewArchiveEndpoint returns an endpoint function that calls method "archive"
// of service "archiver".
func NewArchiveEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ArchivePayload)
		return s.Archive(ctx, p)
	}
}

// NewReadEndpoint returns an endpoint function that calls method "read" of
// service "archiver".
func NewReadEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ReadPayload)
		return s.Read(ctx, p)
	}
}
