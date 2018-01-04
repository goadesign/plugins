// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// archiver client
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/fetcher/archiver/design

package archiversvc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Client is the "archiver" service client.
type Client struct {
	ArchiveEndpoint endpoint.Endpoint
	ReadEndpoint    endpoint.Endpoint
}

// NewClient initializes a "archiver" service client given the endpoints.
func NewClient(archive, read endpoint.Endpoint) *Client {
	return &Client{
		ArchiveEndpoint: archive,
		ReadEndpoint:    read,
	}
}

// Archive calls the "archive" endpoint of the "archiver" service.
func (c *Client) Archive(ctx context.Context, p *ArchivePayload) (res *ArchiveMedia, err error) {
	var ires interface{}
	ires, err = c.ArchiveEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*ArchiveMedia), nil
}

// Read calls the "read" endpoint of the "archiver" service.
// Read can return the following error types:
//	- *Error
//	- *Error
//	- error: generic transport error.
func (c *Client) Read(ctx context.Context, p *ReadPayload) (res *ArchiveMedia, err error) {
	var ires interface{}
	ires, err = c.ReadEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*ArchiveMedia), nil
}
