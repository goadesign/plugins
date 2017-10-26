// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// archiver HTTP server
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/client/archiver/design

package server

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	goahttp "goa.design/goa/http"
	archiver "goa.design/plugins/goakit/examples/client/archiver/gen/archiver"
)

// Server lists the archiver service endpoint HTTP handlers.
type Server struct {
	Archive http.Handler
	Read    http.Handler
}

// New instantiates HTTP handlers for all the archiver service endpoints.
func New(
	e *archiver.Endpoints,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
) *Server {
	return &Server{
		Archive: NewArchiveHandler(e.Archive, mux, dec, enc),
		Read:    NewReadHandler(e.Read, mux, dec, enc),
	}
}

// Mount configures the mux to serve the archiver endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountArchiveHandler(mux, h.Archive)
	MountReadHandler(mux, h.Read)
}

// MountArchiveHandler configures the mux to serve the "archiver" service
// "archive" endpoint.
func MountArchiveHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/archive", f)
}

// NewArchiveHandler creates a HTTP handler which loads the HTTP request and
// calls the "archiver" service "archive" endpoint.
func NewArchiveHandler(
	endpoint endpoint.Endpoint,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
) http.Handler {
	var (
		decodeRequest  = DecodeArchiveRequest(mux, dec)
		encodeResponse = EncodeArchiveResponse(enc)
		encodeError    = goahttp.ErrorEncoder(enc)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		ctx := context.WithValue(r.Context(), goahttp.ContextKeyAcceptType, accept)
		payload, err := decodeRequest(r)
		if err != nil {
			encodeError(ctx, w, err)
			return
		}

		res, err := endpoint(ctx, payload)

		if err != nil {
			encodeError(ctx, w, err)
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			encodeError(ctx, w, err)
		}
	})
}

// MountReadHandler configures the mux to serve the "archiver" service "read"
// endpoint.
func MountReadHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/archive/{id}", f)
}

// NewReadHandler creates a HTTP handler which loads the HTTP request and calls
// the "archiver" service "read" endpoint.
func NewReadHandler(
	endpoint endpoint.Endpoint,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
) http.Handler {
	var (
		decodeRequest  = DecodeReadRequest(mux, dec)
		encodeResponse = EncodeReadResponse(enc)
		encodeError    = EncodeReadError(enc)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		ctx := context.WithValue(r.Context(), goahttp.ContextKeyAcceptType, accept)
		payload, err := decodeRequest(r)
		if err != nil {
			encodeError(ctx, w, err)
			return
		}

		res, err := endpoint(ctx, payload)

		if err != nil {
			encodeError(ctx, w, err)
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			encodeError(ctx, w, err)
		}
	})
}
