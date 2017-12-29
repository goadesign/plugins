// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// sommelier HTTP server
//
// Command:
// $ goa gen goa.design/plugins/cors/examples/cellar/design

package server

import (
	"context"
	"net/http"
	"regexp"

	goa "goa.design/goa"
	goahttp "goa.design/goa/http"
	"goa.design/plugins/cors"
	sommelier "goa.design/plugins/cors/examples/cellar/gen/sommelier"
)

// Server lists the sommelier service endpoint HTTP handlers.
type Server struct {
	Mounts []*MountPoint
	Pick   http.Handler
	CORS   http.Handler
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

// New instantiates HTTP handlers for all the sommelier service endpoints.
func New(
	e *sommelier.Endpoints,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"Pick", "POST", "/sommelier"},
			{"CORS", "OPTIONS", "/sommelier"},
		},
		Pick: NewPickHandler(e.Pick, mux, dec, enc),
		CORS: NewCORSHandler(),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "sommelier" }

// Mount configures the mux to serve the sommelier endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountPickHandler(mux, h.Pick)
	MountCORSHandler(mux, h.CORS)
}

// MountPickHandler configures the mux to serve the "sommelier" service "pick"
// endpoint.
func MountPickHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleSommelierOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/sommelier", f)
}

// NewPickHandler creates a HTTP handler which loads the HTTP request and calls
// the "sommelier" service "pick" endpoint.
func NewPickHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
) http.Handler {
	var (
		decodeRequest  = DecodePickRequest(mux, dec)
		encodeResponse = EncodePickResponse(enc)
		encodeError    = EncodePickError(enc)
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

// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service sommelier.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = handleSommelierOrigin(h)
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("OPTIONS", "/sommelier", f)
}

// NewCORSHandler creates a HTTP handler which returns a simple 200 response.
func NewCORSHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
}

// handleSommelierOrigin applies the CORS response headers corresponding to the
// origin for the service sommelier.
func handleSommelierOrigin(h http.Handler) http.Handler {
	spec0 := regexp.MustCompile(".*localhost*")
	origHndlr := h.(http.HandlerFunc)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			origHndlr(w, r)
			return
		}
		if cors.MatchOriginRegexp(origin, spec0) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Expose-Headers", "X-Time")
			w.Header().Set("Access-Control-Max-Age", "600")
			w.Header().Set("Access-Control-Allow-Credentials", "false")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
			}
			origHndlr(w, r)
			return
		}
		if cors.MatchOrigin(origin, "http://localhost") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Expose-Headers", "X-Time")
			w.Header().Set("Access-Control-Max-Age", "600")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
				w.Header().Set("Access-Control-Allow-Headers", "X-Shared-Secret")
			}
			origHndlr(w, r)
			return
		}
		origHndlr(w, r)
		return
	})
}
