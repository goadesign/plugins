// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// secured_service HTTP server
//
// Command:
// $ goa gen goa.design/plugins/security/example/design

package server

import (
	"context"
	"net/http"

	goa "goa.design/goa"
	goahttp "goa.design/goa/http"
	"goa.design/plugins/security/example/gen/securedservice"
)

// Server lists the secured_service service endpoint HTTP handlers.
type Server struct {
	Signin           http.Handler
	Secure           http.Handler
	DoublySecure     http.Handler
	AlsoDoublySecure http.Handler
}

// New instantiates HTTP handlers for all the secured_service service endpoints.
func New(
	e *securedservice.Endpoints,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
) *Server {
	return &Server{
		Signin:           NewSigninHandler(e.Signin, mux, dec, enc),
		Secure:           NewSecureHandler(e.Secure, mux, dec, enc),
		DoublySecure:     NewDoublySecureHandler(e.DoublySecure, mux, dec, enc),
		AlsoDoublySecure: NewAlsoDoublySecureHandler(e.AlsoDoublySecure, mux, dec, enc),
	}
}

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
		decodeRequest  = DecodeSigninRequest(mux, dec)
		encodeResponse = EncodeSigninResponse(enc)
		encodeError    = EncodeSigninError(enc)
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
		decodeRequest  = DecodeSecureRequest(mux, dec)
		encodeResponse = EncodeSecureResponse(enc)
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

// MountDoublySecureHandler configures the mux to serve the "secured_service"
// service "doubly_secure" endpoint.
func MountDoublySecureHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/secure", f)
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
		decodeRequest  = DecodeDoublySecureRequest(mux, dec)
		encodeResponse = EncodeDoublySecureResponse(enc)
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
		decodeRequest  = DecodeAlsoDoublySecureRequest(mux, dec)
		encodeResponse = EncodeAlsoDoublySecureResponse(enc)
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
