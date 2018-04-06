package testdata

var EndpointInitWithoutRequirementCode = `// NewSecureEndpoints wraps the methods of a EndpointWithoutRequirement service
// with security scheme aware endpoints.
func NewSecureEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Unsecure: NewUnsecureEndpoint(s),
	}
}
`

var EndpointInitWithRequirementsCode = `// NewSecureEndpoints wraps the methods of a EndpointsWithRequirements service
// with security scheme aware endpoints.
func NewSecureEndpoints(s Service, authBasicAuthFn security.AuthorizeBasicAuthFunc, authJWTFn security.AuthorizeJWTFunc) *Endpoints {
	return &Endpoints{
		SecureWithRequirements:       SecureSecureWithRequirements(NewSecureWithRequirementsEndpoint(s), authBasicAuthFn),
		DoublySecureWithRequirements: SecureDoublySecureWithRequirements(NewDoublySecureWithRequirementsEndpoint(s), authBasicAuthFn, authJWTFn),
	}
}
`

var EndpointInitWithServiceRequirementsCode = `// NewSecureEndpoints wraps the methods of a EndpointsWithServiceRequirements
// service with security scheme aware endpoints.
func NewSecureEndpoints(s Service, authBasicAuthFn security.AuthorizeBasicAuthFunc) *Endpoints {
	return &Endpoints{
		SecureWithRequirements:     SecureSecureWithRequirements(NewSecureWithRequirementsEndpoint(s), authBasicAuthFn),
		AlsoSecureWithRequirements: SecureAlsoSecureWithRequirements(NewAlsoSecureWithRequirementsEndpoint(s), authBasicAuthFn),
	}
}
`

var EndpointInitNoSecurityCode = `// NewSecureEndpoints wraps the methods of a EndpointNoSecurity service with
// security scheme aware endpoints.
func NewSecureEndpoints(s Service) *Endpoints {
	return &Endpoints{
		NoSecurity: NewNoSecurityEndpoint(s),
	}
}
`

var EndpointWithRequiredScopesCode = `// SecureSecureWithRequiredScopes returns an endpoint function which
// initializes the context with the security requirements for the method
// "SecureWithRequiredScopes" of service "EndpointWithRequiredScopes".
func SecureSecureWithRequiredScopes(ep goa.Endpoint, authJWTFn security.AuthorizeJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*SecureWithRequiredScopesPayload)
		var err error
		jwtSch := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"api:read", "api:write"},
			RequiredScopes: []string{"api:read", "api:write"},
		}
		ctx, err = authJWTFn(ctx, *p.Token, &jwtSch)
		if err != nil {
			return nil, err
		}
		return ep(ctx, req)
	}
}
`

var EndpointWithAPIKeyOverrideCode = `// SecureSecureWithAPIKeyOverride returns an endpoint function which
// initializes the context with the security requirements for the method
// "SecureWithAPIKeyOverride" of service "EndpointWithAPIKeyOverride".
func SecureSecureWithAPIKeyOverride(ep goa.Endpoint, authAPIKeyFn security.AuthorizeAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*SecureWithAPIKeyOverridePayload)
		var err error
		apiKeySch := security.APIKeyScheme{
			Name: "api_key",
		}
		ctx, err = authAPIKeyFn(ctx, *p.Key, &apiKeySch)
		if err != nil {
			return nil, err
		}
		return ep(ctx, req)
	}
}
`

var EndpointWithOAuth2Code = `// SecureSecureWithOAuth2 returns an endpoint function which initializes the
// context with the security requirements for the method "SecureWithOAuth2" of
// service "EndpointWithOAuth2".
func SecureSecureWithOAuth2(ep goa.Endpoint, authOAuth2Fn security.AuthorizeOAuth2Func) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*SecureWithOAuth2Payload)
		var err error
		oauth2Sch := security.OAuth2Scheme{
			Name:           "authCode",
			Scopes:         []string{"api:write", "api:read"},
			RequiredScopes: []string{},
			Flows: []*security.OAuthFlow{
				&security.OAuthFlow{
					Type:             "authorization_code",
					AuthorizationURL: "/authorization",
					TokenURL:         "/token",
					RefreshURL:       "/refresh",
				},
			},
		}
		ctx, err = authOAuth2Fn(ctx, *p.Token, &oauth2Sch)
		if err != nil {
			return nil, err
		}
		return ep(ctx, req)
	}
}
`

var SingleServiceAuthFuncsCode = `// SingleServiceAuthAPIKeyFn implements the authorization logic for APIKey
// scheme.
func SingleServiceAuthAPIKeyFn(ctx context.Context, key string, s *security.APIKeyScheme) (context.Context, error) {
	// Add authorization logic
	return ctx, fmt.Errorf("internal error")
}
`

var MultipleServicesAuth1FuncsCode = `// ServiceWithAPIKeyAuthAuthAPIKeyFn implements the authorization logic for
// APIKey scheme.
func ServiceWithAPIKeyAuthAuthAPIKeyFn(ctx context.Context, key string, s *security.APIKeyScheme) (context.Context, error) {
	// Add authorization logic
	return ctx, fmt.Errorf("internal error")
}
`

var MultipleServicesAuth2FuncsCode = `// ServiceWithJWTAndAPIKeyAuthAPIKeyFn implements the authorization logic for
// APIKey scheme.
func ServiceWithJWTAndAPIKeyAuthAPIKeyFn(ctx context.Context, key string, s *security.APIKeyScheme) (context.Context, error) {
	// Add authorization logic
	return ctx, fmt.Errorf("internal error")
}

// ServiceWithJWTAndAPIKeyAuthJWTFn implements the authorization logic for JWT
// scheme.
func ServiceWithJWTAndAPIKeyAuthJWTFn(ctx context.Context, token string, s *security.JWTScheme) (context.Context, error) {
	// Add authorization logic
	return ctx, fmt.Errorf("internal error")
}
`
