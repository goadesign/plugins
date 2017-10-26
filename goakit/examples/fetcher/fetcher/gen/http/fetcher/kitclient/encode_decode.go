// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// fetcher go-kit HTTP client encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/client/fetcher/design

package client

import (
	"context"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	goahttp "goa.design/goa/http"
	"goa.design/plugins/goakit/examples/client/fetcher/gen/http/fetcher/client"
)

// DecodeFetchResponse returns a go-kit DecodeResponseFunc suitable for
// decoding fetcher fetch responses.
func DecodeFetchResponse(decoder func(*http.Response) goahttp.Decoder) kithttp.DecodeResponseFunc {
	dec := client.DecodeFetchResponse(decoder, false)
	return func(ctx context.Context, resp *http.Response) (interface{}, error) {
		return dec(resp)
	}
}
