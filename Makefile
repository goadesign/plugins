#! /usr/bin/make
#
# Makefile for goa v3 plugins

GOPATH=$(shell go env GOPATH)
GOA:=$(shell goa version 2> /dev/null)

# Add new plugins here to enable make
PLUGINS=\
	cors \
	docs \
	goakit \
	zaplogger \
	i18n

export GO111MODULE=on

all: check-goa gen lint test

travis: depend all check-freshness

check-goa:
ifdef GOA
	go mod download
	@echo $(GOA)
else
	go get -u goa.design/goa/v3@v3
	go get -u goa.design/goa/v3/...@v3
	go mod download
	@echo $(GOA)
endif

$(GOPATH)/bin/goimports:
	@go get golang.org/x/tools/cmd/goimports

$(GOPATH)/bin/golint:
	@go get golang.org/x/lint/golint

$(GOPATH)/bin/goa:
	@go install goa.design/goa/v3/cmd/goa && goa version

depend: $(GOPATH)/bin/goimports $(GOPATH)/bin/golint $(GOPATH)/bin/goa

tidy:
	@go mod tidy -v

gen:
	@for p in $(PLUGINS) ; do \
		make -C $$p gen || exit 1; \
	done

fmt: $(GOPATH)/bin/goimports
	@files=$$(find . -type f -not -path '*/\.*' -not -path "./vendor/*" -name "*\.go" | grep -Ev '/(gen)/'); \
	$(GOPATH)/bin/goimports -w -l $$files

lint: $(GOPATH)/bin/golint
	@for p in $(PLUGINS) ; do \
		make -C $$p lint || exit 1; \
	done

test:
	@for p in $(PLUGINS) ; do \
		make -C $$p test || exit 1; \
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
