package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/endpoint"
	goahttp "goa.design/goa/v3/http"
	cli "goa.design/plugins/goakit/examples/fetcher/fetcher/gen/http/cli/fetcher"
)

func doHTTP(scheme, host string, timeout int, debug bool) (endpoint.Endpoint, interface{}, error) {
	var (
		doer goahttp.Doer
	)
	{
		doer = &http.Client{Timeout: time.Duration(timeout) * time.Second}
		if debug {
			doer = goahttp.NewDebugDoer(doer)
			doer.(goahttp.DebugDoer).Fprint(os.Stderr)
		}
	}

	return cli.ParseEndpoint(
		scheme,
		host,
		doer,
		goahttp.RequestEncoder,
		goahttp.ResponseDecoder,
		debug,
	)
}
func httpUsageCommands() string {
	return cli.UsageCommands()
}

func httpUsageExamples() string {
	return cli.UsageExamples()
}
