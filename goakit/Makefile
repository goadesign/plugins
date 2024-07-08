#! /usr/bin/make
#
# Makefile for goa v3 goakit plugin
#
# Targets:
# - "gen" generates the goa files for the example services

# include common Makefile content for plugins
PLUGIN_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
include ../plugins.mk

gen: clean
	@goa gen goa.design/plugins/v3/goakit/examples/calc/design -o "$(PLUGIN_DIR)/examples/calc"
	@goa gen goa.design/plugins/v3/goakit/examples/fetcher/fetcher/design -o "$(PLUGIN_DIR)/examples/fetcher/fetcher"
	@goa gen goa.design/plugins/v3/goakit/examples/fetcher/archiver/design -o "$(PLUGIN_DIR)/examples/fetcher/archiver"
	@goa example goa.design/plugins/v3/goakit/examples/calc/design -o "$(PLUGIN_DIR)/examples/calc"
	@goa example goa.design/plugins/v3/goakit/examples/fetcher/fetcher/design -o "$(PLUGIN_DIR)/examples/fetcher/fetcher"
	@goa example goa.design/plugins/v3/goakit/examples/fetcher/archiver/design -o "$(PLUGIN_DIR)/examples/fetcher/archiver"

build-examples:
	@cd "$(PLUGIN_DIR)/examples/calc" && \
		go build ./cmd/calc && go build ./cmd/calc-cli
	@cd "$(PLUGIN_DIR)/examples/fetcher/fetcher" && \
		go build ./cmd/fetcher && go build ./cmd/fetcher-cli
	@cd "$(PLUGIN_DIR)/examples/fetcher/archiver" && \
		go build ./cmd/archiver && go build ./cmd/archiver-cli

clean:
	@cd "$(PLUGIN_DIR)/examples/calc" && \
		rm -rf cmd calc.go
	@cd "$(PLUGIN_DIR)/examples/fetcher/fetcher" && \
		rm -rf cmd/http.go health.go
	@cd "$(PLUGIN_DIR)/examples/fetcher/archiver" && \
		rm -rf cmd health.go
