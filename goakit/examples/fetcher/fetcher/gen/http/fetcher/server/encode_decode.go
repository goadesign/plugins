// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// fetcher HTTP server encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/fetcher/fetcher/design

package server

import (
	"context"
	"net/http"

	goa "goa.design/goa"
	goahttp "goa.design/goa/http"
	fetchersvc "goa.design/plugins/goakit/examples/fetcher/fetcher/gen/fetcher"
)

// EncodeFetchResponse returns an encoder for responses returned by the fetcher
// fetch endpoint.
func EncodeFetchResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*fetchersvc.FetchMedia)
		enc := encoder(ctx, w)
		body := NewFetchResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeFetchRequest returns a decoder for requests sent to the fetcher fetch
// endpoint.
func DecodeFetchRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			url_ string
			err  error

			params = mux.Vars(r)
		)
		url_ = params["url"]
		err = goa.MergeErrors(err, goa.ValidateFormat("url_", url_, goa.FormatURI))

		if err != nil {
			return nil, err
		}

		return NewFetchFetchPayload(url_), nil
	}
}

// EncodeFetchError returns an encoder for errors returned by the fetch fetcher
// endpoint.
func EncodeFetchError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		switch res := v.(type) {
		case *fetchersvc.Error:
			if res.Name == "bad_request" {
				enc := encoder(ctx, w)
				body := NewFetchBadRequestResponseBody(res)
				w.WriteHeader(http.StatusBadRequest)
				return enc.Encode(body)
			}
			if res.Name == "internal_error" {
				enc := encoder(ctx, w)
				body := NewFetchInternalErrorResponseBody(res)
				w.WriteHeader(http.StatusInternalServerError)
				return enc.Encode(body)
			}
		default:
			return encodeError(ctx, w, v)
		}
		return nil
	}
}
