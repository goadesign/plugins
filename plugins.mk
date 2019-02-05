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

PLUGIN_DIR=goa.design/plugins

DEPEND=\
  github.com/sergi/go-diff/diffmatchpatch \
  golang.org/x/lint/golint \
  golang.org/x/tools/cmd/goimports \
	goa.design/goa/...

all: depend gen test lint build-examples clean

GOOS=$(shell go env GOOS)
ifeq ($(GOOS),windows)
GOA_PATH="$(GOPATH)\src\goa.design\goa"
else
GOA_PATH="$(GOPATH)/src/goa.design/goa"
endif
depend:
	@go get -t -v ./...
	@go get -v $(DEPEND)
	@cd $(GOA_PATH)/cmd/goa && go install

test:
	@go test ./...

lint:
	$(eval DIRS := $(shell go list -f {{.Dir}} ./...))
	@for d in $(DIRS) ; do \
		if [ "`goimports -l $$d/*.go | tee /dev/stderr`" ]; then \
			echo "^ - Repo contains improperly formatted go files" && echo && exit 1; \
		fi \
	done
	@if [ "`golint ./... | grep -vf .golint_exclude | tee /dev/stderr`" ]; then \
		echo "^ - Lint errors!" && echo && exit 1; \
	fi
