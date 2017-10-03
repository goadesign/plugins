package security

import (
	"bytes"
	"html/template"
	"strings"

	"goa.design/goa/codegen"
	"goa.design/goa/codegen/service"
	"goa.design/goa/eval"
	"goa.design/plugins/security/design"
	"goa.design/plugins/security/http"
)

// SecuredServiceMethodData contains the data necessry to render
// endpoints.
type SecuredServiceMethodData struct {
	// VarName is the goified name of the method.
	VarName string
	// Requirements lists the security requirements that apply to the
	// secured method.
	Requirements []*design.SecurityExpr
}

// Register the plugin Generator functions.
func init() {
	codegen.RegisterPlugin("gen", Generate)
	codegen.RegisterPlugin("example", Example)
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
			f.SectionTemplates = append(f.SectionTemplates, &codegen.SectionTemplate{
				Source: secureEndpointContextT,
				Data: &SecuredServiceMethodData{
					VarName:      data.VarName,
					Requirements: reqs,
				},
			})
		}
	}
	http.Generate(files)
	return output, nil
}

// anchorComment is used as the marker before which to insert the wrapping code.
const anchorComment = "	// Configure the mux."

var (
	// wrapperTmpl is the template used to render the endpoint wrapper code.
	wrapperTmpl = template.Must(template.New("security-example").Funcs(TemplateFuncs).Parse(wrapperCodeT))

	// TemplateFuncs lists common template helper functions.
	TemplateFuncs = map[string]interface{}{"IsSecured": func(svc, m string) bool {
		return len(design.Requirements(svc, m)) > 0
	}}
)

// Example modified the generated main function so that the secured endpoints
// context gets initialized with the security requirements.
func Example(_ string, _ []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	output := make([]*codegen.File, len(files))
	for i, f := range files {
		output[i] = f
		if s := f.Section("endpoint-method"); s != nil {
			data := s.Data.(*service.EndpointMethodData)
			reqs := design.Requirements(data.ServiceName, data.Name)
			if len(reqs) == 0 {
				continue
			}
		}
		if s := f.Section("service-main"); s != nil {
			var buffer bytes.Buffer
			if err := wrapperTmpl.Execute(&buffer, s.Data); err != nil {
				panic(err) // bug
			}
			s.Source = strings.Replace(s.Source, anchorComment,
				buffer.String()+"\n"+anchorComment, 1)
		}
	}
	return output, nil
}

// input: "service-main" section template input
const wrapperCodeT = `{{- range .Services }}{{ $service := . }}
	{{- range .Methods }}
		{{- if IsSecured $service.Name .Name }}
			{{ $service.VarName }} = Secure{{ .VarName }}Endpoint({{ $service.VarName }}.{{ .VarName }})
		{{- end }}
	{{- end }}
{{- end }}
`

// input: securedServiceMethodData
const secureEndpointContextT = `// {{ printf "Secure%sEndpoint returns an endpoint function which initializes the context with the security requirements and calls ep." .VarName | comment }}
func Secure{{ .VarName }}Endpoint(ep goa.Endpoint) goa.Endpoint {
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
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		ctx = context.WithValue(ctx, security.ContextKey, reqs)
		return ep(ctx, req)
	}
}
`
