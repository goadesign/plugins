package main

import (
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	goahttp "goa.design/goa/v3/http"
	cli "goa.design/plugins/v3/goakit/examples/calc/gen/http/cli/calc"
)

func doHTTP(scheme, host string, timeout int, debug bool) (endpoint.Endpoint, any, error) {
	var (
		doer goahttp.Doer
	)
	{
		doer = &http.Client{Timeout: time.Duration(timeout) * time.Second}
		if debug {
			doer = goahttp.NewDebugDoer(doer)
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
