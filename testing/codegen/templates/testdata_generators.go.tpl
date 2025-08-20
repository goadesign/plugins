{{ printf "TestData provides generators for test data." | comment }}
type TestData struct{}

{{ printf "NewTestData creates a new test data generator." | comment }}
func NewTestData() *TestData {
	return &TestData{}
}


{{- range .Methods }}
{{- $method := . }}
{{- if .PayloadEx }}

{{ printf "Valid%sPayload generates a valid %s payload." (goify .Name true) .Name | comment }}
func (td *TestData) Valid{{ goify .Name true }}Payload() {{ if .PayloadRef }}{{ .PayloadRef }}{{ else }}any{{ end }} {
	{{- if .PayloadInit }}
	var payload {{ if .PayloadRef }}{{ .PayloadRef }}{{ else }}any{{ end }}
	json.Unmarshal([]byte(`{{ .PayloadInit }}`), &payload)
	return payload
	{{- else }}
	return nil
	{{- end }}
}

{{- /* Generate edge case methods for this payload type */ -}}
{{- range .EdgeCases }}

{{ printf "%sPayloadWith%s generates a %s payload %s." (goify $method.Name true) .Name $method.Name .Description | comment }}
func (td *TestData) {{ goify $method.Name true }}PayloadWith{{ .Name }}() {{ if $method.PayloadRef }}{{ $method.PayloadRef }}{{ else }}any{{ end }} {
	{{- if .Init }}
	var payload {{ if $method.PayloadRef }}{{ $method.PayloadRef }}{{ else }}any{{ end }}
	json.Unmarshal([]byte(`{{ .Init }}`), &payload)
	return payload
	{{- else }}
	return td.Valid{{ goify $method.Name true }}Payload()
	{{- end }}
}
{{- end }}

{{ printf "%sPayloadWithAllFields generates a %s payload with all fields populated." (goify .Name true) .Name | comment }}
func (td *TestData) {{ goify .Name true }}PayloadWithAllFields() {{ if .PayloadRef }}{{ .PayloadRef }}{{ else }}any{{ end }} {
	{{- if .PayloadInit }}
	var payload {{ if .PayloadRef }}{{ .PayloadRef }}{{ else }}any{{ end }}
	json.Unmarshal([]byte(`{{ .PayloadInit }}`), &payload)
	return payload
	{{- else }}
	return nil
	{{- end }}
}

{{- if .Payload }}
{{ printf "%sPayloadBuilder provides a fluent interface for building %s payload instances." (goify .Name true) .Name | comment }}
type {{ goify .Name true }}PayloadBuilder struct {
	obj {{ if .PayloadRef }}{{ .PayloadRef }}{{ else }}any{{ end }}
}

{{ printf "New%sPayloadBuilder creates a new builder for %s payload." (goify .Name true) .Name | comment }}
func (td *TestData) New{{ goify .Name true }}PayloadBuilder() *{{ goify .Name true }}PayloadBuilder {
	return &{{ goify .Name true }}PayloadBuilder{
		obj: td.Valid{{ goify .Name true }}Payload(),
	}
}

{{- /* Generate builder methods for payload fields based on service data */ -}}

{{ printf "Build returns the constructed %s payload." .Name | comment }}
func (b *{{ goify .Name true }}PayloadBuilder) Build() {{ if .PayloadRef }}{{ .PayloadRef }}{{ else }}any{{ end }} {
	return b.obj
}
{{- end }}
{{- end }}
{{- end }}