package testdata

var BasicAuthSecureEncoderCode = `// SecureEncodeLoginRequest returns an encoder for requests sent to the
// BasicAuth login endpoint that is security scheme aware.
func SecureEncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeLoginRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*basicauth.LoginPayload)
		req.SetBasicAuth(*payload.User, *payload.Pass)
		return nil
	}
}
`

var OAuth2SecureEncoderCode = `// SecureEncodeLoginRequest returns an encoder for requests sent to the OAuth2
// login endpoint that is security scheme aware.
func SecureEncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeLoginRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*oauth2.LoginPayload)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *payload.Token))
		return nil
	}
}
`

var OAuth2InParamSecureEncoderCode = `// SecureEncodeLoginRequest returns an encoder for requests sent to the
// OAuth2InParam login endpoint that is security scheme aware.
func SecureEncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeLoginRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*oauth2inparam.LoginPayload)
		values := req.URL.Query()
		values.Add("t", *payload.Token)
		req.URL.RawQuery = values.Encode()
		return nil
	}
}
`

var JWTSecureEncoderCode = `// SecureEncodeLoginRequest returns an encoder for requests sent to the JWT
// login endpoint that is security scheme aware.
func SecureEncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeLoginRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*jwt.LoginPayload)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *payload.Token))
		return nil
	}
}
`

var APIKeySecureEncoderCode = `// SecureEncodeLoginRequest returns an encoder for requests sent to the APIKey
// login endpoint that is security scheme aware.
func SecureEncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeLoginRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*apikey.LoginPayload)
		req.Header.Add("Authorization", *payload.Key)
		return nil
	}
}
`

var APIKeyInParamSecureEncoderCode = `// SecureEncodeLoginRequest returns an encoder for requests sent to the
// APIKeyInParam login endpoint that is security scheme aware.
func SecureEncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeLoginRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*apikeyinparam.LoginPayload)
		values := req.URL.Query()
		values.Add("key", *payload.Key)
		req.URL.RawQuery = values.Encode()
		return nil
	}
}
`

var MultipleAndSecureEncoderCode = `// SecureEncodeLoginRequest returns an encoder for requests sent to the
// MultipleAnd login endpoint that is security scheme aware.
func SecureEncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeLoginRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*multipleand.LoginPayload)
		values := req.URL.Query()
		req.SetBasicAuth(*payload.User, *payload.Password)
		values.Add("k", *payload.Key)
		req.URL.RawQuery = values.Encode()
		return nil
	}
}
`

var MultipleOrSecureEncoderCode = `// SecureEncodeLoginRequest returns an encoder for requests sent to the
// MultipleOr login endpoint that is security scheme aware.
func SecureEncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeLoginRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*multipleor.LoginPayload)
		values := req.URL.Query()
		req.SetBasicAuth(*payload.User, *payload.Password)
		values.Add("k", *payload.Key)
		req.URL.RawQuery = values.Encode()
		return nil
	}
}
`
