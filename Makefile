#! /usr/bin/make
#
# Makefile for goa v2 plugins

PLUGINS=$(shell find . -mindepth 1 -maxdepth 1 -not -path "*/\.*" -type d)

all: depend test-plugins

depend:
	@mkdir -p $(GOPATH)/src/golang.org/x
	@git clone https://github.com/golang/lint.git $(GOPATH)/src/golang.org/x/lint | true

gen:
	@for p in $(PLUGINS) ; do \
		make -C $$p gen || exit 1; \
	done

test-plugins:
	@for p in $(PLUGINS) ; do \
		make -C $$p || exit 1; \
	done
