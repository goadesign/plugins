package testdata

var EndpointContextWithRequiredScopesCode = `// SecureSecureWithRequiredScopes returns an endpoint function which
// initializes the context with the security requirements for the method
// "SecureWithRequiredScopes" of service "EndpointWithRequiredScopes".
func SecureSecureWithRequiredScopes(ep goa.Endpoint) goa.Endpoint {
	reqs := []*security.Requirement{
		&security.Requirement{
			RequiredScopes: []string{"api:read", "api:write"},
			Schemes: []*security.Scheme{
				&security.Scheme{
					Kind:   security.SchemeKind(4),
					Name:   "jwt",
					Scopes: []string{"api:read", "api:write"},
				},
			},
		},
	}
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		ctx = context.WithValue(ctx, security.ContextKey, reqs)
		return ep(ctx, req)
	}
}
`

var EndpointContextWithAPIKeyOverrideCode = `// SecureSecureWithAPIKeyOverride returns an endpoint function which
// initializes the context with the security requirements for the method
// "SecureWithAPIKeyOverride" of service "EndpointWithAPIKeyOverride".
func SecureSecureWithAPIKeyOverride(ep goa.Endpoint) goa.Endpoint {
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
`

var EndpointContextWithOAuth2Code = `// SecureSecureWithOAuth2 returns an endpoint function which initializes the
// context with the security requirements for the method "SecureWithOAuth2" of
// service "EndpointWithOAuth2".
func SecureSecureWithOAuth2(ep goa.Endpoint) goa.Endpoint {
	reqs := []*security.Requirement{
		&security.Requirement{
			Schemes: []*security.Scheme{
				&security.Scheme{
					Kind:   security.SchemeKind(1),
					Name:   "implicit",
					Scopes: []string{"api:write", "api:read"},
					Flows: []*security.OAuthFlow{
						&security.OAuthFlow{
							Kind:             security.FlowKind(2),
							AuthorizationURL: "/authorization",
							RefreshURL:       "/refresh",
						},
					},
				},
				&security.Scheme{
					Kind:   security.SchemeKind(1),
					Name:   "authCode",
					Scopes: []string{"api:write", "api:read"},
					Flows: []*security.OAuthFlow{
						&security.OAuthFlow{
							Kind:             security.FlowKind(1),
							AuthorizationURL: "/authorization",
							TokenURL:         "/token",
							RefreshURL:       "/refresh",
						},
					},
				},
				&security.Scheme{
					Kind:   security.SchemeKind(1),
					Name:   "pass",
					Scopes: []string{"api:write", "api:read"},
					Flows: []*security.OAuthFlow{
						&security.OAuthFlow{
							Kind:       security.FlowKind(3),
							TokenURL:   "/token",
							RefreshURL: "/refresh",
						},
					},
				},
				&security.Scheme{
					Kind:   security.SchemeKind(1),
					Name:   "clicred",
					Scopes: []string{"api:write", "api:read"},
					Flows: []*security.OAuthFlow{
						&security.OAuthFlow{
							Kind:       security.FlowKind(4),
							TokenURL:   "/token",
							RefreshURL: "/refresh",
						},
					},
				},
			},
		},
	}
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		ctx = context.WithValue(ctx, security.ContextKey, reqs)
		return ep(ctx, req)
	}
}
`

var EndpointContextNoSecurityCode = `// SecureNoSecurity returns an endpoint function which initializes the context
// with the security requirements for the method "NoSecurity" of service
// "EndpointNoSecurity".
func SecureNoSecurity(ep goa.Endpoint) goa.Endpoint {
	reqs := []*security.Requirement{
		&security.Requirement{
			Schemes: []*security.Scheme{
				&security.Scheme{
					Kind: security.SchemeKind(5),
					Name: "",
				},
			},
		},
	}
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		ctx = context.WithValue(ctx, security.ContextKey, reqs)
		return ep(ctx, req)
	}
}
`
