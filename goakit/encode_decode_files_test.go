package goakit

import (
	"strings"
	"testing"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/expr"
	httpcodegen "goa.design/goa/v3/http/codegen"
	"goa.design/plugins/v3/goakit/testdata"
)

func TestServerEncodeDecode(t *testing.T) {
	cases := map[string]struct {
		DSL    func()
		Code   map[string][]string
		Import string
	}{
		"simple-service": {
			DSL: testdata.SimpleServiceDSL,
			Code: map[string][]string{
				"goakit-response-encoder": {testdata.SimpleMethodGoakitResponseEncoderCode},
				"goakit-request-decoder":  {},
				"goakit-error-encoder":    {},
			},
			Import: "/http/simple_service/server",
		},
		"with-payload": {
			DSL: testdata.WithPayloadDSL,
			Code: map[string][]string{
				"goakit-response-encoder": {testdata.WithPayloadMethodGoakitResponseEncoderCode},
				"goakit-request-decoder":  {testdata.WithPayloadMethodGoakitRequestDecoderCode},
				"goakit-error-encoder":    {},
			},
			Import: "/http/with_payload_service/server",
		},
		"with-error": {
			DSL: testdata.WithErrorDSL,
			Code: map[string][]string{
				"goakit-response-encoder": {testdata.WithErrorMethodGoakitResponseEncoderCode},
				"goakit-request-decoder":  {},
				"goakit-error-encoder":    {testdata.WithErrorMethodGoakitErrorEncoderCode},
			},
			Import: "/http/with_error_service/server",
		},
		"multi-endpoints": {
			DSL: testdata.MultiEndpointDSL,
			Code: map[string][]string{
				"goakit-response-encoder": {testdata.Endpoint1GoakitResponseEncoderCode, testdata.Endpoint2GoakitResponseEncoderCode},
				"goakit-request-decoder":  {testdata.Endpoint1GoakitRequestDecoderCode},
				"goakit-error-encoder":    {testdata.Endpoint1GoakitErrorEncoderCode, testdata.Endpoint2GoakitErrorEncoderCode},
			},
			Import: "/http/multi_endpoint_service/server",
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
					requireImport(t, f, c.Import)
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
		DSL    func()
		Code   map[string][]string
		Import string
	}{
		"simple-service": {
			DSL: testdata.SimpleServiceDSL,
			Code: map[string][]string{
				"goakit-response-decoder": {testdata.SimpleMethodGoakitResponseDecoderCode},
				"goakit-request-encoder":  {},
			},
			Import: "/http/simple_service/client",
		},
		"with-payload": {
			DSL: testdata.WithPayloadDSL,
			Code: map[string][]string{
				"goakit-response-decoder": {testdata.WithPayloadMethodGoakitResponseDecoderCode},
				"goakit-request-encoder":  {testdata.WithPayloadMethodGoakitRequestEncoderCode},
			},
			Import: "/http/with_payload_service/client",
		},
		"with-error": {
			DSL: testdata.WithErrorDSL,
			Code: map[string][]string{
				"goakit-response-decoder": {testdata.WithErrorMethodGoakitResponseDecoderCode},
				"goakit-request-encoder":  {},
			},
			Import: "/http/with_error_service/client",
		},
		"multi-endpoints": {
			DSL: testdata.MultiEndpointDSL,
			Code: map[string][]string{
				"goakit-response-decoder": {testdata.Endpoint1GoakitResponseDecoderCode, testdata.Endpoint2GoakitResponseDecoderCode},
				"goakit-request-encoder":  {testdata.Endpoint1GoakitRequestEncoderCode},
			},
			Import: "/http/multi_endpoint_service/client",
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
					requireImport(t, f, c.Import)
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
