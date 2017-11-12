// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// fetcher HTTP client CLI support package
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/fetcher/fetcher/design

package client

import (
	fetcher "goa.design/plugins/goakit/examples/fetcher/fetcher/gen/fetcher"
)

// BuildFetchFetchPayload builds the payload for the fetcher fetch endpoint
// from CLI flags.
func BuildFetchFetchPayload(fetcherFetchURL string) (*fetcher.FetchPayload, error) {
	var url_ string
	{
		url_ = fetcherFetchURL
	}
	payload := &fetcher.FetchPayload{
		URL: url_,
	}
	return payload, nil
}
