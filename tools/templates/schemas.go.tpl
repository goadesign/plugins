var (
{{- range .Types }}
	{{ .SchemaVar }} = {{ .SchemaLiteral }}
{{- end }}
)
