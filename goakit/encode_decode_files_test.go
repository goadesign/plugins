package goakit

import (
	"strings"
	"testing"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/expr"
	httpcodegen "goa.design/goa/v3/http/codegen"
	"goa.design/plugins/goakit/testdata"
)

func TestServerEncodeDecode(t *testing.T) {
	cases := map[string]struct {
		DSL  func()
		Code map[string][]string
	}{
		"simple-service": {
			DSL: testdata.SimpleServiceDSL,
			Code: map[string][]string{
				"goakit-response-encoder": []string{testdata.SimpleMethodGoakitResponseEncoderCode},
				"goakit-request-decoder":  []string{},
				"goakit-error-encoder":    []string{},
			},
		},
		"with-payload": {
			DSL: testdata.WithPayloadDSL,
			Code: map[string][]string{
				"goakit-response-encoder": []string{testdata.WithPayloadMethodGoakitResponseEncoderCode},
				"goakit-request-decoder":  []string{testdata.WithPayloadMethodGoakitRequestDecoderCode},
				"goakit-error-encoder":    []string{},
			},
		},
		"with-error": {
			DSL: testdata.WithErrorDSL,
			Code: map[string][]string{
				"goakit-response-encoder": []string{testdata.WithErrorMethodGoakitResponseEncoderCode},
				"goakit-request-decoder":  []string{},
				"goakit-error-encoder":    []string{testdata.WithErrorMethodGoakitErrorEncoderCode},
			},
		},
		"multi-endpoints": {
			DSL: testdata.MultiEndpointDSL,
			Code: map[string][]string{
				"goakit-response-encoder": []string{testdata.Endpoint1GoakitResponseEncoderCode, testdata.Endpoint2GoakitResponseEncoderCode},
				"goakit-request-decoder":  []string{testdata.Endpoint1GoakitRequestDecoderCode},
				"goakit-error-encoder":    []string{testdata.Endpoint1GoakitErrorEncoderCode, testdata.Endpoint2GoakitErrorEncoderCode},
			},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			httpcodegen.RunHTTPDSL(t, c.DSL)
			fs := EncodeDecodeFiles("", expr.Root)
			if len(fs) != 2 {
				t.Fatalf("got %d files, expected 2", len(fs))
			}
			var found bool
			for _, f := range fs {
				if strings.Contains(f.Path, "kitserver") {
					found = true
					for sec, secCode := range c.Code {
						testCode(t, f, sec, secCode)
					}
				}
			}
			if !found {
				t.Fatalf("kitserver encode_decode.go file not found")
			}
		})
	}
}

func TestClientEncodeDecode(t *testing.T) {
	cases := map[string]struct {
		DSL  func()
		Code map[string][]string
	}{
		"simple-service": {
			DSL: testdata.SimpleServiceDSL,
			Code: map[string][]string{
				"goakit-response-decoder": []string{testdata.SimpleMethodGoakitResponseDecoderCode},
				"goakit-request-encoder":  []string{},
			},
		},
		"with-payload": {
			DSL: testdata.WithPayloadDSL,
			Code: map[string][]string{
				"goakit-response-decoder": []string{testdata.WithPayloadMethodGoakitResponseDecoderCode},
				"goakit-request-encoder":  []string{testdata.WithPayloadMethodGoakitRequestEncoderCode},
			},
		},
		"with-error": {
			DSL: testdata.WithErrorDSL,
			Code: map[string][]string{
				"goakit-response-decoder": []string{testdata.WithErrorMethodGoakitResponseDecoderCode},
				"goakit-request-encoder":  []string{},
			},
		},
		"multi-endpoints": {
			DSL: testdata.MultiEndpointDSL,
			Code: map[string][]string{
				"goakit-response-decoder": []string{testdata.Endpoint1GoakitResponseDecoderCode, testdata.Endpoint2GoakitResponseDecoderCode},
				"goakit-request-encoder":  []string{testdata.Endpoint1GoakitRequestEncoderCode},
			},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			httpcodegen.RunHTTPDSL(t, c.DSL)
			fs := EncodeDecodeFiles("", expr.Root)
			if len(fs) != 2 {
				t.Fatalf("got %d files, expected 2", len(fs))
			}
			var found bool
			for _, f := range fs {
				if strings.Contains(f.Path, "kitclient") {
					found = true
					for sec, secCode := range c.Code {
						testCode(t, f, sec, secCode)
					}
				}
			}
			if !found {
				t.Fatalf("kitserver encode_decode.go file not found")
			}
		})
	}
}

func testCode(t *testing.T, file *codegen.File, section string, expCode []string) {
	sections := file.Section(section)
	if len(sections) != len(expCode) {
		t.Fatalf("%s: got %d sections, expected %d", section, len(sections), len(expCode))
	}
	for i, c := range expCode {
		code := codegen.SectionCode(t, sections[i])
		if code != c {
			t.Errorf("invalid code, got:\n%s\ngot vs. expected:\n%s", code, codegen.Diff(t, code, c))
		}
	}
}
