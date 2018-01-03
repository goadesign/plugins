// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// secured_service HTTP client encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/security/examples/multi_auth/design

package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	goa "goa.design/goa"
	goahttp "goa.design/goa/http"
	securedservice "goa.design/plugins/security/examples/multi_auth/gen/secured_service"
)

// BuildSigninRequest instantiates a HTTP request object with method and path
// set to call the "secured_service" service "signin" endpoint
func (c *Client) BuildSigninRequest(v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: SigninSecuredServicePath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("secured_service", "signin", u.String(), err)
	}

	return req, nil
}

// EncodeSigninRequest returns an encoder for requests sent to the
// secured_service signin server.
func EncodeSigninRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*securedservice.SigninPayload)
		if !ok {
			return goahttp.ErrInvalidType("secured_service", "signin", "*securedservice.SigninPayload", v)
		}
		body := NewSigninRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("secured_service", "signin", err)
		}
		return nil
	}
}

// DecodeSigninResponse returns a decoder for responses returned by the
// secured_service signin endpoint. restoreBody controls whether the response
// body should be restored after having been read.
// DecodeSigninResponse may return the following error types:
//	- *securedservice.Unauthorized: http.StatusUnauthorized
//	- error: generic transport error.
func DecodeSigninResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusNoContent:
			var (
				authorization string
				err           error
			)
			authorization = resp.Header.Get("Authorization")
			if authorization != "" {
				err = goa.MergeErrors(err, goa.MissingFieldError("Authorization", "header"))
			}
			if err != nil {
				return nil, fmt.Errorf("invalid response: %s", err)
			}
			return nil, nil
		case http.StatusUnauthorized:
			var (
				body SigninUnauthorizedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("secured_service", "signin", err)
			}
			err = body.Validate()
			if err != nil {
				return nil, fmt.Errorf("invalid response: %s", err)
			}

			return nil, NewSigninUnauthorized(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("account", "create", resp.StatusCode, string(body))
		}
	}
}

// BuildSecureRequest instantiates a HTTP request object with method and path
// set to call the "secured_service" service "secure" endpoint
func (c *Client) BuildSecureRequest(v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: SecureSecuredServicePath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("secured_service", "secure", u.String(), err)
	}

	return req, nil
}

// EncodeSecureRequest returns an encoder for requests sent to the
// secured_service secure server.
func EncodeSecureRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*securedservice.SecurePayload)
		if !ok {
			return goahttp.ErrInvalidType("secured_service", "secure", "*securedservice.SecurePayload", v)
		}
		values := req.URL.Query()
		if p.Fail != nil {
			values.Add("fail", fmt.Sprintf("%v", *p.Fail))
		}
		req.URL.RawQuery = values.Encode()
		body := NewSecureRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("secured_service", "secure", err)
		}
		return nil
	}
}

// DecodeSecureResponse returns a decoder for responses returned by the
// secured_service secure endpoint. restoreBody controls whether the response
// body should be restored after having been read.
func DecodeSecureResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body string
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("secured_service", "secure", err)
			}

			return body, nil
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("account", "create", resp.StatusCode, string(body))
		}
	}
}

// BuildDoublySecureRequest instantiates a HTTP request object with method and
// path set to call the "secured_service" service "doubly_secure" endpoint
func (c *Client) BuildDoublySecureRequest(v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: DoublySecureSecuredServicePath()}
	req, err := http.NewRequest("PUT", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("secured_service", "doubly_secure", u.String(), err)
	}

	return req, nil
}

// EncodeDoublySecureRequest returns an encoder for requests sent to the
// secured_service doubly_secure server.
func EncodeDoublySecureRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*securedservice.DoublySecurePayload)
		if !ok {
			return goahttp.ErrInvalidType("secured_service", "doubly_secure", "*securedservice.DoublySecurePayload", v)
		}
		values := req.URL.Query()
		if p.Key != nil {
			values.Add("k", *p.Key)
		}
		req.URL.RawQuery = values.Encode()
		body := NewDoublySecureRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("secured_service", "doubly_secure", err)
		}
		return nil
	}
}

// DecodeDoublySecureResponse returns a decoder for responses returned by the
// secured_service doubly_secure endpoint. restoreBody controls whether the
// response body should be restored after having been read.
func DecodeDoublySecureResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body string
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("secured_service", "doubly_secure", err)
			}

			return body, nil
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("account", "create", resp.StatusCode, string(body))
		}
	}
}

// BuildAlsoDoublySecureRequest instantiates a HTTP request object with method
// and path set to call the "secured_service" service "also_doubly_secure"
// endpoint
func (c *Client) BuildAlsoDoublySecureRequest(v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: AlsoDoublySecureSecuredServicePath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("secured_service", "also_doubly_secure", u.String(), err)
	}

	return req, nil
}

// EncodeAlsoDoublySecureRequest returns an encoder for requests sent to the
// secured_service also_doubly_secure server.
func EncodeAlsoDoublySecureRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*securedservice.AlsoDoublySecurePayload)
		if !ok {
			return goahttp.ErrInvalidType("secured_service", "also_doubly_secure", "*securedservice.AlsoDoublySecurePayload", v)
		}
		if p.Key != nil {
			req.Header.Set("Authorization", *p.Key)
		}
		body := NewAlsoDoublySecureRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("secured_service", "also_doubly_secure", err)
		}
		return nil
	}
}

// DecodeAlsoDoublySecureResponse returns a decoder for responses returned by
// the secured_service also_doubly_secure endpoint. restoreBody controls
// whether the response body should be restored after having been read.
func DecodeAlsoDoublySecureResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body string
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("secured_service", "also_doubly_secure", err)
			}

			return body, nil
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("account", "create", resp.StatusCode, string(body))
		}
	}
}

// SecureEncodeSigninRequest returns an encoder for requests sent to the
// secured_service signin endpoint that is security scheme aware.
func SecureEncodeSigninRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeSigninRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*securedservice.SigninPayload)
		req.SetBasicAuth(*payload.Username, *payload.Password)
		return nil
	}
}

// SecureEncodeSecureRequest returns an encoder for requests sent to the
// secured_service secure endpoint that is security scheme aware.
func SecureEncodeSecureRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeSecureRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*securedservice.SecurePayload)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *payload.Token))
		return nil
	}
}

// SecureEncodeDoublySecureRequest returns an encoder for requests sent to the
// secured_service doubly_secure endpoint that is security scheme aware.
func SecureEncodeDoublySecureRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeDoublySecureRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*securedservice.DoublySecurePayload)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *payload.Token))
		req.URL.Query().Set("k", *payload.Key)
		return nil
	}
}

// SecureEncodeAlsoDoublySecureRequest returns an encoder for requests sent to
// the secured_service also_doubly_secure endpoint that is security scheme
// aware.
func SecureEncodeAlsoDoublySecureRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeAlsoDoublySecureRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*securedservice.AlsoDoublySecurePayload)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *payload.Token))
		req.Header.Set("Authorization", *payload.Key)
		return nil
	}
}
