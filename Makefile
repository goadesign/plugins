#! /usr/bin/make
#
# Makefile for goa v2 plugins

PLUGINS=$(shell find . -mindepth 1 -maxdepth 1 -not -path "*/\.*" -type d)

all: test-plugins test-aliaser

test-plugins:
	@for p in $(PLUGINS) ; do \
		make -C $$p || exit 1; \
	done

test-aliaser:
	@for p in $(PLUGINS) ; do \
		make -C $$p test-aliaser || exit 1; \
	done
