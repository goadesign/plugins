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
	fetcher "goa.design/plugins/goakit/examples/fetcher/fetcher"
	fetchersvc "goa.design/plugins/goakit/examples/fetcher/fetcher/gen/fetcher"
	health "goa.design/plugins/goakit/examples/fetcher/fetcher/gen/health"
	fetchersvcsvr "goa.design/plugins/goakit/examples/fetcher/fetcher/gen/http/fetcher/kitserver"
	healthsvr "goa.design/plugins/goakit/examples/fetcher/fetcher/gen/http/health/kitserver"
)

func main() {
	// Define command line flags, add any other flag required to configure
	// the service.
	var (
		addr         = flag.String("listen", ":8080", "HTTP listen `address`")
		archiverHost = flag.String("archiver", ":8081", "archiver service `host:port`")
	)
	flag.Parse()
	if *archiverHost == "" {
		fmt.Fprintf(os.Stderr, "missing required flag --archiver")
		os.Exit(1)
	}

	// Setup logger.
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
		healths     health.Service
		fetchersvcs fetchersvc.Service
	)
	{
		healths = fetcher.NewHealth(logger)
		fetchersvcs = fetcher.NewFetcher(logger, *archiverHost)
	}

	// Wrap the services in endpoints that can be invoked from other
	// services potentially running in different processes.
	var (
		healthe     *health.Endpoints
		fetchersvce *fetchersvc.Endpoints
	)
	{
		healthe = health.NewEndpoints(healths)
		fetchersvce = fetchersvc.NewEndpoints(fetchersvcs)
	}

	// Provide the transport specific request decoder and response encoder.
	// The goa http package has built-in support for JSON, XML and gob.
	// Other encodings can be used by providing the corresponding functions,
	// see goa.design/encoding.
	var (
		dec = goahttp.RequestDecoder
		enc = goahttp.ResponseEncoder
	)

	// Build the service HTTP request router (a.k.a. mux).
	var mux goahttp.Muxer
	{
		mux = goahttp.NewMuxer()
	}

	// Wrap the endpoints with the transport specific layer.
	var (
		healthShowHandler      *kithttp.Server
		fetchersvcFetchHandler *kithttp.Server
	)
	{
		healthShowHandler = kithttp.NewServer(
			endpoint.Endpoint(healthe.Show),
			func(context.Context, *http.Request) (request interface{}, err error) { return nil, nil },
			healthsvr.EncodeShowResponse(enc),
		)
		fetchersvcFetchHandler = kithttp.NewServer(
			endpoint.Endpoint(fetchersvce.Fetch),
			fetchersvcsvr.DecodeFetchRequest(mux, dec),
			fetchersvcsvr.EncodeFetchResponse(enc),
		)
	}

	// Configure the mux.
	healthsvr.MountShowHandler(mux, healthShowHandler)
	fetchersvcsvr.MountFetchHandler(mux, fetchersvcFetchHandler)

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
	srv := &http.Server{Addr: *addr, Handler: mux}
	go func() {
		logger.Log("listening", *addr)
		errc <- srv.ListenAndServe()
	}()

	// Wait for signal.
	logger.Log("exiting", <-errc)

	// Shutdown gracefully with a 30s timeout.
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	srv.Shutdown(ctx)

	logger.Log("server", "exited")
}
