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
	@cd examples/calc/calc && \
	goa gen goa.design/plugins/security/examples/calc/calc/design
	@cd examples/calc/adder && \
	goa gen goa.design/plugins/security/examples/calc/adder/design

aliases:
	@aliaser -dsl -src goa.design/goa/http/dsl -dest $(PLUGIN_DIR)/$(PLUGIN_NAME)/dsl > /dev/null

test-aliaser: aliases
	@if [ "`git diff */aliases.go | tee /dev/stderr`" ]; then \
		echo "^ - Aliaser tool output not identical!" && echo && exit 1; \
	else \
		echo "Aliaser tool output identical"; \
	fi
