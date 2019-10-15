---
id: building
title: How to Build
---

Mittens is written in Go and uses [Go Modules](https://github.com/golang/go/wiki/Modules). 
We provide a [Makefile](https://github.com/ExpediaGroup/mittens/blob/master/Makefile) which can be used to generate an executable binary and a Dockerfile if you prefer to run using Docker.

### Binary
#### Build binary executable

To build the project run the following:
    
    make build

This will generate a binary executable.

#### Run tests

To run the tests:

    make test
 
### Docker
#### Build image

To build a Docker image named `mittens`:

    make docker
