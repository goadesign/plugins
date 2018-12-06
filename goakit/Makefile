#! /usr/bin/make
#
# Makefile for goa v2 goakit plugin
#
# Targets:
# - "gen" generates the goa files for the example services

# include common Makefile content for plugins
include $(GOPATH)/src/goa.design/plugins/plugins.mk

gen:
	@goa gen goa.design/plugins/goakit/examples/calc/design -o $(GOPATH)/src/goa.design/plugins/goakit/examples/calc
	@goa gen goa.design/plugins/goakit/examples/fetcher/fetcher/design -o $(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/fetcher
	@goa gen goa.design/plugins/goakit/examples/fetcher/archiver/design -o $(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/archiver

example:
	@goa example goa.design/plugins/goakit/examples/calc/design -o $(GOPATH)/src/goa.design/plugins/goakit/examples/calc
	@goa example goa.design/plugins/goakit/examples/fetcher/fetcher/design -o $(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/fetcher
	@goa example goa.design/plugins/goakit/examples/fetcher/archiver/design -o $(GOPATH)/src/goa.design/plugins/goakit/examples/fetcher/archiver
