// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// archiver HTTP client CLI support package
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/fetcher/archiver/design

package cli

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/endpoint"
	goahttp "goa.design/goa/http"
	archiversvcc "goa.design/plugins/goakit/examples/fetcher/archiver/gen/http/archiver/client"
	healthc "goa.design/plugins/goakit/examples/fetcher/archiver/gen/http/health/client"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//    command (subcommand1|subcommand2|...)
//
func UsageCommands() string {
	return `archiver (archive|read)
health show
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` archiver archive --body '{
      "body": "Unde sed nulla.",
      "status": 200
   }'` + "\n" +
		os.Args[0] + ` health show` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(
	scheme, host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restore bool,
) (endpoint.Endpoint, interface{}, error) {
	var (
		archiverFlags = flag.NewFlagSet("archiver", flag.ContinueOnError)

		archiverArchiveFlags    = flag.NewFlagSet("archive", flag.ExitOnError)
		archiverArchiveBodyFlag = archiverArchiveFlags.String("body", "REQUIRED", "")

		archiverReadFlags  = flag.NewFlagSet("read", flag.ExitOnError)
		archiverReadIDFlag = archiverReadFlags.String("id", "REQUIRED", "ID of archive")

		healthFlags = flag.NewFlagSet("health", flag.ContinueOnError)

		healthShowFlags = flag.NewFlagSet("show", flag.ExitOnError)
	)
	archiverFlags.Usage = archiverUsage
	archiverArchiveFlags.Usage = archiverArchiveUsage
	archiverReadFlags.Usage = archiverReadUsage

	healthFlags.Usage = healthUsage
	healthShowFlags.Usage = healthShowUsage

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return nil, nil, err
	}

	if len(os.Args) < flag.NFlag()+3 {
		return nil, nil, fmt.Errorf("not enough arguments")
	}

	var (
		svcn string
		svcf *flag.FlagSet
	)
	{
		svcn = os.Args[1+flag.NFlag()]
		switch svcn {
		case "archiver":
			svcf = archiverFlags
		case "health":
			svcf = healthFlags
		default:
			return nil, nil, fmt.Errorf("unknown service %q", svcn)
		}
	}
	if err := svcf.Parse(os.Args[2+flag.NFlag():]); err != nil {
		return nil, nil, err
	}

	var (
		epn string
		epf *flag.FlagSet
	)
	{
		epn = os.Args[2+flag.NFlag()+svcf.NFlag()]
		switch svcn {
		case "archiver":
			switch epn {
			case "archive":
				epf = archiverArchiveFlags

			case "read":
				epf = archiverReadFlags

			}

		case "health":
			switch epn {
			case "show":
				epf = healthShowFlags

			}

		}
	}
	if epf == nil {
		return nil, nil, fmt.Errorf("unknown %q endpoint %q", svcn, epn)
	}

	// Parse endpoint flags if any
	if len(os.Args) > 2+flag.NFlag()+svcf.NFlag() {
		if err := epf.Parse(os.Args[3+flag.NFlag()+svcf.NFlag():]); err != nil {
			return nil, nil, err
		}
	}

	var (
		data     interface{}
		endpoint endpoint.Endpoint
		err      error
	)
	{
		switch svcn {
		case "archiver":
			c := archiversvcc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "archive":
				endpoint = c.Archive()
				data, err = archiversvcc.BuildArchivePayload(*archiverArchiveBodyFlag)
			case "read":
				endpoint = c.Read()
				data, err = archiversvcc.BuildReadPayload(*archiverReadIDFlag)
			}
		case "health":
			c := healthc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "show":
				endpoint = c.Show()
				data = nil
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// archiverUsage displays the usage of the archiver command and its subcommands.
func archiverUsage() {
	fmt.Fprintf(os.Stderr, `Service is the archiver service interface.
Usage:
    %s [globalflags] archiver COMMAND [flags]

COMMAND:
    archive: Archive HTTP response
    read: Read HTTP response from archive

Additional help:
    %s archiver COMMAND --help
`, os.Args[0], os.Args[0])
}
func archiverArchiveUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] archiver archive -body JSON

Archive HTTP response
    -body JSON: 

Example:
    `+os.Args[0]+` archiver archive --body '{
      "body": "Unde sed nulla.",
      "status": 200
   }'
`, os.Args[0])
}

func archiverReadUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] archiver read -id INT

Read HTTP response from archive
    -id INT: ID of archive

Example:
    `+os.Args[0]+` archiver read --id 7904528753725778865
`, os.Args[0])
}

// healthUsage displays the usage of the health command and its subcommands.
func healthUsage() {
	fmt.Fprintf(os.Stderr, `Service is the health service interface.
Usage:
    %s [globalflags] health COMMAND [flags]

COMMAND:
    show: Health check endpoint

Additional help:
    %s health COMMAND --help
`, os.Args[0], os.Args[0])
}
func healthShowUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] health show

Health check endpoint

Example:
    `+os.Args[0]+` health show
`, os.Args[0])
}
