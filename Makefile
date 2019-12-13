#! /usr/bin/make
#
# Makefile for goa v2 plugins

export GO111MODULE=off

# Add new plugins here to enable make
PLUGINS=\
	cors \
	goakit \
	zaplogger

all: depend gen lint test-plugins

depend:
	@go get -v golang.org/x/lint/golint
	@for p in $(PLUGINS) ; do \
		make -C $$p depend || exit 1; \
	done

gen:
	@for p in $(PLUGINS) ; do \
		make -C $$p gen || exit 1; \
	done

lint:
	@for p in $(PLUGINS) ; do \
		make -C $$p lint || exit 1; \
	done

test-plugins:
	@for p in $(PLUGINS) ; do \
		make -C $$p || exit 1; \
	done
