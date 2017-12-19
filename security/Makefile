#! /usr/bin/make
#
# Makefile for goa v2 security plugin
#
# Targets:
# - "depend" retrieves the Go packages needed to run the linter and tests
# - "lint" runs the linter and checks the code format using goimports
# - "test" runs the tests
#
# Meta targets:
# - "all" is the default target, it runs all the targets in the order above.

# Only list test and build dependencies
# Standard dependencies are installed via go get
DEPEND=\
  github.com/sergi/go-diff/diffmatchpatch \
  github.com/golang/lint/golint \
  golang.org/x/tools/cmd/goimports \
	goa.design/goa/...

all: depend lint aliases gen test

depend:
	@go get -t -v ./...
	@go get -v $(DEPEND)

lint:
	@if [ "`goimports -l *.go | tee /dev/stderr`" ]; then \
		echo "^ - Repo contains improperly formatted go files" && echo && exit 1; \
	fi
	@if [ "`golint ./... | grep -vf .golint_exclude | tee /dev/stderr`" ]; then \
    echo "^ - Lint errors!" && echo && exit 1; \
  fi

aliases:
	@aliaser -src goa.design/goa/http/dsl -dest goa.design/plugins/security/dsl > /dev/null

gen:
	@cd examples/multi_auth && \
	goa gen goa.design/plugins/security/examples/multi_auth/design

test:
	go test ./...

test-aliaser: aliases
	@if [ "`git diff */aliases.go | tee /dev/stderr`" ]; then \
		echo "^ - Aliaser tool output not identical!" && echo && exit 1; \
	 else \
      echo "Aliaser tool output identical"; \
   fi
