#! /usr/bin/make
#
# Makefile for goa v3 plugins

GOOS=$(shell go env GOOS)
GOPATH=$(shell go env GOPATH)
GOA:=$(shell goa version 2> /dev/null)

# Only list test and build dependencies
# Standard dependencies are installed via go get
DEPEND=\
	google.golang.org/protobuf/cmd/protoc-gen-go@latest \
	google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
	honnef.co/go/tools/cmd/staticcheck@latest \
	goa.design/goa/v3/cmd/goa@v3

# Add new plugins here to enable make
PLUGINS=\
	cors \
	docs \
	goakit \
	i18n \
	otel \
	types \
	zaplogger \
	zerologger

PROTOC_VERSION=22.2
UNZIP=unzip
ifeq ($(GOOS),linux)
	PROTOC=protoc-$(PROTOC_VERSION)-linux-x86_64
	PROTOC_EXEC=$(PROTOC)/bin/protoc
endif
ifeq ($(GOOS),darwin)
	PROTOC=protoc-$(PROTOC_VERSION)-osx-x86_64
	PROTOC_EXEC=$(PROTOC)/bin/protoc
endif
ifeq ($(GOOS),windows)
	PROTOC=protoc-$(PROTOC_VERSION)-win32
	PROTOC_EXEC="$(PROTOC)\bin\protoc.exe"
	GOPATH:=$(subst \,/,$(GOPATH))
endif

all: check-goa gen tidy lint test

ci: depend all 

depend:
	@echo INSTALLING DEPENDENCIES...
	@go mod download
	@for package in $(DEPEND); do go install $$package; done
	@go mod tidy -compat=1.19
	@echo INSTALLING PROTOC...
	@mkdir $(PROTOC)
	@cd $(PROTOC); \
	curl -O -L https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/$(PROTOC).zip; \
	$(UNZIP) $(PROTOC).zip
	@cp $(PROTOC_EXEC) $(GOPATH)/bin && \
		rm -rf $(PROTOC) && \
		echo "`protoc --version`"

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

tidy:
	@go mod tidy

gen:
	@for p in $(PLUGINS) ; do \
		make -C $$p gen || exit 1; \
	done

lint:
ifneq ($(GOOS),windows)
	@if [ "`staticcheck ./... | grep -v ".pb.go" | tee /dev/stderr`" ]; then \
		echo "^ - staticcheck errors!" && echo && exit 1; \
	fi
endif

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
