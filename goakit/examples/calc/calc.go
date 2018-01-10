package calc

import (
	"context"

	"github.com/go-kit/kit/log"
	calcsvc "goa.design/plugins/goakit/examples/calc/gen/calc"
)

// calc service example implementation.
// The example methods log the requests and return zero values.
type calcsvcsvc struct {
	logger log.Logger
}

// NewCalc returns the calc service implementation.
func NewCalc(logger log.Logger) calcsvc.Service {
	return &calcsvcsvc{logger}
}

// Add adds up the two integer parameters and returns the results.
func (s *calcsvcsvc) Add(ctx context.Context, p *calcsvc.AddPayload) (int, error) {
	var res int
	s.logger.Log("msg", "calc.add")
	return res, nil
}
