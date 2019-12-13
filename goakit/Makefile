#! /usr/bin/make
#
# Makefile for goa v2 goakit plugin
#
# Targets:
# - "gen" generates the goa files for the example services

# include common Makefile content for plugins
include $(GOPATH)/src/goa.design/plugins/plugins.mk

gen:
	@goa gen goa.design/plugins/goakit/examples/calc/design -o "$(GOPATH)/src/goa.design/plugins/goakit/examples/calc"
	@goa gen goa.design/plugins/goakit/examples/fetcher/fetcher/design -o "$(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/fetcher"
	@goa gen goa.design/plugins/goakit/examples/fetcher/archiver/design -o "$(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/archiver"
	@make example

example:
	@rm -rf "$(GOPATH)/src/goa.design/plugins/goakit/examples/calc/cmd" && \
	rm -rf "$(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/archiver/cmd"
	# rm -rf "$(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/fetcher/cmd"
	# You must manually regenerate examples/fetcher/fetcher/cmd because it contains
	# changes after generation. The changes are pointed by the comment XXX.
	@goa example goa.design/plugins/goakit/examples/calc/design -o "$(GOPATH)/src/goa.design/plugins/goakit/examples/calc"
	@goa example goa.design/plugins/goakit/examples/fetcher/fetcher/design -o "$(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/fetcher"
	@goa example goa.design/plugins/goakit/examples/fetcher/archiver/design -o "$(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/archiver"

build-examples:
	@cd "$(GOPATH)/src/goa.design/plugins/goakit/examples/calc" && \
		go build ./cmd/calc && go build ./cmd/calc-cli
	@cd "$(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/fetcher" && \
		go build ./cmd/fetcher && go build ./cmd/fetcher-cli
	@cd "$(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/archiver" && \
		go build ./cmd/archiver && go build ./cmd/archiver-cli

clean:
	@cd "$(GOPATH)/src/goa.design/plugins/goakit/examples/calc" && \
		rm -f calc calc-cli
	@cd "$(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/fetcher" && \
		rm -f fetcher fetcher-cli
	@cd "$(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/archiver" && \
		rm -f archiver archiver-cli
