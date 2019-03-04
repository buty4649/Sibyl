GO ?= go

default: build

build: get test
	$(GO) build -o sibyl

get:
	$(GO) get ./...

test:
	$(GO) test ./...

.PHONY: build
