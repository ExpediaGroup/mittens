.PHONY: all build test

override GOOS:=$(shell uname)
override GO111MODULE=on

all: build

test:
	@go test ./...

build: test
	@CGO_ENABLED=0 go build

docker:
	@docker build -t mittens .
