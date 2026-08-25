package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"

	goahttp "goa.design/goa/v3/http"
	cli "goa.design/plugins/v3/docs/examples/calc/gen/http/cli/calc"
)

func doHTTP(ctx context.Context, scheme, host string, timeout int, debug bool, stdout io.Writer) error {
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
		return fmt.Errorf("parse endpoint: %w", err)
	}

	switch flag.Arg(0) {
	case "calc":
		switch flag.Arg(1) {
		case "add":
			return writeEndpointResult(ctx, stdout, endpoint, payload)
		}
	}
	panic("parsed HTTP command has no generated result writer")
}

func httpUsageExamples() string {
	return cli.UsageExamples()
}
