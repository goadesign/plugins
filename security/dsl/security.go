package dsl

import (
	"goa.design/goa/dsl"
	"goa.design/goa/eval"
	goadesign "goa.design/goa/http/design"
	"goa.design/plugins/security/design"
)

// Security defines authentication requirements to access an API, a service or a
// service endpoint.
//
// The requirement refers to one or more OAuth2Security, BasicAuthSecurity,
// APIKeySecurity or JWTSecurity security scheme. If the schemes include a
// OAuth2Security or JWTSecurity scheme then required scopes may be listed by
// name in the Security DSL. All the listed schemes must be validated by the
// client for the request to be authorized. Security may appear multiple times
// in the same scope in which case the client may validate any one of the
// requirements for the request to be authorized.
//
// Security must appear in a API, Service or Method expression.
//
// Security accepts an arbitrary number of security schemes as argument
// specified by name or by reference and an optional DSL function as last
// argument.
//
// Examples:
//
//    var _ = API("calc", func() {
//        // All API endpoints are secured via basic auth by default.
//        Security(BasicAuth)
//    })
//
//    var _ = Service("calculator", func() {
//        // Override default API security requirements. Accept either basic
//        // auth or OAuth2 access token with "api:read" scope.
//        Security(BasicAuth)
//        Security("oauth2", func() {
//            Scope("api:read")
//        })
//
//        Method("add", func() {
//            Description("Add two operands")
//
//            // Override default service security requirements. Require
//            // both basic auth and OAuth2 access token with "api:write"
//            // scope.
//            Security(BasicAuth, "oauth2", func() {
//                Scope("api:write")
//            })
//
//            Payload(Operands)
//            Error(ErrBadRequest, ErrorResult)
//        })
//    })
//
func Security(args ...interface{}) {
	var dsl func()
	{
		if d, ok := args[len(args)-1].(func()); ok {
			args = args[:len(args)-1]
			dsl = d
		}
	}

	var schemes []*design.SchemeExpr
	{
		schemes = make([]*design.SchemeExpr, len(args))
		for i, arg := range args {
			switch val := arg.(type) {
			case string:
				for _, s := range design.Root.Schemes {
					if s.SchemeName == val {
						schemes[i] = s
						break
					}
				}
				if schemes[i] == nil {
					eval.ReportError("security scheme %q not found", val)
					return
				}
			case *design.SchemeExpr:
				schemes[i] = val
			default:
				eval.InvalidArgError("security scheme or security scheme name", val)
				return
			}
		}
	}

	security := &design.SecurityExpr{Schemes: schemes}
	if dsl != nil {
		if !eval.Execute(dsl, security) {
			return
		}
	}

	current := eval.Current()
	switch actual := current.(type) {
	case *goadesign.EndpointExpr:
		sec := &design.EndpointSecurityExpr{SecurityExpr: security, Endpoint: actual}
		design.Root.EndpointSecurity = append(design.Root.EndpointSecurity, sec)
	case *goadesign.FileServerExpr:
		sec := &design.FileServerSecurityExpr{SecurityExpr: security, FileServer: actual}
		design.Root.FileServerSecurity = append(design.Root.FileServerSecurity, sec)
	case *goadesign.ServiceExpr:
		sec := &design.ServiceSecurityExpr{SecurityExpr: security, Service: actual}
		design.Root.ServiceSecurity = append(design.Root.ServiceSecurity, sec)
	case *goadesign.RootExpr:
		design.Root.APISecurity = append(design.Root.APISecurity, security)
	default:
		eval.IncompatibleDSL()
		return
	}
}

// NoSecurity resets the security schemes for a Service or a Method.
//
// NoSecurity can be used in Service or Method.
func NoSecurity() {
	expr := &design.SecurityExpr{
		Schemes: []*design.SchemeExpr{&design.SchemeExpr{Kind: design.NoSecurityKind}},
	}

	current := eval.Current()
	switch actual := current.(type) {
	case *design.EndpointExpr:
		actual.Security = def
	case *design.FileServerExpr:
		actual.Security = def
	case *design.ServiceExpr:
		actual.Security = def
	default:
		eval.IncompatibleDSL()
		return
	}
}

// BasicAuthSecurity defines a basic authentication security scheme.
//
// BasicAuthSecurity is a top level DSL.
//
// BasicAuthSecurity takes a name as first argument and an optional DSL as
// second argument.
//
// Example:
//
//     var Basic = BasicAuthSecurity("password", func() {
//         Description("Use your own password!")
//     })
//
func BasicAuthSecurity(name string, dsl ...func()) *design.SchemeExpr {
	if _, ok := eval.Current().(eval.TopExpr); !ok {
		eval.IncompatibleDSL()
		return nil
	}

	if securitySchemeRedefined(name) {
		return nil
	}

	expr := &design.SchemeExpr{
		Kind:       design.BasicAuthSecurityKind,
		SchemeName: name,
	}

	if len(dsl) != 0 {
		if !eval.Execute(dsl[0], expr) {
			return nil
		}
	}

	design.Root.Schemes = append(design.Root.Schemes, expr)

	return expr
}

// Username defines the attribute used to provide the username to a endpoint
// secured with basic authentication. The parameters and usage of Username are
// the same as the goa DSL Attribute function.
//
// The generated code produced by goa uses the value of the corresponding
// payload field to compute the basic authentication Authorization header value.
//
// Example:
//
//    Method("login", func() {
//        Security(Basic)
//        Payload(func() {
//            Username("user")
//            Password("pass")
//        })
//        HTTP(func() {
//            POST("/login")
//        })
//    })
//
func Username(name string, args ...interface{}) {
	args = useDSL(args, func() { dsl.Metadata("security:username") })
	dsl.Attribute(name, args...)
}

// Password defines the attribute used to provide the password to a endpoint
// secured with basic authentication. The parameters and usage of Password are
// the same as the goa DSL Attribute function.
//
// The generated code produced by goa uses the value of the corresponding
// payload field to compute the basic authentication Authorization header value.
//
// Example:
//
//    Method("login", func() {
//        Security(Basic)
//        Payload(func() {
//            Username("user")
//            Password("pass")
//        })
//        HTTP(func() {
//            POST("/login")
//        })
//    })
//
func Password(name string, args ...interface{}) {
	args = useDSL(args, func() { dsl.Metadata("security:password") })
	dsl.Attribute(name, args...)
}

// APIKeySecurity defines a API key security scheme available throughout the API.
//
// APIKeySecurity is a top level DSL.
//
// APIKeySecurity takes a name as first argument and an optional DSL as
// second argument.
//
// Example:
//
//    var APIKey = APIKeySecurity("key", func() {
//          Description("Shared secret")
//    })
//
func APIKeySecurity(name string, dsl ...func()) *design.SchemeExpr {
	if _, ok := eval.Current().(eval.TopExpr); !ok {
		eval.IncompatibleDSL()
		return nil
	}

	if securitySchemeRedefined(name) {
		return nil
	}

	expr := &design.SchemeExpr{
		Kind:       design.APIKeySecurityKind,
		SchemeName: name,
	}

	if len(dsl) != 0 {
		if !eval.Execute(dsl[0], expr) {
			return nil
		}
	}

	design.Root.Schemes = append(design.Root.Schemes, expr)

	return expr
}

// APIKey defines the attribute used to provide the API key to a endpoint
// secured with API keys. The parameters and usage of APIKey are the same as
// the goa DSL Attribute function.
//
// The generated code produced by goa uses the value of the corresponding
// payload field to set the API key value.
//
// Example:
//
//    Method("secured_read", func() {
//        Security(APIKeyAuth)
//        Payload(func() {
//            APIKey("key", String, "API key used to perform authorization")
//            Required("key")
//        })
//        Result(String)
//        HTTP(func() {
//            GET("/")
//            Param("key:api_key") // Provide the key as a query string param
//        })
//    })
//
//    Method("secured_write", func() {
//        Security(APIKeyAuth)
//        Payload(func() {
//            APIKey("key", String, "API key used to perform authorization")
//            Attribute("data", String, "Data to be written")
//            Required("key", "data")
//        })
//        HTTP(func() {
//            POST("/")
//            Header("key:Authorization") // Provide the key as a header (default)
//        })
//    })
//
func APIKey(name string, args ...interface{}) {
	args = useDSL(args, func() { dsl.Metadata("security:apikey") })
	dsl.Attribute(name, args...)
}

// OAuth2Security defines an OAuth2 security scheme. The DSL provided as second
// argument must define one and exactly one flow. One of AccessCodeFlow,
// ImplicitFlow, PasswordFlow or ApplicationFlow. Each flow defines endpoints
// for retrieving OAuth2 authorization codes and/or refresh and access tokens.
// The endpoint URLs may be complete or may be just a path in which case the API
// scheme and host are used to build the full URL. See for example [Aaron
// Parecki's writeup](https://aaronparecki.com/2012/07/29/2/oauth2-simplified)
// for additional details on OAuth2 flows.
//
// The OAuth2 DSL also defines the scopes that may be associated with the
// incoming request tokens.
//
// OAuth2Security is a top level DSL.
//
// OAuth2Security takes a name as first argument and a DSL as second argument.
//
// Example:
//
//    var OAuth2 = OAuth2Security("googAuth", func() {
//        ImplicitFlow("/authorization")
//     // PasswordFlow("/token"...)
//     // ClientCredentialsFlow("/token")
//     // AuthorizationCodeFlow("/authorization", "/token")
//
//        Scope("my_system:write", "Write to the system")
//        Scope("my_system:read", "Read anything in there")
//    })
//
func OAuth2Security(name string, dsl ...func()) *design.SchemeExpr {
	switch eval.Current().(type) {
	case *design.APIExpr, *eval.TopLevelExpr:
	default:
		eval.IncompatibleDSL()
		return nil
	}

	if securitySchemeRedefined(name) {
		return nil
	}

	def := &design.SchemeExpr{
		SchemeName: name,
		Kind:       design.OAuth2SecurityKind,
		Type:       "oauth2",
	}

	if len(dsl) != 0 {
		def.DSLFunc = dsl[0]
	}

	design.Design.Schemes = append(design.Design.Schemes, def)

	return def
}

// AccessToken defines the attribute used to provide the access token to a
// endpoint secured with OAuth2. The parameters and usage of AccessToken are the
// same as the goa DSL Attribute function.
func AccessToken(name string, args ...interface{}) {
	args = useDSL(args, func() { dsl.Metadata("security:accesstoken") })
	dsl.Attribute(name, args...)
}

// JWTSecurity defines an HTTP security scheme with support for Scopes and a
// TokenURL.
//
// Since Scopes and TokenURLs are not compatible with the Swagger specification,
// the swagger generator inserts comments in the description of the different
// elements on which they are defined.
//
// JWTSecurity is a top level DSL.
//
// JWTSecurity takes a name as first argument and a DSL as second argument.
//
// Example:
//
//    var JWT = JWTSecurity("jwt", func() {
//        TokenURL("https://example.com/token")
//        Scope("my_system:write", "Write to the system")
//        Scope("my_system:read", "Read anything in there")
//    })
//
func JWTSecurity(name string, dsl ...func()) *design.SchemeExpr {
	switch eval.Current().(type) {
	case *design.APIExpr, *eval.TopLevelExpr:
	default:
		eval.IncompatibleDSL()
		return nil
	}

	if securitySchemeRedefined(name) {
		return nil
	}

	def := &design.SchemeExpr{
		SchemeName: name,
		Kind:       design.JWTSecurityKind,
		Type:       "apiKey",
	}

	if len(dsl) != 0 {
		def.DSLFunc = dsl[0]
	}

	design.Design.Schemes = append(design.Design.Schemes, def)

	return def
}

// Token defines the attribute used to provide the JWT to a endpoint secured via
// JWT. The parameters and usage of Token are the same as the goa DSL Attribute
// function.
func Token(name string, args ...interface{}) {
	args = useDSL(args, func() { dsl.Metadata("security:token") })
	dsl.Attribute(name, args...)
}

// Scope can be used in: Security, JWTSecurity, OAuth2Security
//
// Scope defines an authorization scope. Used within Scheme, a description may be provided
// explaining what the scope means. Within a Security block, only a scope is needed.
func Scope(name string, desc ...string) {
	switch current := eval.Current().(type) {
	case *design.SecurityExpr:
		if len(desc) >= 1 {
			eval.ReportError("too many arguments")
			return
		}
		current.Scopes = append(current.Scopes, name)
	case *design.SchemeExpr:
		if len(desc) > 1 {
			eval.ReportError("too many arguments")
			return
		}
		if current.FlowExpr == nil {
			current.Scopes = make(map[string]string)
		}
		d := "no description"
		if len(desc) == 1 {
			d = desc[0]
		}
		current.Scopes[name] = d
	default:
		eval.IncompatibleDSL()
	}
}

// inHeader is called by `Header()`, see documentation there.
func inHeader(headerName string) {
	if current, ok := eval.Current().(*design.SchemeExpr); ok {
		if current.Kind == design.APIKeySecurityKind || current.Kind == design.JWTSecurityKind {
			if current.In != "" {
				eval.ReportError("'In' previously defined through Header or Query")
				return
			}
			current.In = "header"
			current.Name = headerName
			return
		}
	}
	eval.IncompatibleDSL()
}

// Query can be used in: APIKeySecurity, JWTSecurity
//
// Query defines that an APIKeySecurity or JWTSecurity implementation must check in the query
// parameter named "parameterName" to get the api key.
func Query(parameterName string) {
	if current, ok := eval.Current().(*design.SchemeExpr); ok {
		if current.Kind == design.APIKeySecurityKind || current.Kind == design.JWTSecurityKind {
			if current.In != "" {
				eval.ReportError("'In' previously defined through Header or Query")
				return
			}
			current.In = "query"
			current.Name = parameterName
			return
		}
	}
	eval.IncompatibleDSL()
}

// AccessCodeFlow can be used in: OAuth2Security
//
// AccessCodeFlow defines an "access code" OAuth2 flow.  Use within an OAuth2Security expression.
func AccessCodeFlow(authorizationURL, tokenURL string) {
	if current, ok := eval.Current().(*design.SchemeExpr); ok {
		if current.Kind == design.OAuth2SecurityKind {
			current.Flow = "accessCode"
			current.AuthorizationURL = authorizationURL
			current.TokenURL = tokenURL
			return
		}
	}
	eval.IncompatibleDSL()
}

// ApplicationFlow can be used in: OAuth2Security
//
// ApplicationFlow defines an "application" OAuth2 flow.  Use within an OAuth2Security expression.
func ApplicationFlow(tokenURL string) {
	if parent, ok := eval.Current().(*design.SchemeExpr); ok {
		if parent.Kind == design.OAuth2SecurityKind {
			parent.Flow = "application"
			parent.TokenURL = tokenURL
			return
		}
	}
	eval.IncompatibleDSL()
}

// PasswordFlow can be used in: OAuth2Security
//
// PasswordFlow defines a "password" OAuth2 flow.  Use within an OAuth2Security expression.
func PasswordFlow(tokenURL string) {
	if parent, ok := eval.Current().(*design.SchemeExpr); ok {
		if parent.Kind == design.OAuth2SecurityKind {
			parent.Flow = "password"
			parent.TokenURL = tokenURL
			return
		}
	}
	eval.IncompatibleDSL()
}

// ImplicitFlow can be used in: OAuth2Security
//
// ImplicitFlow defines an "implicit" OAuth2 flow. Use within an OAuth2Security
// expression.
func ImplicitFlow(authorizationURL string) {
	if parent, ok := eval.Current().(*design.SchemeExpr); ok {
		if parent.Kind == design.OAuth2SecurityKind {
			parent.Flow = "implicit"
			parent.AuthorizationURL = authorizationURL
			return
		}
	}
	eval.IncompatibleDSL()
}

// TokenURL can be used in: JWTSecurity
//
// TokenURL defines a URL to get an access token.  If you are defining OAuth2 flows, use
// `ImplicitFlow`, `PasswordFlow`, `AccessCodeFlow` or `ApplicationFlow` instead. This will set an
// endpoint where you can obtain a JWT with the JWTSecurity scheme. The URL may be a complete URL
// or just a path in which case the API scheme and host are used to build the full URL.
func TokenURL(tokenURL string) {
	if parent, ok := eval.Current().(*design.SchemeExpr); ok {
		if parent.Kind == design.JWTSecurityKind {
			parent.TokenURL = tokenURL
			return
		}
	}
	eval.IncompatibleDSL()
}

func securitySchemeRedefined(name string) bool {
	for _, s := range design.Root.Schemes {
		if s.SchemeName == name {
			eval.ReportError("cannot redefine security scheme with name %q", name)
			return true
		}
	}
	return false
}

// useDSL modifies the Attribute function to use the given function as DSL,
// merging it with any pre-exsiting DSL.
func useDSL(d func(), args []interface{}) []interface{} {
	dsl, ok = args[len(args)-1].(func())
	if ok {
		newdsl = func() { dsl(); d() }
		args = append(args[:len(args)-1], newdsl)
	} else {
		args = append(args, d)
	}
	return args
}
