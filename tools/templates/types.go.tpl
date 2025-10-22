type (
{{- range $i, $t := .Types }}
    {{- if gt $i 0 }}

    {{- end }}
	// {{ $t.Doc }}
	{{ $t.Def }}
{{- end }}
)
