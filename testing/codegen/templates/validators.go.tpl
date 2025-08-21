package {{ .Package }}

import (
	"fmt"
	{{ .PkgName }} "{{ .PkgPath }}"
)

{{ printf "This file contains validator function variables that you can define in your test files." | comment }}
{{ printf "The scenario runner will automatically call these functions if they are defined." | comment }}
{{ printf "" | comment }}
{{ printf "Example usage in your test file:" | comment }}
{{ printf "  func init() {" | comment }}
{{ printf "    // Define a validator for the GetBook method" | comment }}
{{ printf "    ValidateGetBookResult = func(result *bookstore.Book, expected map[string]any) error {" | comment }}
{{ printf "      if result.Title != expected[\"title\"] {" | comment }}
{{ printf "        return fmt.Errorf(\"title mismatch: got %q, want %q\", result.Title, expected[\"title\"])" | comment }}
{{ printf "      }" | comment }}
{{ printf "      return nil" | comment }}
{{ printf "    }" | comment }}
{{ printf "  }" | comment }}

{{- range .Methods }}
{{ printf "Validate%sResult validates the result of the %s method." .VarName .Name | comment }}
{{ printf "Define this function in your test to provide custom validation logic." | comment }}
var Validate{{ .VarName }}Result func(result *{{ $.PkgName }}.{{ .Result }}, expected map[string]any) error

{{- if or .ServerStream .ClientStream }}
{{ printf "Validate%sStream validates streaming results from the %s method." .VarName .Name | comment }}
var Validate{{ .VarName }}Stream func(stream any, expected []map[string]any) error
{{- end }}
{{- end }}

{{ printf "Helper functions for common validation patterns" | comment }}

{{ printf "ValidateFields checks that all expected fields match in the result." | comment }}
func ValidateFields(result any, expected map[string]any, fields ...string) error {
	// Convert result to map
	resultMap := make(map[string]any)
	if err := toMap(result, &resultMap); err != nil {
		return fmt.Errorf("failed to convert result to map: %w", err)
	}
	
	// Check each specified field
	for _, field := range fields {
		expectedValue, ok := expected[field]
		if !ok {
			continue // Field not in expected, skip
		}
		
		actualValue, ok := resultMap[field]
		if !ok {
			return fmt.Errorf("missing field %q in result", field)
		}
		
		if fmt.Sprintf("%v", expectedValue) != fmt.Sprintf("%v", actualValue) {
			return fmt.Errorf("field %q: expected %v, got %v", field, expectedValue, actualValue)
		}
	}
	
	return nil
}

{{ printf "ValidateContains checks that the result contains the expected values." | comment }}
func ValidateContains(result any, expected map[string]any) error {
	// Convert result to map
	resultMap := make(map[string]any)
	if err := toMap(result, &resultMap); err != nil {
		return fmt.Errorf("failed to convert result to map: %w", err)
	}
	
	// Check each expected field exists with the right value
	for key, expectedValue := range expected {
		actualValue, ok := resultMap[key]
		if !ok {
			return fmt.Errorf("missing expected field %q", key)
		}
		
		if fmt.Sprintf("%v", expectedValue) != fmt.Sprintf("%v", actualValue) {
			return fmt.Errorf("field %q: expected %v, got %v", key, expectedValue, actualValue)
		}
	}
	
	return nil
}

func toMap(from any, to *map[string]any) error {
	// This is a simplified version - in production you might use reflection or json marshaling
	// The actual implementation would be provided by the user or a utility library
	return nil
}