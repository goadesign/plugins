// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// adder HTTP client encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/adder/design

package client

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"

	goahttp "goa.design/goa/http"
	addersvc "goa.design/plugins/security/examples/calc/adder/gen/adder"
)

// BuildAddRequest instantiates a HTTP request object with method and path set
// to call the "adder" service "add" endpoint
func (c *Client) BuildAddRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		a int
		b int
	)
	{
		p, ok := v.(*addersvc.AddPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("adder", "add", "*addersvc.AddPayload", v)
		}
		a = p.A
		b = p.B
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: AddAdderPath(a, b)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("adder", "add", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeAddRequest returns an encoder for requests sent to the adder add
// server.
func EncodeAddRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*addersvc.AddPayload)
		if !ok {
			return goahttp.ErrInvalidType("adder", "add", "*addersvc.AddPayload", v)
		}
		values := req.URL.Query()
		values.Add("key", p.Key)
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeAddResponse returns a decoder for responses returned by the adder add
// endpoint. restoreBody controls whether the response body should be restored
// after having been read.
// DecodeAddResponse may return the following error types:
//	- addersvc.InvalidScopes: http.StatusForbidden
//	- addersvc.Unauthorized: http.StatusUnauthorized
//	- error: generic transport error.
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
				return nil, goahttp.ErrDecodingError("adder", "add", err)
			}

			return body, nil
		case http.StatusForbidden:
			var (
				body AddInvalidScopesResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("adder", "add", err)
			}

			return nil, NewAddInvalidScopes(body)
		case http.StatusUnauthorized:
			var (
				body AddUnauthorizedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("adder", "add", err)
			}

			return nil, NewAddUnauthorized(body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("account", "create", resp.StatusCode, string(body))
		}
	}
}

// SecureEncodeAddRequest returns an encoder for requests sent to the adder add
// endpoint that is security scheme aware.
func SecureEncodeAddRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	rawEncoder := EncodeAddRequest(encoder)
	return func(req *http.Request, v interface{}) error {
		if err := rawEncoder(req, v); err != nil {
			return err
		}
		payload := v.(*addersvc.AddPayload)
		req.URL.Query().Set("key", payload.Key)
		return nil
	}
}
