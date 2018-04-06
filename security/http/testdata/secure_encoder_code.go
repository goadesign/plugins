package testdata

var BasicAuthSecureEncoderCode = `// SecureEncodeLoginRequest returns an encoder for requests sent to the
// BasicAuth login endpoint that is security scheme aware.
func SecureEncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		payload := v.(*basicauth.LoginPayload)
		if payload.User == nil {
			user := ""
			payload.User = &user
		}
		if payload.Pass == nil {
			pass := ""
			payload.Pass = &pass
		}
		req.SetBasicAuth(*payload.User, *payload.Pass)
		return nil
	}
}
`

var BasicAuthRequiredSecureEncoderCode = `// SecureEncodeLoginRequest returns an encoder for requests sent to the
// BasicAuthRequired login endpoint that is security scheme aware.
func SecureEncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeLoginRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*basicauthrequired.LoginPayload)
		req.SetBasicAuth(payload.User, payload.Pass)
		return nil
	}
}
`

var BasicAuthRequiredEncoderCode = `// EncodeLoginRequest returns an encoder for requests sent to the
// BasicAuthRequired login server.
func EncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*basicauthrequired.LoginPayload)
		if !ok {
			return goahttp.ErrInvalidType("BasicAuthRequired", "login", "*basicauthrequired.LoginPayload", v)
		}
		values := req.URL.Query()
		if p.ID != nil {
			values.Add("id", *p.ID)
		}
		req.URL.RawQuery = values.Encode()
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
		if payload.Token == nil {
			req.Header.Set("Authorization", "")
		} else if !strings.Contains(*payload.Token, " ") {
			req.Header.Set("Authorization", "Bearer "+*payload.Token)
		}
		return nil
	}
}
`

var OAuth2RequiredSecureEncoderCode = `// SecureEncodeLoginRequest returns an encoder for requests sent to the
// OAuth2Required login endpoint that is security scheme aware.
func SecureEncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeLoginRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*oauth2required.LoginPayload)
		if !strings.Contains(payload.Token, " ") {
			req.Header.Set("Authorization", "Bearer "+payload.Token)
		}
		return nil
	}
}
`

var OAuth2EncoderCode = `// EncodeLoginRequest returns an encoder for requests sent to the OAuth2 login
// server.
func EncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*oauth2.LoginPayload)
		if !ok {
			return goahttp.ErrInvalidType("OAuth2", "login", "*oauth2.LoginPayload", v)
		}
		if p.Token != nil {
			req.Header.Set("Authorization", *p.Token)
		}
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
		if payload.Token == nil {
			values.Set("t", "")
		} else if strings.Contains(*payload.Token, " ") {
			s := strings.SplitN(*payload.Token, " ", 2)[1]
			values.Set("t", s)
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}
`

var OAuth2InParamRequiredSecureEncoderCode = `// SecureEncodeLoginRequest returns an encoder for requests sent to the
// OAuth2InParamRequired login endpoint that is security scheme aware.
func SecureEncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeLoginRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*oauth2inparamrequired.LoginPayload)
		values := req.URL.Query()
		if strings.Contains(payload.Token, " ") {
			s := strings.SplitN(payload.Token, " ", 2)[1]
			values.Set("t", s)
		}
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
		if payload.Token == nil {
			req.Header.Set("Authorization", "")
		} else if !strings.Contains(*payload.Token, " ") {
			req.Header.Set("Authorization", "Bearer "+*payload.Token)
		}
		return nil
	}
}
`

var JWTRequiredSecureEncoderCode = `// SecureEncodeLoginRequest returns an encoder for requests sent to the
// JWTRequired login endpoint that is security scheme aware.
func SecureEncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeLoginRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*jwtrequired.LoginPayload)
		if !strings.Contains(payload.Token, " ") {
			req.Header.Set("Authorization", "Bearer "+payload.Token)
		}
		return nil
	}
}
`

var JWTEncoderCode = `// EncodeLoginRequest returns an encoder for requests sent to the JWT login
// server.
func EncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*jwt.LoginPayload)
		if !ok {
			return goahttp.ErrInvalidType("JWT", "login", "*jwt.LoginPayload", v)
		}
		if p.Token != nil {
			req.Header.Set("Authorization", *p.Token)
		}
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
		if payload.Key == nil {
			req.Header.Set("Authorization", "")
		} else if strings.Contains(*payload.Key, " ") {
			s := strings.SplitN(*payload.Key, " ", 2)[1]
			req.Header.Set("Authorization", s)
		}
		return nil
	}
}
`

var APIKeyRequiredSecureEncoderCode = `// SecureEncodeLoginRequest returns an encoder for requests sent to the
// APIKeyRequired login endpoint that is security scheme aware.
func SecureEncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeLoginRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*apikeyrequired.LoginPayload)
		if strings.Contains(payload.Key, " ") {
			s := strings.SplitN(payload.Key, " ", 2)[1]
			req.Header.Set("Authorization", s)
		}
		return nil
	}
}
`

var APIKeyEncoderCode = `// EncodeLoginRequest returns an encoder for requests sent to the APIKey login
// server.
func EncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*apikey.LoginPayload)
		if !ok {
			return goahttp.ErrInvalidType("APIKey", "login", "*apikey.LoginPayload", v)
		}
		if p.Key != nil {
			req.Header.Set("Authorization", *p.Key)
		}
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
		if payload.Key == nil {
			values.Set("key", "")
		} else if strings.Contains(*payload.Key, " ") {
			s := strings.SplitN(*payload.Key, " ", 2)[1]
			values.Set("key", s)
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}
`

var APIKeyInParamRequiredSecureEncoderCode = `// SecureEncodeLoginRequest returns an encoder for requests sent to the
// APIKeyInParamRequired login endpoint that is security scheme aware.
func SecureEncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeLoginRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*apikeyinparamrequired.LoginPayload)
		values := req.URL.Query()
		if strings.Contains(payload.Key, " ") {
			s := strings.SplitN(payload.Key, " ", 2)[1]
			values.Set("key", s)
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}
`

var APIKeyInParamEncoderCode = `// EncodeLoginRequest returns an encoder for requests sent to the APIKeyInParam
// login server.
func EncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*apikeyinparam.LoginPayload)
		if !ok {
			return goahttp.ErrInvalidType("APIKeyInParam", "login", "*apikeyinparam.LoginPayload", v)
		}
		values := req.URL.Query()
		if p.Key != nil {
			values.Add("key", *p.Key)
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}
`

var APIKeyInBodySecureEncoderCode = `// SecureEncodeLoginRequest returns an encoder for requests sent to the
// APIKeyInBody login endpoint that is security scheme aware.
func SecureEncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeLoginRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		return nil
	}
}
`

var APIKeyInBodyEncoderCode = `// EncodeLoginRequest returns an encoder for requests sent to the APIKeyInBody
// login server.
func EncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*apikeyinbody.LoginPayload)
		if !ok {
			return goahttp.ErrInvalidType("APIKeyInBody", "login", "*apikeyinbody.LoginPayload", v)
		}
		body := p
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("APIKeyInBody", "login", err)
		}
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
		if payload.User == nil {
			user := ""
			payload.User = &user
		}
		if payload.Password == nil {
			password := ""
			payload.Password = &password
		}
		req.SetBasicAuth(*payload.User, *payload.Password)
		if payload.Key == nil {
			values.Set("k", "")
		} else if strings.Contains(*payload.Key, " ") {
			s := strings.SplitN(*payload.Key, " ", 2)[1]
			values.Set("k", s)
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}
`

var MultipleAndRequiredSecureEncoderCode = `// SecureEncodeLoginRequest returns an encoder for requests sent to the
// MultipleAndRequired login endpoint that is security scheme aware.
func SecureEncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeLoginRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*multipleandrequired.LoginPayload)
		values := req.URL.Query()
		req.SetBasicAuth(payload.User, payload.Password)
		if strings.Contains(payload.Key, " ") {
			s := strings.SplitN(payload.Key, " ", 2)[1]
			values.Set("k", s)
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}
`

var MultipleAndEncoderCode = `// EncodeLoginRequest returns an encoder for requests sent to the MultipleAnd
// login server.
func EncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*multipleand.LoginPayload)
		if !ok {
			return goahttp.ErrInvalidType("MultipleAnd", "login", "*multipleand.LoginPayload", v)
		}
		values := req.URL.Query()
		if p.Key != nil {
			values.Add("k", *p.Key)
		}
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
		if payload.User == nil {
			user := ""
			payload.User = &user
		}
		if payload.Password == nil {
			password := ""
			payload.Password = &password
		}
		req.SetBasicAuth(*payload.User, *payload.Password)
		if payload.Key == nil {
			values.Set("k", "")
		} else if strings.Contains(*payload.Key, " ") {
			s := strings.SplitN(*payload.Key, " ", 2)[1]
			values.Set("k", s)
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}
`

var MultipleOrRequiredSecureEncoderCode = `// SecureEncodeLoginRequest returns an encoder for requests sent to the
// MultipleOrRequired login endpoint that is security scheme aware.
func SecureEncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeLoginRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*multipleorrequired.LoginPayload)
		values := req.URL.Query()
		req.SetBasicAuth(payload.User, payload.Password)
		if strings.Contains(payload.Key, " ") {
			s := strings.SplitN(payload.Key, " ", 2)[1]
			values.Set("k", s)
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}
`

var MultipleOrEncoderCode = `// EncodeLoginRequest returns an encoder for requests sent to the MultipleOr
// login server.
func EncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*multipleor.LoginPayload)
		if !ok {
			return goahttp.ErrInvalidType("MultipleOr", "login", "*multipleor.LoginPayload", v)
		}
		values := req.URL.Query()
		if p.Key != nil {
			values.Add("k", *p.Key)
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}
`

var SameSchemeMethod1SecureEncoderCode = `// SecureEncodeMethod1Request returns an encoder for requests sent to the
// SameSchemeMultipleEndpoints method1 endpoint that is security scheme aware.
func SecureEncodeMethod1Request(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeMethod1Request(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*sameschememultipleendpoints.Method1Payload)
		values := req.URL.Query()
		if payload.Key == nil {
			values.Set("k", "")
		} else if strings.Contains(*payload.Key, " ") {
			s := strings.SplitN(*payload.Key, " ", 2)[1]
			values.Set("k", s)
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}
`

var SameSchemeMethod2SecureEncoderCode = `// SecureEncodeMethod2Request returns an encoder for requests sent to the
// SameSchemeMultipleEndpoints method2 endpoint that is security scheme aware.
func SecureEncodeMethod2Request(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeMethod2Request(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*sameschememultipleendpoints.Method2Payload)
		if payload.Key == nil {
			req.Header.Set("Authorization", "")
		} else if strings.Contains(*payload.Key, " ") {
			s := strings.SplitN(*payload.Key, " ", 2)[1]
			req.Header.Set("Authorization", s)
		}
		return nil
	}
}
`

var SchemesInTypeSecureEncoderCode = `// SecureEncodeLoginRequest returns an encoder for requests sent to the
// SchemesInTypeDSL login endpoint that is security scheme aware.
func SecureEncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeLoginRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*schemesintypedsl.Schemes)
		values := req.URL.Query()
		if payload.Key == nil {
			values.Set("k", "")
		} else if strings.Contains(*payload.Key, " ") {
			s := strings.SplitN(*payload.Key, " ", 2)[1]
			values.Set("k", s)
		}
		if payload.Token == nil {
			req.Header.Set("Authorization", "")
		} else if !strings.Contains(*payload.Token, " ") {
			req.Header.Set("Authorization", "Bearer "+*payload.Token)
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}
`

var SchemesInTypeRequiredSecureEncoderCode = `// SecureEncodeLoginRequest returns an encoder for requests sent to the
// SchemesInTypeRequiredDSL login endpoint that is security scheme aware.
func SecureEncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeLoginRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*schemesintyperequireddsl.SchemesRequired)
		values := req.URL.Query()
		if strings.Contains(payload.Key, " ") {
			s := strings.SplitN(payload.Key, " ", 2)[1]
			values.Set("k", s)
		}
		if !strings.Contains(payload.Token, " ") {
			req.Header.Set("Authorization", "Bearer "+payload.Token)
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}
`

var SchemesInTypeEncoderCode = `// EncodeLoginRequest returns an encoder for requests sent to the
// SchemesInTypeDSL login server.
func EncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*schemesintypedsl.Schemes)
		if !ok {
			return goahttp.ErrInvalidType("SchemesInTypeDSL", "login", "*schemesintypedsl.Schemes", v)
		}
		if p.Token != nil {
			req.Header.Set("Authorization", *p.Token)
		}
		values := req.URL.Query()
		if p.Key != nil {
			values.Add("k", *p.Key)
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}
`

var MultipleSchemesWithParamsEncoderCode = `// EncodeLoginRequest returns an encoder for requests sent to the
// MultipleSchemesWithParams login server.
func EncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*multipleschemeswithparams.LoginPayload)
		if !ok {
			return goahttp.ErrInvalidType("MultipleSchemesWithParams", "login", "*multipleschemeswithparams.LoginPayload", v)
		}
		if p.UserAgent != nil {
			req.Header.Set("User-Agent", *p.UserAgent)
		}
		req.Header.Set("Authorization", p.Token)
		values := req.URL.Query()
		if p.Key != nil {
			values.Add("k", *p.Key)
		}
		if p.Name != nil {
			values.Add("name", *p.Name)
		}
		req.URL.RawQuery = values.Encode()
		body := NewLoginRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("MultipleSchemesWithParams", "login", err)
		}
		return nil
	}
}
`
