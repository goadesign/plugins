package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	goahttp "goa.design/goa/http"
	"goa.design/goa/http/middleware"
	fetcher "goa.design/plugins/goakit/examples/fetcher/fetcher"
	fetchersvc "goa.design/plugins/goakit/examples/fetcher/fetcher/gen/fetcher"
	health "goa.design/plugins/goakit/examples/fetcher/fetcher/gen/health"
	fetchersvckitsvr "goa.design/plugins/goakit/examples/fetcher/fetcher/gen/http/fetcher/kitserver"
	fetchersvcsvr "goa.design/plugins/goakit/examples/fetcher/fetcher/gen/http/fetcher/server"
	healthkitsvr "goa.design/plugins/goakit/examples/fetcher/fetcher/gen/http/health/kitserver"
	healthsvr "goa.design/plugins/goakit/examples/fetcher/fetcher/gen/http/health/server"
)

func main() {
	// Define command line flags, add any other flag required to configure
	// the service.
	var (
		addr         = flag.String("listen", "localhost:80", "HTTP listen `address`")
		archiverHost = flag.String("archiver", "localhost:8081", "archiver service `host:port`")
	)
	flag.Parse()
	if *archiverHost == "" {
		fmt.Fprintf(os.Stderr, "missing required flag --archiver")
		os.Exit(1)
	}

	// Setup gokit logger.
	var (
		logger log.Logger
	)
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Create the structs that implement the services.
	var (
		fetcherSvc fetchersvc.Service
		healthSvc  health.Service
	)
	{
		fetcherSvc = fetcher.NewFetcher(logger, *archiverHost)
		healthSvc = fetcher.NewHealth(logger)
	}

	// Wrap the services in endpoints that can be invoked from other
	// services potentially running in different processes.
	var (
		fetcherEndpoints *fetchersvc.Endpoints
		healthEndpoints  *health.Endpoints
	)
	{
		fetcherEndpoints = fetchersvc.NewEndpoints(fetcherSvc)
		healthEndpoints = health.NewEndpoints(healthSvc)
	}

	// Provide the transport specific request decoder and response encoder.
	// The goa http package has built-in support for JSON, XML and gob.
	// Other encodings can be used by providing the corresponding functions,
	// see goa.design/encoding.
	var (
		dec = goahttp.RequestDecoder
		enc = goahttp.ResponseEncoder
	)

	// Build the service HTTP request multiplexer and configure it to serve
	// HTTP requests to the service endpoints.
	var mux goahttp.Muxer
	{
		mux = goahttp.NewMuxer()
	}

	// Wrap the endpoints with the transport specific layers. The generated
	// server packages contains code generated from the design which maps
	// the service input and output data structures to HTTP requests and
	// responses.
	var (
		fetcherFetchHandler *kithttp.Server
		fetcherServer       *fetchersvcsvr.Server
		healthShowHandler   *kithttp.Server
		healthServer        *healthsvr.Server
	)
	{
		eh := ErrorHandler(logger)
		fetcherFetchHandler = kithttp.NewServer(
			endpoint.Endpoint(fetcherEndpoints.Fetch),
			fetchersvckitsvr.DecodeFetchRequest(mux, dec),
			fetchersvckitsvr.EncodeFetchResponse(enc),
		)
		fetcherServer = fetchersvcsvr.New(fetcherEndpoints, mux, dec, enc, eh)
		healthShowHandler = kithttp.NewServer(
			endpoint.Endpoint(healthEndpoints.Show),
			func(context.Context, *http.Request) (request interface{}, err error) { return nil, nil },
			healthkitsvr.EncodeShowResponse(enc),
		)
		healthServer = healthsvr.New(healthEndpoints, mux, dec, enc, eh)
	}

	// Configure the mux.
	fetchersvckitsvr.MountFetchHandler(mux, fetcherFetchHandler)
	healthkitsvr.MountShowHandler(mux, healthShowHandler)

	// Wrap the multiplexer with additional middlewares. Middlewares mounted
	// here apply to all the service endpoints.
	var handler http.Handler = mux
	{
		handler = middleware.Log(logger)(handler)
		handler = middleware.RequestID()(handler)
	}

	// Create channel used by both the signal handler and server goroutines
	// to notify the main goroutine when to stop the server.
	errc := make(chan error)
	// Setup interrupt handler. This optional step configures the process so
	// that SIGINT and SIGTERM signals cause the service to stop gracefully.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		errc <- fmt.Errorf("%s", <-c)
	}()
	// Start HTTP server using default configuration, change the code to
	// configure the server as required by your service.
	srv := &http.Server{Addr: *addr, Handler: handler}
	go func() {
		for _, m := range fetcherServer.Mounts {
			logger.Log("info", fmt.Sprintf("method %q mounted on %s %s", m.Method, m.Verb, m.Pattern))
		}
		for _, m := range healthServer.Mounts {
			logger.Log("info", fmt.Sprintf("method %q mounted on %s %s", m.Method, m.Verb, m.Pattern))
		}
		logger.Log("info", fmt.Sprintf("listening on %s", *addr))
		errc <- srv.ListenAndServe()
	}()

	// Wait for signal.
	logger.Log("info", fmt.Sprintf("exiting (%v)", <-errc))
	// Shutdown gracefully with a 30s timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	logger.Log("info", fmt.Sprintf("exited"))
}

// ErrorHandler returns a function that writes and logs the given error.
// The function also writes and logs the error unique ID so that it's possible
// to correlate.
func ErrorHandler(logger log.Logger) func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		id := ctx.Value(middleware.RequestIDKey).(string)
		w.Write([]byte("[" + id + "] encoding: " + err.Error()))
		logger.Log("info", fmt.Sprintf("[%s] ERROR: %s", id, err.Error()))
	}
}
