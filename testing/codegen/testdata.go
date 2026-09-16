// This file renders valid payload builders and validation-edge examples for
// generated service tests.
package codegen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/expr"
)

// generateTestData generates test data generators for a service.
// This generates payload generators only, not result generators, because:
// - Tests need valid inputs (payloads) to send to the service
// - The service implementation produces the results, not the test
// - Tests verify the actual behavior of the service, not mock results
func generateTestData(genpkg string, svcData *service.Data, root *expr.RootExpr, svc *expr.ServiceExpr) *codegen.File {
	if svcData == nil {
		return nil
	}

	path := filepath.Join(testingPath(svcData), "testdata.go")

	specs := []*codegen.ImportSpec{
		{Path: "encoding/json"},
		{Path: filepath.Join(genpkg, svcData.PathName), Name: svcData.PkgName},
	}

	// Build per-method test data matching the template expectations
	methods := make([]*methodTestData, 0, len(svc.Methods))
	examples := expr.NewExampleGenerator(root.API.RandomizerFactory)
	for _, m := range svc.Methods {
		mtd := &methodTestData{
			Name:      m.Name,
			EdgeCases: []*edgeCaseData{},
		}

		// Payload info (non-streaming)
		if m.Payload != nil && m.Payload.Type != expr.Empty {
			payloadExamples := examples.At(expr.MethodPayloadExampleIdentity(m))
			ref := svcData.Scope.GoFullTypeRef(m.Payload, svcData.PkgName)
			mtd.Payload = true
			mtd.PayloadRef = ref
			// Generate JSON example init if available
			if ex := m.Payload.Example(payloadExamples); ex != nil {
				mtd.PayloadEx = true
				if b, err := json.Marshal(ex); err == nil {
					var buf bytes.Buffer
					_, _ = buf.Write(b)
					mtd.PayloadInit = buf.String()
				}
			}
			// Generate edge cases based on validation rules
			mtd.EdgeCases = generateEdgeCases(m.Name, m.Payload, payloadExamples)
		}

		// Streaming payload info
		if m.StreamingPayload != nil && m.StreamingPayload.Type != expr.Empty {
			streamingExamples := examples.At(expr.MethodStreamingPayloadExampleIdentity(m))
			ref := svcData.Scope.GoFullTypeRef(m.StreamingPayload, svcData.PkgName)
			mtd.Payload = true
			mtd.PayloadRef = ref
			if ex := m.StreamingPayload.Example(streamingExamples); ex != nil {
				mtd.PayloadEx = true
				if b, err := json.Marshal(ex); err == nil {
					var buf bytes.Buffer
					_, _ = buf.Write(b)
					mtd.PayloadInit = buf.String()
				}
			}
			// Generate edge cases for streaming payloads
			if len(mtd.EdgeCases) == 0 {
				mtd.EdgeCases = generateEdgeCases(m.Name, m.StreamingPayload, streamingExamples)
			}
		}

		methods = append(methods, mtd)
	}

	data := struct {
		Service *service.Data
		Methods []*methodTestData
	}{
		Service: svcData,
		Methods: methods,
	}

	sections := []*codegen.SectionTemplate{
		codegen.Header(fmt.Sprintf("Test data generators for %s service", svc.Name), svcData.PathName+"test", specs),
		{
			Name:   "testdata-generators",
			Source: testingTemplates.Read(testdataGeneratorsT),
			FuncMap: template.FuncMap{
				"goify": codegen.Goify,
			},
			Data: data,
		},
	}

	return &codegen.File{
		Path:             path,
		SectionTemplates: sections,
	}
}

// methodTestData contains data for generating test data for a method.
type methodTestData struct {
	Name        string
	Payload     bool
	PayloadEx   bool
	PayloadRef  string
	PayloadInit string
	EdgeCases   []*edgeCaseData
}

// edgeCaseData contains data for generating edge case test data.
type edgeCaseData struct {
	Name        string // e.g., "MinValues", "MaxValues", "ZeroDivisor"
	Description string // e.g., "with minimum valid values"
	Init        string // JSON initialization for the edge case
}

// generateEdgeCases generates edge case test data based on validation rules.
func generateEdgeCases(_ string, attr *expr.AttributeExpr, gen *expr.ExampleGenerator) []*edgeCaseData {
	var cases []*edgeCaseData

	// Note: We intentionally don't generate error-triggering cases here
	// because the design doesn't specify what input values cause errors.
	// For example, even if a method has a "division_by_zero" error defined,
	// we can't know that divisor=0 causes it without looking at the implementation.
	// Users should write explicit test scenarios for business logic errors.

	// Generate min/max value cases based on validation
	minCase := generateMinValues(attr, gen)
	if minCase != nil {
		if b, err := json.Marshal(minCase); err == nil {
			cases = append(cases, &edgeCaseData{
				Name:        "MinValues",
				Description: "with minimum valid values",
				Init:        string(b),
			})
		}
	}

	maxCase := generateMaxValues(attr, gen)
	if maxCase != nil {
		if b, err := json.Marshal(maxCase); err == nil {
			cases = append(cases, &edgeCaseData{
				Name:        "MaxValues",
				Description: "with maximum valid values",
				Init:        string(b),
			})
		}
	}

	// Generate min/max length cases for strings and arrays
	minLengthCase := generateMinLengthValues(attr, gen)
	if minLengthCase != nil {
		if b, err := json.Marshal(minLengthCase); err == nil {
			cases = append(cases, &edgeCaseData{
				Name:        "MinLength",
				Description: "with minimum length strings/arrays",
				Init:        string(b),
			})
		}
	}

	maxLengthCase := generateMaxLengthValues(attr, gen)
	if maxLengthCase != nil {
		if b, err := json.Marshal(maxLengthCase); err == nil {
			cases = append(cases, &edgeCaseData{
				Name:        "MaxLength",
				Description: "with maximum length strings/arrays",
				Init:        string(b),
			})
		}
	}

	return cases
}

// generateMinValues generates an example with all numeric fields at minimum values.
func generateMinValues(attr *expr.AttributeExpr, gen *expr.ExampleGenerator) any {
	// AsObject handles UserTypes automatically
	if obj := expr.AsObject(attr.Type); obj != nil {
		// For objects, check each field for numeric validations
		result := make(map[string]any)
		hasMinValues := false

		for _, att := range *obj {
			if att.Attribute.Validation != nil &&
				(att.Attribute.Validation.Minimum != nil || att.Attribute.Validation.ExclusiveMinimum != nil) {
				// This field has min validation
				if minVal := generateMinValues(att.Attribute, gen); minVal != nil {
					result[att.Name] = minVal
					hasMinValues = true
				} else if example := att.Attribute.Example(gen); example != nil {
					result[att.Name] = example
				}
			} else if nested := generateMinValues(att.Attribute, gen); nested != nil {
				// Check if nested objects have min values
				result[att.Name] = nested
				hasMinValues = true
			} else if example := att.Attribute.Example(gen); example != nil {
				result[att.Name] = example
			}
		}

		if hasMinValues {
			return result
		}
		return nil
	}

	// Handle primitive numeric types
	switch attr.Type.Kind() {
	case expr.IntKind, expr.Int32Kind, expr.Int64Kind,
		expr.UIntKind, expr.UInt32Kind, expr.UInt64Kind,
		expr.Float32Kind, expr.Float64Kind:
		if attr.Validation != nil {
			if attr.Validation.Minimum != nil {
				return *attr.Validation.Minimum
			}
			if attr.Validation.ExclusiveMinimum != nil {
				// Add small epsilon for exclusive minimum
				return *attr.Validation.ExclusiveMinimum + 0.0001
			}
		}
	}

	return nil
}

// generateMaxValues generates an example with all numeric fields at maximum values.
func generateMaxValues(attr *expr.AttributeExpr, gen *expr.ExampleGenerator) any {
	// AsObject handles UserTypes automatically
	if obj := expr.AsObject(attr.Type); obj != nil {
		// For objects, check each field for numeric validations
		result := make(map[string]any)
		hasMaxValues := false

		for _, att := range *obj {
			if att.Attribute.Validation != nil &&
				(att.Attribute.Validation.Maximum != nil || att.Attribute.Validation.ExclusiveMaximum != nil) {
				// This field has max validation
				if maxVal := generateMaxValues(att.Attribute, gen); maxVal != nil {
					result[att.Name] = maxVal
					hasMaxValues = true
				} else if example := att.Attribute.Example(gen); example != nil {
					result[att.Name] = example
				}
			} else if nested := generateMaxValues(att.Attribute, gen); nested != nil {
				// Check if nested objects have max values
				result[att.Name] = nested
				hasMaxValues = true
			} else if example := att.Attribute.Example(gen); example != nil {
				result[att.Name] = example
			}
		}

		if hasMaxValues {
			return result
		}
		return nil
	}

	// Handle primitive numeric types
	switch attr.Type.Kind() {
	case expr.IntKind, expr.Int32Kind, expr.Int64Kind,
		expr.UIntKind, expr.UInt32Kind, expr.UInt64Kind,
		expr.Float32Kind, expr.Float64Kind:
		if attr.Validation != nil {
			if attr.Validation.Maximum != nil {
				return *attr.Validation.Maximum
			}
			if attr.Validation.ExclusiveMaximum != nil {
				// Subtract small epsilon for exclusive maximum
				return *attr.Validation.ExclusiveMaximum - 0.0001
			}
		}
	}

	return nil
}

// generateMinLengthValues generates examples with minimum length strings/arrays.
func generateMinLengthValues(attr *expr.AttributeExpr, gen *expr.ExampleGenerator) any {
	// AsObject handles UserTypes automatically
	if obj := expr.AsObject(attr.Type); obj != nil {
		// For objects, check each field for length validations
		result := make(map[string]any)
		hasMinLength := false

		for _, att := range *obj {
			if att.Attribute.Validation != nil && att.Attribute.Validation.MinLength != nil {
				// This field has min length validation
				if minVal := generateMinLengthValues(att.Attribute, gen); minVal != nil {
					result[att.Name] = minVal
					hasMinLength = true
				} else if example := att.Attribute.Example(gen); example != nil {
					result[att.Name] = example
				}
			} else if nested := generateMinLengthValues(att.Attribute, gen); nested != nil {
				// Check if nested objects have min length values
				result[att.Name] = nested
				hasMinLength = true
			} else if example := att.Attribute.Example(gen); example != nil {
				result[att.Name] = example
			}
		}

		if hasMinLength {
			return result
		}
		return nil
	}

	// Handle primitive string/array types
	switch attr.Type.Kind() {
	case expr.StringKind:
		if attr.Validation != nil && attr.Validation.MinLength != nil {
			// Generate string with exact minimum length
			return strings.Repeat("a", *attr.Validation.MinLength)
		}
	case expr.BytesKind:
		if attr.Validation != nil && attr.Validation.MinLength != nil {
			// Generate bytes with minimum length
			minLen := *attr.Validation.MinLength
			b := make([]byte, minLen)
			for i := range b {
				b[i] = byte('A' + (i % 26))
			}
			return b
		}
	case expr.ArrayKind:
		if attr.Validation != nil && attr.Validation.MinLength != nil {
			// Generate array with minimum number of elements
			arr := attr.Type.(*expr.Array)
			minLen := *attr.Validation.MinLength
			result := make([]any, minLen)
			for i := 0; i < minLen; i++ {
				result[i] = arr.ElemType.Type.Example(gen)
			}
			return result
		}
	}

	return nil
}

// generateMaxLengthValues generates examples with maximum length strings/arrays.
func generateMaxLengthValues(attr *expr.AttributeExpr, gen *expr.ExampleGenerator) any {
	// AsObject handles UserTypes automatically
	if obj := expr.AsObject(attr.Type); obj != nil {
		// For objects, check each field for length validations
		result := make(map[string]any)
		hasMaxLength := false

		for _, att := range *obj {
			if att.Attribute.Validation != nil && att.Attribute.Validation.MaxLength != nil {
				// This field has max length validation
				if maxVal := generateMaxLengthValues(att.Attribute, gen); maxVal != nil {
					result[att.Name] = maxVal
					hasMaxLength = true
				} else if example := att.Attribute.Example(gen); example != nil {
					result[att.Name] = example
				}
			} else if nested := generateMaxLengthValues(att.Attribute, gen); nested != nil {
				// Check if nested objects have max length values
				result[att.Name] = nested
				hasMaxLength = true
			} else if example := att.Attribute.Example(gen); example != nil {
				result[att.Name] = example
			}
		}

		if hasMaxLength {
			return result
		}
		return nil
	}

	// Handle primitive string/array types
	switch attr.Type.Kind() {
	case expr.StringKind:
		if attr.Validation != nil && attr.Validation.MaxLength != nil {
			// Generate string with exact maximum length
			return strings.Repeat("Z", *attr.Validation.MaxLength)
		}
	case expr.BytesKind:
		if attr.Validation != nil && attr.Validation.MaxLength != nil {
			// Generate bytes with maximum length
			maxLen := *attr.Validation.MaxLength
			b := make([]byte, maxLen)
			for i := range b {
				b[i] = byte('z' - (i % 26))
			}
			return b
		}
	case expr.ArrayKind:
		if attr.Validation != nil && attr.Validation.MaxLength != nil {
			// Generate array with maximum number of elements
			arr := attr.Type.(*expr.Array)
			maxLen := *attr.Validation.MaxLength
			result := make([]any, maxLen)
			for i := 0; i < maxLen; i++ {
				result[i] = arr.ElemType.Type.Example(gen)
			}
			return result
		}
	}

	return nil
}
