// This file writes go-kit request and response wrappers using the names chosen
// while Goa planned the generated client and server packages.
package goakit

import (
	"fmt"
	"path/filepath"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/expr"
	httpcodegen "goa.design/goa/v3/http/codegen"
)

type (
	// goakitEndpointData contains the chosen names rendered by go-kit templates
	// alongside the HTTP method data supplied by Goa.
	goakitEndpointData struct {
		*httpcodegen.EndpointData
		MountHandler, RequestDecoder, ResponseEncoder string
		ErrorEncoder, RequestEncoder, ResponseDecoder string
	}
)

// EncodeDecodeFiles produces a set of go-kit transport encoders and decoders
// that wrap the corresponding generated goa functions.
func EncodeDecodeFiles(genpkg string, root *expr.RootExpr) []*codegen.File {
	plan, err := planHTTP(genpkg, root)
	if err != nil {
		panic(err)
	}
	return encodeDecodeFiles(genpkg, plan)
}

// encodeDecodeFiles builds go-kit codecs from the HTTP names chosen for the
// current generation run.
func encodeDecodeFiles(genpkg string, plan *goakitRootPlan) []*codegen.File {
	fw := make([]*codegen.File, 2*len(plan.root.API.HTTP.Services))
	for i, service := range plan.root.API.HTTP.Services {
		data := httpServiceData(plan.http, service)
		fw[i] = serverEncodeDecode(genpkg, service, data, plan.services[service])
	}
	for i, service := range plan.root.API.HTTP.Services {
		data := httpServiceData(plan.http, service)
		fw[i+len(plan.root.API.HTTP.Services)] = clientEncodeDecode(genpkg, service, data, plan.services[service])
	}
	return fw
}

// httpServiceData returns the linked HTTP data for a service already listed
// by the same design root. Its absence means the Goa plan is inconsistent.
func httpServiceData(plan *httpcodegen.Plan, service *expr.HTTPServiceExpr) *httpcodegen.ServiceData {
	data, ok := plan.Service(service)
	if !ok {
		panic(fmt.Sprintf("goakit: HTTP service %q is missing from its plan", service.Name()))
	}
	return data
}

// serverEncodeDecode returns the file defining the go-kit HTTP server encoding
// and decoding logic.
func serverEncodeDecode(
	genpkg string,
	svc *expr.HTTPServiceExpr,
	data *httpcodegen.ServiceData,
	names *goakitServicePlan,
) *codegen.File {
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
		endpoint := plannedEndpointData(e, names.endpoints[e.Method.Name])
		sections = append(sections, &codegen.SectionTemplate{
			Name:   "goakit-response-encoder",
			Source: responseEncoderT,
			Data:   endpoint,
		})

		if e.Payload.Ref != "" {
			sections = append(sections, &codegen.SectionTemplate{
				Name:   "goakit-request-decoder",
				Source: requestDecoderT,
				Data:   endpoint,
			})
		}

		if len(e.Errors) > 0 {
			sections = append(sections, &codegen.SectionTemplate{
				Name:   "goakit-error-encoder",
				Source: errorEncoderT,
				Data:   endpoint,
			})
		}
	}

	return &codegen.File{Path: path, SectionTemplates: sections}
}

// clientEncodeDecode returns the file defining the go-kit HTTP client encoding
// and decoding logic.
func clientEncodeDecode(
	genpkg string,
	svc *expr.HTTPServiceExpr,
	data *httpcodegen.ServiceData,
	names *goakitServicePlan,
) *codegen.File {
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
		endpoint := plannedEndpointData(e, names.endpoints[e.Method.Name])
		if e.RequestEncoder != "" {
			sections = append(sections, &codegen.SectionTemplate{
				Name:   "goakit-request-encoder",
				Source: requestEncoderT,
				Data:   endpoint,
			})
		}
		if e.Result != nil || len(e.Errors) > 0 {
			sections = append(sections, &codegen.SectionTemplate{
				Name:   "goakit-response-decoder",
				Source: responseDecoderT,
				Data:   endpoint,
			})
		}
	}

	return &codegen.File{Path: path, SectionTemplates: sections}
}

// plannedEndpointData combines Goa's linked HTTP values with the wrapper names
// declared by this plugin before generation was frozen.
func plannedEndpointData(endpoint *httpcodegen.EndpointData, names *goakitEndpointNames) *goakitEndpointData {
	data := &goakitEndpointData{EndpointData: endpoint}
	if names.mountHandler != nil {
		data.MountHandler = names.mountHandler.Name()
	}
	if names.requestDecoder != nil {
		data.RequestDecoder = names.requestDecoder.Name()
	}
	if names.responseEncoder != nil {
		data.ResponseEncoder = names.responseEncoder.Name()
	}
	if names.errorEncoder != nil {
		data.ErrorEncoder = names.errorEncoder.Name()
	}
	if names.requestEncoder != nil {
		data.RequestEncoder = names.requestEncoder.Name()
	}
	if names.responseDecoder != nil {
		data.ResponseDecoder = names.responseDecoder.Name()
	}
	return data
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
