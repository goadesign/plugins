#! /usr/bin/make
#
# Makefile for goa v2 plugins

GOOS=$(shell go env GOOS)
ifeq ($(GOOS),windows)
PLUGINS=$(shell /usr/bin/find . -mindepth 1 -maxdepth 1 -not -path "*[/\]\.*" -type d)
else
PLUGINS=$(shell find . -mindepth 1 -maxdepth 1 -not -path "*/\.*" -type d)
endif

all: depend test-plugins

depend:
	@go get -v golang.org/x/lint/golint

gen:
	@for p in $(PLUGINS) ; do \
		make -C $$p gen || exit 1; \
	done

test-plugins: gen
	@for p in $(PLUGINS) ; do \
		make -C $$p || exit 1; \
	done
