package http

import (
	"path/filepath"
	"strings"
	"testing"

	"goa.design/goa/codegen"
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
