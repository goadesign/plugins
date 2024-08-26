package main

import (
	"fmt"
	"net/http"

	goahttp "goa.design/goa/v3/http"
	test "goa.design/plugins/v3/arnz/test"
	genarnz "goa.design/plugins/v3/arnz/test/gen/arnz"
	genarnzhttp "goa.design/plugins/v3/arnz/test/gen/http/arnz/server"
)

func main() {
	server(8080).ListenAndServe()
}

func server(port int) *http.Server {
	svc := &test.Service{}
	endpoints := genarnz.NewEndpoints(svc)
	mux := goahttp.NewMuxer()
	enc := goahttp.ResponseEncoder
	dec := goahttp.RequestDecoder
	api := genarnzhttp.New(endpoints, mux, dec, enc, nil, nil)

	genarnzhttp.Mount(mux, api)
	return &http.Server{Addr: fmt.Sprintf("%d", port), Handler: mux}
}
