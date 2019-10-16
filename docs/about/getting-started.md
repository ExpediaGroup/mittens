---
id: getting-started
title: Getting Started
---

The application receives a number of command-line flags including the requests that will be sent to warm up the main service. Depending on the format of the requests this will invoke REST or/and gRPC calls.

Mittens can also run as a linked container, or even as a sidecar in Kubernetes.

## Usage

    mittens [flags]

## Flags

| Flag                               | Type    | Default value    | Description                                                                                                     |
|:-----------------------------------|:--------|:-----------------|:----------------------------------------------------------------------------------------------------------------|
| --concurrency                      | int     | 2                | Number of concurrent requests for warm up                                                                       |
| --exit-after-warmup                | bool    | false            | If warm up process should finish after completion. This is useful to prevent container restarts.                |
| --grpc-headers                     | strings | N/A              | gRPC headers to be sent with warm up requests.                                                                  |
| --grpc-requests                    | strings | N/A              | gRPC request to be sent. Request is in '<service>/<method>[:message]' format. E.g. health/ping:{"key": "value"} |
| --help, -h                         |         |                  | N/A Help for warmup-sidecar                                                                                     |
| --http-headers                     | strings | N/A              | Http headers to be sent with warm up requests.                                                                  |
| --http-requests                    | strings | N/A              | Http request to be sent. Request is in '<http-method>:<path>[:body]' format. E.g. post:/ping:{"key": "value"}   |
| --probe-liveness-path              | strings | /alive           | Warm up sidecar liveness probe path                                                                             |
| --probe-port                       | int     | /8000            | Warm up sidecar port for liveness and readiness probe                                                           |
| --probe-readiness-path             | string  | /ready           | Warm up sidecar readiness probe path                                                                            |
| --profile-cpu                      | string  | N/A              | Name of the file where to write CPU profile data                                                                |
| --profile-memory                   | string  | N/A              | Name of the file where to write memory profile data                                                             |
| --request-delay-milliseconds       | int     | 50               | Delay in milliseconds between requests                                                                          |
| --target-grpc-host                 | string  | localhost        | gRPC host to warm up                                                                                            |
| --target-grpc-port                 | int     | 50051            | gRPC port for warm up requests                                                                                  |
| --target-http-host                 | string  | http://localhost | Http host to warm up                                                                                            |
| --target-http-port                 | int     | 8080             | Http port for warm up requests                                                                                  |
| --target-insecure                  | bool    | false            | Whether to skip TLS validation                                                                                  |
| --target-readiness-path            | string  | /ready           | The path used for target readiness probe                                                                        |
| --target-readiness-timeout-seconds | int     | -1               | Timeout for target readiness probe                                                                              |
| --timeout-seconds                  | int     | 60               | Time after which warm up will stop making requests                                                              |

### Warmup request
A warmup request can be an HTTP one (over REST) or a gRPC one.

#### HTTP requests

HTTP requests are in the form `method:path[:body]` (`body` is optional).
Host and port are taken from `--target-http-host` and
`--target-http-port` flags.

E.g.:
 - `get:/health`: HTTP GET request.
 - `post:/warmupUrl:{"key":"value"}`: POST request with its url being `/warmupUrl` and its body being `{"key":"value"}`.

#### gRPC requests

gRPC requests are in the form `service/method[:message]` (`message` is optional). Host and port are taken from `--target-grpc-host` and
`--target-grpc-port` flags.
