package calc

import (
	"context"

	calcsvc "goa.design/plugins/v3/zaplogger/examples/calc/gen/calc"
	log "goa.design/plugins/v3/zaplogger/examples/calc/gen/log"
)

// calc service example implementation.
// The example methods log the requests and return zero values.
type calcSvc struct {
	logger *log.Logger
}

// NewCalc returns the calc service implementation.
func NewCalc(logger *log.Logger) calcsvc.Service {
	return &calcSvc{logger}
}

// Add implements add.
func (s *calcSvc) Add(ctx context.Context, p *calcsvc.AddPayload) (res int, err error) {
	s.logger.Info("calc.add")
	return
}
