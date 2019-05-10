#! /usr/bin/make
#
# Makefile for goa v2 plugins

# Add new plugins here to enable make
PLUGINS=\
	cors \
	docs \
	goakit \
	zaplogger

export GO111MODULE=on

all: gen lint test-plugins

travis: depend all check-freshness

depend:
	@env GO111MODULE=off go get -v golang.org/x/lint/golint

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

check-freshness:
	@if [ "`git diff | wc -l`" -gt "0" ]; then \
	        echo "[ERROR] generated code not in-sync with design:"; \
	        echo; \
	        git status -s; \
	        git --no-pager diff; \
	        echo; \
	        exit 1; \
	fi
