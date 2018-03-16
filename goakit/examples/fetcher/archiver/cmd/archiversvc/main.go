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
	archiver "goa.design/plugins/goakit/examples/fetcher/archiver"
	archiversvc "goa.design/plugins/goakit/examples/fetcher/archiver/gen/archiver"
	health "goa.design/plugins/goakit/examples/fetcher/archiver/gen/health"
	archiversvckitsvr "goa.design/plugins/goakit/examples/fetcher/archiver/gen/http/archiver/kitserver"
	archiversvcsvr "goa.design/plugins/goakit/examples/fetcher/archiver/gen/http/archiver/server"
	healthkitsvr "goa.design/plugins/goakit/examples/fetcher/archiver/gen/http/health/kitserver"
	healthsvr "goa.design/plugins/goakit/examples/fetcher/archiver/gen/http/health/server"
)

func main() {
	// Define command line flags, add any other flag required to configure
	// the service.
	var (
		addr = flag.String("listen", ":8081", "HTTP listen `address`")
	)
	flag.Parse()

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
		archiversvcs archiversvc.Service
		healths      health.Service
	)
	{
		archiversvcs = archiver.NewArchiver(logger)
		healths = archiver.NewHealth(logger)
	}

	// Wrap the services in endpoints that can be invoked from other
	// services potentially running in different processes.
	var (
		archiversvce *archiversvc.Endpoints
		healthe      *health.Endpoints
	)
	{
		archiversvce = archiversvc.NewEndpoints(archiversvcs)
		healthe = health.NewEndpoints(healths)
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
		archiversvcArchiveHandler *kithttp.Server
		archiversvcReadHandler    *kithttp.Server
		archiversvcServer         *archiversvcsvr.Server
		healthShowHandler         *kithttp.Server
		healthServer              *healthsvr.Server
	)
	{
		archiversvcArchiveHandler = kithttp.NewServer(
			endpoint.Endpoint(archiversvce.Archive),
			archiversvckitsvr.DecodeArchiveRequest(mux, dec),
			archiversvckitsvr.EncodeArchiveResponse(enc),
		)
		archiversvcReadHandler = kithttp.NewServer(
			endpoint.Endpoint(archiversvce.Read),
			archiversvckitsvr.DecodeReadRequest(mux, dec),
			archiversvckitsvr.EncodeReadResponse(enc),
		)
		archiversvcServer = archiversvcsvr.New(archiversvce, mux, dec, enc)
		healthShowHandler = kithttp.NewServer(
			endpoint.Endpoint(healthe.Show),
			func(context.Context, *http.Request) (request interface{}, err error) { return nil, nil },
			healthkitsvr.EncodeShowResponse(enc),
		)
		healthServer = healthsvr.New(healthe, mux, dec, enc)
	}

	// Configure the mux.
	archiversvckitsvr.MountArchiveHandler(mux, archiversvcArchiveHandler)
	archiversvckitsvr.MountReadHandler(mux, archiversvcReadHandler)
	healthkitsvr.MountShowHandler(mux, healthShowHandler)

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
		for _, m := range archiversvcServer.Mounts {
			logger.Log("info", fmt.Sprintf("service %s method %s mounted on %s %s", archiversvcServer.Service(), m.Method, m.Verb, m.Pattern))
		}
		for _, m := range healthServer.Mounts {
			logger.Log("info", fmt.Sprintf("service %s method %s mounted on %s %s", healthServer.Service(), m.Method, m.Verb, m.Pattern))
		}
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
