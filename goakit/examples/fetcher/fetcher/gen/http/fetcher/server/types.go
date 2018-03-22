// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// fetcher HTTP server types
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/fetcher/fetcher/design

package server

import (
	fetchersvc "goa.design/plugins/goakit/examples/fetcher/fetcher/gen/fetcher"
)

// FetchResponseBody is the type of the "fetcher" service "fetch" endpoint HTTP
// response body.
type FetchResponseBody struct {
	// HTTP status code returned by fetched service
	Status int `form:"status" json:"status" xml:"status"`
	// The href to the corresponding archive in the archiver service
	ArchiveHref string `form:"archive_href" json:"archive_href" xml:"archive_href"`
}

// FetchBadRequestResponseBody is the type of the "fetcher" service "fetch"
// endpoint HTTP response body for the "bad_request" error.
type FetchBadRequestResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary *bool `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"`
	// Is the error a timeout?
	Timeout *bool `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"`
}

// FetchInternalErrorResponseBody is the type of the "fetcher" service "fetch"
// endpoint HTTP response body for the "internal_error" error.
type FetchInternalErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary *bool `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"`
	// Is the error a timeout?
	Timeout *bool `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"`
}

// NewFetchResponseBody builds the HTTP response body from the result of the
// "fetch" endpoint of the "fetcher" service.
func NewFetchResponseBody(res *fetchersvc.FetchMedia) *FetchResponseBody {
	body := &FetchResponseBody{
		Status:      res.Status,
		ArchiveHref: res.ArchiveHref,
	}
	return body
}

// NewFetchBadRequestResponseBody builds the HTTP response body from the result
// of the "fetch" endpoint of the "fetcher" service.
func NewFetchBadRequestResponseBody(res *fetchersvc.Error) *FetchBadRequestResponseBody {
	body := &FetchBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
	}
	return body
}

// NewFetchInternalErrorResponseBody builds the HTTP response body from the
// result of the "fetch" endpoint of the "fetcher" service.
func NewFetchInternalErrorResponseBody(res *fetchersvc.Error) *FetchInternalErrorResponseBody {
	body := &FetchInternalErrorResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
	}
	return body
}

// NewFetchFetchPayload builds a fetcher service fetch endpoint payload.
func NewFetchFetchPayload(url_ string) *fetchersvc.FetchPayload {
	return &fetchersvc.FetchPayload{
		URL: url_,
	}
}
