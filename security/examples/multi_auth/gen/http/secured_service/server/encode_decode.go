// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// secured_service HTTP server encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/security/examples/multi_auth/design

package server

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	goa "goa.design/goa"
	goahttp "goa.design/goa/http"
	securedservice "goa.design/plugins/security/examples/multi_auth/gen/secured_service"
)

// EncodeSigninResponse returns an encoder for responses returned by the
// secured_service signin endpoint.
func EncodeSigninResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(string)
		w.Header().Set("Authorization", res)
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}

// DecodeSigninRequest returns a decoder for requests sent to the
// secured_service signin endpoint.
func DecodeSigninRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			username *string
			password *string
		)
		usernameRaw := r.Header.Get("Authorization")
		if usernameRaw != "" {
			username = &usernameRaw
		}
		passwordRaw := r.Header.Get("Authorization")
		if passwordRaw != "" {
			password = &passwordRaw
		}

		return NewSigninSigninPayload(username, password), nil
	}
}

// EncodeSigninError returns an encoder for errors returned by the signin
// secured_service endpoint.
func EncodeSigninError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		switch res := v.(type) {
		case securedservice.Unauthorized:
			enc := encoder(ctx, w)
			body := NewSigninUnauthorizedResponseBody(res)
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
		return nil
	}
}

// EncodeSecureResponse returns an encoder for responses returned by the
// secured_service secure endpoint.
func EncodeSecureResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(string)
		enc := encoder(ctx, w)
		body := res
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeSecureRequest returns a decoder for requests sent to the
// secured_service secure endpoint.
func DecodeSecureRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			fail  *bool
			token *string
			err   error
		)
		{
			failRaw := r.URL.Query().Get("fail")
			if failRaw != "" {
				v, err2 := strconv.ParseBool(failRaw)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("fail", failRaw, "boolean"))
				}
				fail = &v
			}
		}
		tokenRaw := r.Header.Get("Authorization")
		if tokenRaw != "" {
			token = &tokenRaw
		}
		if err != nil {
			return nil, err
		}

		return NewSecureSecurePayload(fail, token), nil
	}
}

// EncodeSecureError returns an encoder for errors returned by the secure
// secured_service endpoint.
func EncodeSecureError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		switch res := v.(type) {
		case securedservice.Unauthorized:
			enc := encoder(ctx, w)
			body := NewSecureUnauthorizedResponseBody(res)
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
		return nil
	}
}

// EncodeDoublySecureResponse returns an encoder for responses returned by the
// secured_service doubly_secure endpoint.
func EncodeDoublySecureResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(string)
		enc := encoder(ctx, w)
		body := res
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeDoublySecureRequest returns a decoder for requests sent to the
// secured_service doubly_secure endpoint.
func DecodeDoublySecureRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			key   *string
			token *string
		)
		keyRaw := r.URL.Query().Get("k")
		if keyRaw != "" {
			key = &keyRaw
		}
		tokenRaw := r.Header.Get("Authorization")
		if tokenRaw != "" {
			token = &tokenRaw
		}

		return NewDoublySecureDoublySecurePayload(key, token), nil
	}
}

// EncodeDoublySecureError returns an encoder for errors returned by the
// doubly_secure secured_service endpoint.
func EncodeDoublySecureError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		switch res := v.(type) {
		case securedservice.Unauthorized:
			enc := encoder(ctx, w)
			body := NewDoublySecureUnauthorizedResponseBody(res)
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
		return nil
	}
}

// EncodeAlsoDoublySecureResponse returns an encoder for responses returned by
// the secured_service also_doubly_secure endpoint.
func EncodeAlsoDoublySecureResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(string)
		enc := encoder(ctx, w)
		body := res
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeAlsoDoublySecureRequest returns a decoder for requests sent to the
// secured_service also_doubly_secure endpoint.
func DecodeAlsoDoublySecureRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			key        *string
			token      *string
			oauthToken *string
			username   *string
			password   *string
		)
		keyRaw := r.Header.Get("Authorization")
		if keyRaw != "" {
			key = &keyRaw
		}
		tokenRaw := r.Header.Get("Authorization")
		if tokenRaw != "" {
			token = &tokenRaw
		}
		oauthTokenRaw := r.Header.Get("Authorization")
		if oauthTokenRaw != "" {
			oauthToken = &oauthTokenRaw
		}
		usernameRaw := r.Header.Get("Authorization")
		if usernameRaw != "" {
			username = &usernameRaw
		}
		passwordRaw := r.Header.Get("Authorization")
		if passwordRaw != "" {
			password = &passwordRaw
		}

		return NewAlsoDoublySecureAlsoDoublySecurePayload(key, token, oauthToken, username, password), nil
	}
}

// EncodeAlsoDoublySecureError returns an encoder for errors returned by the
// also_doubly_secure secured_service endpoint.
func EncodeAlsoDoublySecureError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		switch res := v.(type) {
		case securedservice.Unauthorized:
			enc := encoder(ctx, w)
			body := NewAlsoDoublySecureUnauthorizedResponseBody(res)
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
		return nil
	}
}

// SecureDecodeSigninRequest returns a decoder for requests sent to the
// secured_service signin endpoint that is security scheme aware.
func SecureDecodeSigninRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeSigninRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*securedservice.SigninPayload)
		user, pass, ok := r.BasicAuth()
		if !ok {
			return p, nil
		}
		payload.Username = &user
		payload.Password = &pass
		return payload, nil
	}
}

// SecureDecodeSecureRequest returns a decoder for requests sent to the
// secured_service secure endpoint that is security scheme aware.
func SecureDecodeSecureRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeSecureRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*securedservice.SecurePayload)
		hJWT := r.Header.Get("Authorization")
		if hJWT == "" {
			return p, nil
		}
		tokenJWT := strings.TrimPrefix(hJWT, "Bearer ")
		payload.Token = &tokenJWT
		return payload, nil
	}
}

// SecureDecodeDoublySecureRequest returns a decoder for requests sent to the
// secured_service doubly_secure endpoint that is security scheme aware.
func SecureDecodeDoublySecureRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeDoublySecureRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*securedservice.DoublySecurePayload)
		hJWT := r.Header.Get("Authorization")
		if hJWT == "" {
			return p, nil
		}
		tokenJWT := strings.TrimPrefix(hJWT, "Bearer ")
		payload.Token = &tokenJWT
		key := r.URL.Query().Get("k")
		if key == "" {
			return p, nil
		}
		payload.Key = &key
		return payload, nil
	}
}

// SecureDecodeAlsoDoublySecureRequest returns a decoder for requests sent to
// the secured_service also_doubly_secure endpoint that is security scheme
// aware.
func SecureDecodeAlsoDoublySecureRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeAlsoDoublySecureRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*securedservice.AlsoDoublySecurePayload)
		hJWT := r.Header.Get("Authorization")
		if hJWT == "" {
			return p, nil
		}
		tokenJWT := strings.TrimPrefix(hJWT, "Bearer ")
		payload.Token = &tokenJWT
		key := r.Header.Get("Authorization")
		if key == "" {
			return p, nil
		}
		payload.Key = &key
		hOAuth2 := r.Header.Get("Authorization")
		if hOAuth2 == "" {
			return p, nil
		}
		tokenOAuth2 := strings.TrimPrefix(hOAuth2, "Bearer ")
		payload.OauthToken = &tokenOAuth2
		user, pass, ok := r.BasicAuth()
		if !ok {
			return p, nil
		}
		payload.Username = &user
		payload.Password = &pass
		return payload, nil
	}
}
