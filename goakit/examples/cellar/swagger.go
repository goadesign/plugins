package cellar

import (
	"github.com/go-kit/kit/log"
	"goa.design/plugins/goakit/examples/cellar/gen/swagger"
)

// swagger service example implementation.
// The example methods log the requests and return zero values.
type swaggersvc struct {
	logger log.Logger
}

// NewSwagger returns the swagger service implementation.
func NewSwagger(logger log.Logger) swagger.Service {
	return &swaggersvc{logger}
}
