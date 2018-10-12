package calc

import (
	"context"

	goalog "goa.design/goa/logging"
	calcsvc "goa.design/plugins/cors/examples/calc/gen/calc"
)

// calc service example implementation.
// The example methods log the requests and return zero values.
type calcSvc struct {
	logger goalog.Logger
}

// Required for compatibility with Service interface
func (s *calcSvc) GetLogger() goalog.Logger {
	return s.logger
}

// NewCalc returns the calc service implementation.
func NewCalc(logger goalog.Logger) calcsvc.Service {
	return &calcSvc{logger: logger}
}

// Add adds up the two integer parameters and returns the results.
func (s *calcSvc) Add(ctx context.Context, p *calcsvc.AddPayload) (res int, err error) {
	s.logger.Debug("calc.add")
	return
}
