// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// archiver HTTP client CLI support package
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/client/archiver/design

package client

import (
	"encoding/json"
	"fmt"
	"strconv"

	goa "goa.design/goa"
	archiver "goa.design/plugins/goakit/examples/client/archiver/gen/archiver"
)

// BuildArchiveArchivePayload builds the payload for the archiver archive
// endpoint from CLI flags.
func BuildArchiveArchivePayload(archiverArchiveBody string) (*archiver.ArchivePayload, error) {
	var err error
	var body ArchiveRequestBody
	{
		err = json.Unmarshal([]byte(archiverArchiveBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, example of valid JSON:\n%s", "'{\n      \"body\": \"Unde sed nulla.\",\n      \"status\": 200\n   }'")
		}
		if body.Status < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.status", body.Status, 0, true))
		}
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	v := &archiver.ArchivePayload{
		Status: body.Status,
		Body:   body.Body,
	}

	return v, nil
}

// BuildReadReadPayload builds the payload for the archiver read endpoint from
// CLI flags.
func BuildReadReadPayload(archiverReadID string) (*archiver.ReadPayload, error) {
	var err error
	var id int
	{
		var v int64
		v, err = strconv.ParseInt(archiverReadID, 10, 64)
		id = int(v)
		if err != nil {
			err = fmt.Errorf("invalid value for id, must be INT")
		}
		if id < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("id", id, 0, true))
		}
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	payload := &archiver.ReadPayload{
		ID: id,
	}
	return payload, nil
}
