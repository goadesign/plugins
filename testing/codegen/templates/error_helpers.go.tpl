{{ printf "ErrorAsserter provides methods for asserting specific errors." | comment }}
type ErrorAsserter struct {
	t *testing.T
}

{{ printf "NewErrorAsserter creates a new error asserter." | comment }}
func NewErrorAsserter(t *testing.T) *ErrorAsserter {
	t.Helper()
	return &ErrorAsserter{t: t}
}

{{- range .Errors }}
{{ printf "Assert%s asserts that the error is %s." (goify .Name true) .Name | comment }}
{{- if .Description }}
// {{ .Description }}
{{- end }}
// This error can be returned by: {{ range $i, $m := .Methods }}{{ if $i }}, {{ end }}{{ $m }}{{ end }}
func (a *ErrorAsserter) Assert{{ goify .Name true }}(err error) {
	a.t.Helper()
	if err == nil {
		a.t.Errorf(`Expected {{ .Name }} error, got nil`)
		return
	}
	
	// Check if it's a Goa error with ErrorName method
	type errorNamer interface{ ErrorName() string }
	if en, ok := err.(errorNamer); ok {
		if en.ErrorName() != "{{ .Name }}" {
			a.t.Errorf(`Expected {{ .Name }} error, got %s`, en.ErrorName())
		}
		return
	}
	
	// For HTTP transport errors, check the error message
	errMsg := err.Error()
	if strings.Contains(errMsg, `"name":"{{ .Name }}"`) {
		// HTTP transport error contains the error name in JSON
		return
	}
	
	// For gRPC errors, check if it's the right error type
	if strings.Contains(errMsg, "{{ .Name }}") {
		return
	}
	
	a.t.Errorf(`Expected {{ .Name }} error, got: %v`, err)
}

{{ printf "Expect%s runs a function and asserts it returns %s error." (goify .Name true) .Name | comment }}
func (a *ErrorAsserter) Expect{{ goify .Name true }}(fn func() error) {
	a.t.Helper()
	err := fn()
	a.Assert{{ goify .Name true }}(err)
}
{{- end }}
