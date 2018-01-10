// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// fetcher service
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/fetcher/fetcher/design

package fetchersvc

import (
	"context"

	"goa.design/goa"
)

// Service is the fetcher service interface.
type Service interface {
	// Fetch makes a GET request to the given URL and stores the results in the
	// archiver service which must be running or the request fails
	Fetch(context.Context, *FetchPayload) (*FetchMedia, error)
}

// FetchPayload is the payload type of the fetcher service fetch method.
type FetchPayload struct {
	// URL to be fetched
	URL string
}

// FetchMedia is the result type of the fetcher service fetch method.
type FetchMedia struct {
	// HTTP status code returned by fetched service
	Status int
	// The href to the corresponding archive in the archiver service
	ArchiveHref string
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

// NewBadRequest initilializes a Error struct reference from a goa.Error
func NewBadRequest(err goa.Error) *Error {
	return &Error{
		ID:      err.ID(),
		Status:  int(err.Status()),
		Code:    "bad_request",
		Message: err.Message(),
	}
}

// NewInternalError initilializes a Error struct reference from a goa.Error
func NewInternalError(err goa.Error) *Error {
	return &Error{
		ID:      err.ID(),
		Status:  int(err.Status()),
		Code:    "internal_error",
		Message: err.Message(),
	}
}
