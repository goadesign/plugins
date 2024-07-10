package calcapi

import (
	"context"

	"go.uber.org/zap"
	calc "goa.design/plugins/v3/zaplogger/examples/calc/gen/calc"
)

// calc service example implementation.
// The example methods log the requests and return zero values.
type calcsrvc struct {
	logger *zap.SugaredLogger
}

// NewCalc returns the calc service implementation.
func NewCalc(logger *zap.SugaredLogger) calc.Service {
	return &calcsrvc{
		logger: logger,
	}
}

// Add implements add.
func (s *calcsrvc) Add(ctx context.Context, p *calc.AddPayload) (res int, err error) {
	s.logger.Info("calc.add")
	return
}
