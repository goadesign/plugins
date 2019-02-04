#! /usr/bin/make
#
# Makefile for goa v2 zaplogger plugin
#
# Targets:
# - "gen" generates the goa files for the example services

PLUGIN_NAME=zaplogger
ALIASER_SRC=goa.design/goa/http/dsl

# include common Makefile content for plugins
include $(GOPATH)/src/goa.design/plugins/plugins.mk

gen:
	goa gen goa.design/plugins/zaplogger/examples/calc/design -o $(GOPATH)/src/goa.design/plugins/zaplogger/examples/calc && \
	make example

example:
	@ rm -rf $(GOPATH)/src/goa.design/plugins/zaplogger/examples/calc/cmd && \
	goa example goa.design/plugins/zaplogger/examples/calc/design -o $(GOPATH)/src/goa.design/plugins/zaplogger/examples/calc
