package main

import (
	"fmt"
	"net/http"
	"strconv"

	goahttp "goa.design/goa/v3/http"
	"goa.design/plugins/v3/arnz/example"
	likehttp "goa.design/plugins/v3/arnz/example/gen/http/like/server"
	matchhttp "goa.design/plugins/v3/arnz/example/gen/http/match/server"
	"goa.design/plugins/v3/arnz/example/gen/like"
	"goa.design/plugins/v3/arnz/example/gen/match"
)

func main() {
	fmt.Println("Starting server on :8080")
	server(8080).ListenAndServe()
}

func server(port int) *http.Server {
	mux := goahttp.NewMuxer()

	crud := &example.Service{}

	likeEndpoints := like.NewEndpoints(crud)
	matchEndpoints := match.NewEndpoints(crud)

	likeApi := likehttp.New(likeEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
	matchApi := matchhttp.New(matchEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)

	likeApi.Mount(mux)
	matchApi.Mount(mux)

	return &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}
}
