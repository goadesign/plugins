package calc

import (
	"context"

	jwt "github.com/dgrijalva/jwt-go"
	"goa.design/plugins/security"
	adder "goa.design/plugins/security/examples/calc/adder/gen/adder"
	"goa.design/plugins/security/examples/calc/calc/gen/calc"
)

var (
	// ErrUnauthorized is the error returned by Login when the request credentials
	// are invalid.
	ErrUnauthorized error = calcsvc.Unauthorized("invalid username and password combination")

	// ErrInvalidToken is the error returned by Login when the request credentials
	// are invalid.
	ErrInvalidToken error = adder.Unauthorized("invalid token")

	// ErrInvalidScopes is the error returned by Login when the scopes provided in
	// the JWT token claims are invalid.
	ErrInvalidScopes error = adder.InvalidScopes("invalid scopes, requires 'calc:add'")

	// Key is the key used in JWT authentication
	Key = []byte("secret")
)

// BasicAuthFunc implements the basic auth scheme.
func BasicAuthFunc(ctx context.Context, username, password string, s *security.BasicAuthScheme) (context.Context, error) {
	if username != "goa" {
		return ctx, ErrUnauthorized
	}
	if password != "rocks" {
		return ctx, ErrUnauthorized
	}
	return ctx, nil
}

// JWTAuthFunc implements the JWT security scheme.
func JWTAuthFunc(ctx context.Context, token string, s *security.JWTScheme) (context.Context, error) {
	claims := make(jwt.MapClaims)

	// authorize request
	// 1. parse JWT token, token key is hardcoded to "secret" in this example
	_, err := jwt.ParseWithClaims(token, claims, func(_ *jwt.Token) (interface{}, error) { return Key, nil })
	if err != nil {
		return ctx, ErrInvalidToken
	}

	// 2. validate provided "scopes" claim
	if claims["scopes"] == nil {
		return ctx, ErrInvalidScopes
	}
	scopes, ok := claims["scopes"].([]interface{})
	if !ok {
		return ctx, ErrInvalidScopes
	}
	hasAddScope := false
	for _, s := range scopes {
		if s == "calc:add" {
			hasAddScope = true
			break
		}
	}
	if !hasAddScope {
		return ctx, ErrInvalidScopes
	}

	return ctx, nil
}
