package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/log"
	"goa.design/clue/debug"
	"goa.design/clue/log"
	goa "goa.design/goa/v3/pkg"
	fetcherapi "goa.design/plugins/v3/goakit/examples/fetcher/fetcher"
	fetcher "goa.design/plugins/v3/goakit/examples/fetcher/fetcher/gen/fetcher"
	health "goa.design/plugins/v3/goakit/examples/fetcher/fetcher/gen/health"
)

func main() {
	// Define command line flags, add any other flag required to configure the
	// service.
	var (
		hostF        = flag.String("host", "localhost", "Server host (valid values: localhost)")
		domainF      = flag.String("domain", "", "Host domain name (overrides host domain specified in service design)")
		httpPortF    = flag.String("http-port", "", "HTTP port (overrides host HTTP port specified in service design)")
		secureF      = flag.Bool("secure", false, "Use secure scheme (https or grpcs)")
		dbgF         = flag.Bool("debug", false, "Log request and response bodies")
		archiverHost = flag.String("archiver", "localhost:8081", "archiver service `host:port`")
	)
	flag.Parse()

	// Setup logger. Replace logger with your own log package of choice.
	format := log.FormatJSON
	if log.IsTerminal() {
		format = log.FormatTerminal
	}
	ctx := log.Context(context.Background(), log.WithFormat(format))
	if *dbgF {
		ctx = log.Context(ctx, log.WithDebug())
		log.Debugf(ctx, "debug logs enabled")
	}
	log.Print(ctx, log.KV{K: "http-port", V: *httpPortF})

	// Initialize the services.
	var (
		fetcherSvc fetcher.Service
		healthSvc  health.Service
	)
	{
		{
			var logger kitlog.Logger
			logger = kitlog.NewLogfmtLogger(os.Stderr)
			logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
			logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
			logger = kitlog.With(logger, "service", "fetcher")
			fetcherSvc = fetcherapi.NewFetcher(logger, *archiverHost)
		}
		{
			var logger kitlog.Logger
			logger = kitlog.NewLogfmtLogger(os.Stderr)
			logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
			logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
			logger = kitlog.With(logger, "service", "health")
			healthSvc = fetcherapi.NewHealth(logger)
		}
	}

	// Wrap the services in endpoints that can be invoked from other services
	// potentially running in different processes.
	var (
		fetcherEndpoints *fetcher.Endpoints
		healthEndpoints  *health.Endpoints
	)
	{
		fetcherEndpoints = fetcher.NewEndpoints(fetcherSvc)
		fetcherEndpoints.Use(wrapMiddleware(debug.LogPayloads()))
		fetcherEndpoints.Use(wrapMiddleware(log.Endpoint))
		healthEndpoints = health.NewEndpoints(healthSvc)
		healthEndpoints.Use(wrapMiddleware(debug.LogPayloads()))
		healthEndpoints.Use(wrapMiddleware(log.Endpoint))
	}

	// Create channel used by both the signal handler and server goroutines
	// to notify the main goroutine when to stop the server.
	errc := make(chan error)

	// Setup interrupt handler. This optional step configures the process so
	// that SIGINT and SIGTERM signals cause the services to stop gracefully.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)

	// Start the servers and send errors (if any) to the error channel.
	switch *hostF {
	case "localhost":
		{
			addr := "http://localhost:80"
			u, err := url.Parse(addr)
			if err != nil {
				log.Fatalf(ctx, err, "invalid URL %#v\n", addr)
			}
			if *secureF {
				u.Scheme = "https"
			}
			if *domainF != "" {
				u.Host = *domainF
			}
			if *httpPortF != "" {
				h, _, err := net.SplitHostPort(u.Host)
				if err != nil {
					log.Fatalf(ctx, err, "invalid URL %#v\n", u.Host)
				}
				u.Host = net.JoinHostPort(h, *httpPortF)
			} else if u.Port() == "" {
				u.Host = net.JoinHostPort(u.Host, "80")
			}
			handleHTTPServer(ctx, u, fetcherEndpoints, healthEndpoints, &wg, errc, *dbgF)
		}

	default:
		log.Fatal(ctx, fmt.Errorf("invalid host argument: %q (valid hosts: localhost)", *hostF))
	}

	// Wait for signal.
	log.Printf(ctx, "exiting (%v)", <-errc)

	// Send cancellation signal to the goroutines.
	cancel()

	wg.Wait()
	log.Printf(ctx, "exited")
}

// Wrap goa middleware into go-kit middleware.
func wrapMiddleware(mw func(goa.Endpoint) goa.Endpoint) func(endpoint.Endpoint) endpoint.Endpoint {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return endpoint.Endpoint(mw(goa.Endpoint(e)))
	}
}
