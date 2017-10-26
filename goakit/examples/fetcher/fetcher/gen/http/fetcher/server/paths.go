// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// HTTP request path constructors for the fetcher service.
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/client/fetcher/design

package server

import (
	"fmt"
)

// FetchFetcherPath returns the URL path to the fetcher service fetch HTTP endpoint.
func FetchFetcherPath(url_ string) string {
	return fmt.Sprintf("/fetch/%v", url_)
}
