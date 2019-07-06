#! /usr/bin/make
#
# Makefile for goa v3 plugins

# Add new plugins here to enable make
PLUGINS=\
	cors \
	docs \
	goakit \
	zaplogger \
	i18n

export GO111MODULE=on

all: gen lint test-plugins

travis: depend all check-freshness

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
