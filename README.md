<h1 align="left">
  <img width="500" alt="Mittens" src="images/mittens_logo.svg">
</h1>

# Mittens
Mittens is a -command line- tool that can be used as a warm-up routine against an http application over HTTP (REST) or/and gRPC.

## Design
The application receives a number of command line arguments including the requests that will be sent to warm up the main service. Depending on the format of the requests this will invoke REST or/and gRPC calls.

An overview of the architecture is shown below:

![Mittens architecture.](images/mittens_architecture.png)

This can also run as a linked container, or even as a sidecar in Kubernetes.

## Usage

    mittens [flags]

### Flags

| Flag                               | Type    | Default value    | Description                                                                                                     |
|:-----------------------------------|:--------|:-----------------|:----------------------------------------------------------------------------------------------------------------|
| --concurrency                      | int     | 2                | Number of concurrent requests for warm up                                                                       |
| --exit-after-warmup                | bool    | false            | If warm up process should finish after completion. This is useful to prevent container restarts.                |
| --grpc-headers                     | strings | N/A              | gRPC headers to be sent with warm up requests.                                                                  |
| --grpc-requests                    | strings | N/A              | gRPC request to be sent. Request is in '\<service\>/\<method\>\[:message\]' format. E.g. health/ping:{"key": "value"} |
| --help, -h                         |         |                  | N/A Help for warmup-sidecar                                                                                     |
| --http-headers                     | strings | N/A              | Http headers to be sent with warm up requests.                                                                  |
| --http-requests                    | strings | N/A              | Http request to be sent. Request is in '\<http-method\>:\<path\>\[:body\]' format. E.g. post:/ping:{"key": "value"}   |
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

HTTP requests are in the form `method:path:body` (`body` is optional).
Host and port are taken from `--target-http-host` and
`--target-http-port` flags.

E.g.:
 - `get:/health`: HTTP GET request.
 - `post:/warmupUrl:{"key":"value"}`: POST request with its url being `/warmupUrl` and its body being `{"key":"value"}`.

#### gRPC requests

gRPC requests are in the form `service/method:message`. Host and port are taken from `--target-grpc-host` and
`--target-grpc-port` flags.

## How to build and run
Mittens is written in Go and the simplest way to run it is as a cmd application. It receives a number of command line arguments (also see "Required arguments") including the requests that will be sent to warm up the main service. Depending on the format of the requests this will invoke REST or/and gRPC calls.

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
        
    ./mittens --target-readiness-path=/health --target-insecure=true --target-grpc-port=6565 --timeout-seconds=60 --concurrency=3 --http-requests=get:/hotel/potatoes --grpc-requests=service/method:"{\"foo\":\"bar\", \"bar\":\"foo\"}"

### Docker
#### Build image

To build a Docker image named `mittens`:

    make docker

#### Run container

To run the container:

    docker run mittens:latest --target-readiness-path=/health --target-insecure=true --target-grpc-port=6565 --timeout-seconds=60 --concurrency=3 --http-requests=get:/hotel/potatoes --grpc-requests=service/method:"{\"foo\":\"bar\", \"bar\":\"foo\"}"

_Note_: If you use Docker for Mac you might need to set `targetHost` to `docker.for.mac.localhost`, or `docker.for.mac.host.internal`, or `host.docker.internal` (depending on your version of Docker) so that your container can resolve localhost.

### Kubernetes deployment (as a sidecar)

```yaml
# for versions after 1.9.0 you can use apps/v1
# for versions before 1.6.0 use extensions/v1beta1
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: mittens
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: mittens
    spec:
      containers:
      # primary container goes here
      # - name: foo
      #   image: lorem/ipsum:1.0
      # side car follows
      - name: mittens
        image: mittens:latest
        resources:
          limits:
            memory: 40Mi
            cpu: 20m
          requests:
            memory: 40Mi
            cpu: 20m
        readinessProbe: 
            exec:
              command:
              - cat
              - ready
            initialDelaySeconds: 10
            periodSeconds: 30
        livenessProbe: 
            exec:
              command:
              - cat
              - alive
            initialDelaySeconds: 10
            periodSeconds: 30
        args: 
        - "-concurrency"
        - "2"
        - "-durationSeconds"
        - "120"
        - "-readinessPath"
        - "/ready"
        - "-warmupRequest"
        - "http:get:/tomatoes/"
        - "-warmupRequest"
        - "http:get:/potatoes"
        - "-warmupRequest"
        - "http:post:/hotel/aubergines:{\"foo\":\"bar\"}
        - "-warmupRequest"
        - "grpc:service/method:{\"foo\":\"bar\"}
```

### Notes about warmup duration

Be aware that setting **timeoutForReadinessProbeSeconds** will change how long the warmup routine will run for.

#### Option 1: setting just durationSeconds

```
readinessPath: /ready
durationSeconds: 90
warmupRequest: someRequest
warmupRequest: anotherRequest
```

With these configs the mittens container will start to call _/ready_.
Let's say that your application takes 30 seconds to start (ie, for _/ready_ to start returning 200).
What happens is that after these initial 30 seconds, mittens will start but it will only run for 60 seconds. This is because we already spent 30 seconds waiting for the app to start.
Note that during the warmup _someRequest_ and _anotherRequest_ will be called randomly and not in any particular order.

If the application is not ready after 90 seconds, we skip the warmup routine.

#### Option 2: setting durationSeconds and timeoutForReadinessProbeSeconds

```
readinessPath: /ready
durationSeconds: 90
timeoutForReadinessProbeSeconds: 60
warmupRequest: someRequest
warmupRequest: anotherRequest
```

With these configs the mittens container will start to call _/ready_.
Let's say that your application takes 30 seconds to start (ie, for _/ready_ to start returning 200).
What happens is that after these initial 30 seconds, the warmup will start but unlike the previous example, this time it will run for a full 90 seconds.
Note that during the warmup _someRequest_ and _anotherRequest_ will be called randomly and not in any particular order.

If the application is not ready after the defined 60 seconds, we skip the warmup routine.

## References
* [Mittens at Docker Hub](https://hub.docker.com/r/expediagroup/mittens/)

## Legal

This project is available under the [Apache 2.0 License](http://www.apache.org/licenses/LICENSE-2.0.html).

Copyright 2019 Expedia, Inc.
