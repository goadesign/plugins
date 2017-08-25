// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// storage go-kit HTTP server encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/cellar/design

package server

import (
	"net/http"

	goahttp "goa.design/goa/http"
)

// MountListHandler configures the mux to serve the "storage" service "list"
// endpoint.
func MountListHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/storage", f)
}

// MountShowHandler configures the mux to serve the "storage" service "show"
// endpoint.
func MountShowHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/storage/{id}", f)
}

// MountAddHandler configures the mux to serve the "storage" service "add"
// endpoint.
func MountAddHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/storage", f)
}

// MountRemoveHandler configures the mux to serve the "storage" service
// "remove" endpoint.
func MountRemoveHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("DELETE", "/storage/{id}", f)
}
