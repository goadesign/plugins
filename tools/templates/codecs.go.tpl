var (
{{- range .Types }}
	// {{ .ExportedCodec }} serializes values of type {{ .FullRef }} to canonical JSON.
	{{ .ExportedCodec }} = tools.JSONCodec[{{ .FullRef }}]{
		ToJSON:   {{ .MarshalFunc }},
		FromJSON: {{ .UnmarshalFunc }},
	}

{{- end }}
{{- range .Types }}
	// {{ .GenericCodec }} provides an untyped codec for {{ .FullRef }}.
	{{ .GenericCodec }} = tools.JSONCodec[any]{
		ToJSON: func(v any) ([]byte, error) {
			typed, ok := v.({{ .FullRef }})
			if !ok {
				return nil, fmt.Errorf("expected {{ .FullRef }}, got %T", v)
			}
			return {{ .MarshalFunc }}(typed)
		},
		FromJSON: func(data []byte) (any, error) {
			return {{ .UnmarshalFunc }}(data)
		},
	}

{{- end }}
)

func PayloadCodec(name string) (*tools.JSONCodec[any], bool) {
	switch name {
{{- range .Tools }}
    {{- if .Payload }}
	case "{{ .Name }}":
		return &{{ .Payload.GenericCodec }}, true
    {{- end }}
{{- end }}
	default:
		return nil, false
	}
}

func ResultCodec(name string) (*tools.JSONCodec[any], bool) {
	switch name {
{{- range .Tools }}
    {{- if .Result }}
	case "{{ .Name }}":
		return &{{ .Result.GenericCodec }}, true
    {{- end }}
{{- end }}
	default:
		return nil, false
	}
}

{{- range .Types }}
func {{ .MarshalFunc }}(v {{ .FullRef }}) ([]byte, error) {
    {{- if .CheckNil }}
	if v == nil {
		return nil, fmt.Errorf("{{ .NilError }}")
	}
    {{- end }}
    {{- if .HasValidation }}
	if err := {{ .ValidateFunc }}({{ .MarshalArg }}); err != nil {
		return nil, fmt.Errorf("{{ .ValidateError }}: %w", err)
	}
    {{- end }}
	return json.Marshal(v)
}

func {{ .UnmarshalFunc }}(data []byte) ({{ .FullRef }}, error) {
    {{- if not .Pointer }}
	var zero {{ .FullRef }}
    {{- end }}
	if len(data) == 0 {
        {{- if .Pointer }}
		return nil, fmt.Errorf("{{ .EmptyError }}")
        {{- else }}
		return zero, fmt.Errorf("{{ .EmptyError }}")
        {{- end }}
	}
	var v {{ .ElemRef }}
	if err := json.Unmarshal(data, &v); err != nil {
        {{- if .Pointer }}
		return nil, fmt.Errorf("{{ .DecodeError }}: %w", err)
        {{- else }}
		return zero, fmt.Errorf("{{ .DecodeError }}: %w", err)
        {{- end }}
	}
    {{- if .HasValidation }}
	if err := {{ .ValidateFunc }}({{ .UnmarshalArg }}); err != nil {
        {{- if .Pointer }}
		return nil, fmt.Errorf("{{ .ValidateError }}: %w", err)
        {{- else }}
		return zero, fmt.Errorf("{{ .ValidateError }}: %w", err)
        {{- end }}
	}
    {{- end }}
    {{- if .Pointer }}
	return &v, nil
    {{- else }}
	return v, nil
    {{- end }}
}

func {{ .ValidateFunc }}(body {{ .FullRef }}) (err error) {
    {{- if .HasValidation }}
        {{- range .ValidationSrc }}
            {{- if eq . "" }}

            {{- else }}
	{{ . }}
            {{- end }}
        {{- end }}
    {{- else }}
	_ = body
    {{- end }}
	return
}

{{- end }}
