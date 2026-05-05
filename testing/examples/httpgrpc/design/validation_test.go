package design

import (
	. "goa.design/goa/v3/dsl"
)

// TestValidation service to test edge case generation based on validation rules
var _ = Service("test-validation", func() {
	Description("Test service with validation rules for edge case generation")

	Method("validate_strings", func() {
		Payload(func() {
			Field(1, "short", String, func() {
				MinLength(3)
				MaxLength(10)
				Description("String with min 3, max 10 chars")
			})
			Field(2, "pattern", String, func() {
				Pattern("^[A-Z][a-z]+$")
				Description("String matching pattern")
			})
			Field(3, "email", String, func() {
				Format(FormatEmail)
				Description("Email format")
			})
			Required("short", "pattern", "email")
		})
		Result(String)
		HTTP(func() { POST("/validate/strings") })
	})

	Method("validate_numbers", func() {
		Payload(func() {
			Field(1, "age", Int, func() {
				Minimum(0)
				Maximum(120)
				Description("Age between 0 and 120")
			})
			Field(2, "price", Float64, func() {
				Minimum(0.01)
				Maximum(999999.99)
				Description("Price between 0.01 and 999999.99")
			})
			Field(3, "quantity", Int, func() {
				Minimum(1)
				Description("Quantity at least 1")
			})
			Required("age", "price", "quantity")
		})
		Result(String)
		HTTP(func() { POST("/validate/numbers") })
	})

	Method("validate_arrays", func() {
		Payload(func() {
			Field(1, "tags", ArrayOf(String), func() {
				MinLength(1)
				MaxLength(5)
				Description("Array with 1 to 5 tags")
			})
			Field(2, "scores", ArrayOf(Int), func() {
				MinLength(3)
				Description("Array with at least 3 scores")
			})
			Required("tags", "scores")
		})
		Result(String)
		HTTP(func() { POST("/validate/arrays") })
	})

	Method("validate_optional", func() {
		Payload(func() {
			Field(1, "required_field", String, func() {
				MinLength(1)
			})
			Field(2, "optional_field", String, func() {
				MinLength(5)
				Description("Optional field with validation when present")
			})
			Field(3, "optional_number", Int, func() {
				Minimum(10)
				Maximum(100)
			})
			Required("required_field")
		})
		Result(String)
		HTTP(func() { POST("/validate/optional") })
	})
})
