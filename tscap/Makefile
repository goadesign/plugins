#! /usr/bin/make
#
# Makefile for goa v3 tscap plugin
#
# Targets:
# - "gen" generates the goa files for the example services
# - "example" generates the example files for the example services

# include common Makefile content for plugins
PLUGIN_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
include ../plugins.mk

gen:
	@goa gen goa.design/plugins/v3/tscap/example/design -o "$(PLUGIN_DIR)/example" && \
	make example

example:
	@rm -rf "$(PLUGIN_DIR)/example/cmd"
	goa example goa.design/plugins/v3/tscap/example/design -o "$(PLUGIN_DIR)/example"

build-examples:
	@cd "$(PLUGIN_DIR)/example" && \
		go build ./cmd/example

clean:
	@cd "$(PLUGIN_DIR)/example" && \
		rm -f example
