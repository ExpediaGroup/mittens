---
id: getting-started
title: Getting Started
---

The application receives a number of command-line flags including the requests that will be sent to warm up the main service. Depending on the format of the requests this will invoke REST or/and gRPC calls.

## Usage

    mittens [flags]

## Flags

| Flag                              | Type    | Default value               | Description                                                                                                                                                                        |
|:----------------------------------|:--------|:----------------------------|:-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| -concurrency                      | int     | 2                           | Number of concurrent requests for warm up                                                                                                                                          |
| -exit-after-warmup                | bool    | false                       | If warm up process should exit after completion                                                                                                                                    |
| -grpc-headers                     | strings | N/A                         | gRPC headers to be sent with warm up requests. To send multiple headers define this flag for each header                                                                           |
| -grpc-requests                    | strings | N/A                         | gRPC requests to be sent. Request is in '\<service\>\<method\>\[:message\]' format. E.g. health/ping:{"key": "value"}. To send multiple requests define this flag for each request |
| -http-headers                     | strings | N/A                         | Http headers to be sent with warm up requests. To send multiple headers define this flag for each header                                                                           |
| -http-requests                    | string  | N/A                         | Http request to be sent. Request is in '\<http-method\>:\<path\>\[:body\]' format. E.g. post:/ping:{"key": "value"}. To send multiple requests define this flag for each request   |
| -file-probe-enabled               | bool    | true                        | If set to true writes files to be used as readiness/liveness probes                                                                                                                |
| -file-probe-liveness-path         | string  | alive                       | File to be used for liveness probe                                                                                                                                                 |
| -file-probe-readiness-path        | string  | ready                       | File to be used for readiness probe                                                                                                                                                |
| -server-probe-enabled             | bool    | false                       | If set to true runs a web server that exposes endpoints to be used as readiness/liveness probes                                                                                    |
| -server-probe-port                | int     | 8000                        | Port on which probe server is running                                                                                                                                              |
| -server-probe-liveness-path       | string  | /alive                      | Probe server endpoint used as liveness probe                                                                                                                                       |
| -server-probe-readiness-path      | string  | /ready                      | Probe server endpoint used as readiness probe                                                                                                                                      |
| -profile-cpu                      | string  | ""                          | Name of the file where to write CPU profile data. If empty no CPU profiling takes place                                                                                            |
| -profile-memory                   | string  | ""                          | Name of the file where to write memory profile data. If empty no memory profiling takes place                                                                                      |
| -request-delay-milliseconds       | int     | 50                          | Delay in milliseconds between requests                                                                                                                                             |
| -target-grpc-host                 | string  | localhost                   | gRPC host to warm up                                                                                                                                                               |
| -target-grpc-port                 | int     | 50051                       | gRPC port for warm up requests                                                                                                                                                     |
| -target-http-host                 | string  | http://localhost            | Http host to warm up                                                                                                                                                               |
| -target-http-port                 | int     | 8080                        | Http port for warm up requests                                                                                                                                                     |
| -target-insecure                  | bool    | false                       | Whether to skip TLS validation                                                                                                                                                     |
| -target-readiness-grpc-method     | string  | grpc.health.v1.Health/Check | The service method used for gRPC target readiness probe                                                                                                                            |
| -target-readiness-http-path       | string  | /ready                      | The path used for target readiness probe                                                                                                                                           |
| -target-readiness-port            | int     | same as -target-http-port   | The port used for target readiness probe                                                                                                                                           |
| -target-readiness-protocol        | string  | http                        | Protocol to be used for readiness check. One of [http, grpc]                                                                                                                       |
| -target-readiness-timeout-seconds | int     | -1                          | Timeout for target readiness probe                                                                                                                                                 |
| -timeout-seconds                  | int     | 60                          | Time after which warm up will stop making requests                                                                                                                                 |

### Warmup request
A warmup request can be an HTTP one (over REST) or a gRPC one.

#### HTTP requests

HTTP requests are in the form `method:path[:body]` (`body` is optional).
Host and port are taken from `target-http-host` and
`target-http-port` flags.

E.g.:
 - `get:/health`: HTTP GET request.
 - `post:/warmupUrl:{"key":"value"}`: POST request with its url being `/warmupUrl` and its body being `{"key":"value"}`.

#### gRPC requests

gRPC requests are in the form `service/method[:message]` (`message` is
optional). Host and port are taken from `target-grpc-host` and
`target-grpc-port` flags.

#### Placeholders for random elements

Mittens allows you to use special keywords if you need to generate randomized urls.
The following are available:
- `{{currentDate|days+x,months+y,years+z}}`: you can adjust the temporal offset by adding or subtracting days, months or years. The offsets are optional and can be removed.
- {{currentTimestamp}} - Time from Unix epoch in milliseconds.
- `{{random|foo,bar,baz}}`: Mittens will randomly select an element from the provided list, eg: one of foo, bar or baz. Special chars are not supported. Valid: [0-9A-Za-z_]
- {{range|min=x,max=y}} - both min and max are required arguments. Range is inclusive.

E.g.:
 - `get:/some-path?date="{{currentDate|days+1,months+1,years+1}}"` 
 - `post:/some-path:{"id": "{{range|min=1,max=5}}", "currentDate": "{{currentDate|days+2,months+1}}"}`

### Liveness/readiness probes

#### File probes
By default Mittens writes files that can be used as liveness/readiness probes. Using files is the suggested way for such probes and is preferred over server probes for the following reasons:
- Running a web server for probes increases memory/cpu consumed by Mittens. This needs to be taken into consideration when setting the resources for this container in Kubernetes.
- Using files is less error-prone; a file is persisted on disk and will be there whenever Kubernetes does a liveness check against the pod. On the other hand, an endpoint could at any point be unavailable for all sorts of reasons.

In case such probes are not needed you can disable this feature by setting `file-probe-enabled` to `false`. 

#### Server probes

Setting `server-probe-enabled` to `true` will start a web server that exposes liveness/readiness endpoints. 
Note that running this web server instead of or in addition to having file probes increases memory and cpu consumption.

### Health checks over HTTP and gRPC

Mittens supports both HTTP and gRPC for application health checks.

By default it uses HTTP to call the `-target-readiness-http-path` endpoint. If your app exposes a health check over gRPC you can set `-target-readiness-protocol` to `grpc` and define the RPC method to be called in `-target-readiness-grpc-method`. Method should be in the form `service/method`.
See [here](https://github.com/grpc/grpc/blob/master/doc/health-checking.md) on how to implement a gRPC health check on your applications. This has already been implemented in many languages including [Java](https://github.com/grpc/grpc-java/blob/master/services/src/main/proto/grpc/health/v1/health.proto) and [Go](https://github.com/grpc/grpc/blob/master/src/proto/grpc/health/v1/health.proto).
Based on the [gRPC Health Checking Protocol](https://github.com/grpc/grpc/blob/master/doc/health-checking.md) the suggested format of service name is `grpc.health.v1.Health` which would translate to `-target-readiness-grpc-method=grpc.health.v1.Health/Check`.
