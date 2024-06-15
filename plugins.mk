# Common Makefile rules for plugins
#
# Include this file in all the plugins Makefile as follows
#
#		include ../plugins.mk
#
# Each plugin Makefile may define the following targets:
#
# - gen: generate the plugin code
# - build-examples: build the plugin examples binaries
# - clean: clean the plugin examples binaries
#
# Targets:
# - "all" calls "gen", "test", "lint", "build-examples" and "clean"
# - "lint" runs the linter and checks the code format using goimports
# - "test" runs the tests

all: gen test lint build-examples clean

test:
	@go test ./...

lint:
	$(eval GO_FILES := $(shell find . -type f -name '*.go'))
	@if [ "`goimports -l $(GO_FILES) | tee /dev/stderr`" ]; then \
		echo "^ - Repo contains improperly formatted go files" && echo && exit 1; \
	fi
	@if [ "`staticcheck ./... | grep -vf .golint_exclude | tee /dev/stderr`" ]; then \
		echo "^ - Lint errors!" && echo && exit 1; \
	fi
