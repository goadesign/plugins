// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// secured_service endpoints
//
// Command:
// $ goa gen goa.design/plugins/security/examples/multi_auth/design

package securedservice

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
	return &Endpoints{
		Signin:           NewSigninEndpoint(s),
		Secure:           NewSecureEndpoint(s),
		DoublySecure:     NewDoublySecureEndpoint(s),
		AlsoDoublySecure: NewAlsoDoublySecureEndpoint(s),
	}
}

// NewSigninEndpoint returns an endpoint function that calls method "signin" of
// service "secured_service".
func NewSigninEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*SigninPayload)
		return s.Signin(ctx, p)
	}
}

// NewSecureEndpoint returns an endpoint function that calls method "secure" of
// service "secured_service".
func NewSecureEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*SecurePayload)
		return s.Secure(ctx, p)
	}
}

// NewDoublySecureEndpoint returns an endpoint function that calls method
// "doubly_secure" of service "secured_service".
func NewDoublySecureEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*DoublySecurePayload)
		return s.DoublySecure(ctx, p)
	}
}

// NewAlsoDoublySecureEndpoint returns an endpoint function that calls method
// "also_doubly_secure" of service "secured_service".
func NewAlsoDoublySecureEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*AlsoDoublySecurePayload)
		return s.AlsoDoublySecure(ctx, p)
	}
}
