// Code generated by goa v3.1.1, DO NOT EDIT.
//
// calc client
//
// Command:
// $ goa gen goa.design/plugins/v3/i18n/examples/calc/design -o
// $(GOPATH)/src/goa.design/plugins/i18n/examples/calc

package calc

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "calc" service client.
type Client struct {
	AddEndpoint goa.Endpoint
}

// NewClient initializes a "calc" service client given the endpoints.
func NewClient(add goa.Endpoint) *Client {
	return &Client{
		AddEndpoint: add,
	}
}

// Add calls the "add" endpoint of the "calc" service.
func (c *Client) Add(ctx context.Context, p *AddPayload) (res int, err error) {
	var ires interface{}
	ires, err = c.AddEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(int), nil
}