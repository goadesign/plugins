// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// storage endpoints
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/cellar/design

package storage

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type (
	// Endpoints wraps the storage service endpoints.
	Endpoints struct {
		List   endpoint.Endpoint
		Show   endpoint.Endpoint
		Add    endpoint.Endpoint
		Remove endpoint.Endpoint
	}
)

// NewEndpoints wraps the methods of a storage service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		List:   NewListEndpoint(s),
		Show:   NewShowEndpoint(s),
		Add:    NewAddEndpoint(s),
		Remove: NewRemoveEndpoint(s),
	}
}

// NewListEndpoint returns an endpoint function that calls method "list" of
// service "storage".
func NewListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.List(ctx)
	}
}

// NewShowEndpoint returns an endpoint function that calls method "show" of
// service "storage".
func NewShowEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ShowPayload)
		return s.Show(ctx, p)
	}
}

// NewAddEndpoint returns an endpoint function that calls method "add" of
// service "storage".
func NewAddEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*Bottle)
		return s.Add(ctx, p)
	}
}

// NewRemoveEndpoint returns an endpoint function that calls method "remove" of
// service "storage".
func NewRemoveEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*RemovePayload)
		return nil, s.Remove(ctx, p)
	}
}
