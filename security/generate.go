package security

import (
	"goa.design/goa/codegen"
	"goa.design/goa/codegen/service"
	"goa.design/goa/eval"
	"goa.design/plugins/security/design"
	"goa.design/plugins/security/http"
)

// SecuredServiceMethodData contains the data necessry to render
// endpoints.
type SecuredServiceMethodData struct {
	*service.EndpointMethodData
	Requirements []*design.SecurityExpr
}

// Register the plugin Generator functions.
func init() {
	codegen.RegisterPlugin("gen", Generate)
}

// Generate produces server code that enforce the security requirements defined
// in the design. Generate also produces client code that makes it possible to
// provide the required security artifacts. Finally Generate also generate code
// that initializes the context given to the service methods with security
// information.
func Generate(_ string, _ []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	output := make([]*codegen.File, len(files))
	for i, f := range files {
		output[i] = f
		if s := f.Section("endpoint-method"); s != nil {
			data := s.Data.(*service.EndpointMethodData)
			reqs := design.Requirements(data.ServiceName, data.Name)
			if len(reqs) == 0 {
				continue
			}
			s.Source = secureEndpointMethodT
			s.Data = &SecuredServiceMethodData{
				EndpointMethodData: data,
				Requirements:       reqs,
			}
		}
	}
	http.Generate(files)
	return output, nil
}

// input: securedServiceMethodData
const secureEndpointMethodT = `// {{ printf "New%sEndpoint returns an endpoint function that calls method %q of service %q." .VarName .Name .ServiceName | comment }}
func New{{ .VarName }}Endpoint(s {{ .ServiceName }}) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		ctx = context.WithValue(ctx, security.ContextKey, reqs)
		reqs := make([]*security.Requirement, {{ len .Requirements }})
		{{- range $i, $req := .Requirements }}
			{{- if $req.Scopes }}
		reqs[$i].RequiredScopes = []string{ {{- range $req.Scopes }}{{ printf "%q" . }}, {{ end }} }
			{{- end }}
		reqs.Schemes := make([]*security.Schemes, {{ len .Schemes }})
			{{- range $i, $scheme := .Schemes }}
		reqs.Schemes[$i] = &security.Scheme{
			Kind: security.SchemeKind({{ $scheme.Kind }}),
			Name: {{ $scheme.SchemeName }},
		}
				{{- if .Scopes }}
			reqs.Schemes[$i].Scopes = []string{ {{- range $scheme.Scopes }}{{ printf "%q" . }}, {{ end }} }
				{{- end }}
				{{- if .In }}
		reqs.Schemes[$i].RequiredKey = &security.Key{
			In: {{ printf "%q" .In }},
			Name: {{ printf "%q" .Name }},
		}
				{{- end }}
				{{- if .Flows }}
		reqs.Schemes[$i].Flows = make([]*security.OAuthFlow, {{ len .Flows }})
					{{- range $j, $flow := .Flows }}
		reqs.Schemes[$i].Flows[$j] = &security.OAuthFlow{
			Kind: security.FlowKind({{ $flow.Kind }}),
		}
						{{- if .AuthorizationURL }}
		reqs.Schemes[$i].Flows[$j].AuthorizationURL = {{ printf "%q" .AuthorizationURL }}
						{{- end }}
						{{- if .TokenURL }}
		reqs.Schemes[$i].Flows[$j].TokenURL = {{ printf "%q" .TokenURL }}
						{{- end }}
						{{- if .RefreshURL }}
		reqs.Schemes[$i].Flows[$j].RefreshURL = {{ printf "%q" .RefreshURL }}
						{{- end }}
					{{- end }}
				{{- end }}
			{{- end }}
		{{- end }}
		{{- if .PayloadRef }}
		p := req.({{ .PayloadRef }})
		{{- end }}
		{{- if .ResultRef }}
		return s.{{ .Name }}(ctx{{ if .PayloadRef }}, p{{ end }})
		{{- else }}
		return nil, s.{{ .Name }}(ctx{{ if .PayloadRef }}, p{{ end }})
		{{- end }}
	}
}
`
