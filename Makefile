.PHONY: all build test

all: build

test:
	GO111MODULE=on go test ./...

build: test
	GO111MODULE=on go build

docker:
	@docker build -t mittens .
