// Code generated by goa v3.6.2, DO NOT EDIT.
//
// calc HTTP server
//
// Command:
// $ goa gen goa.design/plugins/v3/cors/examples/calc/design -o
// $(GOPATH)/src/goa.design/plugins/cors/examples/calc

package server

import (
	"context"
	"net/http"
	"regexp"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
	"goa.design/plugins/v3/cors"
	calc "goa.design/plugins/v3/cors/examples/calc/gen/calc"
)

// Server lists the calc service endpoint HTTP handlers.
type Server struct {
	Mounts    []*MountPoint
	Add       http.Handler
	CORS      http.Handler
	IndexHTML http.Handler
}

// ErrorNamer is an interface implemented by generated error structs that
// exposes the name of the error as defined in the design.
type ErrorNamer interface {
	ErrorName() string
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

// New instantiates HTTP handlers for all the calc service endpoints using the
// provided encoder and decoder. The handlers are mounted on the given mux
// using the HTTP verb and path defined in the design. errhandler is called
// whenever a response fails to be encoded. formatter is used to format errors
// returned by the service methods prior to encoding. Both errhandler and
// formatter are optional and can be nil.
func New(
	e *calc.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
	fileSystemIndexHTML http.FileSystem,
) *Server {
	if fileSystemIndexHTML == nil {
		fileSystemIndexHTML = http.Dir(".")
	}
	return &Server{
		Mounts: []*MountPoint{
			{"Add", "GET", "/add/{a}/{b}"},
			{"CORS", "OPTIONS", "/add/{a}/{b}"},
			{"CORS", "OPTIONS", "/"},
			{"/index.html", "GET", "/"},
		},
		Add:       NewAddHandler(e.Add, mux, decoder, encoder, errhandler, formatter),
		CORS:      NewCORSHandler(),
		IndexHTML: http.FileServer(fileSystemIndexHTML),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "calc" }

// Use wraps the server handlers with the given middleware.
func (s *Server) Use(m func(http.Handler) http.Handler) {
	s.Add = m(s.Add)
	s.CORS = m(s.CORS)
}

// Mount configures the mux to serve the calc endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountAddHandler(mux, h.Add)
	MountCORSHandler(mux, h.CORS)
	MountIndexHTML(mux, h.IndexHTML)
}

// Mount configures the mux to serve the calc endpoints.
func (s *Server) Mount(mux goahttp.Muxer) {
	Mount(mux, s)
}

// MountAddHandler configures the mux to serve the "calc" service "add"
// endpoint.
func MountAddHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := HandleCalcOrigin(h).(http.HandlerFunc)
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
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeAddRequest(mux, decoder)
		encodeResponse = EncodeAddResponse(encoder)
		encodeError    = goahttp.ErrorEncoder(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "add")
		ctx = context.WithValue(ctx, goa.ServiceKey, "calc")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
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

// MountIndexHTML configures the mux to serve GET request made to "/".
func MountIndexHTML(mux goahttp.Muxer, h http.Handler) {
	mux.Handle("GET", "/", HandleCalcOrigin(h).ServeHTTP)
}

// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service calc.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = HandleCalcOrigin(h)
	mux.Handle("OPTIONS", "/add/{a}/{b}", h.ServeHTTP)
	mux.Handle("OPTIONS", "/", h.ServeHTTP)
}

// NewCORSHandler creates a HTTP handler which returns a simple 200 response.
func NewCORSHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
}

// HandleCalcOrigin applies the CORS response headers corresponding to the
// origin for the service calc.
func HandleCalcOrigin(h http.Handler) http.Handler {
	spec0 := regexp.MustCompile(".*localhost.*")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			h.ServeHTTP(w, r)
			return
		}
		if cors.MatchOriginRegexp(origin, spec0) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Expose-Headers", "X-Time, X-Api-Version")
			w.Header().Set("Access-Control-Max-Age", "100")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
			}
			h.ServeHTTP(w, r)
			return
		}
		if cors.MatchOrigin(origin, "http://127.0.0.1") {
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
			h.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r)
		return
	})
}
