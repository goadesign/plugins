package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	goahttp "goa.design/goa/v3/http"
	cli "goa.design/plugins/v3/goakit/examples/fetcher/archiver/gen/http/cli/archiver"
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

	endpoint, payload, err := cli.ParseEndpoint(
		scheme,
		host,
		doer,
		goahttp.RequestEncoder,
		goahttp.ResponseDecoder,
		debug,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("parse endpoint: %w", err)
	}
	return endpoint, payload, nil
}

func httpUsageCommands() []string {
	return cli.UsageCommands()
}

func httpUsageExamples() string {
	return cli.UsageExamples()
}
