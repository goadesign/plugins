package calcapi

import (
	"context"

	calc "goa.design/plugins/v3/docs/examples/calc/gen/calc"
)

// calc service example implementation.
// The example methods log the requests and return zero values.
type calcsrvc struct{}

// NewCalc returns the calc service implementation.
func NewCalc() calc.Service {
	return &calcsrvc{}
}

// Add implements add.
func (s *calcsrvc) Add(ctx context.Context, p *calc.AddPayload, stream calc.AddServerStream) error {
	defer stream.Close()
	if err := stream.Send(p.Left + p.Right); err != nil {
		return err
	}
	data, err := stream.Recv()
	for err != nil {
		return err
	}
	return stream.Send(data.A + data.B)
}
