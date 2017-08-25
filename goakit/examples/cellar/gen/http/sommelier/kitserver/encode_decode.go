// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// sommelier go-kit HTTP server encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/cellar/design

package server

import (
	"context"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	goahttp "goa.design/goa/http"
	"goa.design/plugins/goakit/examples/cellar/gen/http/sommelier/server"
)

// EncodePickResponse returns a go-kit EncodeResponseFunc suitable for encoding
// sommelier pick responses.
func EncodePickResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) kithttp.EncodeResponseFunc {
	return server.EncodePickResponse(encoder)
}

// DecodePickRequest returns a go-kit DecodeRequestFunc suitable for decoding
// sommelier pick requests.
func DecodePickRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) kithttp.DecodeRequestFunc {
	dec := server.DecodePickRequest(mux, decoder)
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		r = r.WithContext(ctx)
		return dec(r)
	}
}

// EncodePickResponse returns a go-kit EncodeResponseFunc suitable for encoding
// errors returned by the sommelier pick endpoint.
func EncodePickError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) kithttp.EncodeResponseFunc {
	enc := server.EncodePickError(encoder)
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		enc(ctx, w, v.(error))
		return nil
	}
}
