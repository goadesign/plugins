#! /usr/bin/make
#
# Makefile for goa v2 CORS plugin
#
# Targets:
# - "gen" generates the goa files for the example services
# - "example" generates the example files for the example services

# include common Makefile content for plugins
include $(GOPATH)/src/goa.design/plugins/plugins.mk

gen:
	@goa gen goa.design/plugins/cors/examples/calc/design -o goa.design/plugins/cors/examples/calc

example:
	@goa example goa.design/plugins/cors/examples/calc/design -o goa.design/plugins/cors/examples/calc
