// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// HTTP request path constructors for the archiver service.
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/client/archiver/design

package server

import (
	"fmt"
)

// ArchiveArchiverPath returns the URL path to the archiver service archive HTTP endpoint.
func ArchiveArchiverPath() string {
	return "/archive"
}

// ReadArchiverPath returns the URL path to the archiver service read HTTP endpoint.
func ReadArchiverPath(id int) string {
	return fmt.Sprintf("/archive/%v", id)
}
