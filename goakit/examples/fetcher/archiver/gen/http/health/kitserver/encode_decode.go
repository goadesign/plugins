// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// health go-kit HTTP server encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/fetcher/archiver/design

package server

import (
	"context"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	goahttp "goa.design/goa/http"
	"goa.design/plugins/goakit/examples/fetcher/archiver/gen/http/health/server"
)

// EncodeShowResponse returns a go-kit EncodeResponseFunc suitable for encoding
// health show responses.
func EncodeShowResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) kithttp.EncodeResponseFunc {
	return server.EncodeShowResponse(encoder)
}
