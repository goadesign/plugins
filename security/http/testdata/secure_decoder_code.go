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
			user = ""
			pass = ""
		}
		payload.User = &user
		payload.Pass = &pass
		return payload, nil
	}
}
`

var BasicAuthRequiredSecureDecoderCode = `// SecureDecodeLoginRequest returns a decoder for requests sent to the
// BasicAuthRequired login endpoint that is security scheme aware.
func SecureDecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeLoginRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*basicauthrequired.LoginPayload)
		user, pass, ok := r.BasicAuth()
		if !ok {
			return nil, goa.MissingFieldError("Authorization", "header")
		}
		payload.User = user
		payload.Pass = pass
		return payload, nil
	}
}
`

var BasicAuthDecoderCode = `// DecodeLoginRequest returns a decoder for requests sent to the BasicAuth
// login endpoint.
func DecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {

		return NewLoginLoginPayload(), nil
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
		if payload.Token == nil {
			token := ""
			payload.Token = &token
		} else if strings.Contains(*payload.Token, " ") {
			payload.Token = &(strings.SplitN(*payload.Token, " ", 2)[1])
		}
		return payload, nil
	}
}
`

var OAuth2RequiredSecureDecoderCode = `// SecureDecodeLoginRequest returns a decoder for requests sent to the
// OAuth2Required login endpoint that is security scheme aware.
func SecureDecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeLoginRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*oauth2required.LoginPayload)
		if strings.Contains(payload.Token, " ") {
			payload.Token = strings.SplitN(payload.Token, " ", 2)[1]
		}
		return payload, nil
	}
}
`

var OAuth2DecoderCode = `// DecodeLoginRequest returns a decoder for requests sent to the OAuth2 login
// endpoint.
func DecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			token *string
		)
		tokenRaw := r.Header.Get("Authorization")
		if tokenRaw != "" {
			token = &tokenRaw
		}

		return NewLoginLoginPayload(token), nil
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
		if payload.Token == nil {
			token := ""
			payload.Token = &token
		} else if strings.Contains(*payload.Token, " ") {
			payload.Token = &(strings.SplitN(*payload.Token, " ", 2)[1])
		}
		return payload, nil
	}
}
`

var JWTRequiredSecureDecoderCode = `// SecureDecodeLoginRequest returns a decoder for requests sent to the
// JWTRequired login endpoint that is security scheme aware.
func SecureDecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeLoginRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*jwtrequired.LoginPayload)
		if strings.Contains(payload.Token, " ") {
			payload.Token = strings.SplitN(payload.Token, " ", 2)[1]
		}
		return payload, nil
	}
}
`

var JWTDecoderCode = `// DecodeLoginRequest returns a decoder for requests sent to the JWT login
// endpoint.
func DecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			token *string
		)
		tokenRaw := r.Header.Get("Authorization")
		if tokenRaw != "" {
			token = &tokenRaw
		}

		return NewLoginLoginPayload(token), nil
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
		if payload.Key == nil {
			key := ""
			payload.Key = &key
		} else if strings.Contains(*payload.Key, " ") {
			payload.Key = &(strings.SplitN(*payload.Key, " ", 2)[1])
		}
		return payload, nil
	}
}
`

var APIKeyRequiredSecureDecoderCode = `// SecureDecodeLoginRequest returns a decoder for requests sent to the
// APIKeyRequired login endpoint that is security scheme aware.
func SecureDecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeLoginRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*apikeyrequired.LoginPayload)
		if strings.Contains(payload.Key, " ") {
			payload.Key = strings.SplitN(payload.Key, " ", 2)[1]
		}
		return payload, nil
	}
}
`

var APIKeyDecoderCode = `// DecodeLoginRequest returns a decoder for requests sent to the APIKey login
// endpoint.
func DecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			key *string
		)
		keyRaw := r.Header.Get("Authorization")
		if keyRaw != "" {
			key = &keyRaw
		}

		return NewLoginLoginPayload(key), nil
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
		if payload.Key == nil {
			key := ""
			payload.Key = &key
		} else if strings.Contains(*payload.Key, " ") {
			payload.Key = &(strings.SplitN(*payload.Key, " ", 2)[1])
		}
		return payload, nil
	}
}
`

var APIKeyInParamRequiredSecureDecoderCode = `// SecureDecodeLoginRequest returns a decoder for requests sent to the
// APIKeyInParamRequired login endpoint that is security scheme aware.
func SecureDecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeLoginRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*apikeyinparamrequired.LoginPayload)
		if strings.Contains(payload.Key, " ") {
			payload.Key = strings.SplitN(payload.Key, " ", 2)[1]
		}
		return payload, nil
	}
}
`

var APIKeyInParamDecoderCode = `// DecodeLoginRequest returns a decoder for requests sent to the APIKeyInParam
// login endpoint.
func DecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			key *string
		)
		keyRaw := r.URL.Query().Get("key")
		if keyRaw != "" {
			key = &keyRaw
		}

		return NewLoginLoginPayload(key), nil
	}
}
`

var APIKeyInBodySecureDecoderCode = `// SecureDecodeLoginRequest returns a decoder for requests sent to the
// APIKeyInBody login endpoint that is security scheme aware.
func SecureDecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeLoginRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*apikeyinbody.LoginPayload)
		if payload.Key == nil {
			key := ""
			payload.Key = &key
		} else if strings.Contains(*payload.Key, " ") {
			payload.Key = &(strings.SplitN(*payload.Key, " ", 2)[1])
		}
		return payload, nil
	}
}
`

var APIKeyInBodyDecoderCode = `// DecodeLoginRequest returns a decoder for requests sent to the APIKeyInBody
// login endpoint.
func DecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			body string
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}

		return NewLoginLoginPayload(body), nil
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
			user = ""
			pass = ""
		}
		payload.User = &user
		payload.Password = &pass
		if payload.Key == nil {
			key := ""
			payload.Key = &key
		} else if strings.Contains(*payload.Key, " ") {
			payload.Key = &(strings.SplitN(*payload.Key, " ", 2)[1])
		}
		return payload, nil
	}
}
`

var MultipleAndRequiredSecureDecoderCode = `// SecureDecodeLoginRequest returns a decoder for requests sent to the
// MultipleAndRequired login endpoint that is security scheme aware.
func SecureDecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeLoginRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*multipleandrequired.LoginPayload)
		user, pass, ok := r.BasicAuth()
		if !ok {
			return nil, goa.MissingFieldError("Authorization", "header")
		}
		payload.User = user
		payload.Password = pass
		if strings.Contains(payload.Key, " ") {
			payload.Key = strings.SplitN(payload.Key, " ", 2)[1]
		}
		return payload, nil
	}
}
`

var MultipleAndDecoderCode = `// DecodeLoginRequest returns a decoder for requests sent to the MultipleAnd
// login endpoint.
func DecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			key *string
		)
		keyRaw := r.URL.Query().Get("k")
		if keyRaw != "" {
			key = &keyRaw
		}

		return NewLoginLoginPayload(key), nil
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
			user = ""
			pass = ""
		}
		payload.User = &user
		payload.Password = &pass
		if payload.Key == nil {
			key := ""
			payload.Key = &key
		} else if strings.Contains(*payload.Key, " ") {
			payload.Key = &(strings.SplitN(*payload.Key, " ", 2)[1])
		}
		return payload, nil
	}
}
`

var MultipleOrRequiredSecureDecoderCode = `// SecureDecodeLoginRequest returns a decoder for requests sent to the
// MultipleOrRequired login endpoint that is security scheme aware.
func SecureDecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeLoginRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*multipleorrequired.LoginPayload)
		user, pass, ok := r.BasicAuth()
		if !ok {
			return nil, goa.MissingFieldError("Authorization", "header")
		}
		payload.User = user
		payload.Password = pass
		if strings.Contains(payload.Key, " ") {
			payload.Key = strings.SplitN(payload.Key, " ", 2)[1]
		}
		return payload, nil
	}
}
`

var MultipleOrDecoderCode = `// DecodeLoginRequest returns a decoder for requests sent to the MultipleOr
// login endpoint.
func DecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			key *string
		)
		keyRaw := r.URL.Query().Get("k")
		if keyRaw != "" {
			key = &keyRaw
		}

		return NewLoginLoginPayload(key), nil
	}
}
`

var SchemesInTypeSecureDecoderCode = `// SecureDecodeLoginRequest returns a decoder for requests sent to the
// SchemesInTypeDSL login endpoint that is security scheme aware.
func SecureDecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeLoginRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*schemesintypedsl.Schemes)
		if payload.Key == nil {
			key := ""
			payload.Key = &key
		} else if strings.Contains(*payload.Key, " ") {
			payload.Key = &(strings.SplitN(*payload.Key, " ", 2)[1])
		}
		if payload.Token == nil {
			token := ""
			payload.Token = &token
		} else if strings.Contains(*payload.Token, " ") {
			payload.Token = &(strings.SplitN(*payload.Token, " ", 2)[1])
		}
		return payload, nil
	}
}
`

var SchemesInTypeRequiredSecureDecoderCode = `// SecureDecodeLoginRequest returns a decoder for requests sent to the
// SchemesInTypeRequiredDSL login endpoint that is security scheme aware.
func SecureDecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	rawDecoder := DecodeLoginRequest(mux, decoder)
	return func(r *http.Request) (interface{}, error) {
		p, err := rawDecoder(r)
		if err != nil {
			return nil, err
		}
		payload := p.(*schemesintyperequireddsl.SchemesRequired)
		if strings.Contains(payload.Key, " ") {
			payload.Key = strings.SplitN(payload.Key, " ", 2)[1]
		}
		if strings.Contains(payload.Token, " ") {
			payload.Token = strings.SplitN(payload.Token, " ", 2)[1]
		}
		return payload, nil
	}
}
`

var SchemesInTypeDecoderCode = `// DecodeLoginRequest returns a decoder for requests sent to the
// SchemesInTypeDSL login endpoint.
func DecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
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

		return NewLoginSchemes(key, token), nil
	}
}
`

var MultipleSchemesWithParamsDecoderCode = `// DecodeLoginRequest returns a decoder for requests sent to the
// MultipleSchemesWithParams login endpoint.
func DecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			body LoginRequestBody
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
			id        int
			key       *string
			name      *string
			userAgent *string
			token     string

			params = mux.Vars(r)
		)
		{
			idRaw := params["id"]
			v, err2 := strconv.ParseInt(idRaw, 10, strconv.IntSize)
			if err2 != nil {
				err = goa.MergeErrors(err, goa.InvalidFieldTypeError("id", idRaw, "integer"))
			}
			id = int(v)
		}
		keyRaw := r.URL.Query().Get("k")
		if keyRaw != "" {
			key = &keyRaw
		}
		nameRaw := r.URL.Query().Get("name")
		if nameRaw != "" {
			name = &nameRaw
		}
		userAgentRaw := r.Header.Get("User-Agent")
		if userAgentRaw != "" {
			userAgent = &userAgentRaw
		}
		token = r.Header.Get("Authorization")
		if token == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("Authorization", "header"))
		}
		if err != nil {
			return nil, err
		}

		return NewLoginLoginPayload(&body, id, key, name, userAgent, token), nil
	}
}
`
