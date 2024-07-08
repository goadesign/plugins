#! /usr/bin/make
#
# Makefile for Goa zaplogger plugin
#
# Targets:
# - "gen" generates the goa files for the example services

# include common Makefile content for plugins
GOPATH=$(shell go env GOPATH)
include ../plugins.mk

gen: clean
	goa gen goa.design/plugins/v3/zaplogger/examples/calc/design -o "$(GOPATH)/src/goa.design/plugins/zaplogger/examples/calc" && \
	make example

example:
	@ rm -rf "$(GOPATH)/src/goa.design/plugins/zaplogger/examples/calc/cmd" && \
	goa example goa.design/plugins/v3/zaplogger/examples/calc/design -o "$(GOPATH)/src/goa.design/plugins/zaplogger/examples/calc"

build-examples:
	@cd "$(GOPATH)/src/goa.design/plugins/zaplogger/examples/calc" && \
		go build ./cmd/calc && go build ./cmd/calc-cli

clean:
	@cd "$(GOPATH)/src/goa.design/plugins/zaplogger/examples/calc" && \
		rm -rf cmd calc.go
