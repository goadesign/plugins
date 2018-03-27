// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// fetcher HTTP client types
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/fetcher/fetcher/design

package client

import (
	goa "goa.design/goa"
	fetchersvc "goa.design/plugins/goakit/examples/fetcher/fetcher/gen/fetcher"
)

// FetchResponseBody is the type of the "fetcher" service "fetch" endpoint HTTP
// response body.
type FetchResponseBody struct {
	// HTTP status code returned by fetched service
	Status *int `form:"status,omitempty" json:"status,omitempty" xml:"status,omitempty"`
	// The href to the corresponding archive in the archiver service
	ArchiveHref *string `form:"archive_href,omitempty" json:"archive_href,omitempty" xml:"archive_href,omitempty"`
}

// FetchBadRequestResponseBody is the type of the "fetcher" service "fetch"
// endpoint HTTP response body for the "bad_request" error.
type FetchBadRequestResponseBody struct {
	// Name is the name of this class of errors.
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// Is the error temporary?
	Temporary *bool `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"`
	// Is the error a timeout?
	Timeout *bool `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"`
}

// FetchInternalErrorResponseBody is the type of the "fetcher" service "fetch"
// endpoint HTTP response body for the "internal_error" error.
type FetchInternalErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// Is the error temporary?
	Temporary *bool `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"`
	// Is the error a timeout?
	Timeout *bool `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"`
}

// NewFetchFetchMediaOK builds a "fetcher" service "fetch" endpoint result from
// a HTTP "OK" response.
func NewFetchFetchMediaOK(body *FetchResponseBody) *fetchersvc.FetchMedia {
	v := &fetchersvc.FetchMedia{
		Status:      *body.Status,
		ArchiveHref: *body.ArchiveHref,
	}
	return v
}

// NewFetchBadRequest builds a fetcher service fetch endpoint bad_request error.
func NewFetchBadRequest(body *FetchBadRequestResponseBody) *fetchersvc.Error {
	v := &fetchersvc.Error{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
	}
	return v
}

// NewFetchInternalError builds a fetcher service fetch endpoint internal_error
// error.
func NewFetchInternalError(body *FetchInternalErrorResponseBody) *fetchersvc.Error {
	v := &fetchersvc.Error{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
	}
	return v
}

// Validate runs the validations defined on FetchResponseBody
func (body *FetchResponseBody) Validate() (err error) {
	if body.Status == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("status", "body"))
	}
	if body.ArchiveHref == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("archive_href", "body"))
	}
	if body.Status != nil {
		if *body.Status < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.status", *body.Status, 0, true))
		}
	}
	if body.ArchiveHref != nil {
		err = goa.MergeErrors(err, goa.ValidatePattern("body.archive_href", *body.ArchiveHref, "^/archive/[0-9]+$"))
	}
	return
}

// Validate runs the validations defined on FetchBadRequestResponseBody
func (body *FetchBadRequestResponseBody) Validate() (err error) {
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Message == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("message", "body"))
	}
	if body.Temporary == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("temporary", "body"))
	}
	if body.Timeout == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("timeout", "body"))
	}
	return
}

// Validate runs the validations defined on FetchInternalErrorResponseBody
func (body *FetchInternalErrorResponseBody) Validate() (err error) {
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Message == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("message", "body"))
	}
	if body.Temporary == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("temporary", "body"))
	}
	if body.Timeout == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("timeout", "body"))
	}
	return
}
