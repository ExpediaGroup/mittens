.PHONY: build unit-tests integration-tests

override GOOS:=$(shell uname)
override GO111MODULE=on

unit-tests:
	@CGO_ENABLED=0 go build && go test ./...

integration-tests:
	@go test -tags=integration ./...

docker:
	@docker build -t expediagroup/mittens:latest .
