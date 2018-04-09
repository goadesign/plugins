package calc

import (
	"context"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"goa.design/plugins/security"
	"goa.design/plugins/security/examples/calc/calc/gen/calc"
)

var (
	// ErrLoginFailed is the error returned by Login when the request credentials
	// are invalid.
	ErrLoginFailed error = calcsvc.Unauthorized("invalid username and password combination")

	// ErrForbidden is the error returned by Add when the JWT token doesn't contain
	// the required scopes.
	ErrForbidden error = calcsvc.MakeForbidden(fmt.Errorf("forbidden, requires 'calc:add'"))

	// ErrInvalidToken is the error returned by Login when the request credentials
	// are invalid.
	ErrInvalidToken error = calcsvc.Unauthorized("invalid token")

	// ErrInvalidScopes is the error returned by Add when the scopes provided in
	// the JWT token claims are invalid.
	ErrInvalidScopes error = calcsvc.Unauthorized("invalid scopes")

	// Key is the key used in JWT authentication
	Key = []byte("secret")
)

// BasicAuthFunc implements the basic auth scheme.
func BasicAuthFunc(ctx context.Context, username, password string, s *security.BasicAuthScheme) (context.Context, error) {
	if username != "goa" {
		return ctx, ErrLoginFailed
	}
	if password != "rocks" {
		return ctx, ErrLoginFailed
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
		return ctx, ErrForbidden
	}

	return ctx, nil
}
