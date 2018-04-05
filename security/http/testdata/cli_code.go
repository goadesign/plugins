package testdata

var BasicAuthBuildCode = `// BuildLoginPayload builds the payload for the BasicAuth login endpoint from
// CLI flags.
func BuildLoginPayload(basicAuthLoginUser string, basicAuthLoginPass string) (*basicauth.LoginPayload, error) {
	var user *string
	{
		if basicAuthLoginUser != "" {
			user = &basicAuthLoginUser
		}
	}
	var pass *string
	{
		if basicAuthLoginPass != "" {
			pass = &basicAuthLoginPass
		}
	}
	payload := &basicauth.LoginPayload{
		User: user,
		Pass: pass,
	}
	return payload, nil
}
`

var BasicAuthRequiredParseCode = `// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(
	scheme, host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restore bool,
) (goa.Endpoint, interface{}, error) {
	var (
		basicAuthRequiredFlags = flag.NewFlagSet("basic-auth-required", flag.ContinueOnError)

		basicAuthRequiredLoginFlags    = flag.NewFlagSet("login", flag.ExitOnError)
		basicAuthRequiredLoginUserFlag = basicAuthRequiredLoginFlags.String("user", "REQUIRED", "")
		basicAuthRequiredLoginPassFlag = basicAuthRequiredLoginFlags.String("pass", "REQUIRED", "")
		basicAuthRequiredLoginIDFlag   = basicAuthRequiredLoginFlags.String("id", "", "")
	)
	basicAuthRequiredFlags.Usage = basicAuthRequiredUsage
	basicAuthRequiredLoginFlags.Usage = basicAuthRequiredLoginUsage

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
		case "basic-auth-required":
			svcf = basicAuthRequiredFlags
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
		case "basic-auth-required":
			switch epn {
			case "login":
				epf = basicAuthRequiredLoginFlags

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
		case "basic-auth-required":
			c := basicauthrequiredc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "login":
				endpoint = c.Login()
				data, err = basicauthrequiredc.BuildLoginPayload(*basicAuthRequiredLoginUserFlag, *basicAuthRequiredLoginPassFlag, *basicAuthRequiredLoginIDFlag)
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}
`

var BasicAuthRequiredBuildCode = `// BuildLoginPayload builds the payload for the BasicAuthRequired login
// endpoint from CLI flags.
func BuildLoginPayload(basicAuthRequiredLoginUser string, basicAuthRequiredLoginPass string, basicAuthRequiredLoginID string) (*basicauthrequired.LoginPayload, error) {
	var user string
	{
		user = basicAuthRequiredLoginUser
	}
	var pass string
	{
		pass = basicAuthRequiredLoginPass
	}
	var id *string
	{
		if basicAuthRequiredLoginID != "" {
			id = &basicAuthRequiredLoginID
		}
	}
	payload := &basicauthrequired.LoginPayload{
		User: user,
		Pass: pass,
		ID:   id,
	}
	return payload, nil
}
`

var JWTBuildCode = `// BuildLoginPayload builds the payload for the JWT login endpoint from CLI
// flags.
func BuildLoginPayload(jwtLoginToken string) (*jwt.LoginPayload, error) {
	var token *string
	{
		if jwtLoginToken != "" {
			token = &jwtLoginToken
		}
	}
	payload := &jwt.LoginPayload{
		Token: token,
	}
	return payload, nil
}
`
