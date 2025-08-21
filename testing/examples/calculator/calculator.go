package calculatorapi

import (
	"context"

	"goa.design/clue/log"
	calculator "goa.design/plugins/v3/testing/examples/calculator/gen/calculator"
)

// calculator service example implementation.
// The example methods log the requests and return zero values.
type calculatorsrvc struct{}

// NewCalculator returns the calculator service implementation.
func NewCalculator() calculator.Service {
	return &calculatorsrvc{}
}

// Add two numbers
func (s *calculatorsrvc) Add(ctx context.Context, p *calculator.AddPayload) (res *calculator.AddResult, err error) {
	log.Printf(ctx, "calculator.add: %v + %v", p.A, p.B)
	return &calculator.AddResult{
		Result:    p.A + p.B,
		Operation: "add",
	}, nil
}

// Divide two numbers
func (s *calculatorsrvc) Divide(ctx context.Context, p *calculator.DividePayload) (res *calculator.DivideResult, err error) {
	log.Printf(ctx, "calculator.divide: %v / %v", p.Dividend, p.Divisor)
	if p.Divisor == 0 {
		return nil, &calculator.CalculatorError{Message: "division by zero", Code: "division_by_zero"}
	}
	return &calculator.DivideResult{
		Result:    p.Dividend / p.Divisor,
		Operation: "divide",
	}, nil
}

// Calculate factorial of a number
func (s *calculatorsrvc) Factorial(ctx context.Context, p *calculator.FactorialPayload) (res *calculator.FactorialResult, err error) {
	log.Printf(ctx, "calculator.factorial: %v", p.N)
	var result int64 = 1
	for i := 2; i <= p.N; i++ {
		result *= int64(i)
	}
	return &calculator.FactorialResult{
		Result:            result,
		Operation:         "factorial",
		ComputationTimeMs: 0,
	}, nil
}

// Calculate statistics for a list of numbers
func (s *calculatorsrvc) Statistics(ctx context.Context, p *calculator.StatisticsPayload) (res *calculator.StatisticsResult, err error) {
	log.Printf(ctx, "calculator.statistics: %v", p.Numbers)
	if len(p.Numbers) == 0 {
		return &calculator.StatisticsResult{}, nil
	}

	var sum, min, max float64
	min = p.Numbers[0]
	max = p.Numbers[0]

	for _, n := range p.Numbers {
		sum += n
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}

	mean := sum / float64(len(p.Numbers))

	// Calculate median
	sorted := make([]float64, len(p.Numbers))
	copy(sorted, p.Numbers)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	var median float64
	if len(sorted)%2 == 0 {
		median = (sorted[len(sorted)/2-1] + sorted[len(sorted)/2]) / 2
	} else {
		median = sorted[len(sorted)/2]
	}

	return &calculator.StatisticsResult{
		Mean:   mean,
		Median: median,
		Min:    min,
		Max:    max,
		Count:  len(p.Numbers),
		Sum:    sum,
	}, nil
}

// Add multiple pairs of numbers using streaming
func (s *calculatorsrvc) BatchAdd(ctx context.Context, stream calculator.BatchAddServerStream) (err error) {
	log.Printf(ctx, "calculator.batch_add")
	index := 0
	for {
		payload, err := stream.Recv()
		if err != nil {
			return err
		}
		result := &calculator.BatchAddResult{
			Result: payload.A + payload.B,
			Index:  index,
		}
		index++
		if err := stream.Send(result); err != nil {
			return err
		}
	}
}
