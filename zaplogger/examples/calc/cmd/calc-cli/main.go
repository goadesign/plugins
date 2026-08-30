package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	goa "goa.design/goa/v3/pkg"
)

func main() {
	var (
		hostF    = flag.String("host", "development", "Server host (valid values: development, production)")
		addrF    = flag.String("url", "", "URL to service host")
		versionF = flag.String("version", "v1", "API version")

		verboseF = flag.Bool("verbose", false, "Print request and response details")
		vF       = flag.Bool("v", false, "Print request and response details")
		timeoutF = flag.Int("timeout", 30, "Maximum number of seconds to wait for response")
	)
	flag.Usage = usage
	flag.Parse()

	var (
		addr    string
		timeout int
		debug   bool
	)
	{
		addr = *addrF
		if addr == "" {
			switch *hostF {
			case "development":
				addr = "http://localhost:8000/calc"
			case "production":
				addr = "https://{version}.goa.design/calc"
				addr = strings.ReplaceAll(addr, "{version}", *versionF)
			default:
				fmt.Fprintf(os.Stderr, "invalid host argument: %q (valid hosts: development|production)\n", *hostF)
				os.Exit(1)
			}
		}
		timeout = *timeoutF
		debug = *verboseF || *vF
	}

	var (
		scheme string
		host   string
	)
	{
		u, err := url.Parse(addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid URL %#v: %s\n", addr, err)
			os.Exit(1)
		}
		scheme = u.Scheme
		host = u.Host
	}

	var (
		err error
	)
	{
		switch scheme {
		case "http", "https":
			err = doHTTP(context.Background(), scheme, host, timeout, debug, os.Stdout)
		default:
			fmt.Fprintf(os.Stderr, "invalid scheme: %q (valid schemes: http|https)\n", scheme)
			os.Exit(1)
		}
	}
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		fmt.Fprintln(os.Stderr, err.Error())
		fmt.Fprintln(os.Stderr, "run '"+os.Args[0]+" --help' for detailed usage.")
		os.Exit(1)
	}

}

// writeEndpointResult calls one normal endpoint and writes its result as JSON.
func writeEndpointResult(ctx context.Context, stdout io.Writer, endpoint goa.Endpoint, payload any) error {
	data, err := endpoint(ctx, payload)
	if err != nil {
		return err
	}
	return writeJSON(stdout, data)
}

// writeJSON writes one indented JSON value followed by a newline.
func writeJSON(stdout io.Writer, data any) error {
	if data == nil {
		return nil
	}
	encoded, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("encode result: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, string(encoded)); err != nil {
		return fmt.Errorf("write result: %w", err)
	}
	return nil
}

func usage() {
	usageCommands := []string{
		"calc add",
	}
	fmt.Fprintf(os.Stderr, `%s is a command line client for the calc API.

Usage:
    %s [-host HOST][-url URL][-timeout SECONDS][-verbose|-v][-version VERSION] SERVICE ENDPOINT [flags]

    -host HOST:  server host (development). valid values: development, production
    -url URL:    specify service URL overriding host URL (http://localhost:8000/calc)
    -timeout:    maximum number of seconds to wait for response (30)
    -verbose|-v: print request and response details (false)
    -version:    API version (v1)

Commands:
%s
Additional help:
    %s SERVICE [ENDPOINT] --help

Example:
%s
`, os.Args[0], os.Args[0], indent(strings.Join(usageCommands, "\n")), os.Args[0], indent(httpUsageExamples()))
}

func indent(s string) string {
	if s == "" {
		return ""
	}
	return "    " + strings.ReplaceAll(s, "\n", "\n    ")
}
