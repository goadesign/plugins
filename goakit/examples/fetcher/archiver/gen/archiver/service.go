// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// archiver service
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/fetcher/archiver/design

package archiversvc

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

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "archiver"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [2]string{"archive", "read"}

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
	// Name is the name of this class of errors.
	Name string
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string
	// Is the error temporary?
	Temporary bool
	// Is the error a timeout?
	Timeout bool
}

// Error returns "error".
func (e *Error) Error() string {
	return "error"
}

// MakeNotFound builds a Error from an error.
func MakeNotFound(err error) *Error {
	return &Error{
		Name:    "not_found",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeBadRequest builds a Error from an error.
func MakeBadRequest(err error) *Error {
	return &Error{
		Name:    "bad_request",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}
