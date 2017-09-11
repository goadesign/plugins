// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// secured HTTP client CLI support package
//
// Command:
// $ goa gen goa.design/plugins/security/example/design

package cli

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	goa "goa.design/goa"
	goahttp "goa.design/goa/http"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//    command (subcommand1|subcommand2|...)
//
func UsageCommands() string {
	return `securedService (signin|secure|doublySecure|alsoDoublySecure)
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` securedService signin --p "user:password"` + "\n" +
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
		securedServiceFlags = flag.NewFlagSet("securedService", flag.ContinueOnError)

		securedServiceSigninFlags = flag.NewFlagSet("signin", flag.ExitOnError)
		securedServiceSigninPFlag = securedServiceSigninFlags.String("p", "REQUIRED", "Credentials used to authenticate to retrieve JWT token")

		securedServiceSecureFlags    = flag.NewFlagSet("secure", flag.ExitOnError)
		securedServiceSecureBodyFlag = securedServiceSecureFlags.String("body", "REQUIRED", "")
		securedServiceSecureFailFlag = securedServiceSecureFlags.String("fail", "", "")

		securedServiceDoublySecureFlags    = flag.NewFlagSet("doublySecure", flag.ExitOnError)
		securedServiceDoublySecureBodyFlag = securedServiceDoublySecureFlags.String("body", "REQUIRED", "")
		securedServiceDoublySecureKeyFlag  = securedServiceDoublySecureFlags.String("key", "", "")

		securedServiceAlsoDoublySecureFlags    = flag.NewFlagSet("alsoDoublySecure", flag.ExitOnError)
		securedServiceAlsoDoublySecureBodyFlag = securedServiceAlsoDoublySecureFlags.String("body", "REQUIRED", "")
		securedServiceAlsoDoublySecureKeyFlag  = securedServiceAlsoDoublySecureFlags.String("key", "", "")
	)
	securedServiceFlags.Usage = securedServiceUsage
	securedServiceSigninFlags.Usage = securedServiceSigninUsage
	securedServiceSecureFlags.Usage = securedServiceSecureUsage
	securedServiceDoublySecureFlags.Usage = securedServiceDoublySecureUsage
	securedServiceAlsoDoublySecureFlags.Usage = securedServiceAlsoDoublySecureUsage

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
		case "securedService":
			svcf = securedServiceFlags
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
		case "securedService":
			switch epn {
			case "signin":
				epf = securedServiceSigninFlags

			case "secure":
				epf = securedServiceSecureFlags

			case "doublySecure":
				epf = securedServiceDoublySecureFlags

			case "alsoDoublySecure":
				epf = securedServiceAlsoDoublySecureFlags

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
		case "securedService":
			c := securedservicec.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "signin":
				endpoint = c.Signin()
				data = nil
			case "secure":
				endpoint = c.Secure()
				data, err = securedservicec.BuildSecurePayload(*securedServiceSecureBodyFlag, *securedServiceSecureFailFlag)
			case "doublySecure":
				endpoint = c.DoublySecure()
				data, err = securedservicec.BuildDoublySecurePayload(*securedServiceDoublySecureBodyFlag, *securedServiceDoublySecureKeyFlag)
			case "alsoDoublySecure":
				endpoint = c.AlsoDoublySecure()
				data, err = securedservicec.BuildAlsoDoublySecurePayload(*securedServiceAlsoDoublySecureBodyFlag, *securedServiceAlsoDoublySecureKeyFlag)
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// securedServiceUsage displays the usage of the securedService command and its
// subcommands.
func securedServiceUsage() {
	fmt.Fprintf(os.Stderr, `The secured service exposes endpoints that require valid authorization credentials.
Usage:
    %s [globalflags] securedService COMMAND [flags]

COMMAND:
    signin: Creates a valid JWT
    secure: This action is secured with the jwt scheme
    doublySecure: This action is secured with the jwt scheme and also requires an API key query string.
    alsoDoublySecure: This action is secured with the jwt scheme and also requires an API key header.

Additional help:
    %s securedService COMMAND --help
`, os.Args[0], os.Args[0])
}
func securedServiceSigninUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] securedService signin -p STRING

Creates a valid JWT
    -p STRING: Credentials used to authenticate to retrieve JWT token

Example:
    `+os.Args[0]+` securedService signin --p "user:password"
`, os.Args[0])
}

func securedServiceSecureUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] securedService secure -body JSON -fail BOOL

This action is secured with the jwt scheme
    -body JSON: 
    -fail BOOL: 

Example:
    `+os.Args[0]+` securedService secure --body '{
      "token": "Voluptatem architecto consequatur fuga nisi veritatis."
   }' --fail false
`, os.Args[0])
}

func securedServiceDoublySecureUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] securedService doublySecure -body JSON -key STRING

This action is secured with the jwt scheme and also requires an API key query string.
    -body JSON: 
    -key STRING: 

Example:
    `+os.Args[0]+` securedService doubly-secure --body '{
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"
   }' --key "abcdef12345"
`, os.Args[0])
}

func securedServiceAlsoDoublySecureUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] securedService alsoDoublySecure -body JSON -key STRING

This action is secured with the jwt scheme and also requires an API key header.
    -body JSON: 
    -key STRING: 

Example:
    `+os.Args[0]+` securedService also-doubly-secure --body '{
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"
   }' --key "abcdef12345"
`, os.Args[0])
}
