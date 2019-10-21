<h1 align="left">
  <img width="500" alt="Mittens" src="images/mittens_logo.svg">
</h1>

# Mittens
Mittens is a tool that can be used to warm up an http application over REST or gRPC.

## Features

Mittens can run as a standalone [command-line tool](https://expediagroup.github.io/mittens/docs/installation/running#run-as-a-cmd-application), as a [linked Docker container](https://expediagroup.github.io/mittens/docs/installation/running#run-as-a-linked-docker-container), or even as a [sidecar in Kubernetes](https://expediagroup.github.io/mittens/docs/installation/running#run-as-a-sidecar-on-kubernetes).

Its main features are summarised below:
- Sends requests continuously for X seconds
- Supports REST and gRPC
- Supports HTTP and gRPC headers
- Supports concurrent requests
- Provides files or/and endpoints that can be used as liveness/readiness probes in Kubernetes

## Usage
The application receives a number of command-line flags. Read the [documentation](https://expediagroup.github.io/mittens/docs/about/getting-started) for more context.

## How to build and run
Mittens is written in Go and the simplest way to run it is as a cmd application. It receives a number of command line arguments (see [Flags](https://github.com/ExpediaGroup/mittens#flags)).

The project uses [Go Modules](https://github.com/golang/go/wiki/Modules).
We provide a [Makefile](Makefile) which can be used to generate an executable binary and a Dockerfile if you prefer to run using Docker.

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
   
#### Run the executable

To run the binary:
        
    ./mittens --target-readiness-path=/health --target-grpc-port=6565 --timeout-seconds=60 --concurrency=3 --http-requests=get:/hotel/potatoes --grpc-requests=service/method:"{\"foo\":\"bar\", \"bar\":\"foo\"}"

### Docker
#### Build image

To build a Docker image named `mittens`:

    make docker

#### Run container

To run the container:

    docker run mittens:latest --target-readiness-path=/health --target-grpc-port=6565 --timeout-seconds=60 --concurrency=3 --http-requests=get:/hotel/potatoes --grpc-requests=service/method:"{\"foo\":\"bar\", \"bar\":\"foo\"}"

_Note_: If you use Docker for Mac you might need to set `targetHost` to `docker.for.mac.localhost`, or `docker.for.mac.host.internal`, or `host.docker.internal` (depending on your version of Docker) so that your container can resolve localhost.

## Contributing

Please refer to our [CONTRIBUTING](./CONTRIBUTING.md) file.

## Use Cases

* [Hotels.com](https://www.hotels.com/) - Used in the production environment as a [linked Docker container](https://expediagroup.github.io/mittens/docs/installation/running#run-as-a-linked-docker-container) and as a [Kubernetes sidecar](https://expediagroup.github.io/mittens/docs/installation/running#run-as-a-sidecar-on-kubernetes) to eliminate cold starts. 

## References

* [Mittens Documentation](https://expediagroup.github.io/mittens/docs/about/getting-started)
* [Mittens at Docker Hub](https://hub.docker.com/r/expediagroup/mittens/)

## Legal

This project is available under the [Apache 2.0 License](http://www.apache.org/licenses/LICENSE-2.0.html).

Copyright 2019 Expedia, Inc.
