// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// fetcher HTTP client encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/client/fetcher/design

package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	goahttp "goa.design/goa/http"
	fetcher "goa.design/plugins/goakit/examples/client/fetcher/gen/fetcher"
)

// BuildFetchRequest instantiates a HTTP request object with method and path
// set to call the "fetcher" service "fetch" endpoint
func (c *Client) BuildFetchRequest(v interface{}) (*http.Request, error) {
	var (
		url_ string
	)
	{
		p, ok := v.(*fetcher.FetchPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("fetcher", "fetch", "*fetcher.FetchPayload", v)
		}
		url_ = p.URL
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: FetchFetcherPath(url_)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("fetcher", "fetch", u.String(), err)
	}

	return req, nil
}

// DecodeFetchResponse returns a decoder for responses returned by the fetcher
// fetch endpoint. restoreBody controls whether the response body should be
// restored after having been read.
func DecodeFetchResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
				body FetchResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("fetcher", "fetch", err)
			}
			err = body.Validate()
			if err != nil {
				return nil, fmt.Errorf("invalid response: %s", err)
			}

			return NewFetchFetchMediaOK(&body), nil
		case http.StatusBadRequest:
			var (
				body FetchBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("fetcher", "fetch", err)
			}
			err = body.Validate()
			if err != nil {
				return nil, fmt.Errorf("invalid response: %s", err)
			}

			return NewFetchBadRequest(&body), nil
		case http.StatusInternalServerError:
			var (
				body FetchInternalErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("fetcher", "fetch", err)
			}
			err = body.Validate()
			if err != nil {
				return nil, fmt.Errorf("invalid response: %s", err)
			}

			return NewFetchInternalError(&body), nil
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("account", "create", resp.StatusCode, string(body))
		}
	}
}
