package calc

import (
	"context"
	"log"

	"goa.design/plugins/security"
	adder "goa.design/plugins/security/examples/calc/adder/gen/adder"
)

// adder service example implementation.
// The example methods log the requests and return zero values.
type adderSvc struct {
	logger *log.Logger
}

// NewAdder returns the adder service implementation.
func NewAdder(logger *log.Logger) adder.Service {
	return &adderSvc{logger}
}

// This action returns the sum of two integers and is secured with the API key
// scheme.
func (s *adderSvc) Add(ctx context.Context, p *adder.AddPayload) (int, error) {
	return p.A + p.B, nil
}

// AuthAPIKeyFn implements the authorization logic for APIKey scheme.
func AuthAPIKeyFn(ctx context.Context, key string, s *security.APIKeyScheme) (context.Context, error) {
	// Add authorization logic
	return ctx, nil
}
