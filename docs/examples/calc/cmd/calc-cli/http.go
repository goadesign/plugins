package main

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
	cli "goa.design/plugins/v3/docs/examples/calc/gen/http/cli/calc"
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

	var (
		dialer *websocket.Dialer
	)
	{
		dialer = websocket.DefaultDialer
	}

	return cli.ParseEndpoint(
		scheme,
		host,
		doer,
		goahttp.RequestEncoder,
		goahttp.ResponseDecoder,
		debug,
		dialer,
		nil,
	)
}

func httpUsageCommands() string {
	return cli.UsageCommands()
}

func httpUsageExamples() string {
	return cli.UsageExamples()
}
