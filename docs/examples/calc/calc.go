package calcapi

import (
	"context"
	"log"

	calc "goa.design/plugins/v3/docs/examples/calc/gen/calc"
)

// calc service example implementation.
// The example methods log the requests and return zero values.
type calcsrvc struct {
	logger *log.Logger
}

// NewCalc returns the calc service implementation.
func NewCalc(logger *log.Logger) calc.Service {
	return &calcsrvc{logger}
}

// Add implements add.
func (s *calcsrvc) Add(ctx context.Context, p *calc.AddPayload, stream calc.AddServerStream) error {
	stream.Send(p.Left + p.Right)
	data, err := stream.Recv()
	for err == nil {
		stream.Send(data.A + data.B)
	}
	return stream.Close()
}
