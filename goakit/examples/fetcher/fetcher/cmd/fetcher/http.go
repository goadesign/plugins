package main

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"goa.design/clue/debug"
	"goa.design/clue/log"
	goahttp "goa.design/goa/v3/http"
	fetcher "goa.design/plugins/v3/goakit/examples/fetcher/fetcher/gen/fetcher"
	health "goa.design/plugins/v3/goakit/examples/fetcher/fetcher/gen/health"
	fetcherkitsvr "goa.design/plugins/v3/goakit/examples/fetcher/fetcher/gen/http/fetcher/kitserver"
	fetchersvr "goa.design/plugins/v3/goakit/examples/fetcher/fetcher/gen/http/fetcher/server"
	healthkitsvr "goa.design/plugins/v3/goakit/examples/fetcher/fetcher/gen/http/health/kitserver"
	healthsvr "goa.design/plugins/v3/goakit/examples/fetcher/fetcher/gen/http/health/server"
)

// handleHTTPServer starts configures and starts a HTTP server on the given
// URL. It shuts down the server if any error is received in the error channel.
func handleHTTPServer(ctx context.Context, u *url.URL, fetcherEndpoints *fetcher.Endpoints, healthEndpoints *health.Endpoints, wg *sync.WaitGroup, errc chan error, dbg bool) {

	// Provide the transport specific request decoder and response encoder.
	// The goa http package has built-in support for JSON, XML and gob.
	// Other encodings can be used by providing the corresponding functions,
	// see goa.design/implement/encoding.
	var (
		dec = goahttp.RequestDecoder
		enc = goahttp.ResponseEncoder
	)

	// Build the service HTTP request multiplexer and mount debug and profiler
	// endpoints in debug mode.
	var mux goahttp.Muxer
	{
		mux = goahttp.NewMuxer()
		if dbg {
			// Mount pprof handlers for memory profiling under /debug/pprof.
			debug.MountPprofHandlers(debug.Adapt(mux))
			// Mount /debug endpoint to enable or disable debug logs at runtime.
			debug.MountDebugLogEnabler(debug.Adapt(mux))
		}
	}

	// Wrap the endpoints with the transport specific layers. The generated
	// server packages contains code generated from the design which maps
	// the service input and output data structures to HTTP requests and
	// responses.
	var (
		fetcherFetchHandler *kithttp.Server
		fetcherServer       *fetchersvr.Server
		healthShowHandler   *kithttp.Server
		healthServer        *healthsvr.Server
	)
	{
		eh := errorHandler(ctx)
		fetcherFetchHandler = kithttp.NewServer(
			endpoint.Endpoint(fetcherEndpoints.Fetch),
			fetcherkitsvr.DecodeFetchRequest(mux, dec),
			fetcherkitsvr.EncodeFetchResponse(enc),
			kithttp.ServerErrorEncoder(fetcherkitsvr.EncodeFetchError(enc, nil)),
		)
		fetcherServer = fetchersvr.New(fetcherEndpoints, mux, dec, enc, eh, nil)
		healthShowHandler = kithttp.NewServer(
			endpoint.Endpoint(healthEndpoints.Show),
			func(context.Context, *http.Request) (request interface{}, err error) { return nil, nil },
			healthkitsvr.EncodeShowResponse(enc),
		)
		healthServer = healthsvr.New(healthEndpoints, mux, dec, enc, eh, nil)
	}

	// Configure the mux.
	fetcherkitsvr.MountFetchHandler(mux, fetcherFetchHandler)
	healthkitsvr.MountShowHandler(mux, healthShowHandler)

	var handler http.Handler = mux
	if dbg {
		// Log query and response bodies if debug logs are enabled.
		handler = debug.HTTP()(handler)
	}
	handler = log.HTTP(ctx)(handler)

	// Start HTTP server using default configuration, change the code to
	// configure the server as required by your service.
	srv := &http.Server{Addr: u.Host, Handler: handler, ReadHeaderTimeout: time.Second * 60}
	for _, m := range fetcherServer.Mounts {
		log.Printf(ctx, "HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range healthServer.Mounts {
		log.Printf(ctx, "HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}

	(*wg).Add(1)
	go func() {
		defer (*wg).Done()

		// Start HTTP server in a separate goroutine.
		go func() {
			log.Printf(ctx, "HTTP server listening on %q", u.Host)
			errc <- srv.ListenAndServe()
		}()

		<-ctx.Done()
		log.Printf(ctx, "shutting down HTTP server at %q", u.Host)

		// Shutdown gracefully with a 30s timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			log.Printf(ctx, "failed to shutdown: %v", err)
		}
	}()
}

// errorHandler returns a function that writes and logs the given error.
// The function also writes and logs the error unique ID so that it's possible
// to correlate.
func errorHandler(logCtx context.Context) func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		log.Printf(logCtx, "ERROR: %s", err.Error())
	}
}
