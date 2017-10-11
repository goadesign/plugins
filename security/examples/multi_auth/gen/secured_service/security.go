// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// secured_service service security
//
// Command:
// $ goa gen goa.design/plugins/security/examples/multi_auth/design

package securedservice

import (
	"context"

	"goa.design/goa"
	"goa.design/plugins/security"
)

// NewSecureEndpoints wraps the methods of a secured_service service with
// security scheme aware endpoints.
func NewSecureEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Signin:           SecureSignin(NewSigninEndpoint(s)),
		Secure:           SecureSecure(NewSecureEndpoint(s)),
		DoublySecure:     SecureDoublySecure(NewDoublySecureEndpoint(s)),
		AlsoDoublySecure: SecureAlsoDoublySecure(NewAlsoDoublySecureEndpoint(s)),
	}
}

// SecureSignin returns an endpoint function which initializes the context with
// the security requirements for the method "signin" of service
// "secured_service".
func SecureSignin(ep goa.Endpoint) goa.Endpoint {
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

// SecureSecure returns an endpoint function which initializes the context with
// the security requirements for the method "secure" of service
// "secured_service".
func SecureSecure(ep goa.Endpoint) goa.Endpoint {
	reqs := make([]*security.Requirement, 1)
	reqs[0].RequiredScopes = []string{"api:read"}
	reqs[0].Schemes = make([]*security.Scheme, 1)
	reqs[0].Schemes[0] = &security.Scheme{
		Kind: security.SchemeKind(4),
		Name: "jwt",
	}
	reqs[0].Schemes[0].Scopes = []string{"api:read", "api:write"}
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		ctx = context.WithValue(ctx, security.ContextKey, reqs)
		return ep(ctx, req)
	}
}

// SecureDoublySecure returns an endpoint function which initializes the
// context with the security requirements for the method "doubly_secure" of
// service "secured_service".
func SecureDoublySecure(ep goa.Endpoint) goa.Endpoint {
	reqs := make([]*security.Requirement, 1)
	reqs[0].RequiredScopes = []string{"api:read", "api:write"}
	reqs[0].Schemes = make([]*security.Scheme, 2)
	reqs[0].Schemes[0] = &security.Scheme{
		Kind: security.SchemeKind(4),
		Name: "jwt",
	}
	reqs[0].Schemes[0].Scopes = []string{"api:read", "api:write"}
	reqs[0].Schemes[1] = &security.Scheme{
		Kind: security.SchemeKind(3),
		Name: "api_key",
	}
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		ctx = context.WithValue(ctx, security.ContextKey, reqs)
		return ep(ctx, req)
	}
}

// SecureAlsoDoublySecure returns an endpoint function which initializes the
// context with the security requirements for the method "also_doubly_secure"
// of service "secured_service".
func SecureAlsoDoublySecure(ep goa.Endpoint) goa.Endpoint {
	reqs := make([]*security.Requirement, 1)
	reqs[0].RequiredScopes = []string{"api:read", "api:write"}
	reqs[0].Schemes = make([]*security.Scheme, 2)
	reqs[0].Schemes[0] = &security.Scheme{
		Kind: security.SchemeKind(4),
		Name: "jwt",
	}
	reqs[0].Schemes[0].Scopes = []string{"api:read", "api:write"}
	reqs[0].Schemes[1] = &security.Scheme{
		Kind: security.SchemeKind(3),
		Name: "api_key",
	}
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		ctx = context.WithValue(ctx, security.ContextKey, reqs)
		return ep(ctx, req)
	}
}
