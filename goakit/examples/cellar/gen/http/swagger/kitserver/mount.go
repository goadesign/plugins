// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// swagger go-kit HTTP server encoders and decoders
//
// Command:
// $ goa gen goa.design/plugins/goakit/examples/cellar/design

package server

import (
	"net/http"

	goahttp "goa.design/goa/http"
)

// MountGenHTTPOpenapiJSON configures the mux to serve GET request made to
// "/swagger/swagger.json".
func MountGenHTTPOpenapiJSON(mux goahttp.Muxer, h http.Handler) {
	mux.Handle("GET", "/swagger/swagger.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../gen/http/openapi.json")
	}))
}
