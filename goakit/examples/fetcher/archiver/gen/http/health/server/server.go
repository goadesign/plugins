// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// health HTTP server
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/fetcher/archiver/design

package server

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	goahttp "goa.design/goa/http"
	health "goa.design/plugins/goakit/examples/fetcher/archiver/gen/health"
)

// Server lists the health service endpoint HTTP handlers.
type Server struct {
	Show http.Handler
}

// New instantiates HTTP handlers for all the health service endpoints.
func New(
	e *health.Endpoints,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
) *Server {
	return &Server{
		Show: NewShowHandler(e.Show, mux, dec, enc),
	}
}

// Mount configures the mux to serve the health endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountShowHandler(mux, h.Show)
}

// MountShowHandler configures the mux to serve the "health" service "show"
// endpoint.
func MountShowHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/health", f)
}

// NewShowHandler creates a HTTP handler which loads the HTTP request and calls
// the "health" service "show" endpoint.
func NewShowHandler(
	endpoint endpoint.Endpoint,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
) http.Handler {
	var (
		encodeResponse = EncodeShowResponse(enc)
		encodeError    = goahttp.ErrorEncoder(enc)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		ctx := context.WithValue(r.Context(), goahttp.ContextKeyAcceptType, accept)
		res, err := endpoint(ctx, nil)

		if err != nil {
			encodeError(ctx, w, err)
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			encodeError(ctx, w, err)
		}
	})
}
