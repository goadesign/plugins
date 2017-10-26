// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// fetcher HTTP server types
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/client/fetcher/design

package server

import (
	fetcher "goa.design/plugins/goakit/examples/client/fetcher/gen/fetcher"
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
	// a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// the HTTP status code applicable to this problem, expressed as a string value.
	Status *string `form:"status,omitempty" json:"status,omitempty" xml:"status,omitempty"`
	// an application-specific error code, expressed as a string value.
	Code string `form:"code" json:"code" xml:"code"`
	// a human-readable explanation specific to this occurrence of the problem.
	Detail string `form:"detail" json:"detail" xml:"detail"`
	// a meta object containing non-standard meta-information about the error.
	Meta map[string]interface{} `form:"meta,omitempty" json:"meta,omitempty" xml:"meta,omitempty"`
}

// FetchInternalErrorResponseBody is the type of the "fetcher" service "fetch"
// endpoint HTTP response body for the "internal_error" error.
type FetchInternalErrorResponseBody struct {
	// a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// the HTTP status code applicable to this problem, expressed as a string value.
	Status *string `form:"status,omitempty" json:"status,omitempty" xml:"status,omitempty"`
	// an application-specific error code, expressed as a string value.
	Code string `form:"code" json:"code" xml:"code"`
	// a human-readable explanation specific to this occurrence of the problem.
	Detail string `form:"detail" json:"detail" xml:"detail"`
	// a meta object containing non-standard meta-information about the error.
	Meta map[string]interface{} `form:"meta,omitempty" json:"meta,omitempty" xml:"meta,omitempty"`
}

// NewFetchResponseBody builds the HTTP response body from the result of the
// "fetch" endpoint of the "fetcher" service.
func NewFetchResponseBody(res *fetcher.FetchMedia) *FetchResponseBody {
	body := &FetchResponseBody{
		Status:      res.Status,
		ArchiveHref: res.ArchiveHref,
	}
	return body
}

// NewFetchBadRequestResponseBody builds the HTTP response body from the result
// of the "fetch" endpoint of the "fetcher" service.
func NewFetchBadRequestResponseBody(res *fetcher.Error) *FetchBadRequestResponseBody {
	body := &FetchBadRequestResponseBody{
		ID:     res.ID,
		Status: res.Status,
		Code:   res.Code,
		Detail: res.Detail,
	}
	if res.Meta != nil {
		body.Meta = make(map[string]interface{}, len(res.Meta))
		for key, val := range res.Meta {
			tk := key
			tv := val
			body.Meta[tk] = tv
		}
	}
	return body
}

// NewFetchInternalErrorResponseBody builds the HTTP response body from the
// result of the "fetch" endpoint of the "fetcher" service.
func NewFetchInternalErrorResponseBody(res *fetcher.Error) *FetchInternalErrorResponseBody {
	body := &FetchInternalErrorResponseBody{
		ID:     res.ID,
		Status: res.Status,
		Code:   res.Code,
		Detail: res.Detail,
	}
	if res.Meta != nil {
		body.Meta = make(map[string]interface{}, len(res.Meta))
		for key, val := range res.Meta {
			tk := key
			tv := val
			body.Meta[tk] = tv
		}
	}
	return body
}

// NewFetchFetchPayload builds a fetcher service fetch endpoint payload.
func NewFetchFetchPayload(url_ string) *fetcher.FetchPayload {
	return &fetcher.FetchPayload{
		URL: url_,
	}
}
