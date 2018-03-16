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
func NewSecureEndpoints(s Service, authBasicAuthFn security.AuthorizeBasicAuthFunc, authJWTFn security.AuthorizeJWTFunc, authAPIKeyFn security.AuthorizeAPIKeyFunc, authOAuth2Fn security.AuthorizeOAuth2Func) *Endpoints {
	return &Endpoints{
		Signin:           SecureSignin(NewSigninEndpoint(s), authBasicAuthFn),
		Secure:           SecureSecure(NewSecureEndpoint(s), authJWTFn),
		DoublySecure:     SecureDoublySecure(NewDoublySecureEndpoint(s), authJWTFn, authAPIKeyFn),
		AlsoDoublySecure: SecureAlsoDoublySecure(NewAlsoDoublySecureEndpoint(s), authJWTFn, authAPIKeyFn, authOAuth2Fn, authBasicAuthFn),
	}
}

// SecureSignin returns an endpoint function which initializes the context with
// the security requirements for the method "signin" of service
// "secured_service".
func SecureSignin(ep goa.Endpoint, authBasicAuthFn security.AuthorizeBasicAuthFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*SigninPayload)
		var err error
		basicAuthSch := security.BasicAuthScheme{
			Name: "basic",
		}
		ctx, err = authBasicAuthFn(ctx, *p.Username, *p.Password, &basicAuthSch)
		if err != nil {
			return nil, err
		}
		return ep(ctx, req)
	}
}

// SecureSecure returns an endpoint function which initializes the context with
// the security requirements for the method "secure" of service
// "secured_service".
func SecureSecure(ep goa.Endpoint, authJWTFn security.AuthorizeJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*SecurePayload)
		var err error
		jwtSch := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"api:read", "api:write"},
			RequiredScopes: []string{"api:read"},
		}
		ctx, err = authJWTFn(ctx, *p.Token, &jwtSch)
		if err != nil {
			return nil, err
		}
		return ep(ctx, req)
	}
}

// SecureDoublySecure returns an endpoint function which initializes the
// context with the security requirements for the method "doubly_secure" of
// service "secured_service".
func SecureDoublySecure(ep goa.Endpoint, authJWTFn security.AuthorizeJWTFunc, authAPIKeyFn security.AuthorizeAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*DoublySecurePayload)
		var err error
		jwtSch := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"api:read", "api:write"},
			RequiredScopes: []string{"api:read", "api:write"},
		}
		ctx, err = authJWTFn(ctx, *p.Token, &jwtSch)
		if err == nil {
			apiKeySch := security.APIKeyScheme{
				Name: "api_key",
			}
			ctx, err = authAPIKeyFn(ctx, *p.Key, &apiKeySch)
		}
		if err != nil {
			return nil, err
		}
		return ep(ctx, req)
	}
}

// SecureAlsoDoublySecure returns an endpoint function which initializes the
// context with the security requirements for the method "also_doubly_secure"
// of service "secured_service".
func SecureAlsoDoublySecure(ep goa.Endpoint, authJWTFn security.AuthorizeJWTFunc, authAPIKeyFn security.AuthorizeAPIKeyFunc, authOAuth2Fn security.AuthorizeOAuth2Func, authBasicAuthFn security.AuthorizeBasicAuthFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*AlsoDoublySecurePayload)
		var err error
		jwtSch := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"api:read", "api:write"},
			RequiredScopes: []string{"api:read", "api:write"},
		}
		ctx, err = authJWTFn(ctx, *p.Token, &jwtSch)
		if err == nil {
			apiKeySch := security.APIKeyScheme{
				Name: "api_key",
			}
			ctx, err = authAPIKeyFn(ctx, *p.Key, &apiKeySch)
		}
		if err != nil {
			oauth2Sch := security.OAuth2Scheme{
				Name:           "oauth2",
				Scopes:         []string{"api:read", "api:write"},
				RequiredScopes: []string{"api:read", "api:write"},
				Flows: []*security.OAuthFlow{
					&security.OAuthFlow{
						Type:             "authorization_code",
						AuthorizationURL: "http://localhost:8080/authorization",
						TokenURL:         "http://localhost:8080/token",
						RefreshURL:       "http://localhost:8080/refresh",
					},
				},
			}
			ctx, err = authOAuth2Fn(ctx, *p.OauthToken, &oauth2Sch)
			if err == nil {
				basicAuthSch := security.BasicAuthScheme{
					Name: "basic",
				}
				ctx, err = authBasicAuthFn(ctx, *p.Username, *p.Password, &basicAuthSch)
			}
		}
		if err != nil {
			return nil, err
		}
		return ep(ctx, req)
	}
}
