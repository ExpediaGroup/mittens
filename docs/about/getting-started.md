---
id: getting-started
title: Getting Started
---

The application receives a number of command-line flags including the requests that will be sent to warm up the main service. Depending on the format of the requests this will invoke REST or/and gRPC calls.

## Usage

    mittens [flags]

## Flags

| Flag                              | Type    | Default value               | Description                                                                                                                                                                                                                                                                             |
|:----------------------------------|:--------|:----------------------------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| -concurrency                      | int     | 2                           | Number of concurrent requests for warm up                                                                                                                                                                                                                                               |
| -exit-after-warmup                | bool    | false                       | If mittens should exit after completion of warm up                                                                                                                                                                                                                                      |
| -http-headers                     | strings | N/A                         | Http headers to be sent with warm up requests. To send multiple headers define this flag for each header                                                                                                                                                                                |
| -grpc-requests                    | strings | N/A                         | gRPC requests to be sent. Request is in '\<service\>\<method\>\[:message\]' format. E.g. health/ping:{"key": "value"}. To send multiple requests, simply repeat this flag for each request. Use the notation `:file/xyz.json` if you want to use an external file for the request body. |
| -http-requests                    | string  | N/A                         | Http request to be sent. Request is in `<http-method>:<path>[:body]` format. E.g. `post:/ping:{"key": "value"}`. To send multiple requests, simply repeat this flag for each request. Use the notation `:file/xyz.json` if you want to use an external file for the request body.       |
| -fail-readiness                   | bool    | false                       | If set to true readiness will fail if the target did not became ready in time                                                                                                                                                                                                           |
| -file-probe-enabled               | bool    | true                        | If set to true writes files that can be used as readiness/liveness probes. a file with the name `alive` is created when Mittens starts and a file named `ready` is created when the warmup completes                                                                                    |
| -request-delay-milliseconds       | int     | 500                         | Delay in milliseconds between requests                                                                                                                                                                                                                                                  |
| -target-grpc-host                 | string  | localhost                   | gRPC host to warm up                                                                                                                                                                                                                                                                    |
| -target-grpc-port                 | int     | 50051                       | gRPC port for warm up requests                                                                                                                                                                                                                                                          |
| -target-http-host                 | string  | http://localhost            | Http host to warm up                                                                                                                                                                                                                                                                    |
| -target-http-port                 | int     | 8080                        | Http port for warm up requests                                                                                                                                                                                                                                                          |
| -target-insecure                  | bool    | false                       | Whether to skip TLS validation                                                                                                                                                                                                                                                          |
| -target-readiness-grpc-method     | string  | grpc.health.v1.Health/Check | The service method used for gRPC target readiness probe                                                                                                                                                                                                                                 |
| -target-readiness-http-path       | string  | /ready                      | The path used for target readiness probe                                                                                                                                                                                                                                                |
| -target-readiness-port            | int     | same as -target-http-port   | The port used for target readiness probe                                                                                                                                                                                                                                                |
| -target-readiness-protocol        | string  | http                        | Protocol to be used for readiness check. One of [`http`, `grpc`]                                                                                                                                                                                                                        |
| -max-duration-seconds             | int     | 60                          | Maximum duration in seconds after which warm up will stop making requests                                                                                                                                                                                                               |
| -concurrency-target-seconds       | int     | 0                           | Time taken to reach expected concurrency. This is useful to ramp up traffic.                                                                                                                                                                                                            |

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

### Placeholders for random elements

Mittens allows you to use special keywords if you need to generate randomized urls.
The following are available:
- `{$currentDate|days+x,months+y,years+z,format=yyyy-MM-dd}`: you can adjust the temporal offset by adding or subtracting days, months, or years. The offsets are optional and can be removed, but their order cannot change (i.e. `days` is always first, or `years` always last). You can optionally specify a custom format using yyyy or yy to represent the year, MM or MMM for the month and dd or d for the day.
- `{$currentTimestamp}`: Time from Unix epoch in milliseconds.
- `{$random|foo,bar,baz}`: Mittens will randomly select an element from the provided list, eg: one of foo, bar or baz. Special chars are not supported. Valid: [0-9A-Za-z_]
- `{$range|min=x,max=y}`: both min and max are required arguments. Range is inclusive.

E.g.:
 - `get:/some-path?date="{$currentDate|days+1,months+1,years+1}"` 
 - `post:/some-path:{"id": "{$range|min=1,max=5}", "currentDate": "{$currentDate|days+2,months+1}"}`

### File probes
Mittens writes files that can be used as liveness and readiness probes. These files are written to disk as `alive` and `ready` respectively.
If you run mittens as a sidecar you can then define a [liveness command](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#define-a-liveness-command) as follows:

```
...
livenessProbe:
  exec:
    command:
    - "cat"
    - "ready"
    ...
```

In case such probes are not needed you can disable this feature by setting `file-probe-enabled` to `false`. 

#### Fail Mittens readiness

Setting `fail-readiness` to true will cause Mittens readiness to fail in case no requests were sent.

### Health checks over HTTP and gRPC

Mittens supports both HTTP and gRPC for application health checks.

By default it uses HTTP to call the `-target-readiness-http-path` endpoint. If your app exposes a health check over gRPC you can set `-target-readiness-protocol` to `grpc` and define the RPC method to be called in `-target-readiness-grpc-method`. Method should be in the form `service/method`.

See [here](https://github.com/grpc/grpc/blob/master/doc/health-checking.md) on how to implement a gRPC health check on your applications. This has already been implemented in many languages including [Java](https://github.com/grpc/grpc-java/blob/master/services/src/main/proto/grpc/health/v1/health.proto) and [Go](https://github.com/grpc/grpc/blob/master/src/proto/grpc/health/v1/health.proto).

Based on the [gRPC Health Checking Protocol](https://github.com/grpc/grpc/blob/master/doc/health-checking.md) the suggested format for the service name is `grpc.health.v1.Health
` which would translate to `-target-readiness-grpc-method=grpc.health.v1.Health/Check`.
