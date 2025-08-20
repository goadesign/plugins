package design

import (
	. "goa.design/goa/v3/dsl"
	_ "goa.design/plugins/v3/testing"
)

// Calculator service demonstrates the testing plugin features with a simple API
var _ = Service("calculator", func() {
	Description("A simple calculator service to demonstrate testing plugin features")

	// Simple addition - demonstrates basic testing
	Method("add", func() {
		Description("Add two numbers")
		Payload(func() {
			Field(1, "a", Float64, "First number")
			Field(2, "b", Float64, "Second number")
			Required("a", "b")
		})
		Result(func() {
			Field(1, "result", Float64, "Sum of a and b")
			Field(2, "operation", String, "Operation performed")
			Required("result", "operation")
		})
		HTTP(func() {
			POST("/add")
		})
		GRPC(func() {})
	})

	// Division - demonstrates error handling
	Method("divide", func() {
		Description("Divide two numbers")
		Payload(func() {
			Field(1, "dividend", Float64, "Number to divide")
			Field(2, "divisor", Float64, "Number to divide by")
			Required("dividend", "divisor")
		})
		Result(func() {
			Field(1, "result", Float64, "Division result")
			Field(2, "operation", String, "Operation performed")
			Required("result", "operation")
		})
		Error("division_by_zero", CalculatorError, "Cannot divide by zero")
		HTTP(func() {
			POST("/divide")
			Response("division_by_zero", StatusBadRequest)
		})
		GRPC(func() {})
	})

	// Factorial - demonstrates timeout scenarios (can be slow for large numbers)
	Method("factorial", func() {
		Description("Calculate factorial of a number")
		Payload(func() {
			Field(1, "n", Int, "Number to calculate factorial of", func() {
				Minimum(0)
				Maximum(20) // Prevent overflow
			})
			Required("n")
		})
		Result(func() {
			Field(1, "result", Int64, "Factorial result")
			Field(2, "operation", String, "Operation performed")
			Field(3, "computation_time_ms", Int, "Time taken in milliseconds")
			Required("result", "operation", "computation_time_ms")
		})
		Error("invalid_input", CalculatorError, "Input must be between 0 and 20")
		HTTP(func() {
			POST("/factorial")
			Response("invalid_input", StatusBadRequest)
		})
		GRPC(func() {})
	})

	// Statistics - demonstrates complex validation scenarios
	Method("statistics", func() {
		Description("Calculate statistics for a list of numbers")
		Payload(func() {
			Field(1, "numbers", ArrayOf(Float64), "List of numbers", func() {
				MinLength(1) // At least one number required
			})
			Required("numbers")
		})
		Result(func() {
			Field(1, "mean", Float64, "Average of numbers")
			Field(2, "median", Float64, "Median value")
			Field(3, "min", Float64, "Minimum value")
			Field(4, "max", Float64, "Maximum value")
			Field(5, "count", Int, "Number of values")
			Field(6, "sum", Float64, "Sum of all values")
			Required("mean", "median", "min", "max", "count", "sum")
		})
		Error("empty_list", CalculatorError, "List cannot be empty")
		HTTP(func() {
			POST("/statistics")
			Response("empty_list", StatusBadRequest)
		})
		GRPC(func() {})
	})

	// Batch operations - demonstrates streaming
	Method("batch_add", func() {
		Description("Add multiple pairs of numbers using streaming")
		StreamingPayload(func() {
			Field(1, "a", Float64, "First number")
			Field(2, "b", Float64, "Second number")
			Required("a", "b")
		})
		StreamingResult(func() {
			Field(1, "result", Float64, "Sum")
			Field(2, "index", Int, "Operation index")
			Required("result", "index")
		})
		HTTP(func() {
			POST("/batch/add")
			Response(StatusOK)
		})
		GRPC(func() {})
	})
})

// CalculatorError is the error result type for calculator operations
var CalculatorError = Type("CalculatorError", func() {
	Description("Calculator operation error")
	Field(1, "message", String, "Error message")
	Field(2, "code", String, "Error code")
	Required("message", "code")
})