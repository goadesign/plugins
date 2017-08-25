// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// storage go-kit HTTP server encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/cellar/design

package server

import (
	"context"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	goahttp "goa.design/goa/http"
	"goa.design/plugins/goakit/examples/cellar/gen/http/storage/server"
)

// EncodeListResponse returns a go-kit EncodeResponseFunc suitable for encoding
// storage list responses.
func EncodeListResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) kithttp.EncodeResponseFunc {
	return server.EncodeListResponse(encoder)
}

// EncodeShowResponse returns a go-kit EncodeResponseFunc suitable for encoding
// storage show responses.
func EncodeShowResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) kithttp.EncodeResponseFunc {
	return server.EncodeShowResponse(encoder)
}

// DecodeShowRequest returns a go-kit DecodeRequestFunc suitable for decoding
// storage show requests.
func DecodeShowRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) kithttp.DecodeRequestFunc {
	dec := server.DecodeShowRequest(mux, decoder)
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		r = r.WithContext(ctx)
		return dec(r)
	}
}

// EncodeShowResponse returns a go-kit EncodeResponseFunc suitable for encoding
// errors returned by the storage show endpoint.
func EncodeShowError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) kithttp.EncodeResponseFunc {
	enc := server.EncodeShowError(encoder)
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		enc(ctx, w, v.(error))
		return nil
	}
}

// EncodeAddResponse returns a go-kit EncodeResponseFunc suitable for encoding
// storage add responses.
func EncodeAddResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) kithttp.EncodeResponseFunc {
	return server.EncodeAddResponse(encoder)
}

// DecodeAddRequest returns a go-kit DecodeRequestFunc suitable for decoding
// storage add requests.
func DecodeAddRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) kithttp.DecodeRequestFunc {
	dec := server.DecodeAddRequest(mux, decoder)
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		r = r.WithContext(ctx)
		return dec(r)
	}
}

// EncodeRemoveResponse returns a go-kit EncodeResponseFunc suitable for
// encoding storage remove responses.
func EncodeRemoveResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) kithttp.EncodeResponseFunc {
	return server.EncodeRemoveResponse(encoder)
}

// DecodeRemoveRequest returns a go-kit DecodeRequestFunc suitable for decoding
// storage remove requests.
func DecodeRemoveRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) kithttp.DecodeRequestFunc {
	dec := server.DecodeRemoveRequest(mux, decoder)
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		r = r.WithContext(ctx)
		return dec(r)
	}
}
