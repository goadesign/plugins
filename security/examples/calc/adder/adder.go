package calc

import (
	"context"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
	adder "goa.design/plugins/security/examples/calc/adder/gen/adder"
)

// ErrUnauthorized is the error returned by Login when the request credentials
// are invalid.
var ErrUnauthorized error = &adder.Unauthorized{"invalid token"}

// ErrInvalidScopes is the error returned by Login when the scopes provided in
// the JWT token claims are invalid.
var ErrInvalidScopes error = &adder.InvalidScopes{"invalid scopes, requires 'calc:add'"}

// adder service example implementation.
// The example methods log the requests and return zero values.
type adderSvc struct {
	logger *log.Logger
}

// NewAdder returns the adder service implementation.
func NewAdder(logger *log.Logger) adder.Service {
	return &adderSvc{logger}
}

// This action returns the sum of two integers and is secured with the API key
// scheme.
func (s *adderSvc) Add(ctx context.Context, p *adder.AddPayload) (int, error) {
	claims := make(jwt.MapClaims)

	// authorize request
	// 1. parse JWT token, token key is hardcoded to "secret" in this example
	t, err := jwt.ParseWithClaims(p.Key, claims, func(_ *jwt.Token) (interface{}, error) { return "secret", nil })
	if err != nil {
		return 0, ErrUnauthorized
	}

	// 2. validate provided "scopes" claim
	if claims["scopes"] == nil {
		return 0, ErrInvalidScopes
	}
	scopes, ok := claims["scopes"].([]interface{})
	if !ok {
		return 0, ErrInvalidScopes
	}
	hasAddScope := false
	for _, s := range scopes {
		if s == "calc:add" {
			hasAddScope = true
			break
		}
	}
	if !hasAddScope {
		return 0, ErrInvalidScopes
	}

	// you shouldn't do that in production...
	s.logger.Printf("Add request authorized with token '%s'", t.Raw)

	return p.A + p.B, nil
}
