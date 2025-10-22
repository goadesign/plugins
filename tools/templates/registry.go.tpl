// ToolRegistry enumerates all generated tool specifications.
var ToolRegistry = []tools.ToolSpec{
{{- range .Tools }}
	{
		Name:    "{{ .Name }}",
    {{- if .Service }}
		Service: "{{ .Service }}",
    {{- end }}
    {{- if .Set }}
		Set:     "{{ .Set }}",
    {{- end }}
    {{- if .Payload }}
		Payload: tools.TypeSpec{
			Name:   "{{ .Payload.TypeName }}",
        {{- if .Payload.SchemaVar }}
			Schema: {{ .Payload.SchemaVar }},
        {{- end }}
			Codec:  {{ .Payload.GenericCodec }},
		},
    {{- end }}
    {{- if .Result }}
		Result: tools.TypeSpec{
			Name:   "{{ .Result.TypeName }}",
        {{- if .Result.SchemaVar }}
			Schema: {{ .Result.SchemaVar }},
        {{- end }}
			Codec:  {{ .Result.GenericCodec }},
		},
    {{- end }}
	},
{{- end }}
}

var toolIndex map[string]*tools.ToolSpec

func init() {
	toolIndex = make(map[string]*tools.ToolSpec, len(ToolRegistry))
	for i := range ToolRegistry {
		spec := &ToolRegistry[i]
		toolIndex[spec.Name] = spec
	}
}

// Names returns the names of all tools in the registry.
func Names() []string {
	names := make([]string, 0, len(toolIndex))
	for name := range toolIndex {
		names = append(names, name)
	}
	return names
}

// Spec returns the specification for the named tool if it exists.
func Spec(name string) (*tools.ToolSpec, bool) {
	spec, ok := toolIndex[name]
	return spec, ok
}

// PayloadSchema returns the JSON schema for the payload of the named tool.
func PayloadSchema(name string) ([]byte, bool) {
	spec, ok := toolIndex[name]
	if !ok {
		return nil, false
	}
	return spec.Payload.Schema, true
}

// ResultSchema returns the JSON schema for the result of the named tool.
func ResultSchema(name string) ([]byte, bool) {
	spec, ok := toolIndex[name]
	if !ok {
		return nil, false
	}
	return spec.Result.Schema, true
}
