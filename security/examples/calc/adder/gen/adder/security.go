// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// adder service security
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/adder/design

package addersvc

import (
	"context"

	"goa.design/goa"
	"goa.design/plugins/security"
)

// NewSecureEndpoints wraps the methods of a adder service with security scheme
// aware endpoints.
func NewSecureEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Add: SecureAdd(NewAddEndpoint(s)),
	}
}

// SecureAdd returns an endpoint function which initializes the context with
// the security requirements for the method "add" of service "adder".
func SecureAdd(ep goa.Endpoint) goa.Endpoint {
	reqs := []*security.Requirement{
		&security.Requirement{
			Schemes: []*security.Scheme{
				&security.Scheme{
					Kind: security.SchemeKind(3),
					Name: "api_key",
				},
			},
		},
	}
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		ctx = context.WithValue(ctx, security.ContextKey, reqs)
		return ep(ctx, req)
	}
}
