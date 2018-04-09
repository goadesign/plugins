package multiauth

import (
	"context"
	"fmt"
	"log"

	"goa.design/plugins/security"
	securedservice "goa.design/plugins/security/examples/multi_auth/gen/secured_service"
)

// secured_service service example implementation.
// The example methods log the requests and return zero values.
type securedServiceSvc struct {
	logger *log.Logger
}

// NewSecuredService returns the secured_service service implementation.
func NewSecuredService(logger *log.Logger) securedservice.Service {
	return &securedServiceSvc{logger}
}

// Creates a valid JWT
func (s *securedServiceSvc) Signin(ctx context.Context, p *securedservice.SigninPayload) error {
	s.logger.Print("securedService.signin")
	return nil
}

// This action is secured with the jwt scheme
func (s *securedServiceSvc) Secure(ctx context.Context, p *securedservice.SecurePayload) (string, error) {
	var res string
	s.logger.Print("securedService.secure")
	return res, nil
}

// This action is secured with the jwt scheme and also requires an API key
// query string.
func (s *securedServiceSvc) DoublySecure(ctx context.Context, p *securedservice.DoublySecurePayload) (string, error) {
	var res string
	s.logger.Print("securedService.doubly_secure")
	return res, nil
}

// This action is secured with the jwt scheme and also requires an API key
// header.
func (s *securedServiceSvc) AlsoDoublySecure(ctx context.Context, p *securedservice.AlsoDoublySecurePayload) (string, error) {
	var res string
	s.logger.Print("securedService.also_doubly_secure")
	return res, nil
}

// SecuredServiceAuthBasicAuthFn implements the authorization logic for
// BasicAuth scheme. It must return one of the following errors
// * securedservice.MakeUnauthorized
// * securedservice.Forbidden
// * error
func SecuredServiceAuthBasicAuthFn(ctx context.Context, user, pass string, s *security.BasicAuthScheme) (context.Context, error) {
	// Add authorization logic
	return ctx, fmt.Errorf("not implemented")
}

// SecuredServiceAuthJWTFn implements the authorization logic for JWT scheme.
// It must return one of the following errors
// * securedservice.MakeUnauthorized
// * securedservice.Forbidden
// * error
func SecuredServiceAuthJWTFn(ctx context.Context, token string, s *security.JWTScheme) (context.Context, error) {
	// Add authorization logic
	return ctx, fmt.Errorf("not implemented")
}

// SecuredServiceAuthAPIKeyFn implements the authorization logic for APIKey
// scheme. It must return one of the following errors
// * securedservice.MakeUnauthorized
// * securedservice.Forbidden
// * error
func SecuredServiceAuthAPIKeyFn(ctx context.Context, key string, s *security.APIKeyScheme) (context.Context, error) {
	// Add authorization logic
	return ctx, fmt.Errorf("not implemented")
}

// SecuredServiceAuthOAuth2Fn implements the authorization logic for OAuth2
// scheme. It must return one of the following errors
// * securedservice.MakeUnauthorized
// * securedservice.Forbidden
// * error
func SecuredServiceAuthOAuth2Fn(ctx context.Context, token string, s *security.OAuth2Scheme) (context.Context, error) {
	// Add authorization logic
	return ctx, fmt.Errorf("not implemented")
}
