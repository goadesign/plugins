#! /usr/bin/make
#
# Makefile for goa v2 CORS plugin
#
# Targets:
# - "gen" generates the goa files for the example services

PLUGIN_NAME=cors
ALIASER_SRC=goa.design/goa/http/dsl

# include common Makefile content for plugins
include $(GOPATH)/src/goa.design/plugins/plugins.mk

gen:
	@cd examples/calc && \
	goa gen goa.design/plugins/cors/examples/calc/design
