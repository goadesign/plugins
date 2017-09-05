// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// secured_service endpoints
//
// Command:
// $ goa gen goa.design/plugins/security/example/design

package securedService

import (
	"context"

	goa "goa.design/goa"
)

type (
	// Endpoints wraps the secured_service service endpoints.
	Endpoints struct {
		Signin           goa.Endpoint
		Secure           goa.Endpoint
		DoublySecure     goa.Endpoint
		AlsoDoublySecure goa.Endpoint
	}
)

// NewEndpoints wraps the methods of a secured_service service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	ep := new(Endpoints)

	ep.Signin = func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(string)
		return s.Signin(ctx, p)
	}

	ep.Secure = func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*SecurePayload)
		return s.Secure(ctx, p)
	}

	ep.DoublySecure = func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*DoublySecurePayload)
		return s.DoublySecure(ctx, p)
	}

	ep.AlsoDoublySecure = func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*AlsoDoublySecurePayload)
		return s.AlsoDoublySecure(ctx, p)
	}

	return ep
}
