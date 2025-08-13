package calcapi

import (
	"context"

	calc "goa.design/plugins/v3/docs/examples/calc/gen/calc"
)

// calc service example implementation.
// The example methods log the requests and return zero values.
type calcsrvc struct{}

// NewCalc returns the calc service implementation.
func NewCalc() calc.Service {
	return &calcsrvc{}
}

// Add implements add.
func (s *calcsrvc) Add(ctx context.Context, p *calc.AddPayload) (int, error) {
	return p.Left + p.Right, nil
}
