#! /usr/bin/make
#
# Makefile for goa v2 goakit plugin
#
# Targets:
# - "gen" generates the goa files for the example services

PLUGIN_NAME=goakit
ALIASER_SRC=goa.design/goa/http/dsl

# include common Makefile content for plugins
include $(GOPATH)/src/goa.design/plugins/plugins.mk

gen:
	@cd examples/calc && \
	goa gen goa.design/plugins/goakit/examples/calc/design
	@cd examples/fetcher/fetcher && \
	goa gen goa.design/plugins/goakit/examples/fetcher/fetcher/design
	@cd examples/fetcher/archiver && \
	goa gen goa.design/plugins/goakit/examples/fetcher/archiver/design

aliases:
	@echo "Nothing to alias"

test-aliases: aliases
