// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// calc HTTP client encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/calc/design

package client

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	goahttp "goa.design/goa/http"
	calcsvc "goa.design/plugins/security/examples/calc/calc/gen/calc"
)

// BuildLoginRequest instantiates a HTTP request object with method and path
// set to call the "calc" service "login" endpoint
func (c *Client) BuildLoginRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: LoginCalcPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("calc", "login", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeLoginResponse returns a decoder for responses returned by the calc
// login endpoint. restoreBody controls whether the response body should be
// restored after having been read.
// DecodeLoginResponse may return the following error types:
//	- calcsvc.Unauthorized: http.StatusUnauthorized
//	- error: generic transport error.
func DecodeLoginResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
				return nil, goahttp.ErrDecodingError("calc", "login", err)
			}

			return body, nil
		case http.StatusUnauthorized:
			var (
				body LoginUnauthorizedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("calc", "login", err)
			}

			return nil, NewLoginUnauthorized(body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("account", "create", resp.StatusCode, string(body))
		}
	}
}

// BuildAddRequest instantiates a HTTP request object with method and path set
// to call the "calc" service "add" endpoint
func (c *Client) BuildAddRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		a int
		b int
	)
	{
		p, ok := v.(*calcsvc.AddPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("calc", "add", "*calcsvc.AddPayload", v)
		}
		a = p.A
		b = p.B
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: AddCalcPath(a, b)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("calc", "add", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeAddRequest returns an encoder for requests sent to the calc add server.
func EncodeAddRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*calcsvc.AddPayload)
		if !ok {
			return goahttp.ErrInvalidType("calc", "add", "*calcsvc.AddPayload", v)
		}
		req.Header.Set("Authorization", p.Token)
		return nil
	}
}

// DecodeAddResponse returns a decoder for responses returned by the calc add
// endpoint. restoreBody controls whether the response body should be restored
// after having been read.
func DecodeAddResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
				body int
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("calc", "add", err)
			}

			return body, nil
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("account", "create", resp.StatusCode, string(body))
		}
	}
}

// SecureEncodeLoginRequest returns an encoder for requests sent to the calc
// login endpoint that is security scheme aware.
func SecureEncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		payload := v.(*calcsvc.LoginPayload)
		req.SetBasicAuth(payload.User, payload.Password)
		return nil
	}
}

// SecureEncodeAddRequest returns an encoder for requests sent to the calc add
// endpoint that is security scheme aware.
func SecureEncodeAddRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeAddRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*calcsvc.AddPayload)
		if !strings.Contains(payload.Token, " ") {
			req.Header.Set("Authorization", "Bearer "+payload.Token)
		}
		return nil
	}
}
