---
id: getting-started
title: Getting Started
---

The application receives a number of command line arguments including the requests that will be sent to warm up the main service. Depending on the format of the requests this will invoke REST or/and gRPC calls.

## Required arguments

| Argument                        | Default value   | Description                                                                                                                                                                                            |
|:--------------------------------|:----------------|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| readinessPath                   | /ready          | Path used as an http readiness probe; we consider the main app to be ready when a GET request to this path returns a 200 http status code.                                                             |
| timeoutForReadinessProbeSeconds | durationSeconds | How long we should wait for the readiness probe before running Mittens (time in seconds). When missing it waits for "durationSeconds".                                                                 |
| durationSeconds                 | 60              | Duration of the warmup routine. If "timeoutForReadinessProbeSeconds" is set this is the duration of the warmup step only. If "timeoutForReadinessProbeSeconds" is not set, this is the total duration. |
| httpTimeoutSeconds              | 10              | Timeout for the http client (for both readiness checks and for the actual warmup routine).                                                                                                             |
| concurrency                     | 2               | Number of concurrent requests during the warmup routine.                                                                                                                                               |
| httpHeader                      | N/A             | Http headers to be sent during the warmup routine. Format: "Header=Value".                                                                                                                             |
| grpcHeader                      | N/A             | Grpc headers to be sent during the warmup routine. Format: "Header:Value".                                                                                                                             |
| warmupRequest                   | N/A             | Request to be used during the warmup routine. See the subsection below on how to format your request.                                                                                                  |
| targetProtocol                  | http            | Protocol for target server. Possible values are http (default) or https.                                                                                                                               |
| targetHost                      | localhost       | The target server.                                                                                                                                                                                     |
| targetHttpPort                  | 8080            | Port for target http server.                                                                                                                                                                           |
| targetGrpcPort                  | 50051           | Port for target grpc server.                                                                                                                                                                           |
| requestDelayMilliseconds        | 0               | Adds a delay between requests.                                                                                                                                                                         |

## Warmup request format
A warmup request can be an HTTP one (over REST) or a gRPC one.

### HTTP (REST) requests
HTTP requests are in the form `http:method:url:body` where `method` is one of `get`, `post`, or `put`, `url` is the url where the request will be sent, and `body` is a properly escaped JSON-formatted string.

Indicative examples are shown below:
- `http:get:/health`: HTTP GET request.
- `http:post:/warmupUrl:{"key":"value"}`: POST request with its url being `/warmupUrl` and its body being `{"key":"value"}`.
- `http:put:/warmupUrl:{"key":"value"}`: PUT request with its url being `/warmupUrl` its body being `{"key":"value"}`.

### gRPC requests
gRPC requests are in the form `grpc:service/method:message` where `service` and `method` are the names of the gRPC service and method respectively, and `message` is a properly escaped JSON-formatted string.

_Note_: For both HTTP and gRPC requests you can use {today} and {today+n} to obtain the date for today or today +/- n days in YYYY-MM-DD format. For HTTP requests the date templating works for both URLs and body.
