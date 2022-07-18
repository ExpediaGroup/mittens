---
id: building
title: How to Build
---

Mittens is written in Go and uses [Go Modules](https://github.com/golang/go/wiki/Modules). 
We provide a [Makefile](https://github.com/ExpediaGroup/mittens/blob/main/Makefile) which can be used to generate an executable binary and a Dockerfile if you prefer to run using Docker.

### Binary

To build the binary make sure you've installed [Go 1.16](https://golang.org/dl/).

#### Build binary executable & run unit tests

To build the project run the following:
    
    make unit-tests

This will run the unit tests and generate a binary executable.
    
#### Run integration tests

To run the integration tests:

    make integration-tests
 
### Docker
#### Build image

To build a Docker image named `mittens`:

    make docker
