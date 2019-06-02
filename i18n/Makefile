#! /usr/bin/make
#
# Makefile for goa v3 I18n plugin
#
# Targets:
# - "gen" generates the goa files for the example services
# - "example" generates the example files for the example services

# include common Makefile content for plugins
include $(GOPATH)/src/goa.design/plugins/plugins.mk

gen:
	@GOA_I18N=en,nl goa gen goa.design/plugins/v3/i18n/examples/calc/design -o "$(GOPATH)/src/goa.design/plugins/i18n/examples/calc" && \
	make example

example:
	@ rm -rf "$(GOPATH)/src/goa.design/plugins/i18n/examples/calc/cmd"
	goa example goa.design/plugins/v3/i18n/examples/calc/design -o "$(GOPATH)/src/goa.design/plugins/i18n/examples/calc"

build-examples:
	@cd "$(GOPATH)/src/goa.design/plugins/i18n/examples/calc" && \
		go build ./cmd/calc && go build ./cmd/calc-cli

clean:
	@cd "$(GOPATH)/src/goa.design/plugins/i18n/examples/calc" && \
		rm -f calc calc-cli
