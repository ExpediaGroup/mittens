<h1 align="left">
  <img width="500" alt="Mittens" src="images/mittens_logo.svg">
</h1>

# Mittens
Mittens is a -command line- tool that can be used to warm up an http application over REST or gRPC.

## Design
The application receives a number of command-line flags including the requests that will be sent to warm up the main service. Depending on the format of the requests this will invoke REST or/and gRPC calls.

Mittens can run as a standalone command-line tool, as a [linked Docker container](https://expediagroup.github.io/mittens/docs/installation/running#run-as-a-linked-docker-container) or even as a [Kubernetes sidecar](https://expediagroup.github.io/mittens/docs/installation/running#run-as-a-sidecar-on-kubernetes).

## Usage

    mittens [flags]

### Flags

| Flag                               | Type    | Default value    | Description                                                                                                           |
|:-----------------------------------|:--------|:-----------------|:----------------------------------------------------------------------------------------------------------------------|
| --concurrency                      | int     | 2                | Number of concurrent requests for warm up                                                                             |
| --exit-after-warmup                | bool    | false            | If warm up process should finish after completion. This is useful to prevent container restarts.                      |
| --grpc-headers                     | strings | N/A              | gRPC headers to be sent with warm up requests.                                                                        |
| --grpc-requests                    | strings | N/A              | gRPC request to be sent. Request is in '\<service\>/\<method\>\[:message\]' format. E.g. health/ping:{"key": "value"} |
| --help, -h                         |         | N/A              | Help for warmup-sidecar                                                                                               |
| --http-headers                     | strings | N/A              | Http headers to be sent with warm up requests.                                                                        |
| --http-requests                    | strings | N/A              | Http request to be sent. Request is in '\<http-method\>:\<path\>\[:body\]' format. E.g. post:/ping:{"key": "value"}   |
| --probe-liveness-path              | strings | /alive           | Warm up sidecar liveness probe path                                                                                   |
| --probe-port                       | int     | /8000            | Warm up sidecar port for liveness and readiness probe                                                                 |
| --probe-readiness-path             | string  | /ready           | Warm up sidecar readiness probe path                                                                                  |
| --profile-cpu                      | string  | N/A              | Name of the file where to write CPU profile data                                                                      |
| --profile-memory                   | string  | N/A              | Name of the file where to write memory profile data                                                                   |
| --request-delay-milliseconds       | int     | 50               | Delay in milliseconds between requests                                                                                |
| --target-grpc-host                 | string  | localhost        | gRPC host to warm up                                                                                                  |
| --target-grpc-port                 | int     | 50051            | gRPC port for warm up requests                                                                                        |
| --target-http-host                 | string  | http://localhost | Http host to warm up                                                                                                  |
| --target-http-port                 | int     | 8080             | Http port for warm up requests                                                                                        |
| --target-insecure                  | bool    | false            | Whether to skip TLS validation                                                                                        |
| --target-readiness-path            | string  | /ready           | The path used for target readiness probe                                                                              |
| --target-readiness-timeout-seconds | int     | -1               | Timeout for target readiness probe                                                                                    |
| --timeout-seconds                  | int     | 60               | Time after which warm up will stop making requests                                                                    |

### Warmup request
A warmup request can be an HTTP one (over REST) or a gRPC one.

#### HTTP requests

HTTP requests are in the form `method:path\[:body\]` (`body` is optional).
Host and port are taken from `--target-http-host` and
`--target-http-port` flags.

E.g.:
 - `get:/health`: HTTP GET request.
 - `post:/warmupUrl:{"key":"value"}`: POST request with its url being `/warmupUrl` and its body being `{"key":"value"}`.

#### gRPC requests

gRPC requests are in the form `service/method:\[message\]` (`message` is optional). Host and port are taken from `--target-grpc-host` and
`--target-grpc-port` flags.

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

#### Run tests

To run the tests:

    make test
   
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
