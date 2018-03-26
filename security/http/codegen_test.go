package http

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/go-openapi/loads"

	"goa.design/goa/codegen"
	goadesign "goa.design/goa/design"
	"goa.design/goa/eval"
	httpcodegen "goa.design/goa/http/codegen"
	httpdesign "goa.design/goa/http/design"
	"goa.design/plugins/security/http/testdata"
)

func TestSecureDecoder(t *testing.T) {
	cases := []struct {
		Name string
		DSL  func()
		Code string
	}{
		{"basic-auth", testdata.BasicAuthDSL, testdata.BasicAuthSecureDecoderCode},
		{"oauth2", testdata.OAuth2DSL, testdata.OAuth2SecureDecoderCode},
		{"oauth2-in-param", testdata.OAuth2InParamDSL, testdata.OAuth2InParamSecureDecoderCode},
		{"jwt", testdata.JWTDSL, testdata.JWTSecureDecoderCode},
		{"api-key", testdata.APIKeyDSL, testdata.APIKeySecureDecoderCode},
		{"api-key-in-param", testdata.APIKeyInParamDSL, testdata.APIKeyInParamSecureDecoderCode},
		{"multiple-and", testdata.MultipleAndDSL, testdata.MultipleAndSecureDecoderCode},
		{"multiple-or", testdata.MultipleOrDSL, testdata.MultipleOrSecureDecoderCode},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			RunHTTPDSL(t, c.DSL, testdata.TopLevelSchemes)
			fs := httpcodegen.ServerFiles("", httpdesign.Root)
			if len(fs) != 2 {
				t.Fatalf("got %d files, expected two", len(fs))
			}
			Generate("", []eval.Root{httpdesign.Root}, fs)
			for _, f := range fs {
				switch filepath.Base(f.Path) {
				case "encode_decode.go":
					sections := f.Section("secure-request-decoder")
					if len(sections) < 1 {
						t.Fatalf("got %d sections, expected at least 1", len(sections))
					}
					code := codegen.SectionCode(t, sections[0])
					if code != c.Code {
						t.Errorf("invalid code, got:\n%s\ngot vs. expected:\n%s", code, codegen.Diff(t, code, c.Code))
					}
				case "server.go":
					sections := f.Section("server-handler-init")
					if len(sections) < 1 {
						t.Fatalf("got %d sections, expected at least 1", len(sections))
					}
					code := codegen.SectionCode(t, sections[0])
					data := sections[0].Data.(*httpcodegen.EndpointData)
					if !strings.Contains(code, "Secure"+data.RequestDecoder) {
						t.Errorf("invalid code, got:\n%s\n expected %s in the code", strings.TrimSpace(code), "Secure"+data.RequestDecoder)
					}
				}
			}
		})
	}
}

func TestSecureEncoder(t *testing.T) {
	cases := []struct {
		Name string
		DSL  func()
		Code string
	}{
		{"basic-auth", testdata.BasicAuthDSL, testdata.BasicAuthSecureEncoderCode},
		{"oauth2", testdata.OAuth2DSL, testdata.OAuth2SecureEncoderCode},
		{"oauth2-in-param", testdata.OAuth2InParamDSL, testdata.OAuth2InParamSecureEncoderCode},
		{"jwt", testdata.JWTDSL, testdata.JWTSecureEncoderCode},
		{"api-key", testdata.APIKeyDSL, testdata.APIKeySecureEncoderCode},
		{"api-key-in-param", testdata.APIKeyInParamDSL, testdata.APIKeyInParamSecureEncoderCode},
		{"multiple-and", testdata.MultipleAndDSL, testdata.MultipleAndSecureEncoderCode},
		{"multiple-or", testdata.MultipleOrDSL, testdata.MultipleOrSecureEncoderCode},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			RunHTTPDSL(t, c.DSL, testdata.TopLevelSchemes)
			fs := httpcodegen.ClientFiles("", httpdesign.Root)
			if len(fs) != 2 {
				t.Fatalf("got %d files, expected two", len(fs))
			}
			Generate("", []eval.Root{httpdesign.Root}, fs)
			for _, f := range fs {
				switch filepath.Base(f.Path) {
				case "encode_decode.go":
					sections := f.Section("secure-request-encoder")
					if len(sections) < 1 {
						t.Fatalf("got %d sections, expected at least 1", len(sections))
					}
					code := codegen.SectionCode(t, sections[0])
					if code != c.Code {
						t.Errorf("invalid code, got:\n%s\ngot vs. expected:\n%s", code, codegen.Diff(t, code, c.Code))
					}
				case "client.go":
					sections := f.Section("client-endpoint-init")
					if len(sections) < 1 {
						t.Fatalf("got %d sections, expected at least 1", len(sections))
					}
					code := codegen.SectionCode(t, sections[0])
					data := sections[0].Data.(*httpcodegen.EndpointData)
					if !strings.Contains(code, "Secure"+data.RequestEncoder) {
						t.Errorf("invalid code, got:\n%s\n expected %s in the code", strings.TrimSpace(code), "Secure"+data.RequestEncoder)
					}
				}
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
		{"basic-auth", testdata.BasicAuthDSL},
		{"oauth2", testdata.OAuth2DSL},
		{"oauth2-in-param", testdata.OAuth2InParamDSL},
		{"jwt", testdata.JWTDSL},
		{"api-key", testdata.APIKeyDSL},
		{"api-key-in-param", testdata.APIKeyInParamDSL},
		{"multiple-and", testdata.MultipleAndDSL},
		{"multiple-or", testdata.MultipleOrDSL},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			RunHTTPDSL(t, c.DSL, testdata.TopLevelSchemes)
			httpdesign.Root.Design.API = a
			f, err := httpcodegen.OpenAPIFiles(httpdesign.Root)
			if err != nil {
				t.Fatalf("error generating openapi file: %s", err)
			}
			for i := 0; i < len(f); i++ {
				OpenAPIV2(httpdesign.Root, f[i])
				s := f[i].SectionTemplates
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
			}
		})
	}
}

func TestExample(t *testing.T) {
	cases := []struct {
		Name     string
		DSL      func()
		Snippets []string
	}{
		{"single-service", testdata.SingleServiceDSL, []string{
			"singleservice.NewSecureEndpoints(singleserviceSvc, testapi.SingleServiceAuthAPIKeyFn)"}},
		{"multiple-services", testdata.MultipleServicesDSL, []string{
			"servicewithapikeyauth.NewSecureEndpoints(servicewithapikeyauthSvc, testapi.ServiceWithAPIKeyAuthAuthAPIKeyFn)",
			"servicewithjwtandbasicauth.NewSecureEndpoints(servicewithjwtandbasicauthSvc, testapi.ServiceWithJWTAndBasicAuthAuthBasicAuthFn, testapi.ServiceWithJWTAndBasicAuthAuthJWTFn)",
			"servicewithnosecurity.NewSecureEndpoints(servicewithnosecuritySvc)"}},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			httpcodegen.RunHTTPDSL(t, c.DSL)
			if len(goadesign.Root.Services) == 0 {
				t.Fatalf("expected at least 1 service")
			}
			fs := httpcodegen.ExampleServerFiles("", httpdesign.Root)
			Example("", []eval.Root{goadesign.Root, httpdesign.Root}, fs)
			for _, f := range fs {
				if filepath.Base(f.Path) != "main.go" {
					continue
				}
				sections := f.Section("service-main")
				if len(sections) < 1 {
					t.Fatalf("service-main: expected at least 1")
				}
				code := codegen.SectionCode(t, sections[0])
				for _, s := range c.Snippets {
					if !strings.Contains(code, s) {
						t.Errorf("invalid code, code:\n%s\ndoes not contain expected snippet:\n%s", code, s)
					}
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
