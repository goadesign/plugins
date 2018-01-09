package calc

import (
	"context"
	"log"

	calcsvc "goa.design/plugins/cors/examples/calc/gen/calc"
)

// calc service example implementation.
// The example methods log the requests and return zero values.
type calcsvcSvc struct {
	logger *log.Logger
}

// NewCalc returns the calc service implementation.
func NewCalc(logger *log.Logger) calcsvc.Service {
	return &calcsvcSvc{logger}
}

// Add adds up the two integer parameters and returns the results.
func (s *calcsvcSvc) Add(ctx context.Context, p *calcsvc.AddPayload) (int, error) {
	var res int
	s.logger.Print("calc.add")
	return res, nil
}
