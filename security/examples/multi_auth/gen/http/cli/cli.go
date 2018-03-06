// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// multi_auth HTTP client CLI support package
//
// Command:
// $ goa gen goa.design/plugins/security/examples/multi_auth/design

package cli

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	goa "goa.design/goa"
	goahttp "goa.design/goa/http"
	securedservicec "goa.design/plugins/security/examples/multi_auth/gen/http/secured_service/client"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//    command (subcommand1|subcommand2|...)
//
func UsageCommands() string {
	return `secured_service (signin|secure|doubly_secure|also_doubly_secure)
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` secured_service signin --body '{
      "password": "password",
      "username": "user"
   }'` + "\n" +
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
		securedServiceFlags = flag.NewFlagSet("secured_service", flag.ContinueOnError)

		securedServiceSigninFlags    = flag.NewFlagSet("signin", flag.ExitOnError)
		securedServiceSigninBodyFlag = securedServiceSigninFlags.String("body", "REQUIRED", "")

		securedServiceSecureFlags    = flag.NewFlagSet("secure", flag.ExitOnError)
		securedServiceSecureBodyFlag = securedServiceSecureFlags.String("body", "REQUIRED", "")
		securedServiceSecureFailFlag = securedServiceSecureFlags.String("fail", "", "")

		securedServiceDoublySecureFlags    = flag.NewFlagSet("doubly_secure", flag.ExitOnError)
		securedServiceDoublySecureBodyFlag = securedServiceDoublySecureFlags.String("body", "REQUIRED", "")
		securedServiceDoublySecureKeyFlag  = securedServiceDoublySecureFlags.String("key", "", "")

		securedServiceAlsoDoublySecureFlags    = flag.NewFlagSet("also_doubly_secure", flag.ExitOnError)
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
		case "secured_service":
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
		case "secured_service":
			switch epn {
			case "signin":
				epf = securedServiceSigninFlags

			case "secure":
				epf = securedServiceSecureFlags

			case "doubly_secure":
				epf = securedServiceDoublySecureFlags

			case "also_doubly_secure":
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
		case "secured_service":
			c := securedservicec.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "signin":
				endpoint = c.Signin()
				data, err = securedservicec.BuildSigninPayload(*securedServiceSigninBodyFlag)
			case "secure":
				endpoint = c.Secure()
				data, err = securedservicec.BuildSecurePayload(*securedServiceSecureBodyFlag, *securedServiceSecureFailFlag)
			case "doubly_secure":
				endpoint = c.DoublySecure()
				data, err = securedservicec.BuildDoublySecurePayload(*securedServiceDoublySecureBodyFlag, *securedServiceDoublySecureKeyFlag)
			case "also_doubly_secure":
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

// secured_serviceUsage displays the usage of the secured_service command and
// its subcommands.
func securedServiceUsage() {
	fmt.Fprintf(os.Stderr, `The secured service exposes endpoints that require valid authorization credentials.
Usage:
    %s [globalflags] secured_service COMMAND [flags]

COMMAND:
    signin: Creates a valid JWT
    secure: This action is secured with the jwt scheme
    doubly_secure: This action is secured with the jwt scheme and also requires an API key query string.
    also_doubly_secure: This action is secured with the jwt scheme and also requires an API key header.

Additional help:
    %s secured_service COMMAND --help
`, os.Args[0], os.Args[0])
}
func securedServiceSigninUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] secured_service signin -body JSON

Creates a valid JWT
    -body JSON: 

Example:
    `+os.Args[0]+` secured_service signin --body '{
      "password": "password",
      "username": "user"
   }'
`, os.Args[0])
}

func securedServiceSecureUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] secured_service secure -body JSON -fail BOOL

This action is secured with the jwt scheme
    -body JSON: 
    -fail BOOL: 

Example:
    `+os.Args[0]+` secured_service secure --body '{
      "token": "Dignissimos reiciendis itaque enim quibusdam."
   }' --fail true
`, os.Args[0])
}

func securedServiceDoublySecureUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] secured_service doubly_secure -body JSON -key STRING

This action is secured with the jwt scheme and also requires an API key query string.
    -body JSON: 
    -key STRING: 

Example:
    `+os.Args[0]+` secured_service doubly-secure --body '{
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"
   }' --key "abcdef12345"
`, os.Args[0])
}

func securedServiceAlsoDoublySecureUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] secured_service also_doubly_secure -body JSON -key STRING

This action is secured with the jwt scheme and also requires an API key header.
    -body JSON: 
    -key STRING: 

Example:
    `+os.Args[0]+` secured_service also-doubly-secure --body '{
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"
   }' --key "abcdef12345"
`, os.Args[0])
}
