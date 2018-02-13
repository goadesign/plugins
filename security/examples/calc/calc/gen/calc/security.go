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
func NewSecureEndpoints(s Service, authBasicAuthFn security.AuthorizeBasicAuthFunc, authJWTFn security.AuthorizeJWTFunc) *Endpoints {
	return &Endpoints{
		Login: SecureLogin(NewLoginEndpoint(s), authBasicAuthFn),
		Add:   SecureAdd(NewAddEndpoint(s), authJWTFn),
	}
}

// SecureLogin returns an endpoint function which initializes the context with
// the security requirements for the method "login" of service "calc".
func SecureLogin(ep goa.Endpoint, authBasicAuthFn security.AuthorizeBasicAuthFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*LoginPayload)
		var err error
		basicAuthSch := security.BasicAuthScheme{
			Name: "basic",
		}
		ctx, err = authBasicAuthFn(ctx, p.User, p.Password, &basicAuthSch)
		if err != nil {
			return nil, err
		}
		return ep(ctx, req)
	}
}

// SecureAdd returns an endpoint function which initializes the context with
// the security requirements for the method "add" of service "calc".
func SecureAdd(ep goa.Endpoint, authJWTFn security.AuthorizeJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*AddPayload)
		var err error
		jwtSch := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"calc:add"},
			RequiredScopes: []string{"calc:add"},
		}
		ctx, err = authJWTFn(ctx, p.Token, &jwtSch)
		if err != nil {
			return nil, err
		}
		return ep(ctx, req)
	}
}
