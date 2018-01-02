// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// calc service security
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/calc/design

package calcsvc

import (
	"context"

	"goa.design/goa"
	"goa.design/plugins/security"
)

// NewSecureEndpoints wraps the methods of a calc service with security scheme
// aware endpoints.
func NewSecureEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Login: SecureLogin(NewLoginEndpoint(s)),
		Add:   SecureAdd(NewAddEndpoint(s)),
	}
}

// SecureLogin returns an endpoint function which initializes the context with
// the security requirements for the method "login" of service "calc".
func SecureLogin(ep goa.Endpoint) goa.Endpoint {
	reqs := make([]*security.Requirement, 1)
	reqs[0].Schemes = make([]*security.Scheme, 1)
	reqs[0].Schemes[0] = &security.Scheme{
		Kind: security.SchemeKind(2),
		Name: "basic",
	}
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		ctx = context.WithValue(ctx, security.ContextKey, reqs)
		return ep(ctx, req)
	}
}

// SecureAdd returns an endpoint function which initializes the context with
// the security requirements for the method "add" of service "calc".
func SecureAdd(ep goa.Endpoint) goa.Endpoint {
	reqs := make([]*security.Requirement, 1)
	reqs[0].RequiredScopes = []string{"calc:add"}
	reqs[0].Schemes = make([]*security.Scheme, 1)
	reqs[0].Schemes[0] = &security.Scheme{
		Kind: security.SchemeKind(4),
		Name: "jwt",
	}
	reqs[0].Schemes[0].Scopes = []string{"calc:add"}
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		ctx = context.WithValue(ctx, security.ContextKey, reqs)
		return ep(ctx, req)
	}
}
