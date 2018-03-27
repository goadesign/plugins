// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// archiver HTTP client types
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/fetcher/archiver/design

package client

import (
	goa "goa.design/goa"
	archiversvc "goa.design/plugins/goakit/examples/fetcher/archiver/gen/archiver"
)

// ArchiveRequestBody is the type of the "archiver" service "archive" endpoint
// HTTP request body.
type ArchiveRequestBody struct {
	// HTTP status
	Status int `form:"status" json:"status" xml:"status"`
	// HTTP response body content
	Body string `form:"body" json:"body" xml:"body"`
}

// ArchiveResponseBody is the type of the "archiver" service "archive" endpoint
// HTTP response body.
type ArchiveResponseBody struct {
	// The archive resouce href
	Href *string `form:"href,omitempty" json:"href,omitempty" xml:"href,omitempty"`
	// HTTP status
	Status *int `form:"status,omitempty" json:"status,omitempty" xml:"status,omitempty"`
	// HTTP response body content
	Body *string `form:"body,omitempty" json:"body,omitempty" xml:"body,omitempty"`
}

// ReadResponseBody is the type of the "archiver" service "read" endpoint HTTP
// response body.
type ReadResponseBody struct {
	// The archive resouce href
	Href *string `form:"href,omitempty" json:"href,omitempty" xml:"href,omitempty"`
	// HTTP status
	Status *int `form:"status,omitempty" json:"status,omitempty" xml:"status,omitempty"`
	// HTTP response body content
	Body *string `form:"body,omitempty" json:"body,omitempty" xml:"body,omitempty"`
}

// ReadNotFoundResponseBody is the type of the "archiver" service "read"
// endpoint HTTP response body for the "not_found" error.
type ReadNotFoundResponseBody struct {
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

// ReadBadRequestResponseBody is the type of the "archiver" service "read"
// endpoint HTTP response body for the "bad_request" error.
type ReadBadRequestResponseBody struct {
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

// NewArchiveRequestBody builds the HTTP request body from the payload of the
// "archive" endpoint of the "archiver" service.
func NewArchiveRequestBody(p *archiversvc.ArchivePayload) *ArchiveRequestBody {
	body := &ArchiveRequestBody{
		Status: p.Status,
		Body:   p.Body,
	}
	return body
}

// NewArchiveArchiveMediaOK builds a "archiver" service "archive" endpoint
// result from a HTTP "OK" response.
func NewArchiveArchiveMediaOK(body *ArchiveResponseBody) *archiversvc.ArchiveMedia {
	v := &archiversvc.ArchiveMedia{
		Href:   *body.Href,
		Status: *body.Status,
		Body:   *body.Body,
	}
	return v
}

// NewReadArchiveMediaOK builds a "archiver" service "read" endpoint result
// from a HTTP "OK" response.
func NewReadArchiveMediaOK(body *ReadResponseBody) *archiversvc.ArchiveMedia {
	v := &archiversvc.ArchiveMedia{
		Href:   *body.Href,
		Status: *body.Status,
		Body:   *body.Body,
	}
	return v
}

// NewReadNotFound builds a archiver service read endpoint not_found error.
func NewReadNotFound(body *ReadNotFoundResponseBody) *archiversvc.Error {
	v := &archiversvc.Error{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
	}
	return v
}

// NewReadBadRequest builds a archiver service read endpoint bad_request error.
func NewReadBadRequest(body *ReadBadRequestResponseBody) *archiversvc.Error {
	v := &archiversvc.Error{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
	}
	return v
}

// Validate runs the validations defined on ArchiveResponseBody
func (body *ArchiveResponseBody) Validate() (err error) {
	if body.Href == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("href", "body"))
	}
	if body.Status == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("status", "body"))
	}
	if body.Body == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("body", "body"))
	}
	if body.Href != nil {
		err = goa.MergeErrors(err, goa.ValidatePattern("body.href", *body.Href, "^/archive/[0-9]+$"))
	}
	if body.Status != nil {
		if *body.Status < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.status", *body.Status, 0, true))
		}
	}
	return
}

// Validate runs the validations defined on ReadResponseBody
func (body *ReadResponseBody) Validate() (err error) {
	if body.Href == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("href", "body"))
	}
	if body.Status == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("status", "body"))
	}
	if body.Body == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("body", "body"))
	}
	if body.Href != nil {
		err = goa.MergeErrors(err, goa.ValidatePattern("body.href", *body.Href, "^/archive/[0-9]+$"))
	}
	if body.Status != nil {
		if *body.Status < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.status", *body.Status, 0, true))
		}
	}
	return
}

// Validate runs the validations defined on ReadNotFoundResponseBody
func (body *ReadNotFoundResponseBody) Validate() (err error) {
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

// Validate runs the validations defined on ReadBadRequestResponseBody
func (body *ReadBadRequestResponseBody) Validate() (err error) {
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
