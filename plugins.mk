# Common Makefile rules for plugins
#
# Include this file in all the plugins Makefile as follows
#
#		include $(GOPATH)/src/goa.design/plugins/plugins.mk
#
# In the plugin Makefile set the following variables
# - PLUGIN_NAME - Name of the plugin (the directory name)
# - ALIASER_SRC - The src DSL for aliaser command
#
# Targets:
# - "depend" retrieves the Go packages needed to run the linter and tests
# - "lint" runs the linter and checks the code format using goimports
# - "aliases" runs the aliaser command for the plugin DSL
# - "test" runs the tests
# - "test-aliaser" checks if there are any uncommitted changes from aliaser command

PLUGIN_DIR=goa.design/plugins

DEPEND=\
  github.com/sergi/go-diff/diffmatchpatch \
  github.com/golang/lint/golint \
  golang.org/x/tools/cmd/goimports \
	goa.design/goa/...

all: depend test lint

depend:
	@go get -t -v ./...
	@go get -v $(DEPEND)

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

aliases:
	@if [ -z "$(PLUGIN_NAME)" ]; then \
		echo "PLUGIN_NAME not set in Makefile!" && exit 1; \
	fi
	@if [ -z "$(ALIASER_SRC)" ]; then \
		echo "ALIASER_SRC not set in Makefile!" && exit 1; \
	fi
	@aliaser -src $(ALIASER_SRC) -dest $(PLUGIN_DIR)/$(PLUGIN_NAME)/dsl > /dev/null; \

test-aliaser: aliases
	@if [ "`git diff */aliases.go | tee /dev/stderr`" ]; then \
		echo "^ - Aliaser tool output not identical!" && echo && exit 1; \
	else \
		echo "Aliaser tool output identical"; \
	fi
