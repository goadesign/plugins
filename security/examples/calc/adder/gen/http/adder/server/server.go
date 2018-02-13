// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// adder HTTP server
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/adder/design

package server

import (
	"context"
	"net/http"

	goa "goa.design/goa"
	goahttp "goa.design/goa/http"
	addersvc "goa.design/plugins/security/examples/calc/adder/gen/adder"
)

// Server lists the adder service endpoint HTTP handlers.
type Server struct {
	Mounts []*MountPoint
	Add    http.Handler
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

// New instantiates HTTP handlers for all the adder service endpoints.
func New(
	e *addersvc.Endpoints,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"Add", "GET", "/add/{a}/{b}"},
		},
		Add: NewAddHandler(e.Add, mux, dec, enc),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "adder" }

// Mount configures the mux to serve the adder endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountAddHandler(mux, h.Add)
}

// MountAddHandler configures the mux to serve the "adder" service "add"
// endpoint.
func MountAddHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/add/{a}/{b}", f)
}

// NewAddHandler creates a HTTP handler which loads the HTTP request and calls
// the "adder" service "add" endpoint.
func NewAddHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
) http.Handler {
	var (
		decodeRequest  = SecureDecodeAddRequest(mux, dec)
		encodeResponse = EncodeAddResponse(enc)
		encodeError    = EncodeAddError(enc)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		ctx := context.WithValue(r.Context(), goahttp.ContextKeyAcceptType, accept)
		ctx = context.WithValue(ctx, goa.ContextKeyMethod, "add")
		ctx = context.WithValue(ctx, goa.ContextKeyService, "adder")
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
