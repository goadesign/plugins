// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// calc HTTP server
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/calc/design

package server

import (
	"context"
	"net/http"

	goa "goa.design/goa"
	goahttp "goa.design/goa/http"
	calcsvc "goa.design/plugins/security/examples/calc/calc/gen/calc"
)

// Server lists the calc service endpoint HTTP handlers.
type Server struct {
	Mounts []*MountPoint
	Login  http.Handler
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

// New instantiates HTTP handlers for all the calc service endpoints.
func New(
	e *calcsvc.Endpoints,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"Login", "POST", "/login"},
			{"Add", "GET", "/add/{a}/{b}"},
		},
		Login: NewLoginHandler(e.Login, mux, dec, enc),
		Add:   NewAddHandler(e.Add, mux, dec, enc),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "calc" }

// Mount configures the mux to serve the calc endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountLoginHandler(mux, h.Login)
	MountAddHandler(mux, h.Add)
}

// MountLoginHandler configures the mux to serve the "calc" service "login"
// endpoint.
func MountLoginHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/login", f)
}

// NewLoginHandler creates a HTTP handler which loads the HTTP request and
// calls the "calc" service "login" endpoint.
func NewLoginHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
) http.Handler {
	var (
		decodeRequest  = SecureDecodeLoginRequest(mux, dec)
		encodeResponse = EncodeLoginResponse(enc)
		encodeError    = EncodeLoginError(enc)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, accept)
		ctx = context.WithValue(ctx, goa.MethodKey, "login")
		ctx = context.WithValue(ctx, goa.ServiceKey, "calc")
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

// MountAddHandler configures the mux to serve the "calc" service "add"
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
// the "calc" service "add" endpoint.
func NewAddHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
) http.Handler {
	var (
		decodeRequest  = SecureDecodeAddRequest(mux, dec)
		encodeResponse = EncodeAddResponse(enc)
		encodeError    = goahttp.ErrorEncoder(enc)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, accept)
		ctx = context.WithValue(ctx, goa.MethodKey, "add")
		ctx = context.WithValue(ctx, goa.ServiceKey, "calc")
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
