// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// secured_service HTTP server encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/security/examples/multi_auth/design

package server

import (
	"context"
	"io"
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
			body SigninRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}

		return NewSigninSigninPayload(&body), nil
	}
}

// EncodeSigninError returns an encoder for errors returned by the signin
// secured_service endpoint.
func EncodeSigninError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		switch res := v.(type) {
		case *securedservice.Unauthorized:
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
			body SecureRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}

		var (
			fail *bool
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
		if err != nil {
			return nil, err
		}

		return NewSecureSecurePayload(&body, fail), nil
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
			body DoublySecureRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}

		var (
			key *string
		)
		keyRaw := r.URL.Query().Get("k")
		if keyRaw != "" {
			key = &keyRaw
		}

		return NewDoublySecureDoublySecurePayload(&body, key), nil
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
			body AlsoDoublySecureRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}

		var (
			key *string
		)
		keyRaw := r.Header.Get("Authorization")
		if keyRaw != "" {
			key = &keyRaw
		}

		return NewAlsoDoublySecureAlsoDoublySecurePayload(&body, key), nil
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
