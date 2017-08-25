// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// sommelier go-kit HTTP client encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/cellar/design

package client

import (
	"context"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	goahttp "goa.design/goa/http"
	"goa.design/plugins/goakit/examples/cellar/gen/http/sommelier/client"
)

// EncodePickRequest returns a go-kit EncodeRequestFunc suitable for encoding
// sommelier pick requests.
func EncodePickRequest(encoder func(*http.Request) goahttp.Encoder) kithttp.EncodeRequestFunc {
	enc := client.EncodePickRequest(encoder)
	return func(_ context.Context, r *http.Request, v interface{}) error {
		return enc(r, v)
	}
}

// DecodePickResponse returns a go-kit DecodeResponseFunc suitable for decoding
// sommelier pick responses.
func DecodePickResponse(decoder func(*http.Response) goahttp.Decoder) kithttp.DecodeResponseFunc {
	dec := client.DecodePickResponse(decoder, false)
	return func(ctx context.Context, resp *http.Response) (interface{}, error) {
		return dec(resp)
	}
}
