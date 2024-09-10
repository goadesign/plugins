// Code generated by goa v3.19.0, DO NOT EDIT.
//
// Arnz HTTP server encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/v3/arnz/example/design -o
// $(GOPATH)/src/goa.design/plugins/arnz//example

package server

import (
	"context"
	"net/http"

	goahttp "goa.design/goa/v3/http"
	arnz "goa.design/plugins/v3/arnz/example/gen/arnz"
)

// EncodeCreateResponse returns an encoder for responses returned by the Arnz
// create endpoint.
func EncodeCreateResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, any) error {
	return func(ctx context.Context, w http.ResponseWriter, v any) error {
		res, _ := v.(*arnz.ResponseBody)
		enc := encoder(ctx, w)
		body := NewCreateResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// EncodeReadResponse returns an encoder for responses returned by the Arnz
// read endpoint.
func EncodeReadResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, any) error {
	return func(ctx context.Context, w http.ResponseWriter, v any) error {
		res, _ := v.(*arnz.ResponseBody)
		enc := encoder(ctx, w)
		body := NewReadResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// EncodeUpdateResponse returns an encoder for responses returned by the Arnz
// update endpoint.
func EncodeUpdateResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, any) error {
	return func(ctx context.Context, w http.ResponseWriter, v any) error {
		res, _ := v.(*arnz.ResponseBody)
		enc := encoder(ctx, w)
		body := NewUpdateResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// EncodeDeleteResponse returns an encoder for responses returned by the Arnz
// delete endpoint.
func EncodeDeleteResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, any) error {
	return func(ctx context.Context, w http.ResponseWriter, v any) error {
		res, _ := v.(*arnz.ResponseBody)
		enc := encoder(ctx, w)
		body := NewDeleteResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// EncodeHealthResponse returns an encoder for responses returned by the Arnz
// health endpoint.
func EncodeHealthResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, any) error {
	return func(ctx context.Context, w http.ResponseWriter, v any) error {
		res, _ := v.(*arnz.ResponseBody)
		enc := encoder(ctx, w)
		body := NewHealthResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}
