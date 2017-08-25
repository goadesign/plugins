// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// storage go-kit HTTP client encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/cellar/design

package client

import (
	"context"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	goahttp "goa.design/goa/http"
	"goa.design/plugins/goakit/examples/cellar/gen/http/storage/client"
)

// DecodeListResponse returns a go-kit DecodeResponseFunc suitable for decoding
// storage list responses.
func DecodeListResponse(decoder func(*http.Response) goahttp.Decoder) kithttp.DecodeResponseFunc {
	dec := client.DecodeListResponse(decoder, false)
	return func(ctx context.Context, resp *http.Response) (interface{}, error) {
		return dec(resp)
	}
}

// DecodeShowResponse returns a go-kit DecodeResponseFunc suitable for decoding
// storage show responses.
func DecodeShowResponse(decoder func(*http.Response) goahttp.Decoder) kithttp.DecodeResponseFunc {
	dec := client.DecodeShowResponse(decoder, false)
	return func(ctx context.Context, resp *http.Response) (interface{}, error) {
		return dec(resp)
	}
}

// EncodeAddRequest returns a go-kit EncodeRequestFunc suitable for encoding
// storage add requests.
func EncodeAddRequest(encoder func(*http.Request) goahttp.Encoder) kithttp.EncodeRequestFunc {
	enc := client.EncodeAddRequest(encoder)
	return func(_ context.Context, r *http.Request, v interface{}) error {
		return enc(r, v)
	}
}

// DecodeAddResponse returns a go-kit DecodeResponseFunc suitable for decoding
// storage add responses.
func DecodeAddResponse(decoder func(*http.Response) goahttp.Decoder) kithttp.DecodeResponseFunc {
	dec := client.DecodeAddResponse(decoder, false)
	return func(ctx context.Context, resp *http.Response) (interface{}, error) {
		return dec(resp)
	}
}

// DecodeRemoveResponse returns a go-kit DecodeResponseFunc suitable for
// decoding storage remove responses.
func DecodeRemoveResponse(decoder func(*http.Response) goahttp.Decoder) kithttp.DecodeResponseFunc {
	dec := client.DecodeRemoveResponse(decoder, false)
	return func(ctx context.Context, resp *http.Response) (interface{}, error) {
		return dec(resp)
	}
}
