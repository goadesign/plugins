// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// secured_service HTTP server
//
// Command:
// $ goa gen goa.design/plugins/security/examples/multi_auth/design

package server

import (
	"context"
	"net/http"

	goa "goa.design/goa"
	goahttp "goa.design/goa/http"
	securedservice "goa.design/plugins/security/examples/multi_auth/gen/secured_service"
)

// Server lists the secured_service service endpoint HTTP handlers.
type Server struct {
	Mounts           []*MountPoint
	Signin           http.Handler
	Secure           http.Handler
	DoublySecure     http.Handler
	AlsoDoublySecure http.Handler
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

// New instantiates HTTP handlers for all the secured_service service endpoints.
func New(
	e *securedservice.Endpoints,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"Signin", "POST", "/signin"},
			{"Secure", "GET", "/secure"},
			{"DoublySecure", "PUT", "/secure"},
			{"AlsoDoublySecure", "POST", "/secure"},
		},
		Signin:           NewSigninHandler(e.Signin, mux, dec, enc),
		Secure:           NewSecureHandler(e.Secure, mux, dec, enc),
		DoublySecure:     NewDoublySecureHandler(e.DoublySecure, mux, dec, enc),
		AlsoDoublySecure: NewAlsoDoublySecureHandler(e.AlsoDoublySecure, mux, dec, enc),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "secured_service" }

// Mount configures the mux to serve the secured_service endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountSigninHandler(mux, h.Signin)
	MountSecureHandler(mux, h.Secure)
	MountDoublySecureHandler(mux, h.DoublySecure)
	MountAlsoDoublySecureHandler(mux, h.AlsoDoublySecure)
}

// MountSigninHandler configures the mux to serve the "secured_service" service
// "signin" endpoint.
func MountSigninHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/signin", f)
}

// NewSigninHandler creates a HTTP handler which loads the HTTP request and
// calls the "secured_service" service "signin" endpoint.
func NewSigninHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
) http.Handler {
	var (
		decodeRequest  = SecureDecodeSigninRequest(mux, dec)
		encodeResponse = EncodeSigninResponse(enc)
		encodeError    = EncodeSigninError(enc)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, accept)
		ctx = context.WithValue(ctx, goa.MethodKey, "signin")
		ctx = context.WithValue(ctx, goa.ServiceKey, "secured_service")
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

// MountSecureHandler configures the mux to serve the "secured_service" service
// "secure" endpoint.
func MountSecureHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/secure", f)
}

// NewSecureHandler creates a HTTP handler which loads the HTTP request and
// calls the "secured_service" service "secure" endpoint.
func NewSecureHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
) http.Handler {
	var (
		decodeRequest  = SecureDecodeSecureRequest(mux, dec)
		encodeResponse = EncodeSecureResponse(enc)
		encodeError    = goahttp.ErrorEncoder(enc)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, accept)
		ctx = context.WithValue(ctx, goa.MethodKey, "secure")
		ctx = context.WithValue(ctx, goa.ServiceKey, "secured_service")
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

// MountDoublySecureHandler configures the mux to serve the "secured_service"
// service "doubly_secure" endpoint.
func MountDoublySecureHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("PUT", "/secure", f)
}

// NewDoublySecureHandler creates a HTTP handler which loads the HTTP request
// and calls the "secured_service" service "doubly_secure" endpoint.
func NewDoublySecureHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
) http.Handler {
	var (
		decodeRequest  = SecureDecodeDoublySecureRequest(mux, dec)
		encodeResponse = EncodeDoublySecureResponse(enc)
		encodeError    = goahttp.ErrorEncoder(enc)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, accept)
		ctx = context.WithValue(ctx, goa.MethodKey, "doubly_secure")
		ctx = context.WithValue(ctx, goa.ServiceKey, "secured_service")
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

// MountAlsoDoublySecureHandler configures the mux to serve the
// "secured_service" service "also_doubly_secure" endpoint.
func MountAlsoDoublySecureHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/secure", f)
}

// NewAlsoDoublySecureHandler creates a HTTP handler which loads the HTTP
// request and calls the "secured_service" service "also_doubly_secure"
// endpoint.
func NewAlsoDoublySecureHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
) http.Handler {
	var (
		decodeRequest  = SecureDecodeAlsoDoublySecureRequest(mux, dec)
		encodeResponse = EncodeAlsoDoublySecureResponse(enc)
		encodeError    = goahttp.ErrorEncoder(enc)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, accept)
		ctx = context.WithValue(ctx, goa.MethodKey, "also_doubly_secure")
		ctx = context.WithValue(ctx, goa.ServiceKey, "secured_service")
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
