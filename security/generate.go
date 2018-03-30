package security

import (
	"fmt"
	"path/filepath"
	"strings"

	"goa.design/goa/codegen"
	goadesign "goa.design/goa/design"
	"goa.design/goa/eval"
	"goa.design/goa/http/codegen/openapi"
	seccodegen "goa.design/plugins/security/codegen"
	"goa.design/plugins/security/design"

	// Initializes the HTTP generator
	_ "goa.design/plugins/security/http"
)

type (
	// AuthFuncsData contains data necessary to render the dummy authorization
	// functions in the example service file.
	AuthFuncsData struct {
		// Schemes is the unique security schemes defined in the API.
		Schemes []*design.SchemeExpr
		// ServiceName is the name of the service.
		ServiceName string
		// ServicePkg is the service package name.
		ServicePkg string
		// Security is the security package name.
		SecurityPkg string
	}
)

// Register the plugin Generator functions.
func init() {
	codegen.RegisterPlugin("gen", design.Root, Generate)
	codegen.RegisterPlugin("example", design.Root, Example)
}

// Generate produces server code that enforce the security requirements defined
// in the design. Generate also produces client code that makes it possible to
// provide the required security artifacts. Finally Generate also generate code
// that initializes the context given to the service methods with security
// information.
func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, root := range roots {
		switch r := root.(type) {
		case *goadesign.RootExpr:
			for _, s := range r.Services {
				if f := SecureEndpointFile(genpkg, s); f != nil {
					files = append(files, f)
				}
			}
		case *design.RootExpr:
			for _, f := range files {
				OpenAPIV2(r, f)
			}
		}
	}
	return files, nil
}

// Example modified the generated main function so that the secured endpoints
// context gets initialized with the security requirements.
func Example(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	var data []*AuthFuncsData
	for _, root := range roots {
		switch r := root.(type) {
		case *goadesign.RootExpr:
			for _, s := range r.Services {
				sd := seccodegen.BuildSecureServiceData(s, "")
				data = append(data, &AuthFuncsData{
					Schemes:     sd.Schemes,
					SecurityPkg: "security",
					ServicePkg:  sd.PkgName,
					ServiceName: codegen.Goify(sd.Name, true)})
			}
		}
	}
	for _, f := range files {
		if s := f.Section("dummy-endpoint"); len(s) > 0 {
			for _, h := range f.Section("source-header") {
				codegen.AddImport(h, codegen.SimpleImport("goa.design/plugins/security"))
			}
			f.SectionTemplates = append(f.SectionTemplates, &codegen.SectionTemplate{
				Name:   "dummy-authorize-funcs",
				Source: dummyAuthFuncsT,
				Data:   data,
			})
		}
	}
	return files, nil
}

// OpenAPIV2 updates the openapi.json file with the security definitions.
func OpenAPIV2(r *design.RootExpr, f *codegen.File) {
	for _, s := range f.Section("openapi") {
		spec := s.Data.(*openapi.V2)
		spec.SecurityDefinitions = buildV2SecurityDefinitions(r.Schemes)
		s.Data = spec
	}
}

// SecureEndpointFile returns the file containing the secure endpoint
// definitions.
func SecureEndpointFile(genpkg string, svc *goadesign.ServiceExpr) *codegen.File {
	data := seccodegen.BuildSecureServiceData(svc, "")
	path := filepath.Join(codegen.Gendir, codegen.SnakeCase(svc.Name), "security.go")
	header := codegen.Header(
		svc.Name+" service security",
		data.PkgName,
		[]*codegen.ImportSpec{
			{Path: "context"},
			{Path: "goa.design/goa"},
			{Path: "goa.design/plugins/security"},
		})
	sections := []*codegen.SectionTemplate{header}
	sections = append(sections, &codegen.SectionTemplate{
		Name:    "secure-endpoint-init",
		Source:  secureEndpointsInitT,
		Data:    data,
		FuncMap: codegen.TemplateFuncs(),
	})
	for _, m := range data.Methods {
		if len(m.Requirements) == 0 {
			continue
		}
		sections = append(sections, &codegen.SectionTemplate{
			Name:   "secure-endpoint",
			Source: secureEndpointT,
			Data:   m,
		})
	}
	if len(sections) == 1 {
		return nil
	}
	return &codegen.File{
		Path:             path,
		SectionTemplates: sections,
	}
}

func buildV2SecurityDefinitions(schemes []*design.SchemeExpr) map[string]*openapi.SecurityDefinition {
	sds := make(map[string]*openapi.SecurityDefinition)
	for _, s := range schemes {
		sd := openapi.SecurityDefinition{
			Description: s.Description,
			Extensions:  openapi.ExtensionsFromExpr(s.Metadata),
		}
		switch s.Kind {
		case design.BasicAuthKind:
			sd.Type = "basic"
		case design.APIKeyKind:
			sd.Type = "apiKey"
			sd.In = s.In
			if sd.In == "" {
				sd.In = "header"
			}
			sd.Name = s.Name
			if sd.Name == "" {
				sd.Name = "key"
			}
		case design.JWTKind:
			sd.Type = "apiKey"
			sd.In = s.In
			if sd.In == "" {
				sd.In = "header"
			}
			sd.Name = s.Name
			if sd.Name == "" {
				sd.Name = "token"
			}
			// OpenAPI V2 spec does not support JWT scheme. Hence we add the scheme
			// information to the description.
			lines := []string{}
			for _, scope := range s.Scopes {
				lines = append(lines, fmt.Sprintf("  * `%s`: %s", scope.Name, scope.Description))
			}
			sd.Description += fmt.Sprintf("\n**Security Scopes**:\n%s", strings.Join(lines, "\n"))
		case design.OAuth2Kind:
			sd.Type = "oauth2"
			if scopesLen := len(s.Scopes); scopesLen > 0 {
				scopes := make(map[string]string, scopesLen)
				for _, scope := range s.Scopes {
					scopes[scope.Name] = scope.Description
				}
				sd.Scopes = scopes
			}
		}
		if len(s.Flows) > 0 {
			switch s.Flows[0].Kind {
			case design.AuthorizationCodeFlowKind:
				sd.Flow = "accessCode"
			case design.ImplicitFlowKind:
				sd.Flow = "implicit"
			case design.PasswordFlowKind:
				sd.Flow = "password"
			case design.ClientCredentialsFlowKind:
				sd.Flow = "application"
			}
			sd.AuthorizationURL = s.Flows[0].AuthorizationURL
			sd.TokenURL = s.Flows[0].TokenURL
		}
		sds[s.SchemeName] = &sd
	}
	return sds
}

// input: securedServiceData
const secureEndpointsInitT = `{{ printf "NewSecure%s wraps the methods of a %s service with security scheme aware endpoints." .EndpointsVarName .Name | comment }}
	func NewSecure{{ .EndpointsVarName }}(s {{ .VarName }}{{ range .Schemes }}, auth{{ .Type }}Fn {{ $.SecurityPkgName }}.Authorize{{ .Type }}Func{{ end }}) *{{ .EndpointsVarName }} {
		return &{{ .EndpointsVarName }}{
			{{- range .Methods }}
			{{ .MethodData.VarName }}: {{ if .Requirements }}{{ .VarName }}({{ end }}{{ .NonSecureVarName }}(s){{ if .Requirements }}{{ range .Schemes }}, auth{{ .Scheme.Type }}Fn{{ end }}){{ end }},
			{{- end }}
		}
	}
`

// input: securedServiceMethodData
const secureEndpointT = `{{ printf "%s returns an endpoint function which initializes the context with the security requirements for the method %q of service %q." .VarName .Name .ServiceName | comment }}
func {{ .VarName }}(ep goa.Endpoint{{ range .Schemes }}, auth{{ .Scheme.Type }}Fn {{ $.SecurityPkgName }}.Authorize{{ .Scheme.Type }}Func{{ end }}) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.({{ if .PayloadRef }}*{{ end }}{{ .Payload }})
		var err error
		{{- range $ridx, $r := .Requirements }}
			{{- if ne $ridx 0 }}
		if err != nil {
			{{- end }}
			{{- range $sidx, $s := .Schemes }}
				{{- $scheme := $.SchemeData $s }}
				{{- if ne $sidx 0 }}
			if err == nil {
				{{- end }}
				{{- if eq .Type "BasicAuth" }}
				basicAuthSch := {{ $.SecurityPkgName }}.BasicAuthScheme{
					Name: {{ printf "%q" .SchemeName }},
				}
				ctx, err = auth{{ .Type }}Fn(ctx, {{ if $scheme.UsernamePointer }}*{{ end }}p.{{ $scheme.UsernameField }}, {{ if $scheme.PasswordPointer }}*{{ end }}p.{{ $scheme.PasswordField }}, &basicAuthSch)

				{{- else if eq .Type "APIKey" }}
				apiKeySch := {{ $.SecurityPkgName }}.APIKeyScheme{
					Name: {{ printf "%q" .SchemeName }},
				}
				ctx, err = auth{{ .Type }}Fn(ctx, {{ if $scheme.CredPointer }}*{{ end }}p.{{ $scheme.CredField }}, &apiKeySch)

				{{- else if eq .Type "JWT" }}
				jwtSch := {{ $.SecurityPkgName }}.JWTScheme{
					Name: {{ printf "%q" .SchemeName }},
					Scopes: []string{ {{- range .Scopes }}{{ printf "%q" .Name }}, {{ end }} },
					RequiredScopes: []string{ {{- range $r.Scopes }}{{ printf "%q" . }}, {{ end }} },
				}
				ctx, err = auth{{ .Type }}Fn(ctx, {{ if $scheme.CredPointer }}*{{ end }}p.{{ $scheme.CredField }}, &jwtSch)

				{{- else if eq .Type "OAuth2" }}
				oauth2Sch := {{ $.SecurityPkgName }}.OAuth2Scheme{
					Name: {{ printf "%q" .SchemeName }},
					Scopes: []string{ {{- range .Scopes }}{{ printf "%q" .Name }}, {{ end }} },
					RequiredScopes: []string{ {{- range $r.Scopes }}{{ printf "%q" . }}, {{ end }} },
					{{- if .Flows }}
					Flows: []*security.OAuthFlow{
						{{- range .Flows }}
						&security.OAuthFlow{
							Type: "{{ .Type }}",
							{{- if .AuthorizationURL }}
							AuthorizationURL: {{ printf "%q" .AuthorizationURL }},
							{{- end }}
							{{- if .TokenURL }}
							TokenURL: {{ printf "%q" .TokenURL }},
							{{- end }}
							{{- if .RefreshURL }}
							RefreshURL: {{ printf "%q" .RefreshURL }},
							{{- end }}
						},
						{{- end }}
					},
					{{- end }}
				}
				ctx, err = auth{{ .Type }}Fn(ctx, {{ if $scheme.CredPointer }}*{{ end }}p.{{ $scheme.CredField }}, &oauth2Sch)

				{{- else }}
				{{ printf "unsupported scheme type %q" .Type | comment }}

				{{- end }}
				{{- if ne $sidx 0 }}
				}
				{{- end }}
			{{- end }}
			{{- if ne $ridx 0 }}
		}
			{{- end }}
		{{- end }}
		if err != nil {
			return nil, err
		}
		return ep(ctx, req)
	}
}
`

// data: AuthFuncsData
const dummyAuthFuncsT = `{{ range $sd := . }}
{{- range .Schemes }}

{{ printf "%sAuth%sFn implements the authorization logic for %s scheme." $sd.ServiceName .Type .Type | comment }}
func {{ $sd.ServiceName }}Auth{{ .Type }}Fn(ctx context.Context, {{ if eq .Type "BasicAuth" }}user, pass{{ else if eq .Type "APIKey" }}key{{ else }}token{{ end }} string, s *{{ $sd.SecurityPkg }}.{{ .Type }}Scheme) (context.Context, error) {
	// Add authorization logic
	{{- if eq .Type "BasicAuth" }}
	if user == "" {
		return ctx, &{{ $sd.ServicePkg }}.Unauthorized{"invalid username"}
	}
	if pass == "" {
		return ctx, &{{ $sd.ServicePkg }}.Unauthorized{"invalid password"}
	}
	{{- else if eq .Type "APIKey" }}
	if key == "" {
		return ctx, &{{ $sd.ServicePkg }}.Unauthorized{"invalid key"}
	}
	{{- else }}
	if token == "" {
		return ctx, &{{ $sd.ServicePkg }}.Unauthorized{"invalid token"}
	}
	{{- end }}
	return ctx, nil
}
{{- end }}
{{- end }}
`
