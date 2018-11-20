// Code generated by goa v2.0.0-wip, DO NOT EDIT.
//
// archiver views
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/fetcher/archiver/design -o
// $(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/archiver

package views

import (
	goa "goa.design/goa"
)

// ArchiveMedia is the viewed result type that is projected based on a view.
type ArchiveMedia struct {
	// Type to project
	Projected *ArchiveMediaView
	// View to render
	View string
}

// ArchiveMediaView is a type that runs validations on a projected type.
type ArchiveMediaView struct {
	// The archive resouce href
	Href *string
	// HTTP status
	Status *int
	// HTTP response body content
	Body *string
}

// Validate runs the validations defined on the viewed result type ArchiveMedia.
func (result *ArchiveMedia) Validate() (err error) {
	switch result.View {
	case "default", "":
		err = result.Projected.Validate()
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"default"})
	}
	return
}

// Validate runs the validations defined on ArchiveMediaView using the
// "default" view.
func (result *ArchiveMediaView) Validate() (err error) {
	if result.Href == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("href", "result"))
	}
	if result.Status == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("status", "result"))
	}
	if result.Body == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("body", "result"))
	}
	if result.Href != nil {
		err = goa.MergeErrors(err, goa.ValidatePattern("result.href", *result.Href, "^/archive/[0-9]+$"))
	}
	if result.Status != nil {
		if *result.Status < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("result.status", *result.Status, 0, true))
		}
	}
	return
}
