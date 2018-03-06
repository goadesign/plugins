// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// calc client HTTP transport
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/calc/design

package client

import (
	"context"
	"net/http"

	goa "goa.design/goa"
	goahttp "goa.design/goa/http"
)

// Client lists the calc service endpoint HTTP clients.
type Client struct {
	// Login Doer is the HTTP client used to make requests to the login endpoint.
	LoginDoer goahttp.Doer

	// Add Doer is the HTTP client used to make requests to the add endpoint.
	AddDoer goahttp.Doer

	// RestoreResponseBody controls whether the response bodies are reset after
	// decoding so they can be read again.
	RestoreResponseBody bool

	scheme  string
	host    string
	encoder func(*http.Request) goahttp.Encoder
	decoder func(*http.Response) goahttp.Decoder
}

// NewClient instantiates HTTP clients for all the calc service servers.
func NewClient(
	scheme string,
	host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restoreBody bool,
) *Client {
	return &Client{
		LoginDoer:           doer,
		AddDoer:             doer,
		RestoreResponseBody: restoreBody,
		scheme:              scheme,
		host:                host,
		decoder:             dec,
		encoder:             enc,
	}
}

// Login returns an endpoint that makes HTTP requests to the calc service login
// server.
func (c *Client) Login() goa.Endpoint {
	var (
		encodeRequest  = SecureEncodeLoginRequest(c.encoder)
		decodeResponse = DecodeLoginResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildLoginRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}

		resp, err := c.LoginDoer.Do(req)

		if err != nil {
			return nil, goahttp.ErrRequestError("calc", "login", err)
		}
		return decodeResponse(resp)
	}
}

// Add returns an endpoint that makes HTTP requests to the calc service add
// server.
func (c *Client) Add() goa.Endpoint {
	var (
		encodeRequest  = SecureEncodeAddRequest(c.encoder)
		decodeResponse = DecodeAddResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildAddRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}

		resp, err := c.AddDoer.Do(req)

		if err != nil {
			return nil, goahttp.ErrRequestError("calc", "add", err)
		}
		return decodeResponse(resp)
	}
}
