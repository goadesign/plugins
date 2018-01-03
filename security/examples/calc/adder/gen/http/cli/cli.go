// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// adder HTTP client CLI support package
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/adder/design

package cli

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	goa "goa.design/goa"
	goahttp "goa.design/goa/http"
	addersvcc "goa.design/plugins/security/examples/calc/adder/gen/http/adder/client"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//    command (subcommand1|subcommand2|...)
//
func UsageCommands() string {
	return `adder add
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` adder add --a 2 --b 3 --key "abcdef12345"` + "\n" +
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
) (goa.Endpoint, interface{}, error) {
	var (
		adderFlags = flag.NewFlagSet("adder", flag.ContinueOnError)

		adderAddFlags   = flag.NewFlagSet("add", flag.ExitOnError)
		adderAddAFlag   = adderAddFlags.String("a", "REQUIRED", "Left operand")
		adderAddBFlag   = adderAddFlags.String("b", "REQUIRED", "Right operand")
		adderAddKeyFlag = adderAddFlags.String("key", "REQUIRED", "")
	)
	adderFlags.Usage = adderUsage
	adderAddFlags.Usage = adderAddUsage

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
		case "adder":
			svcf = adderFlags
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
		case "adder":
			switch epn {
			case "add":
				epf = adderAddFlags

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
		endpoint goa.Endpoint
		err      error
	)
	{
		switch svcn {
		case "adder":
			c := addersvcc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "add":
				endpoint = c.Add()
				data, err = addersvcc.BuildAddAddPayload(*adderAddAFlag, *adderAddBFlag, *adderAddKeyFlag)
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// adderUsage displays the usage of the adder command and its subcommands.
func adderUsage() {
	fmt.Fprintf(os.Stderr, `The adder service exposes an add method secured via API keys.
Usage:
    %s [globalflags] adder COMMAND [flags]

COMMAND:
    add: This action returns the sum of two integers and is secured with the API key scheme

Additional help:
    %s adder COMMAND --help
`, os.Args[0], os.Args[0])
}
func adderAddUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] adder add -a INT -b INT -key STRING

This action returns the sum of two integers and is secured with the API key scheme
    -a INT: Left operand
    -b INT: Right operand
    -key STRING: 

Example:
    `+os.Args[0]+` adder add --a 2 --b 3 --key "abcdef12345"
`, os.Args[0])
}
