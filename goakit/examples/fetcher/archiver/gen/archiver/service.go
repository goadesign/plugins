// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// archiver service
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/client/archiver/design

package archiver

import "context"

// Service is the archiver service interface.
type Service interface {
	// Archive HTTP response
	Archive(context.Context, *ArchivePayload) (*ArchiveMedia, error)
	// Read HTTP response from archive
	Read(context.Context, *ReadPayload) (*ArchiveMedia, error)
}

// ArchivePayload is the payload type of the archiver service archive method.
type ArchivePayload struct {
	// HTTP status
	Status int
	// HTTP response body content
	Body string
}

// ArchiveMedia is the result type of the archiver service archive method.
type ArchiveMedia struct {
	// The archive resouce href
	Href string
	// HTTP status
	Status int
	// HTTP response body content
	Body string
}

// ReadPayload is the payload type of the archiver service read method.
type ReadPayload struct {
	// ID of archive
	ID int
}

// Error response result type
type Error struct {
	// a unique identifier for this particular occurrence of the problem.
	ID string
	// the HTTP status code applicable to this problem, expressed as a string value.
	Status *string
	// an application-specific error code, expressed as a string value.
	Code string
	// a human-readable explanation specific to this occurrence of the problem.
	Detail string
	// a meta object containing non-standard meta-information about the error.
	Meta map[string]interface{}
}

// Error returns "error".
func (e *Error) Error() string {
	return "error"
}
