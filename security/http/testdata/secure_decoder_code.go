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
		hOAuth2 := r.Header.Get("Authorization")
		if hOAuth2 == "" {
			return p, nil
		}
		tokenOAuth2 := strings.TrimPrefix(hOAuth2, "Bearer ")
		payload.Token = &tokenOAuth2
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
		tokenOAuth2 := r.URL.Query().Get("t")
		if tokenOAuth2 == "" {
			return p, nil
		}
		payload.Token = &tokenOAuth2
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
		hJWT := r.Header.Get("Authorization")
		if hJWT == "" {
			return p, nil
		}
		tokenJWT := strings.TrimPrefix(hJWT, "Bearer ")
		payload.Token = &tokenJWT
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

var SameSchemeMethod1DecoderCode = `// SecureDecodeMethod1Request returns a decoder for requests sent to the
// SameSchemeMultipleEndpoints method1 endpoint that is security scheme aware.
func SecureDecodeMethod1Request(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeMethod1Request(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*sameschememultipleendpoints.Method1Payload)
		key := r.URL.Query().Get("k")
		if key == "" {
			return p, nil
		}
		payload.Key = &key
		return payload, nil
	}
}
`

var SameSchemeMethod2DecoderCode = `// SecureDecodeMethod2Request returns a decoder for requests sent to the
// SameSchemeMultipleEndpoints method2 endpoint that is security scheme aware.
func SecureDecodeMethod2Request(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeMethod2Request(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*sameschememultipleendpoints.Method2Payload)
		key := r.Header.Get("Authorization")
		if key == "" {
			return p, nil
		}
		payload.Key = &key
		return payload, nil
	}
}
`
