// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// archiver HTTP server encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/fetcher/archiver/design

package server

import (
	"context"
	"io"
	"net/http"
	"strconv"

	goa "goa.design/goa"
	goahttp "goa.design/goa/http"
	archiversvc "goa.design/plugins/goakit/examples/fetcher/archiver/gen/archiver"
)

// EncodeArchiveResponse returns an encoder for responses returned by the
// archiver archive endpoint.
func EncodeArchiveResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*archiversvc.ArchiveMedia)
		enc := encoder(ctx, w)
		body := NewArchiveResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeArchiveRequest returns a decoder for requests sent to the archiver
// archive endpoint.
func DecodeArchiveRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			body ArchiveRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}
		err = body.Validate()
		if err != nil {
			return nil, err
		}

		return NewArchiveArchivePayload(&body), nil
	}
}

// EncodeReadResponse returns an encoder for responses returned by the archiver
// read endpoint.
func EncodeReadResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*archiversvc.ArchiveMedia)
		enc := encoder(ctx, w)
		body := NewReadResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeReadRequest returns a decoder for requests sent to the archiver read
// endpoint.
func DecodeReadRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			id  int
			err error

			params = mux.Vars(r)
		)
		{
			idRaw := params["id"]
			v, err2 := strconv.ParseInt(idRaw, 10, strconv.IntSize)
			if err2 != nil {
				err = goa.MergeErrors(err, goa.InvalidFieldTypeError("id", idRaw, "integer"))
			}
			id = int(v)
		}
		if id < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("id", id, 0, true))
		}
		if err != nil {
			return nil, err
		}

		return NewReadReadPayload(id), nil
	}
}

// EncodeReadError returns an encoder for errors returned by the read archiver
// endpoint.
func EncodeReadError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		switch res := v.(type) {
		case *archiversvc.Error:
			if res.Name == "not_found" {
				enc := encoder(ctx, w)
				body := NewReadNotFoundResponseBody(res)
				w.WriteHeader(http.StatusNotFound)
				return enc.Encode(body)
			}
			if res.Name == "bad_request" {
				enc := encoder(ctx, w)
				body := NewReadBadRequestResponseBody(res)
				w.WriteHeader(http.StatusBadRequest)
				return enc.Encode(body)
			}
		default:
			return encodeError(ctx, w, v)
		}
		return nil
	}
}
