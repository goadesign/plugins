package {{ .Package }}

import (
	{{ .PkgName }} "{{ .PkgPath }}"
)

{{ printf "Define your custom validator functions here." | comment }}
{{ printf "The scenario runner will call these functions when specified in scenarios.yaml" | comment }}

{{- range .Methods }}
{{- if .ResultRef }}

{{ printf "Validate%s validates results for the %s method." .VarName .Name | comment }}
{{ printf "This function is called when scenarios.yaml specifies:" | comment }}
{{ printf "  validator: Validate%s" .VarName | comment }}
func Validate{{ .VarName }}(result *{{ $.PkgName }}.{{ .Result }}, expected map[string]any) error {
	// Implement your validation logic here
	// Example:
	// if expected["id"] != nil && result.ID != expected["id"] {
	//     return fmt.Errorf("ID mismatch")
	// }
	return nil
}
{{- end }}
{{- end }}

{{ printf "You can also define custom validators with any name and reference them in YAML:" | comment }}
{{ printf "Example:" | comment }}
/*
func MyCustomValidator(result *myservice.Result, expected map[string]any) error {
    // Custom validation logic
    return nil
}

// Then in scenarios.yaml:
// expect:
//   result: { ... }
//   validator: MyCustomValidator
*/