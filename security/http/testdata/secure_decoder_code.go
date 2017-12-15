package testdata

var BasicAuthSecureDecoderCode = `// SecureDecodeLoginRequest returns a decoder for requests sent to the
// BasicAuth login endpoint that is security scheme aware.
func SecureDecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeLoginRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*basicauth.LoginPayload)
		user, pass, ok := r.BasicAuth()
		if !ok {
			return p, nil
		}
		payload.User = &user
		payload.Pass = &pass
		return payload, nil
	}
}
`

var OAuth2SecureDecoderCode = `// SecureDecodeLoginRequest returns a decoder for requests sent to the OAuth2
// login endpoint that is security scheme aware.
func SecureDecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeLoginRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*oauth2.LoginPayload)
		h := r.Header.Get("Authorization")
		if h == "" {
			return p, nil
		}
		token := strings.TrimPrefix(h, "Bearer ")
		payload.Token = &token
		return payload, nil
	}
}
`

var OAuth2InParamSecureDecoderCode = `// SecureDecodeLoginRequest returns a decoder for requests sent to the
// OAuth2InParam login endpoint that is security scheme aware.
func SecureDecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeLoginRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*oauth2inparam.LoginPayload)
		token := r.URL.Query().Get("t")
		if token == "" {
			return p, nil
		}
		payload.Token = &token
		return payload, nil
	}
}
`

var JWTSecureDecoderCode = `// SecureDecodeLoginRequest returns a decoder for requests sent to the JWT
// login endpoint that is security scheme aware.
func SecureDecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeLoginRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*jwt.LoginPayload)
		h := r.Header.Get("Authorization")
		if h == "" {
			return p, nil
		}
		token := strings.TrimPrefix(h, "Bearer ")
		payload.Token = &token
		return payload, nil
	}
}
`

var APIKeySecureDecoderCode = `// SecureDecodeLoginRequest returns a decoder for requests sent to the APIKey
// login endpoint that is security scheme aware.
func SecureDecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeLoginRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*apikey.LoginPayload)
		key := r.Header.Get("Authorization")
		if key == "" {
			return p, nil
		}
		payload.Key = &key
		return payload, nil
	}
}
`

var APIKeyInParamSecureDecoderCode = `// SecureDecodeLoginRequest returns a decoder for requests sent to the
// APIKeyInParam login endpoint that is security scheme aware.
func SecureDecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeLoginRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*apikeyinparam.LoginPayload)
		key := r.URL.Query().Get("key")
		if key == "" {
			return p, nil
		}
		payload.Key = &key
		return payload, nil
	}
}
`

var MultipleAndSecureDecoderCode = `// SecureDecodeLoginRequest returns a decoder for requests sent to the
// MultipleAnd login endpoint that is security scheme aware.
func SecureDecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeLoginRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*multipleand.LoginPayload)
		user, pass, ok := r.BasicAuth()
		if !ok {
			return p, nil
		}
		payload.User = &user
		payload.Password = &pass
		key := r.URL.Query().Get("k")
		if key == "" {
			return p, nil
		}
		payload.Key = &key
		return payload, nil
	}
}
`

var MultipleOrSecureDecoderCode = `// SecureDecodeLoginRequest returns a decoder for requests sent to the
// MultipleOr login endpoint that is security scheme aware.
func SecureDecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeLoginRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*multipleor.LoginPayload)
		user, pass, ok := r.BasicAuth()
		if !ok {
			return p, nil
		}
		payload.User = &user
		payload.Password = &pass
		key := r.URL.Query().Get("k")
		if key == "" {
			return p, nil
		}
		payload.Key = &key
		return payload, nil
	}
}
`
