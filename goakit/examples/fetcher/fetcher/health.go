package fetcher

import (
	"context"

	"github.com/go-kit/kit/log"
	"goa.design/plugins/goakit/examples/client/fetcher/gen/health"
)

// health service example implementation.
// The example methods log the requests and return zero values.
type healthsvc struct {
	logger log.Logger
}

// NewHealth returns the health service implementation.
func NewHealth(logger log.Logger) health.Service {
	return &healthsvc{logger}
}

// Health check endpoint
func (s *healthsvc) Show(ctx context.Context) (string, error) {
	var res string
	s.logger.Log("msg", "health.show")
	return res, nil
}
