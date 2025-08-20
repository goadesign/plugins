package testvalidators

import (
	"math"
	"testing"

	calculator "goa.design/plugins/v3/testing/examples/calculator/gen/calculator"
)

// ValidatePrecision checks if the division result has acceptable floating-point precision
func ValidatePrecision(t *testing.T, result *calculator.DivideResult, expected map[string]any) {
	// For 100 / 3, we expect approximately 33.333...
	expectedValue := 100.0 / 3.0
	tolerance := 0.0001

	if math.Abs(result.Result-expectedValue) > tolerance {
		t.Errorf("precision error: expected ~%f, got %f", expectedValue, result.Result)
	}
}

// ValidateFactorialTiming checks if factorial computation time is reasonable
func ValidateFactorialTiming(t *testing.T, result *calculator.FactorialResult, expected map[string]any) {
	// For factorial of 15, computation should be very fast (< 100ms)
	if result.ComputationTimeMs > 100 {
		t.Errorf("computation took too long: %d ms", result.ComputationTimeMs)
	}

	// Verify the factorial is correct
	expectedFactorial := int64(1307674368000) // 15!
	if result.Result != expectedFactorial {
		t.Errorf("incorrect factorial: expected %d, got %d", expectedFactorial, result.Result)
	}
}

// ValidateStatisticsAccuracy validates the statistical calculations
func ValidateStatisticsAccuracy(t *testing.T, result *calculator.StatisticsResult, expected map[string]any) {
	// For numbers 1-10, validate the statistics
	expectedMean := 5.5
	expectedMedian := 5.5
	expectedSum := 55.0

	tolerance := 0.0001

	if math.Abs(result.Mean-expectedMean) > tolerance {
		t.Errorf("mean error: expected %f, got %f", expectedMean, result.Mean)
	}

	if math.Abs(result.Median-expectedMedian) > tolerance {
		t.Errorf("median error: expected %f, got %f", expectedMedian, result.Median)
	}

	if math.Abs(result.Sum-expectedSum) > tolerance {
		t.Errorf("sum error: expected %f, got %f", expectedSum, result.Sum)
	}

	if result.Count != 10 {
		t.Errorf("count error: expected 10, got %d", result.Count)
	}
}

// ValidateOutliers checks if statistics are calculated correctly with outliers
func ValidateOutliers(t *testing.T, result *calculator.StatisticsResult, expected map[string]any) {
	// For [1, 2, 3, 4, 100], the median should be 3 (middle value)
	expectedMedian := 3.0

	if result.Median != expectedMedian {
		t.Errorf("median with outlier incorrect: expected %f, got %f", expectedMedian, result.Median)
	}

	// The mean should be (1+2+3+4+100)/5 = 22
	expectedMean := 22.0
	if math.Abs(result.Mean-expectedMean) > 0.0001 {
		t.Errorf("mean with outlier incorrect: expected %f, got %f", expectedMean, result.Mean)
	}
}

// ValidateFloatingPoint validates floating-point arithmetic edge cases
func ValidateFloatingPoint(t *testing.T, result *calculator.AddResult, expected map[string]any) {
	// The classic 0.1 + 0.2 != 0.3 problem
	// In Go, 0.1 + 0.2 will be approximately 0.30000000000000004

	// Check if the result is within acceptable tolerance of 0.3
	expectedValue := 0.3
	tolerance := 0.000001

	if math.Abs(result.Result-expectedValue) > tolerance {
		t.Errorf("floating point precision: expected ~%f, got %f", expectedValue, result.Result)
	}
}
