#! /usr/bin/make
#
# Makefile for testing plugin
#
# Targets:
# - "gen" generates the goa files for the example services

# Include the common plugin makefile
PLUGIN_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
PLUGIN_CUSTOM_TEST := 1
include ../plugins.mk

gen:
	@cd "$(PLUGIN_DIR)/examples/calculator" && \
		goa gen goa.design/plugins/v3/testing/examples/calculator/design -o "$(PLUGIN_DIR)/examples/calculator"
	@cd "$(PLUGIN_DIR)/examples/httpgrpc" && \
		goa gen goa.design/plugins/v3/testing/examples/httpgrpc/design -o "$(PLUGIN_DIR)/examples/httpgrpc"
	@cd "$(PLUGIN_DIR)/examples/jsonrpc" && \
		goa gen goa.design/plugins/v3/testing/examples/jsonrpc/design -o "$(PLUGIN_DIR)/examples/jsonrpc"
	make example

example:
	@cd "$(PLUGIN_DIR)/examples/calculator" && \
		goa example goa.design/plugins/v3/testing/examples/calculator/design -o "$(PLUGIN_DIR)/examples/calculator"
	@cd "$(PLUGIN_DIR)/examples/httpgrpc" && \
		goa example goa.design/plugins/v3/testing/examples/httpgrpc/design -o "$(PLUGIN_DIR)/examples/httpgrpc"
	@cd "$(PLUGIN_DIR)/examples/jsonrpc" && \
		goa example goa.design/plugins/v3/testing/examples/jsonrpc/design -o "$(PLUGIN_DIR)/examples/jsonrpc"

test:
	@cd "$(PLUGIN_DIR)/examples/calculator" && \
		go test -v ./...
	@cd "$(PLUGIN_DIR)/examples/httpgrpc" && \
		go test -v ./...
	@cd "$(PLUGIN_DIR)/examples/jsonrpc" && \
		go test -v ./...

build-examples:
	@cd "$(PLUGIN_DIR)/examples/calculator" && \
		go build ./cmd/calculator && go build ./cmd/calculator-cli
	@cd "$(PLUGIN_DIR)/examples/httpgrpc" && \
		go build ./cmd/test_http_grpc && go build ./cmd/test_http_grpc-cli
	@cd "$(PLUGIN_DIR)/examples/jsonrpc" && \
		go build ./cmd/test_jsonrpc && go build ./cmd/test_jsonrpc-cli

clean:
	@cd "$(PLUGIN_DIR)/examples/calculator" && \
		rm -rf calculator calculator-cli
	@cd "$(PLUGIN_DIR)/examples/httpgrpc" && \
		rm -rf test_http_grpc test_http_grpc-cli
	@cd "$(PLUGIN_DIR)/examples/jsonrpc" && \
		rm -rf test_jsonrpc test_jsonrpc-cli
