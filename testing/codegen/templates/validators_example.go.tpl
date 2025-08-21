package {{ .Package }}

import (
	"fmt"
	{{ .PkgName }} "{{ .PkgPath }}"
)

{{ printf "This file shows examples of custom validators you can define in your test files." | comment }}
{{ printf "Copy these examples to your test files and customize them for your needs." | comment }}

{{ printf "Example: Custom result validator for a method" | comment }}
{{ printf "Define this in your test file to override default validation:" | comment }}
/*
func ValidateGetBookResult(result *bookstore.Book, expected map[string]any) error {
	// Custom validation logic
	if title, ok := expected["title"].(string); ok && result.Title != title {
		return fmt.Errorf("title mismatch: got %q, want %q", result.Title, title)
	}
	
	if author, ok := expected["author"].(string); ok && result.Author != author {
		return fmt.Errorf("author mismatch: got %q, want %q", result.Author, author)
	}
	
	// Check nested fields
	if price, ok := expected["price"].(float64); ok && result.Price != price {
		return fmt.Errorf("price mismatch: got %v, want %v", result.Price, price)
	}
	
	return nil
}
*/

{{ printf "Example: Stream validator for server-sent events" | comment }}
/*
func ValidateStreamUpdatesStream(stream any, expected []map[string]any) error {
	// Cast to the appropriate stream type
	sseStream, ok := stream.(HTTPServerStreamSseClientStream)
	if !ok {
		return fmt.Errorf("invalid stream type")
	}
	
	// Read and validate each expected event
	ctx := context.Background()
	for i, expectedEvent := range expected {
		event, err := sseStream.Recv(ctx)
		if err != nil {
			return fmt.Errorf("failed to receive event %d: %w", i, err)
		}
		
		// Validate event fields
		if id, ok := expectedEvent["id"].(string); ok && event.ID != id {
			return fmt.Errorf("event %d: ID mismatch", i)
		}
	}
	
	return nil
}
*/

{{ printf "Helper function for partial field matching" | comment }}
func ValidateFields(result any, expected map[string]any, fields ...string) error {
	// Convert result to map using JSON marshaling
	resultMap := make(map[string]any)
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}
	if err := json.Unmarshal(data, &resultMap); err != nil {
		return fmt.Errorf("failed to unmarshal result: %w", err)
	}
	
	// Check only specified fields
	for _, field := range fields {
		expectedValue, hasExpected := expected[field]
		actualValue, hasActual := resultMap[field]
		
		if hasExpected && !hasActual {
			return fmt.Errorf("missing field %q in result", field)
		}
		
		if hasExpected {
			expectedJSON, _ := json.Marshal(expectedValue)
			actualJSON, _ := json.Marshal(actualValue)
			if string(expectedJSON) != string(actualJSON) {
				return fmt.Errorf("field %q: expected %s, got %s", field, expectedJSON, actualJSON)
			}
		}
	}
	
	return nil
}

{{ printf "The scenario runner will look for these validator functions by convention:" | comment }}
{{- range .Methods }}
{{ printf "  - Validate%sResult for %s method results" .VarName .Name | comment }}
{{- if or .ServerStream .ClientStream }}
{{ printf "  - Validate%sStream for %s method streams" .VarName .Name | comment }}
{{- end }}
{{- end }}