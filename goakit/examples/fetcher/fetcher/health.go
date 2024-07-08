package fetcherapi

import (
	"context"

	"github.com/go-kit/log"
	health "goa.design/plugins/v3/goakit/examples/fetcher/fetcher/gen/health"
)

// health service example implementation.
// The example methods log the requests and return zero values.
type healthsrvc struct {
	logger log.Logger
}

// NewHealth returns the health service implementation.
func NewHealth(logger log.Logger) health.Service {
	return &healthsrvc{
		logger: logger,
	}
}

// Health check endpoint
func (s *healthsrvc) Show(ctx context.Context) (res string, err error) {
	s.logger.Log("service", "health", "method", "show")
	return
}
