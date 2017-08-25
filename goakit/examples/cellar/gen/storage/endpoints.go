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
	ep := new(Endpoints)

	ep.List = func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.List(ctx)
	}

	ep.Show = func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ShowPayload)
		return s.Show(ctx, p)
	}

	ep.Add = func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*Bottle)
		return s.Add(ctx, p)
	}

	ep.Remove = func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*RemovePayload)
		return nil, s.Remove(ctx, p)
	}

	return ep
}
