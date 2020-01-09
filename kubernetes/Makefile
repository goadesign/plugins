#! /usr/bin/make
#
# Makefile for Goa v3 Kubernetes plugin
#
# Targets:
# - "gen" generates the goa files for the example services
# - "example" generates the example files for the example services

# include common Makefile content for plugins
include $(GOPATH)/src/goa.design/plugins/plugins.mk

gen:
	@goa gen goa.design/plugins/v3/kubernetes/examples/calc/design -o "$(GOPATH)/src/goa.design/plugins/kubernetes/examples/calc" && \
	make example

example:
	@ rm -rf "$(GOPATH)/src/goa.design/plugins/kubernetes/examples/calc/cmd"
	goa example goa.design/plugins/v3/kubernetes/examples/calc/design -o "$(GOPATH)/src/goa.design/plugins/kubernetes/examples/calc"

build-examples:
	@cd "$(GOPATH)/src/goa.design/plugins/kubernetes/examples/calc" && \
		go build ./cmd/calc && go build ./cmd/calc-cli

clean:
	@cd "$(GOPATH)/src/goa.design/plugins/kubernetes/examples/calc" && \
		rm -f calc calc-cli
