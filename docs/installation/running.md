---
id: running
title: How to Run
---

The simplest way to run Mittens is as a cmd application. It receives a number of command line arguments (see [Flags](https://github.com/ExpediaGroup/mittens#flags)).
You can also run it as a linked Docker container or even as a sidecar in Kubernetes.

## Run as a cmd application

You can run the binary executable as follows:
        
    ./mittens -target-readiness-http-path=/health -target-grpc-port=6565 -timeout-seconds=60 -concurrency=3 -http-request=get:/hotel/potatoes -grpc-requests=service/method:"{\"foo\":\"bar\", \"bar\":\"foo\"}"

To read the above configs from file:

    ./mittens -config=configs.json

where `configs.json`:

```json
{
  "target-readiness-http-path": "/health",
  "target-grpc-port": 6565,
  "timeout-seconds": 60,
  "concurrency": 3,
  "http-requests": ["get:/hotel/potatoes"],
  "grpc-requests": ["service/method:'{\"foo\":\"bar\", \"bar\":\"foo\"}'"]
}
```

## Run as a linked Docker container

    version: "2"

    services:
    
      foo:
        image: lorem/ipsum:1.0
        ports:
          - "8080:8080"
    
      mittens:
        image: expediagroup/mittens:latest
        links:
          - app
        command: "-target-readiness-http-path=/health -target-grpc-port=6565 -timeout-seconds=60 -concurrency=3 -http-requests=get:/hotel/potatoes -grpc-requests=service/method:{\"foo\":\"bar\", \"bar\":\"foo\"}"

_Note_: If you use Docker for Mac you might need to set the target host (`target-http-host`, `target-grpc-host`) to `docker.for.mac.localhost`, or `docker.for.mac.host.internal`, or `host.docker.internal` (depending on your version of Docker) so that your container can resolve localhost.

## Run as a sidecar on Kubernetes

```yaml
# for versions before 1.9.0 use apps/v1beta1
# for versions before 1.6.0 use extensions/v1beta1
apiVersion: apps/v1
kind: Deployment
metadata:
  name: foo
spec:
  replicas: 2
  selector:
    matchLabels:
      app: foo
  template:
    metadata:
      labels:
        app: foo
    spec:
      containers:
      # primary container goes here
      # - name: foo
      #   image: lorem/ipsum:1.0
      # sidecar follows
      - name: mittens
        image: mittens:latest
        resources:
          limits:
            memory: 50Mi
            cpu: 50m
          requests:
            memory: 50Mi
            cpu: 50m
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
            scheme: HTTP
          initialDelaySeconds: 10
          periodSeconds: 30
        livenessProbe: 
          httpGet:
            path: /live
            port: 8080
            scheme: HTTP
          initialDelaySeconds: 10
          periodSeconds: 30
        args:
        - "-concurrency=3"
        - "-timeout-seconds=60"
        - "-target-readiness-http-path=/health"
        - "-target-grpc-port=6565"
        - "-http-requests=get:/health"
        - "-http-requests=post:/hotel/aubergines:{\"foo\":\"bar\"}"
        - "-grpc-requests=service/method:{\"foo\":\"bar\",\"bar\":\"foo\"}"
```

### gRPC health checks on Kubernetes

Kubernetes does not natively support gRPC health checks.

This leaves you with a couple of options which are documented [here](https://kubernetes.io/blog/2018/10/01/health-checking-grpc-servers-on-kubernetes/).

## Notes about warm-up duration

Be aware that setting **target-readiness-timeout-seconds** will change how long the warmup routine will run for.

### Option 1: setting just -timeout-seconds

```
"-server-probe-readiness-path": /ready
"-timeout-seconds": 90
"-http-requests": someRequest
"-http-requests": anotherRequest
```

With these configs the mittens container will start to call _/ready_.
Let's say that your application takes 30 seconds to start (ie, for _/ready_ to start returning 200).
What happens is that after these initial 30 seconds, mittens will start but it will only run for 60 seconds. This is because we already spent 30 seconds waiting for the app to start.
Note that during the warmup _someRequest_ and _anotherRequest_ will be called randomly and not in any particular order.

If the application is not ready after 90 seconds, we skip the warmup routine.

### Option 2: setting -timeout-seconds and -target-readiness-timeout-seconds

```
"-server-probe-readiness-path": /ready
"-timeout-seconds": 90
"-target-readiness-timeout-seconds": 60
"-http-requests": someRequest
"-http-requests": anotherRequest
```

With these configs the mittens container will start to call _/ready_.
Let's say that your application takes 30 seconds to start (ie, for _/ready_ to start returning 200).
What happens is that after these initial 30 seconds, the warmup will start but unlike the previous example, this time it will run for a full 90 seconds.
Note that during the warmup _someRequest_ and _anotherRequest_ will be called randomly and not in any particular order.

If the application is not ready after the defined 60 seconds, we skip the warmup routine.
