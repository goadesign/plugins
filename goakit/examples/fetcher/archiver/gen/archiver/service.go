// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// archiver service
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/fetcher/archiver/design

package archiver

import (
	"context"

	"goa.design/goa"
)

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
	// the HTTP status code applicable to this problem.
	Status int
	// an application-specific error code, expressed as a string value.
	Code string
	// a human-readable explanation specific to this occurrence of the problem.
	Message string
}

// Error returns "error".
func (e *Error) Error() string {
	return "error"
}

// NewNotFound initilializes a Error struct reference from a goa.Error
func NewNotFound(err goa.Error) *Error {
	return &Error{
		ID:      err.ID(),
		Status:  int(err.Status()),
		Code:    "not_found",
		Message: err.Message(),
	}
}

// NewBadRequest initilializes a Error struct reference from a goa.Error
func NewBadRequest(err goa.Error) *Error {
	return &Error{
		ID:      err.ID(),
		Status:  int(err.Status()),
		Code:    "bad_request",
		Message: err.Message(),
	}
}
