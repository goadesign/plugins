package calcapi

import (
	"context"

	"github.com/go-kit/log"
	calc "goa.design/plugins/v3/goakit/examples/calc/gen/calc"
)

// calc service example implementation.
// The example methods log the requests and return zero values.
type calcsrvc struct {
	logger log.Logger
}

// NewCalc returns the calc service implementation.
func NewCalc(logger log.Logger) calc.Service {
	return &calcsrvc{
		logger: logger,
	}
}

// Add adds up the two integer parameters and returns the results.
func (s *calcsrvc) Add(ctx context.Context, p *calc.AddPayload) (res int, err error) {
	s.logger.Log("service", "calc", "method", "add")
	return
}
