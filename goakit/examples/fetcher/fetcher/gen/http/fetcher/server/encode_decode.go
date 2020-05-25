// Code generated by goa v3.1.3, DO NOT EDIT.
//
// fetcher HTTP server encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/v3/goakit/examples/fetcher/fetcher/design -o
// $(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/fetcher

package server

import (
	"context"
	"net/http"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
	fetcherviews "goa.design/plugins/v3/goakit/examples/fetcher/fetcher/gen/fetcher/views"
)

// EncodeFetchResponse returns an encoder for responses returned by the fetcher
// fetch endpoint.
func EncodeFetchResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*fetcherviews.FetchMedia)
		enc := encoder(ctx, w)
		body := NewFetchResponseBody(res.Projected)
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
		payload := NewFetchPayload(url_)

		return payload, nil
	}
}

// EncodeFetchError returns an encoder for errors returned by the fetch fetcher
// endpoint.
func EncodeFetchError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		en, ok := v.(ErrorNamer)
		if !ok {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "bad_request":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewFetchBadRequestResponseBody(res)
			}
			w.Header().Set("goa-error", "bad_request")
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		case "internal_error":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewFetchInternalErrorResponseBody(res)
			}
			w.Header().Set("goa-error", "internal_error")
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}
