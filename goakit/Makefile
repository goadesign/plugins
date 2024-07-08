#! /usr/bin/make
#
# Makefile for goa v3 goakit plugin
#
# Targets:
# - "gen" generates the goa files for the example services

# include common Makefile content for plugins
GOPATH=$(shell go env GOPATH)
include ../plugins.mk

gen: clean
	@goa gen goa.design/plugins/v3/goakit/examples/calc/design -o "$(GOPATH)/src/goa.design/plugins/goakit/examples/calc"
	@goa gen goa.design/plugins/v3/goakit/examples/fetcher/fetcher/design -o "$(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/fetcher"
	@goa gen goa.design/plugins/v3/goakit/examples/fetcher/archiver/design -o "$(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/archiver"
	@goa example goa.design/plugins/v3/goakit/examples/calc/design -o "$(GOPATH)/src/goa.design/plugins/goakit/examples/calc"
	@goa example goa.design/plugins/v3/goakit/examples/fetcher/fetcher/design -o "$(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/fetcher"
	@goa example goa.design/plugins/v3/goakit/examples/fetcher/archiver/design -o "$(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/archiver"

build-examples:
	@cd "$(GOPATH)/src/goa.design/plugins/goakit/examples/calc" && \
		go build ./cmd/calc && go build ./cmd/calc-cli
	@cd "$(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/fetcher" && \
		go build ./cmd/fetcher && go build ./cmd/fetcher-cli
	@cd "$(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/archiver" && \
		go build ./cmd/archiver && go build ./cmd/archiver-cli

clean:
	@cd "$(GOPATH)/src/goa.design/plugins/goakit/examples/calc" && \
		rm -rf cmd calc.go
	@cd "$(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/fetcher" && \
		rm -rf cmd health.go
	@cd "$(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/archiver" && \
		rm -rf cmd health.go
