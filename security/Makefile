#! /usr/bin/make
#
# Makefile for goa v2 security plugin
#
# Targets:
# - "gen" generates the goa files for the example services

PLUGIN_NAME=security
ALIASER_SRC=goa.design/goa/http/dsl

# include common Makefile content for plugins
include $(GOPATH)/src/goa.design/plugins/plugins.mk

gen:
	@cd examples/multi_auth && \
	goa gen goa.design/plugins/security/examples/multi_auth/design
