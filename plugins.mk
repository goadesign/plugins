# Common Makefile rules for plugins
#
# Include this file in all the plugins Makefile as follows
#
#		include $(GOPATH)/src/goa.design/plugins/plugins.mk
#
# Targets:
# - "depend" retrieves the Go packages needed to run the linter and tests
# - "lint" runs the linter and checks the code format using goimports
# - "test" runs the tests
export GO111MODULE=on

all: gen test lint build-examples clean

test:
	@go test ./...

lint:
	$(eval GO_FILES := $(shell find . -type f -name '*.go'))
	@if [ "`goimports -l $(GO_FILES) | tee /dev/stderr`" ]; then \
		echo "^ - Repo contains improperly formatted go files" && echo && exit 1; \
	fi
	@if [ "`golint ./... | grep -vf .golint_exclude | tee /dev/stderr`" ]; then \
		echo "^ - Lint errors!" && echo && exit 1; \
	fi
