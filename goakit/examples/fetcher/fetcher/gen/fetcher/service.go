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

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "fetcher"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [1]string{"fetch"}

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

// MakeBadRequest builds a Error from an error.
func MakeBadRequest(err error) *Error {
	return &Error{
		Name:    "bad_request",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeInternalError builds a Error from an error.
func MakeInternalError(err error) *Error {
	return &Error{
		Name:    "internal_error",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}
