// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// sommelier go-kit HTTP server encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/cellar/design

package server

import (
	"net/http"

	goahttp "goa.design/goa/http"
)

// MountPickHandler configures the mux to serve the "sommelier" service "pick"
// endpoint.
func MountPickHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/sommelier", f)
}
