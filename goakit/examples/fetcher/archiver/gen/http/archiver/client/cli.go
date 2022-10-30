// Code generated by goa v3.10.2, DO NOT EDIT.
//
// archiver HTTP client CLI support package
//
// Command:
// $ goa gen goa.design/plugins/v3/goakit/examples/fetcher/archiver/design -o
// $(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/archiver

package client

import (
	"encoding/json"
	"fmt"
	"strconv"

	goa "goa.design/goa/v3/pkg"
	archiver "goa.design/plugins/v3/goakit/examples/fetcher/archiver/gen/archiver"
)

// BuildArchivePayload builds the payload for the archiver archive endpoint
// from CLI flags.
func BuildArchivePayload(archiverArchiveBody string) (*archiver.ArchivePayload, error) {
	var err error
	var body ArchiveRequestBody
	{
		err = json.Unmarshal([]byte(archiverArchiveBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"body\": \"Unde sed nulla.\",\n      \"status\": 200\n   }'")
		}
		if body.Status < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.status", body.Status, 0, true))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &archiver.ArchivePayload{
		Status: body.Status,
		Body:   body.Body,
	}

	return v, nil
}

// BuildReadPayload builds the payload for the archiver read endpoint from CLI
// flags.
func BuildReadPayload(archiverReadID string) (*archiver.ReadPayload, error) {
	var err error
	var id int
	{
		var v int64
		v, err = strconv.ParseInt(archiverReadID, 10, strconv.IntSize)
		id = int(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for id, must be INT")
		}
		if id < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("id", id, 0, true))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &archiver.ReadPayload{}
	v.ID = id

	return v, nil
}
