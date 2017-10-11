package multiauth

import (
	"context"
	"log"

	securedservice "goa.design/plugins/security/examples/multi_auth/gen/secured_service"
)

// secured_service service example implementation.
// The example methods log the requests and return zero values.
type securedserviceSvc struct {
	logger *log.Logger
}

// NewSecuredService returns the secured_service service implementation.
func NewSecuredService(logger *log.Logger) securedservice.Service {
	return &securedserviceSvc{logger}
}

// Creates a valid JWT
func (s *securedserviceSvc) Signin(ctx context.Context, p *securedservice.SigninPayload) (string, error) {
	var res string
	s.logger.Print("secured_service.signin")
	return res, nil
}

// This action is secured with the jwt scheme
func (s *securedserviceSvc) Secure(ctx context.Context, p *securedservice.SecurePayload) (string, error) {
	var res string
	s.logger.Print("secured_service.secure")
	return res, nil
}

// This action is secured with the jwt scheme and also requires an API key
// query string.
func (s *securedserviceSvc) DoublySecure(ctx context.Context, p *securedservice.DoublySecurePayload) (string, error) {
	var res string
	s.logger.Print("secured_service.doubly_secure")
	return res, nil
}

// This action is secured with the jwt scheme and also requires an API key
// header.
func (s *securedserviceSvc) AlsoDoublySecure(ctx context.Context, p *securedservice.AlsoDoublySecurePayload) (string, error) {
	var res string
	s.logger.Print("secured_service.also_doubly_secure")
	return res, nil
}
