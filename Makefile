.PHONY: test

override GOOS:=$(shell uname)
override GO111MODULE=on

test:
	@CGO_ENABLED=0 go build && go test ./...

docker:
	@docker build -t expediagroup/mittens:latest .
