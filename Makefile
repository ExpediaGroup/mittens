.PHONY: all build unit-tests integration-tests

override GOOS:=$(shell uname)
override GO111MODULE=on

all: build

unit-tests:
	@go test ./...

integration-tests:
	@go test -tags=integration ./...

build: unit-tests
	@CGO_ENABLED=0 go build

docker:
	@docker build -t mittens .
