package main

import (
	"fmt"
	"net/http"

	goahttp "goa.design/goa/v3/http"
	"goa.design/plugins/v3/arnz/example"
	genarnz "goa.design/plugins/v3/arnz/example/gen/arnz"
	arnzhttp "goa.design/plugins/v3/arnz/example/gen/http/arnz/server"
)

func main() {
	fmt.Println("Starting server on :8080")
	server(8080).ListenAndServe()
}

func server(port int) *http.Server {
	mux := goahttp.NewMuxer()
	svc := &example.Service{}
	endpoints := genarnz.NewEndpoints(svc)
	api := arnzhttp.New(endpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
	api.Mount(mux)

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
}
