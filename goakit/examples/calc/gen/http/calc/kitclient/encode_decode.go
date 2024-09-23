// Code generated by goa v3.19.1, DO NOT EDIT.
//
// calc go-kit HTTP client encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/v3/goakit/examples/calc/design -o
// $(GOPATH)/src/goa.design/plugins/goakit//examples/calc

package client

import (
	"context"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	goahttp "goa.design/goa/v3/http"
	"goa.design/plugins/v3/goakit/examples/calc/gen/http/calc/client"
)

// DecodeAddResponse returns a go-kit DecodeResponseFunc suitable for decoding
// calc add responses.
func DecodeAddResponse(decoder func(*http.Response) goahttp.Decoder) kithttp.DecodeResponseFunc {
	dec := client.DecodeAddResponse(decoder, false)
	return func(ctx context.Context, resp *http.Response) (interface{}, error) {
		return dec(resp)
	}
}
