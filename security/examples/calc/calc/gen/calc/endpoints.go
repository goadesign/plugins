// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// calc endpoints
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/calc/design

package calcsvc

import (
	"context"

	goa "goa.design/goa"
)

// Endpoints wraps the "calc" service endpoints.
type Endpoints struct {
	Login goa.Endpoint
	Add   goa.Endpoint
}

// NewEndpoints wraps the methods of the "calc" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Login: NewLoginEndpoint(s),
		Add:   NewAddEndpoint(s),
	}
}

// NewLoginEndpoint returns an endpoint function that calls the method "login"
// of service "calc".
func NewLoginEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*LoginPayload)
		return s.Login(ctx, p)
	}
}

// NewAddEndpoint returns an endpoint function that calls the method "add" of
// service "calc".
func NewAddEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*AddPayload)
		return s.Add(ctx, p)
	}
}
