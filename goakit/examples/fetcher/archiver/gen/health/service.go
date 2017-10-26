// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// health service
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/client/archiver/design

package health

import "context"

// Service is the health service interface.
type Service interface {
	// Health check endpoint
	Show(context.Context) (string, error)
}
