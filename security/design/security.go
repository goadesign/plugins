package design

import (
	"fmt"
	"net/url"

	"goa.design/goa/design"
	"goa.design/goa/eval"
	httpdesign "goa.design/goa/http/design"
)

// SchemeKind is a type of security scheme.
type SchemeKind int

const (
	// OAuth2Kind identifies a "OAuth2" security scheme.
	OAuth2Kind SchemeKind = iota + 1
	// BasicAuthKind means "basic" security scheme.
	BasicAuthKind
	// APIKeyKind means "apiKey" security scheme.
	APIKeyKind
	// JWTKind means an "apiKey" security scheme, with support for
	// TokenPath and Scopes.
	JWTKind
	// NoKind means to have no security for this endpoint.
	NoKind
)

// FlowKind is a type of OAuth2 flow.
type FlowKind int

const (
	// AuthorizationCodeFlowKind identifies a OAuth2 authorization code
	// flow.
	AuthorizationCodeFlowKind FlowKind = iota + 1
	// ImplicitFlowKind identifiers a OAuth2 implicit flow.
	ImplicitFlowKind
	// PasswordFlowKind identifies a Resource Owner Password flow.
	PasswordFlowKind
	// ClientCredentialsFlowKind identifies a OAuth Client Credentials flow.
	ClientCredentialsFlowKind
)

type (
	// SecurityExpr defines a security requirement.
	SecurityExpr struct {
		// Schemes is the list of security schemes used for this
		// requirement.
		Schemes []*SchemeExpr
		// Scopes list the required scopes if any.
		Scopes []string
	}

	// ServiceSecurityExpr defines a security requirement that applies to
	// a service.
	ServiceSecurityExpr struct {
		*SecurityExpr
		// Service is the service that the security requirements applies
		// to.
		Service *design.ServiceExpr
	}

	// EndpointSecurityExpr defines a security requirement that applies to
	// an endpoint.
	EndpointSecurityExpr struct {
		*SecurityExpr
		// Method is the method that the security requirements
		// applies to.
		Method *design.MethodExpr
	}

	// SchemeExpr defines a security scheme used to authenticate against the
	// method being designed.
	SchemeExpr struct {
		// Kind is the sort of security scheme this object represents.
		Kind SchemeKind
		// SchemeName is the name of the security scheme, e.g. "googAuth",
		// "my_big_token", "jwt".
		SchemeName string
		// Description describes the security scheme e.g. "Google OAuth2"
		Description string
		// In determines the location of the API key, one of "header" or
		// "query".
		In string
		// Name refers to a header or parameter name, based on In's
		// value.
		Name string
		// Scopes lists the JWT or OAuth2 scopes.
		Scopes []*ScopeExpr
		// Flows determine the oauth2 flows supported by this scheme.
		Flows []*FlowExpr
		// Metadata is a list of key/value pairs
		Metadata design.MetadataExpr
	}

	// FlowExpr describes a specific OAuth2 flow.
	FlowExpr struct {
		// Kind is the kind of flow.
		Kind FlowKind
		// AuthorizationURL to be used for implicit or authorizationCode
		// flows.
		AuthorizationURL string
		// TokenURL to be used for password, clientCredentials or
		// authorizationCode flows.
		TokenURL string
		// RefreshURL to be used for obtaining refresh token.
		RefreshURL string
	}

	// A ScopeExpr defines a scope name and description.
	ScopeExpr struct {
		// Name of the scope.
		Name string
		// Description is the description of the scope.
		Description string
	}
)

// EvalName returns the generic definition name used in error messages.
func (s *SecurityExpr) EvalName() string {
	var suffix string
	if len(s.Schemes) > 0 && len(s.Schemes[0].SchemeName) > 0 {
		suffix = "scheme " + s.Schemes[0].SchemeName
	}
	return "Security" + suffix
}

// EvalName returns the generic definition name used in error messages.
func (s *EndpointSecurityExpr) EvalName() string {
	return "security for endpoint " + s.Method.Name
}

// Validate validates the security schemes.
func (s *EndpointSecurityExpr) Validate() error {
	verr := new(eval.ValidationErrors)
	for _, sch := range s.Schemes {
		switch sch.Kind {
		case BasicAuthKind:
			if !hasTaggedField(s.Method.Payload, "security:username") {
				verr.Add(s, "payload of method %q of service %q does not define a username attribute, use Username to define one.", s.Method.Name, s.Method.Service.Name)
			}
			if !hasTaggedField(s.Method.Payload, "security:password") {
				verr.Add(s, "payload of method %q of service %q does not define a password attribute, use Password to define one.", s.Method.Name, s.Method.Service.Name)
			}
		case APIKeyKind:
			if !hasTaggedField(s.Method.Payload, "security:apikey:"+sch.SchemeName) {
				verr.Add(s, "payload of method %q of service %q does not define an API key attribute, use APIKey to define one.", s.Method.Name, s.Method.Service.Name)
			}
		case JWTKind:
			if !hasTaggedField(s.Method.Payload, "security:token") {
				verr.Add(s, "payload of method %q of service %q does not define a JWT attribute, use Token to define one.", s.Method.Name, s.Method.Service.Name)
			}
		case OAuth2Kind:
			if !hasTaggedField(s.Method.Payload, "security:accesstoken") {
				verr.Add(s, "payload of method %q of service %q does not define a OAuth2 access token attribute, use AccessToken to define one.", s.Method.Name, s.Method.Service.Name)
			}
		default:
			panic(fmt.Sprintf("unknown kind %#v", sch.Kind)) // bug
		}
		if err := sch.Validate(); err != nil {
			verr.Merge(err)
		}
	}
	return verr
}

// Finalize ensures the scheme expressions is set with the location and
// parameter/header name.
func (s *EndpointSecurityExpr) Finalize() {
	if svc := httpdesign.Root.Service(s.Method.Service.Name); svc != nil {
		ep := svc.Endpoint(s.Method.Name)
		for _, sch := range s.Schemes {
			sch.Finalize()
			var field string
			switch sch.Kind {
			case BasicAuthKind:
				userField := securityField(s.Method.Payload, "security:username")
				passField := securityField(s.Method.Payload, "security:password")
				ep.Body.Delete(userField)
				ep.Body.Delete(passField)
				if isEmpty(ep.Body) {
					ep.Body = &design.AttributeExpr{Type: design.Empty}
				}
				continue
			case APIKeyKind:
				field = securityField(s.Method.Payload, "security:apikey:"+sch.SchemeName)
			case JWTKind:
				field = securityField(s.Method.Payload, "security:token")
			case OAuth2Kind:
				field = securityField(s.Method.Payload, "security:accesstoken")
			}
			sch.Name, sch.In = findKey(ep, field)
			if sch.Name == "" {
				sch.Name = "Authorization"
				addHeaderAttr(ep, field, sch.Name)
				// Recompute body to remove the security fields now that they are added
				// to the header
				if ep.Body.Find(field) != nil {
					ep.Body = nil
					ep.Body = httpdesign.RequestBody(ep)
				}
			}
		}
	}
}

// DupScheme creates a copy of the given scheme expression.
func DupScheme(sch *SchemeExpr) *SchemeExpr {
	dup := SchemeExpr{
		Kind:        sch.Kind,
		SchemeName:  sch.SchemeName,
		Description: sch.Description,
		In:          sch.In,
		Scopes:      sch.Scopes,
		Flows:       sch.Flows,
		Metadata:    sch.Metadata,
	}
	return &dup
}

// Type returns the type of the scheme.
func (s *SchemeExpr) Type() string {
	switch s.Kind {
	case OAuth2Kind:
		return "OAuth2"
	case BasicAuthKind:
		return "BasicAuth"
	case APIKeyKind:
		return "APIKey"
	case JWTKind:
		return "JWT"
	default:
		panic(fmt.Sprintf("unknown scheme kind: %#v", s.Kind)) // bug
	}
}

// EvalName returns the generic definition name used in error messages.
func (s *SchemeExpr) EvalName() string {
	return s.Type() + "Security"
}

// Validate ensures that the method payload contains attributes required
// by the scheme.
func (s *SchemeExpr) Validate() *eval.ValidationErrors {
	verr := new(eval.ValidationErrors)
	for _, f := range s.Flows {
		if err := f.Validate(); err != nil {
			verr.Merge(err)
		}
	}
	return verr
}

// Finalize makes the TokenURL and AuthorizationURL complete if needed.
func (s *SchemeExpr) Finalize() {
	for _, f := range s.Flows {
		f.Finalize()
	}
}

// EvalName returns the name of the expression used in error messages.
func (s *FlowExpr) EvalName() string {
	if s.TokenURL != "" {
		return fmt.Sprintf("flow with token URL %q", s.TokenURL)
	}
	return fmt.Sprintf("flow with refresh URL %q", s.RefreshURL)
}

// Validate ensures that TokenURL and AuthorizationURL are valid URLs.
func (s *FlowExpr) Validate() *eval.ValidationErrors {
	verr := new(eval.ValidationErrors)
	_, err := url.Parse(s.TokenURL)
	if err != nil {
		verr.Add(s, "invalid token URL %q: %s", s.TokenURL, err)
	}
	_, err = url.Parse(s.AuthorizationURL)
	if err != nil {
		verr.Add(s, "invalid authorization URL %q: %s", s.AuthorizationURL, err)
	}
	_, err = url.Parse(s.RefreshURL)
	if err != nil {
		verr.Add(s, "invalid refresh URL %q: %s", s.RefreshURL, err)
	}
	return verr
}

// Finalize makes the TokenURL and AuthorizationURL complete if needed.
func (s *FlowExpr) Finalize() {
	tu, _ := url.Parse(s.TokenURL)         // validated in Validate
	au, _ := url.Parse(s.AuthorizationURL) // validated in Validate
	ru, _ := url.Parse(s.RefreshURL)       // validated in Validate
	tokenOK := s.TokenURL == "" || tu.IsAbs()
	authOK := s.AuthorizationURL == "" || au.IsAbs()
	refreshOK := s.RefreshURL == "" || ru.IsAbs()
	if tokenOK && authOK && refreshOK {
		return
	}
	var (
		scheme string
		host   string
	)
	if len(httpdesign.Root.Design.API.Servers) > 0 {
		u, _ := url.Parse(httpdesign.Root.Design.API.Servers[0].URL)
		scheme = u.Scheme
		host = u.Host
	}
	if !tokenOK {
		tu.Scheme = scheme
		tu.Host = host
		s.TokenURL = tu.String()
	}
	if !authOK {
		au.Scheme = scheme
		au.Host = host
		s.AuthorizationURL = au.String()
	}
	if !refreshOK {
		ru.Scheme = scheme
		ru.Host = host
		s.RefreshURL = ru.String()
	}
}

// Type returns the grant type of the OAuth2 grant.
func (s *FlowExpr) Type() string {
	switch s.Kind {
	case AuthorizationCodeFlowKind:
		return "authorization_code"
	case ImplicitFlowKind:
		return "implicit"
	case PasswordFlowKind:
		return "password"
	case ClientCredentialsFlowKind:
		return "client_credentials"
	default:
		panic(fmt.Sprintf("unknown flow kind: %#v", s.Kind)) // bug
	}
}

// hasTaggedField returns true if the given attribute is an object that has an
// attribute with the given tag.
func hasTaggedField(att *design.AttributeExpr, tag string) bool {
	// note: this is called before Finalize and therefore there is no guarantee
	// that the payload attribute has been set to Empty if nil.
	if att == nil {
		return false
	}
	return securityField(att, tag) != ""
}

// securityField returns the security attribute name with the given tag.
func securityField(att *design.AttributeExpr, tag string) string {
	obj := design.AsObject(att.Type)
	if obj == nil {
		return ""
	}
	for _, at := range *obj {
		if _, ok := at.Attribute.Metadata[tag]; ok {
			return at.Name
		}
	}
	return ""
}

// findKey finds the given key in the endpoint expression and returns the
// transport element name and the position (header, query, or body).
func findKey(e *httpdesign.EndpointExpr, keyAtt string) (string, string) {
	if n, exists := e.AllParams().FindKey(keyAtt); exists {
		return n, "query"
	} else if n, exists := e.MappedHeaders().FindKey(keyAtt); exists {
		return n, "header"
	} else if _, ok := e.Body.Metadata["http:body"]; ok {
		if e.Body.Find(keyAtt) != nil {
			return keyAtt, "body"
		}
		if m, ok := e.Body.Metadata["origin:attribute"]; ok && m[0] == keyAtt {
			return keyAtt, "body"
		}
	}
	return "", "header"
}

func addHeaderAttr(ep *httpdesign.EndpointExpr, name, suffix string) {
	headers := ep.Headers()
	hObj := design.AsObject(headers.Type)
	if hObj == nil {
		return
	}
	attName := name
	if suffix != "" {
		attName = attName + ":" + suffix
	}
	attr := ep.MethodExpr.Payload.Find(name)
	hObj.Set(attName, attr)
	if ep.MethodExpr.Payload.IsRequired(name) {
		if headers.Validation == nil {
			headers.Validation = &design.ValidationExpr{}
		}
		headers.Validation.AddRequired(name)
	}
}

func isEmpty(a *design.AttributeExpr) bool {
	if a.Type == design.Empty {
		return true
	}
	obj := design.AsObject(a.Type)
	if obj != nil {
		return len(*obj) == 0
	}
	return false
}
