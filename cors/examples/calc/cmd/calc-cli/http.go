package main

import (
	"fmt"
	"net/http"
	"time"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
	cli "goa.design/plugins/v3/cors/examples/calc/gen/http/cli/calc"
)

func doHTTP(scheme, host string, timeout int, debug bool) (goa.Endpoint, any, error) {
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
