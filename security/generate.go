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
		needsImport := false
		for _, s := range f.Section("endpoint-method") {
			needsImport = true
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
		if needsImport {
			codegen.AddImport(f.SectionTemplates[0],
				&codegen.ImportSpec{Path: "goa.design/plugins/security"})
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
		needsImport := false
		output[i] = f
		for _, s := range f.Section("endpoint-method") {
			data := s.Data.(*service.EndpointMethodData)
			reqs := design.Requirements(data.ServiceName, data.Name)
			if len(reqs) == 0 {
				continue
			}
			needsImport = true
		}
		if needsImport {
			codegen.AddImport(f.SectionTemplates[0],
				&codegen.ImportSpec{Path: "goa.design/plugins/security"})
		}
		for _, s := range f.Section("service-main") {
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
const secureEndpointContextT = `{{ printf "Secure%sEndpoint returns an endpoint function which initializes the context with the security requirements and calls ep." .VarName | comment }}
func Secure{{ .VarName }}Endpoint(ep goa.Endpoint) goa.Endpoint {
	reqs := make([]*security.Requirement, {{ len .Requirements }})
	{{- range $i, $req := .Requirements }}
		{{- if $req.Scopes }}
	reqs[{{ $i }}].RequiredScopes = []string{ {{- range $req.Scopes }}{{ printf "%q" . }}, {{ end }} }
		{{- end }}
	reqs[{{ $i }}].Schemes = make([]*security.Scheme, {{ len .Schemes }})
		{{- range $j, $scheme := .Schemes }}
	reqs[{{ $i }}].Schemes[{{ $j }}] = &security.Scheme{
		Kind: security.SchemeKind({{ $scheme.Kind }}),
		Name: {{ printf "%q" $scheme.SchemeName }},
	}
			{{- if .Scopes }}
	reqs[{{ $i }}].Schemes[{{ $j }}].Scopes = []string{ {{- range $scheme.Scopes }}{{ printf "%q" .Name }}, {{ end }} }
			{{- end }}
			{{- if .In }}
	reqs[{{ $i }}].Schemes[{{ $j }}].RequiredKey = &security.Key{
		In: {{ printf "%q" .In }},
		Name: {{ printf "%q" .Name }},
	}
			{{- end }}
			{{- if .Flows }}
	reqs[{{ $i }}].Schemes[{{ $j }}].Flows = make([]*security.OAuthFlow, {{ len .Flows }})
				{{- range $j, $flow := .Flows }}
	reqs[{{ $i }}].Schemes[{{ $j }}].Flows[{{ $j }}] = &security.OAuthFlow{
		Kind: security.FlowKind({{ $flow.Kind }}),
	}
					{{- if .AuthorizationURL }}
	reqs[{{ $i }}].Schemes[{{ $j }}].Flows[{{ $j }}].AuthorizationURL = {{ printf "%q" .AuthorizationURL }}
					{{- end }}
					{{- if .TokenURL }}
	reqs[{{ $i }}].Schemes[{{ $j }}].Flows[{{ $j }}].TokenURL = {{ printf "%q" .TokenURL }}
					{{- end }}
					{{- if .RefreshURL }}
	reqs[{{ $i }}].Schemes[{{ $j }}].Flows[{{ $j }}].RefreshURL = {{ printf "%q" .RefreshURL }}
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
