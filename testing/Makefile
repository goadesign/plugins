#! /usr/bin/make
#
# Makefile for testing plugin
#

PLUGIN_NAME=testing
PLUGIN_SRC=testing
GO_MODULE=goa.design/plugins/v3/testing

# Include the common plugin makefile
include ../plugins.mk

.PHONY: test-example
test-example: ## Test the example
	@echo "Testing example..."
	@cd example && goa gen design && go test -v ./...