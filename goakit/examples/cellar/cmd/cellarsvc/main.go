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
	cellar "goa.design/plugins/goakit/examples/cellar"
	sommeliersvr "goa.design/plugins/goakit/examples/cellar/gen/http/sommelier/kitserver"
	storagesvr "goa.design/plugins/goakit/examples/cellar/gen/http/storage/kitserver"
	swaggersvr "goa.design/plugins/goakit/examples/cellar/gen/http/swagger/kitserver"
	"goa.design/plugins/goakit/examples/cellar/gen/sommelier"
	"goa.design/plugins/goakit/examples/cellar/gen/storage"
)

func main() {
	// Define command line flags, add any other flag required to configure
	// the service.
	var (
		addr = flag.String("listen", ":8080", "HTTP listen `address`")
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
		sommeliers sommelier.Service
		storages   storage.Service
	)
	{
		sommeliers = cellar.NewSommelier(logger)
		storages = cellar.NewStorage(logger)
	}

	// Wrap the services in endpoints that can be invoked from other
	// services potentially running in different processes.
	var (
		sommeliere *sommelier.Endpoints
		storagee   *storage.Endpoints
	)
	{
		sommeliere = sommelier.NewEndpoints(sommeliers)
		storagee = storage.NewEndpoints(storages)
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
		sommelierPickHandler *kithttp.Server
		storageListHandler   *kithttp.Server
		storageShowHandler   *kithttp.Server
		storageAddHandler    *kithttp.Server
		storageRemoveHandler *kithttp.Server
	)
	{
		sommelierPickHandler = kithttp.NewServer(
			endpoint.Endpoint(sommeliere.Pick),
			sommeliersvr.DecodePickRequest(mux, dec),
			sommeliersvr.EncodePickResponse(enc),
		)
		storageListHandler = kithttp.NewServer(
			endpoint.Endpoint(storagee.List),
			func(context.Context, *http.Request) (request interface{}, err error) { return nil, nil },
			storagesvr.EncodeListResponse(enc),
		)
		storageShowHandler = kithttp.NewServer(
			endpoint.Endpoint(storagee.Show),
			storagesvr.DecodeShowRequest(mux, dec),
			storagesvr.EncodeShowResponse(enc),
		)
		storageAddHandler = kithttp.NewServer(
			endpoint.Endpoint(storagee.Add),
			storagesvr.DecodeAddRequest(mux, dec),
			storagesvr.EncodeAddResponse(enc),
		)
		storageRemoveHandler = kithttp.NewServer(
			endpoint.Endpoint(storagee.Remove),
			storagesvr.DecodeRemoveRequest(mux, dec),
			storagesvr.EncodeRemoveResponse(enc),
		)
	}

	// Configure the mux.
	sommeliersvr.MountPickHandler(mux, sommelierPickHandler)
	storagesvr.MountListHandler(mux, storageListHandler)
	storagesvr.MountShowHandler(mux, storageShowHandler)
	storagesvr.MountAddHandler(mux, storageAddHandler)
	storagesvr.MountRemoveHandler(mux, storageRemoveHandler)
	swaggersvr.MountGenHTTPOpenapiJSON(mux)

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
