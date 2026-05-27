package calculatorapi

import (
	"context"
	"io"
	"testing"
	"time"

	calculator "goa.design/plugins/v3/testing/examples/calculator/gen/calculator"
	calculatortest "goa.design/plugins/v3/testing/examples/calculator/gen/calculator/calculatortest"
)

// RunCalculatorHarness exercises the generated harness against your service
// implementation.
// Call this helper from your test, passing your service implementation.
func RunCalculatorHarness(t *testing.T, svc calculator.Service) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	h := calculatortest.NewHarness(t, svc)
	defer h.Close()

	td := calculatortest.NewTestData()
	t.Run("add", func(t *testing.T) {
		result, err := h.Client.Add(ctx, td.ValidAddPayload())
		if err != nil {
			t.Errorf("add failed: %v", err)
		}
		if result == nil {
			t.Error("add returned nil result")
		}
	})
	t.Run("divide", func(t *testing.T) {
		result, err := h.Client.Divide(ctx, td.ValidDividePayload())
		if err != nil {
			t.Errorf("divide failed: %v", err)
		}
		if result == nil {
			t.Error("divide returned nil result")
		}
	})
	t.Run("factorial", func(t *testing.T) {
		result, err := h.Client.Factorial(ctx, td.ValidFactorialPayload())
		if err != nil {
			t.Errorf("factorial failed: %v", err)
		}
		if result == nil {
			t.Error("factorial returned nil result")
		}
	})
	t.Run("statistics", func(t *testing.T) {
		result, err := h.Client.Statistics(ctx, td.ValidStatisticsPayload())
		if err != nil {
			t.Errorf("statistics failed: %v", err)
		}
		if result == nil {
			t.Error("statistics returned nil result")
		}
	})
	t.Run("batch_add_Stream", func(t *testing.T) {
		stream, err := h.Client.BatchAdd(ctx)
		if err != nil {
			t.Errorf("Failed to create batch_add stream: %v", err)
		}
		if stream == nil {
			t.Fatal("batch_add returned nil stream")
		}
		// Bidirectional stream - send and receive
		// Send a test message
		payload := td.ValidBatchAddPayload()
		if err := stream.Send(payload); err != nil {
			t.Errorf("batch_add send failed: %v", err)
		}

		// Try to receive a response
		_, err = stream.Recv()
		if err != nil && err != io.EOF {
			t.Errorf("batch_add recv failed: %v", err)
		}

		// Close the stream
		if err := stream.Close(); err != nil {
			t.Errorf("batch_add close failed: %v", err)
		}
	})
}
