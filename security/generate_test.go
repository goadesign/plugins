package security

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/go-openapi/loads"

	codegen "goa.design/goa/codegen"
	goadesign "goa.design/goa/design"
	"goa.design/goa/eval"
	httpcodegen "goa.design/goa/http/codegen"
	httpdesign "goa.design/goa/http/design"
	"goa.design/plugins/security/design"
	"goa.design/plugins/security/testdata"
)

func TestSecureEndpointInit(t *testing.T) {
	cases := []struct {
		Name string
		DSL  func()
		Code string
	}{
		{"endpoint-without-requirement", testdata.EndpointWithoutRequirementDSL, testdata.EndpointInitWithoutRequirementCode},
		{"endpoints-with-requirements", testdata.EndpointsWithRequirementsDSL, testdata.EndpointInitWithRequirementsCode},
		{"endpoints-with-service-requirements", testdata.EndpointsWithServiceRequirementsDSL, testdata.EndpointInitWithServiceRequirementsCode},
		{"endpoints-no-security", testdata.EndpointNoSecurityDSL, testdata.EndpointInitNoSecurityCode},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			codegen.RunDSL(t, c.DSL)
			if len(goadesign.Root.Services) != 1 {
				t.Fatalf("got %d services, expected 1", len(goadesign.Root.Services))
			}
			fs := SecureEndpointFile("", goadesign.Root.Services[0])
			if fs == nil {
				t.Fatalf("got nil file, expected not nil")
			}
			sections := fs.SectionTemplates
			if len(sections) < 2 {
				t.Fatalf("got %d sections, expected at least 2", len(sections))
			}
			code := codegen.SectionCode(t, sections[1])
			if code != c.Code {
				t.Errorf("invalid code, got:\n%s\ngot vs. expected:\n%s", code, codegen.Diff(t, code, c.Code))
			}
		})
	}
}

func TestSecureEndpoint(t *testing.T) {
	cases := []struct {
		Name string
		DSL  func()
		Code string
	}{
		{"with-required-scopes", testdata.EndpointWithRequiredScopesDSL, testdata.EndpointWithRequiredScopesCode},
		{"with-api-key-override", testdata.EndpointWithAPIKeyOverrideDSL, testdata.EndpointWithAPIKeyOverrideCode},
		{"with-oauth2", testdata.EndpointWithOAuth2DSL, testdata.EndpointWithOAuth2Code},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			codegen.RunDSL(t, c.DSL)
			if len(goadesign.Root.Services) != 1 {
				t.Fatalf("got %d services, expected 1", len(goadesign.Root.Services))
			}
			fs := SecureEndpointFile("", goadesign.Root.Services[0])
			if fs == nil {
				t.Fatalf("got nil file, expected not nil")
			}
			sections := fs.SectionTemplates
			code := codegen.SectionCode(t, sections[2])
			if code != c.Code {
				t.Errorf("invalid code, got:\n%s\ngot vs. expected:\n%s", code, codegen.Diff(t, code, c.Code))
			}
		})
	}
}

func TestOpenAPIV2(t *testing.T) {
	a := &goadesign.APIExpr{
		Name:    "test",
		Servers: []*goadesign.ServerExpr{{URL: "https://goa.design"}},
	}
	cases := []struct {
		Name string
		DSL  func()
	}{
		{"endpoint-without-requirement", testdata.EndpointWithoutRequirementDSL},
		{"endpoints-with-requirements", testdata.EndpointsWithRequirementsDSL},
		{"endpoints-with-service-requirements", testdata.EndpointsWithServiceRequirementsDSL},
		{"endpoints-no-security", testdata.EndpointNoSecurityDSL},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			codegen.RunDSL(t, c.DSL)
			httpdesign.Root.Design.API = a
			f, err := httpcodegen.OpenAPIFile(httpdesign.Root)
			if err != nil {
				t.Fatalf("error generating openapi file: %s", err)
			}
			OpenAPIV2(design.Root, f)
			s := f.SectionTemplates
			if len(s) != 1 {
				t.Fatalf("%s: expected 1 section, got %d", c.Name, len(s))
			}
			if s[0].Source == "" {
				t.Fatalf("%s: empty section template", c.Name)
			}
			if s[0].Data == nil {
				t.Fatalf("%s: nil data", c.Name)
			}
			var buf bytes.Buffer
			tmpl := template.Must(template.New("openapi").Funcs(s[0].FuncMap).Parse(s[0].Source))
			err = tmpl.Execute(&buf, s[0].Data)
			if err != nil {
				t.Fatalf("%s: failed to render template: %s", c.Name, err)
			}
			validateSwagger(t, c.Name, buf.Bytes())
		})
	}
}

func TestExample(t *testing.T) {
	cases := []struct {
		Name string
		DSL  func()
		Code string
	}{
		{"single-service", testdata.SingleServiceDSL, testdata.SingleServiceAuthFuncsCode},
		{"multiple-services", testdata.MultipleServicesDSL, testdata.MultipleServicesAuthFuncsCode},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			httpcodegen.RunHTTPDSL(t, c.DSL)
			if len(goadesign.Root.Services) == 0 {
				t.Fatalf("expected at least 1 service")
			}
			fs := httpcodegen.ExampleServerFiles("", httpdesign.Root)
			Example("", []eval.Root{goadesign.Root}, fs)
			for _, f := range fs {
				if filepath.Base(f.Path) == "main.go" {
					continue
				}
				sections := f.Section("dummy-authorize-funcs")
				if len(sections) < 1 {
					t.Fatalf("service-main: expected at least 1")
				}
				code := codegen.SectionCode(t, sections[0])
				if code != c.Code {
					t.Errorf("invalid code, got:\n%s\ngot vs. expected:\n%s", code, codegen.Diff(t, code, c.Code))
				}
			}
		})
	}
}

// validateSwagger asserts that the given bytes contain a valid Swagger spec.
func validateSwagger(t *testing.T, title string, b []byte) {
	doc, err := loads.Analyzed(json.RawMessage(b), "")
	if err != nil {
		t.Errorf("%s: invalid swagger: %s", title, err)
	}
	if doc == nil {
		t.Errorf("%s: nil swagger", title)
	}
}
