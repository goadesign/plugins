package calc

import (
	"context"
	"log"

	calcsvc "goa.design/plugins/cors/examples/calc/gen/calc"
)

// calc service example implementation.
// The example methods log the requests and return zero values.
type calcsrvc struct {
	logger *log.Logger
}

// NewCalc returns the calc service implementation.
func NewCalc(logger *log.Logger) calcsvc.Service {
	return &calcsrvc{logger}
}

// Add adds up the two integer parameters and returns the results.
func (s *calcsrvc) Add(ctx context.Context, p *calcsvc.AddPayload) (res int, err error) {
	s.logger.Print("calc.add")
	return
}
