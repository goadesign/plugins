package security

import (
	"fmt"
	"strings"
)

const (
	// ContextKey is the key used to store the request alternate security
	// requirements in the context.
	ContextKey = "SecurityRequirements"
)

type (
	// Requirement lists security schemes that requests must fullfill to be
	// authorized.
	Requirement struct {
		// RequiredScopes lists the JWT or OAuth2 scopes that the
		// request credentials must provide for authorization.
		RequiredScopes []string
		// Schemes list the security schemes. The request must satisfy
		// all the schemes to be authorized.
		Schemes []*Scheme
	}

	// Scheme defines a security scheme: basic auth, API key, JWT or OAuth2.
	Scheme struct {
		// Kind identifies the type of security scheme.
		Kind SchemeKind
		// Name is the scheme name as defined in the design.
		Name string
		// Scopes list the scopes that are defined by the scheme.
		// Applies to JWT and OAuth2 schemes only. Note: this lists the
		// scopes that are possible - not the scopes actually provided
		// in the request. The latter is provided as part of the
		// requirements stored in the request context.
		Scopes []string
		// Flows provide information specific to the OAuth2 security
		// scheme.
		Flows []*OAuthFlow
	}

	// OAuthFlow provide information about a OAuth2 security scheme.
	OAuthFlow struct {
		// Kind identifies the OAuth2 flow as defined in RFC 6749, one
		// of Authorization Code, Implicit, Resource Owner Password
		// Credentials or Client Credentials.
		Kind FlowKind
		// AuthorizationURL is the URL to the endpoint providing
		// the authorization codes.
		AuthorizationURL string
		// TokenURL is the URL to the endpoint providing the access
		// tokens.
		TokenURL string
		// RefreshURL is the URL to the endpoint providing the refresh
		// tokens.
		RefreshURL string
	}

	// SchemeKind defines the type of security scheme.
	SchemeKind int

	// FlowKind defines the type of OAuth2 flow.
	FlowKind int
)

const (
	// OAuth2Kind identifies a OAuth2 security scheme.
	OAuth2Kind SchemeKind = iota + 1
	// BasicAuthKind identifies a basic auth security scheme.
	BasicAuthKind
	// APIKeyKind identifies an API key security scheme.
	APIKeyKind
	// JWTKind identifies an API key security scheme using JWT tokens and
	// with support for scopes.
	JWTKind
)

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

// Validate returns a non-nil error if scopes does not contain all of req
// required scopes.
func (req *Requirement) Validate(scopes []string) error {
	if len(req.RequiredScopes) == 0 {
		return nil
	}
	var missing []string
	for _, r := range req.RequiredScopes {
		found := false
		for _, s := range scopes {
			if s == r {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, r)
		}
	}
	if len(missing) == 0 {
		return nil
	}
	return fmt.Errorf("missing sopes: %s", strings.Join(missing, ", "))
}
