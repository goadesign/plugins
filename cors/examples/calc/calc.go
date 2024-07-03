package calc

import (
	"context"

	"goa.design/clue/log"

	calcsvc "goa.design/plugins/v3/cors/examples/calc/gen/calc"
)

// calc service example implementation.
// The example methods log the requests and return zero values.
type calcsrvc struct {
}

// NewCalc returns the calc service implementation.
func NewCalc() calcsvc.Service {
	return &calcsrvc{}
}

// Add adds up the two integer parameters and returns the results.
func (s *calcsrvc) Add(ctx context.Context, p *calcsvc.AddPayload) (res int, err error) {
	log.Printf(ctx, "calc.add")
	return
}
