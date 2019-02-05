#! /usr/bin/make
#
# Makefile for goa v2 plugins

# Add new plugins here to enable make
PLUGINS=\
	cors \
	goakit \
	zaplogger

all: depend gen test-plugins

depend:
	@go get -v golang.org/x/lint/golint
	@for p in $(PLUGINS) ; do \
		make -C $$p depend || exit 1; \
	done

gen:
	@for p in $(PLUGINS) ; do \
		make -C $$p gen || exit 1; \
	done

test-plugins:
	@for p in $(PLUGINS) ; do \
		make -C $$p || exit 1; \
	done
