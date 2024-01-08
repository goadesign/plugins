// Code generated by goa v3.14.4, DO NOT EDIT.
//
// health HTTP server
//
// Command:
// $ goa gen goa.design/plugins/v3/goakit/examples/fetcher/archiver/design -o
// $(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/archiver

package server

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
	health "goa.design/plugins/v3/goakit/examples/fetcher/archiver/gen/health"
)

// Server lists the health service endpoint HTTP handlers.
type Server struct {
	Mounts []*MountPoint
	Show   http.Handler
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

// New instantiates HTTP handlers for all the health service endpoints using
// the provided encoder and decoder. The handlers are mounted on the given mux
// using the HTTP verb and path defined in the design. errhandler is called
// whenever a response fails to be encoded. formatter is used to format errors
// returned by the service methods prior to encoding. Both errhandler and
// formatter are optional and can be nil.
func New(
	e *health.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"Show", "GET", "/health"},
		},
		Show: NewShowHandler(e.Show, mux, decoder, encoder, errhandler, formatter),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "health" }

// Use wraps the server handlers with the given middleware.
func (s *Server) Use(m func(http.Handler) http.Handler) {
	s.Show = m(s.Show)
}

// MethodNames returns the methods served.
func (s *Server) MethodNames() []string { return health.MethodNames[:] }

// Mount configures the mux to serve the health endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountShowHandler(mux, h.Show)
}

// Mount configures the mux to serve the health endpoints.
func (s *Server) Mount(mux goahttp.Muxer) {
	Mount(mux, s)
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
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) http.Handler {
	var (
		encodeResponse = EncodeShowResponse(encoder)
		encodeError    = goahttp.ErrorEncoder(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "show")
		ctx = context.WithValue(ctx, goa.ServiceKey, "health")
		var err error
		res, err := endpoint(ctx, nil)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}
