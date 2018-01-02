// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// adder HTTP server encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/adder/design

package server

import (
	"context"
	"net/http"
	"strconv"

	goa "goa.design/goa"
	goahttp "goa.design/goa/http"
	addersvc "goa.design/plugins/security/examples/calc/adder/gen/adder"
)

// EncodeAddResponse returns an encoder for responses returned by the adder add
// endpoint.
func EncodeAddResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(int)
		enc := encoder(ctx, w)
		body := res
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeAddRequest returns a decoder for requests sent to the adder add
// endpoint.
func DecodeAddRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			a   int
			b   int
			key string
			err error

			params = mux.Vars(r)
		)
		{
			aRaw := params["a"]
			v, err2 := strconv.ParseInt(aRaw, 10, strconv.IntSize)
			if err2 != nil {
				err = goa.MergeErrors(err, goa.InvalidFieldTypeError("a", aRaw, "integer"))
			}
			a = int(v)
		}
		{
			bRaw := params["b"]
			v, err2 := strconv.ParseInt(bRaw, 10, strconv.IntSize)
			if err2 != nil {
				err = goa.MergeErrors(err, goa.InvalidFieldTypeError("b", bRaw, "integer"))
			}
			b = int(v)
		}
		key = r.URL.Query().Get("key")
		if key == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("key", "query string"))
		}
		if err != nil {
			return nil, err
		}

		return NewAddAddPayload(a, b, key), nil
	}
}

// EncodeAddError returns an encoder for errors returned by the add adder
// endpoint.
func EncodeAddError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, error) {
	encodeError := goahttp.ErrorEncoder(encoder)
	return func(ctx context.Context, w http.ResponseWriter, v error) {
		switch res := v.(type) {
		case *addersvc.InvalidScopes:
			enc := encoder(ctx, w)
			body := NewAddInvalidScopesResponseBody(res)
			w.WriteHeader(http.StatusForbidden)
			if err := enc.Encode(body); err != nil {
				encodeError(ctx, w, err)
			}
		case *addersvc.Unauthorized:
			enc := encoder(ctx, w)
			body := NewAddUnauthorizedResponseBody(res)
			w.WriteHeader(http.StatusUnauthorized)
			if err := enc.Encode(body); err != nil {
				encodeError(ctx, w, err)
			}
		default:
			encodeError(ctx, w, v)
		}
	}
}
