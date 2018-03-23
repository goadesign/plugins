package adder

import (
	"context"
	"log"

	"goa.design/plugins/security"
	addersvc "goa.design/plugins/security/examples/calc/adder/gen/adder"
)

// adder service example implementation.
// The example methods log the requests and return zero values.
type addersvcSvc struct {
	logger *log.Logger
}

// NewAdder returns the adder service implementation.
func NewAdder(logger *log.Logger) addersvc.Service {
	return &addersvcSvc{logger}
}

// This action returns the sum of two integers and is secured with the API key
// scheme
func (s *addersvcSvc) Add(ctx context.Context, p *addersvc.AddPayload) (int, error) {
	return p.A + p.B, nil
}

// AdderAuthAPIKeyFn implements the authorization logic for APIKey scheme.
func AdderAuthAPIKeyFn(ctx context.Context, key string, s *security.APIKeyScheme) (context.Context, error) {
	// Add authorization logic
	if key == "" {
		return ctx, addersvc.Unauthorized("invalid key")
	}
	return ctx, nil
}
