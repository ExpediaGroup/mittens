---
id: building
title: How to Build
---

Mittens is written in Go and uses [Go Modules](https://github.com/golang/go/wiki/Modules). 
We provide a [Makefile](https://github.com/ExpediaGroup/mittens/blob/master/Makefile) which can be used to generate an executable binary and a Dockerfile if you prefer to run using Docker.

### Binary

To build the binary make sure you've installed [Go 1.13](https://golang.org/dl/).

#### Build binary executable

To build the project run the following:
    
    make build

This will generate a binary executable.

#### Run unit tests

To run the tests:

    make tests
    
#### Run integration tests

To run the integration tests:

    make integration-tests
 
### Docker
#### Build image

To build a Docker image named `mittens`:

    make docker
