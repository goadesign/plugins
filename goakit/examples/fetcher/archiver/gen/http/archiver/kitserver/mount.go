// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// archiver go-kit HTTP server encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/client/archiver/design

package server

import (
	"net/http"

	goahttp "goa.design/goa/http"
)

// MountArchiveHandler configures the mux to serve the "archiver" service
// "archive" endpoint.
func MountArchiveHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/archive", f)
}

// MountReadHandler configures the mux to serve the "archiver" service "read"
// endpoint.
func MountReadHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/archive/{id}", f)
}
