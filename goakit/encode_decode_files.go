package goakit

import (
	"fmt"
	"path/filepath"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/expr"
	httpcodegen "goa.design/goa/v3/http/codegen"
)

// EncodeDecodeFiles produces a set of go-kit transport encoders and decoders
// that wrap the corresponding generated goa functions.
func EncodeDecodeFiles(genpkg string, root *expr.RootExpr) []*codegen.File {
	fw := make([]*codegen.File, 2*len(root.API.HTTP.Services))
	for i, r := range root.API.HTTP.Services {
		fw[i] = serverEncodeDecode(genpkg, r)
	}
	for i, r := range root.API.HTTP.Services {
		fw[i+len(root.API.HTTP.Services)] = clientEncodeDecode(genpkg, r)
	}
	return fw
}

// serverEncodeDecode returns the file defining the go-kit HTTP server encoding
// and decoding logic.
func serverEncodeDecode(genpkg string, svc *expr.HTTPServiceExpr) *codegen.File {
	data := httpcodegen.HTTPServices.Get(svc.Name())
	svcName := data.Service.PathName
	path := filepath.Join(codegen.Gendir, "http", svcName, "kitserver", "encode_decode.go")
	title := fmt.Sprintf("%s go-kit HTTP server encoders and decoders", svc.Name())
	sections := []*codegen.SectionTemplate{
		codegen.Header(title, "server", []*codegen.ImportSpec{
			{Path: "context"},
			{Path: "net/http"},
			{Path: "strings"},
			{Path: "github.com/go-kit/kit/transport/http", Name: "kithttp"},
			{Path: "goa.design/goa/v3", Name: "goa"},
			{Path: "goa.design/goa/v3/http", Name: "goahttp"},
			{Path: genpkg + "/http/" + svcName + "/server"},
		}),
	}

	for _, e := range data.Endpoints {
		sections = append(sections, &codegen.SectionTemplate{
			Name:   "goakit-response-encoder",
			Source: responseEncoderT,
			Data:   e,
		})

		if e.Payload.Ref != "" {
			sections = append(sections, &codegen.SectionTemplate{
				Name:   "goakit-request-decoder",
				Source: requestDecoderT,
				Data:   e,
			})
		}

		if len(e.Errors) > 0 {
			sections = append(sections, &codegen.SectionTemplate{
				Name:   "goakit-error-encoder",
				Source: errorEncoderT,
				Data:   e,
			})
		}
	}

	return &codegen.File{Path: path, SectionTemplates: sections}
}

// clientEncodeDecode returns the file defining the go-kit HTTP client encoding
// and decoding logic.
func clientEncodeDecode(genpkg string, svc *expr.HTTPServiceExpr) *codegen.File {
	data := httpcodegen.HTTPServices.Get(svc.Name())
	svcName := data.Service.PathName
	path := filepath.Join(codegen.Gendir, "http", svcName, "kitclient", "encode_decode.go")
	title := fmt.Sprintf("%s go-kit HTTP client encoders and decoders", svc.Name())
	sections := []*codegen.SectionTemplate{
		codegen.Header(title, "client", []*codegen.ImportSpec{
			{Path: "context"},
			{Path: "net/http"},
			{Path: "strings"},
			{Path: "github.com/go-kit/kit/transport/http", Name: "kithttp"},
			{Path: "goa.design/goa/v3", Name: "goa"},
			{Path: "goa.design/goa/v3/http", Name: "goahttp"},
			{Path: genpkg + "/http/" + svcName + "/client"},
		}),
	}

	for _, e := range data.Endpoints {
		if e.RequestEncoder != "" {
			sections = append(sections, &codegen.SectionTemplate{
				Name:   "goakit-request-encoder",
				Source: requestEncoderT,
				Data:   e,
			})
		}
		if e.Result != nil || len(e.Errors) > 0 {
			sections = append(sections, &codegen.SectionTemplate{
				Name:   "goakit-response-decoder",
				Source: responseDecoderT,
				Data:   e,
			})
		}
	}

	return &codegen.File{Path: path, SectionTemplates: sections}
}

// input: EndpointData
const requestEncoderT = `{{ printf "%s returns a go-kit EncodeRequestFunc suitable for encoding %s %s requests." .RequestEncoder .ServiceName .Method.Name | comment }}
func {{ .RequestEncoder }}(encoder func(*http.Request) goahttp.Encoder) kithttp.EncodeRequestFunc {
	enc := client.{{ .RequestEncoder }}(encoder)
	return func(_ context.Context, r *http.Request, v interface{}) error {
		return enc(r, v)
	}
}
`

// input: EndpointData
const requestDecoderT = `{{ printf "%s returns a go-kit DecodeRequestFunc suitable for decoding %s %s requests." .RequestDecoder .ServiceName .Method.Name | comment }}
func {{ .RequestDecoder }}(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) kithttp.DecodeRequestFunc {
	dec := server.{{ .RequestDecoder }}(mux, decoder)
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		r = r.WithContext(ctx)
		return dec(r)
	}
}
`

// input: EndpointData
const responseEncoderT = `{{ printf "%s returns a go-kit EncodeResponseFunc suitable for encoding %s %s responses." .ResponseEncoder .ServiceName .Method.Name | comment }}
 func {{ .ResponseEncoder }}(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) kithttp.EncodeResponseFunc {
 	return server.{{ .ResponseEncoder }}(encoder)
 }
`

// input: EndpointData
const errorEncoderT = `{{ printf "%s returns a go-kit EncodeResponseFunc suitable for encoding errors returned by the %s %s endpoint." .ErrorEncoder .ServiceName .Method.Name | comment }}
 func {{ .ErrorEncoder }}(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(ctx context.Context, err error) goahttp.Statuser) kithttp.ErrorEncoder {
 	enc := server.{{ .ErrorEncoder }}(encoder, formatter)
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		enc(ctx, w, err)
	}
}
`

// input: EndpointData
const responseDecoderT = `{{ printf "%s returns a go-kit DecodeResponseFunc suitable for decoding %s %s responses." .ResponseDecoder .ServiceName .Method.Name | comment }}
func {{ .ResponseDecoder }}(decoder func(*http.Response) goahttp.Decoder) kithttp.DecodeResponseFunc {
	dec := client.{{ .ResponseDecoder }}(decoder, false)
	return func(ctx context.Context, resp *http.Response) (interface{}, error) {
		return dec(resp)
	}
}
`
