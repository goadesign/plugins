package calc

import (
	"context"

	"github.com/go-kit/log"
	calcsvc "goa.design/plugins/v3/goakit/examples/calc/gen/calc"
)

// calc service example implementation.
// The example methods log the requests and return zero values.
type calcSvc struct {
	logger log.Logger
}

// NewCalc returns the calc service implementation.
func NewCalc(logger log.Logger) calcsvc.Service {
	return &calcSvc{logger}
}

// Add adds up the two integer parameters and returns the results.
func (s *calcSvc) Add(ctx context.Context, p *calcsvc.AddPayload) (res int, err error) {
	s.logger.Log("info", "calc.add")
	return
}
