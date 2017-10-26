// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// fetcher HTTP client types
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/client/fetcher/design

package client

import (
	goa "goa.design/goa"
	fetcher "goa.design/plugins/goakit/examples/client/fetcher/gen/fetcher"
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
	// a unique identifier for this particular occurrence of the problem.
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// the HTTP status code applicable to this problem, expressed as a string value.
	Status *string `form:"status,omitempty" json:"status,omitempty" xml:"status,omitempty"`
	// an application-specific error code, expressed as a string value.
	Code *string `form:"code,omitempty" json:"code,omitempty" xml:"code,omitempty"`
	// a human-readable explanation specific to this occurrence of the problem.
	Detail *string `form:"detail,omitempty" json:"detail,omitempty" xml:"detail,omitempty"`
	// a meta object containing non-standard meta-information about the error.
	Meta map[string]interface{} `form:"meta,omitempty" json:"meta,omitempty" xml:"meta,omitempty"`
}

// FetchInternalErrorResponseBody is the type of the "fetcher" service "fetch"
// endpoint HTTP response body for the "internal_error" error.
type FetchInternalErrorResponseBody struct {
	// a unique identifier for this particular occurrence of the problem.
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// the HTTP status code applicable to this problem, expressed as a string value.
	Status *string `form:"status,omitempty" json:"status,omitempty" xml:"status,omitempty"`
	// an application-specific error code, expressed as a string value.
	Code *string `form:"code,omitempty" json:"code,omitempty" xml:"code,omitempty"`
	// a human-readable explanation specific to this occurrence of the problem.
	Detail *string `form:"detail,omitempty" json:"detail,omitempty" xml:"detail,omitempty"`
	// a meta object containing non-standard meta-information about the error.
	Meta map[string]interface{} `form:"meta,omitempty" json:"meta,omitempty" xml:"meta,omitempty"`
}

// NewFetchFetchMediaOK builds a "fetcher" service "fetch" endpoint result from
// a HTTP "OK" response.
func NewFetchFetchMediaOK(body *FetchResponseBody) *fetcher.FetchMedia {
	v := &fetcher.FetchMedia{
		Status:      *body.Status,
		ArchiveHref: *body.ArchiveHref,
	}
	return v
}

// NewFetchBadRequest builds a fetcher service fetch endpoint bad_request error.
func NewFetchBadRequest(body *FetchBadRequestResponseBody) *fetcher.Error {
	v := &fetcher.Error{
		ID:     *body.ID,
		Status: body.Status,
		Code:   *body.Code,
		Detail: *body.Detail,
	}
	v.Meta = make(map[string]interface{}, len(body.Meta))
	for key, val := range body.Meta {
		tk := key
		tv := val
		v.Meta[tk] = tv
	}
	return v
}

// NewFetchInternalError builds a fetcher service fetch endpoint internal_error
// error.
func NewFetchInternalError(body *FetchInternalErrorResponseBody) *fetcher.Error {
	v := &fetcher.Error{
		ID:     *body.ID,
		Status: body.Status,
		Code:   *body.Code,
		Detail: *body.Detail,
	}
	v.Meta = make(map[string]interface{}, len(body.Meta))
	for key, val := range body.Meta {
		tk := key
		tv := val
		v.Meta[tk] = tv
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
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Code == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("code", "body"))
	}
	if body.Detail == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("detail", "body"))
	}
	return
}

// Validate runs the validations defined on FetchInternalErrorResponseBody
func (body *FetchInternalErrorResponseBody) Validate() (err error) {
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Code == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("code", "body"))
	}
	if body.Detail == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("detail", "body"))
	}
	return
}
