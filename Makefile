.PHONY: all build test

override GOOS:=$(shell uname)
override GO111MODULE=on

all: build

dependencies:
	@go mod download

test:
	@go test ./...

build: test
	@mkdir -p build
	@CGO_ENABLED=0 go build -a -installsuffix cgo -o build/mittens .

docker:
	@docker build -t mittens .
