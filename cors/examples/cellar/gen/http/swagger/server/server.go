// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// swagger HTTP server
//
// Command:
// $ goa gen goa.design/plugins/cors/examples/cellar/design

package server

import (
	"context"
	"net/http"

	goahttp "goa.design/goa/http"
	"goa.design/plugins/cors"
	swagger "goa.design/plugins/cors/examples/cellar/gen/swagger"
)

// Server lists the swagger service endpoint HTTP handlers.
type Server struct {
	Mounts []*MountPoint
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

// New instantiates HTTP handlers for all the swagger service endpoints.
func New(
	e *swagger.Endpoints,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"CORS", "OPTIONS", "/swagger/swagger.json"},
			{"../../gen/http/openapi.json", "GET", "/swagger/swagger.json"},
		},
		CORS: NewCORSHandler(),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "swagger" }

// Mount configures the mux to serve the swagger endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountCORSHandler(mux, h.CORS)
	MountGenHTTPOpenapiJSON(mux, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../gen/http/openapi.json")
	}))
}

// MountGenHTTPOpenapiJSON configures the mux to serve GET request made to
// "/swagger/swagger.json".
func MountGenHTTPOpenapiJSON(mux goahttp.Muxer, h http.Handler) {
	mux.Handle("GET", "/swagger/swagger.json", handleSwaggerOrigin(h).ServeHTTP)
}

// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service swagger.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = handleSwaggerOrigin(h)
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("OPTIONS", "/swagger/swagger.json", f)
}

// NewCORSHandler creates a HTTP handler which returns a simple 200 response.
func NewCORSHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
}

// handleSwaggerOrigin applies the CORS response headers corresponding to the
// origin for the service swagger.
func handleSwaggerOrigin(h http.Handler) http.Handler {
	origHndlr := h.(http.HandlerFunc)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
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
