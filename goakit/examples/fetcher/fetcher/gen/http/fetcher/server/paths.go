// Code generated by goa v3.16.0, DO NOT EDIT.
//
// HTTP request path constructors for the fetcher service.
//
// Command:
// $ goa gen goa.design/plugins/v3/goakit/examples/fetcher/fetcher/design -o
// $(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/fetcher

package server

import (
	"fmt"
)

// FetchFetcherPath returns the URL path to the fetcher service fetch HTTP endpoint.
func FetchFetcherPath(url_ string) string {
	return fmt.Sprintf("/fetch/%v", url_)
}
