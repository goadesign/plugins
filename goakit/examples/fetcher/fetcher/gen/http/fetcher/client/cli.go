// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// fetcher HTTP client CLI support package
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/fetcher/fetcher/design

package client

import (
	fetchersvc "goa.design/plugins/goakit/examples/fetcher/fetcher/gen/fetcher"
)

// BuildFetchFetchPayload builds the payload for the fetcher fetch endpoint
// from CLI flags.
func BuildFetchFetchPayload(fetcherFetchURL string) (*fetchersvc.FetchPayload, error) {
	var url_ string
	{
		url_ = fetcherFetchURL
	}
	payload := &fetchersvc.FetchPayload{
		URL: url_,
	}
	return payload, nil
}
