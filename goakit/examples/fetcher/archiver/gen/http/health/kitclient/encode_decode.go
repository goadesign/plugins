// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// health go-kit HTTP client encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/client/archiver/design

package client

import (
	"context"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	goahttp "goa.design/goa/http"
	"goa.design/plugins/goakit/examples/client/archiver/gen/http/health/client"
)

// DecodeShowResponse returns a go-kit DecodeResponseFunc suitable for decoding
// health show responses.
func DecodeShowResponse(decoder func(*http.Response) goahttp.Decoder) kithttp.DecodeResponseFunc {
	dec := client.DecodeShowResponse(decoder, false)
	return func(ctx context.Context, resp *http.Response) (interface{}, error) {
		return dec(resp)
	}
}
