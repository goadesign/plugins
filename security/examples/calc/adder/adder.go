package adder

import (
	"context"
	"log"

	"goa.design/plugins/security"
	addersvc "goa.design/plugins/security/examples/calc/adder/gen/adder"
)

// adder service example implementation.
// The example methods log the requests and return zero values.
type adderSvc struct {
	logger *log.Logger
}

// NewAdder returns the adder service implementation.
func NewAdder(logger *log.Logger) addersvc.Service {
	return &adderSvc{logger}
}

// This action returns the sum of two integers and is secured with the API key
// scheme
func (s *adderSvc) Add(ctx context.Context, p *addersvc.AddPayload) (int, error) {
	return p.A + p.B, nil
}

// AdderAuthAPIKeyFn implements the authorization logic for APIKey scheme. It
// must return one of the following errors
// * addersvc.Unauthorized
// * error
func AdderAuthAPIKeyFn(ctx context.Context, key string, s *security.APIKeyScheme) (context.Context, error) {
	if key == "" {
		return ctx, addersvc.Unauthorized("invalid key")
	}
	return ctx, nil
}
