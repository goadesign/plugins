// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// archiver HTTP server
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/fetcher/archiver/design

package server

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	goa "goa.design/goa"
	goahttp "goa.design/goa/http"
	archiversvc "goa.design/plugins/goakit/examples/fetcher/archiver/gen/archiver"
)

// Server lists the archiver service endpoint HTTP handlers.
type Server struct {
	Mounts  []*MountPoint
	Archive http.Handler
	Read    http.Handler
}

// MountPoint holds information about the mounted endpoints.
type MountPoint struct {
	// Method is the name of the service method served by the mounted HTTP handler.
	Method string
	// Verb is the HTTP method used to match requests to the mounted handler.
	Verb string
	// Pattern is the HTTP request path pattern used to match requests to the
	// mounted handler.
	Pattern string
}

// New instantiates HTTP handlers for all the archiver service endpoints.
func New(
	e *archiversvc.Endpoints,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
	eh func(context.Context, http.ResponseWriter, error),
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"Archive", "POST", "/archive"},
			{"Read", "GET", "/archive/{id}"},
		},
		Archive: NewArchiveHandler(e.Archive, mux, dec, enc, eh),
		Read:    NewReadHandler(e.Read, mux, dec, enc, eh),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "archiver" }

// Use wraps the server handlers with the given middleware.
func (s *Server) Use(m func(http.Handler) http.Handler) {
	s.Archive = m(s.Archive)
	s.Read = m(s.Read)
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
	eh func(context.Context, http.ResponseWriter, error),
) http.Handler {
	var (
		decodeRequest  = DecodeArchiveRequest(mux, dec)
		encodeResponse = EncodeArchiveResponse(enc)
		encodeError    = goahttp.ErrorEncoder(enc)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "archive")
		ctx = context.WithValue(ctx, goa.ServiceKey, "archiver")
		payload, err := decodeRequest(r)
		if err != nil {
			eh(ctx, w, err)
			return
		}

		res, err := endpoint(ctx, payload)

		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				eh(ctx, w, err)
				return
			}
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			eh(ctx, w, err)
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
	eh func(context.Context, http.ResponseWriter, error),
) http.Handler {
	var (
		decodeRequest  = DecodeReadRequest(mux, dec)
		encodeResponse = EncodeReadResponse(enc)
		encodeError    = EncodeReadError(enc)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "read")
		ctx = context.WithValue(ctx, goa.ServiceKey, "archiver")
		payload, err := decodeRequest(r)
		if err != nil {
			eh(ctx, w, err)
			return
		}

		res, err := endpoint(ctx, payload)

		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				eh(ctx, w, err)
				return
			}
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			eh(ctx, w, err)
		}
	})
}
