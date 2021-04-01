---
id: running
title: How to Run
---

The simplest way to run Mittens is as a cmd application. It receives a number of command line arguments (see [flags](https://expediagroup.github.io/mittens/docs/about/getting-started#flags)).
You can also run it as a linked Docker container or even as a sidecar in Kubernetes.

## Run as a cmd application

You can run the binary executable as follows:
        
    ./mittens -target-readiness-http-path=/health -target-grpc-port=6565 -max-duration-seconds=60 -concurrency=3 -http-request=get:/hotel/potatoes -grpc-requests=service/method:"{\"foo\":\"bar\", \"bar\":\"foo\"}"

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
        command: "-target-readiness-http-path=/health -target-grpc-port=6565 -max-duration-seconds=60 -concurrency=3 -http-requests=get:/hotel/potatoes -grpc-requests=service/method:{\"foo\":\"bar\", \"bar\":\"foo\"}"

_Note_: If you use Docker for Mac/Windows you might need to set the target host (`target-http-host`, `target-grpc-host`) to `host.docker.internal` so that your container can resolve localhost. If you use an older version of Docker (< 18.03), the value will depend on your Operating System, e.g. `docker.for.mac.host.internal` or `docker.for.win.host.internal`. For `target-http-host`, you need to prefix the host with the scheme, e.g. `http://host.docker.internal`.

## Run as a sidecar on Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: foo
spec:
  replicas: 1
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
        image: expediagroup/mittens:latest
        resources:
          limits:
            memory: 50Mi
            cpu: 50m
          requests:
            memory: 50Mi
            cpu: 50m
        readinessProbe:
          exec:
            command:
            - "cat"
            - "ready"
          initialDelaySeconds: 10
          periodSeconds: 30
        livenessProbe:
          exec:
            command:
            - "cat"
            - "alive"
          initialDelaySeconds: 10
          periodSeconds: 30
        args:
        - "--concurrency=3"
        - "--max-duration-seconds=60"
        - "--target-readiness-http-path=/health"
        - "--target-grpc-port=6565"
        - "--http-requests=get:/health"
        - "--http-requests=post:/hotel/aubergines:{\"foo\":\"bar\"}"
        - "--grpc-requests=service/method:{\"foo\":\"bar\",\"bar\":\"foo\"}"
```

#### Using Config Maps for complex POST queries

If your aplication is using complex POST requests, you can move them to a separate file.
To do that in Kubernetes, you can make use of ConfigMaps and define the body of your requests there, and then use a volume mount to make them accessible to Mittens.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: foo
spec:
  replicas: 1
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
        image: expediagroup/mittens:latest
        resources:
          limits:
            memory: 50Mi
            cpu: 50m
          requests:
            memory: 50Mi
            cpu: 50m
        readinessProbe:
          exec:
            command:
            - "cat"
            - "ready"
          initialDelaySeconds: 10
          periodSeconds: 30
        livenessProbe:
          exec:
            command:
            - "cat"
            - "alive"
          initialDelaySeconds: 10
          periodSeconds: 30
        args:
        - "--concurrency=3"
        - "--max-duration-seconds=60"
        - "--target-readiness-http-path=/health"
        - "--target-grpc-port=6565"
        - "--http-requests=get:/health"
        # Inlined request body
        - "--http-requests=post:/hotel/aubergines:{\"foo\":\"bar\"}"
        # Request body comes from a json file (defined in a config map)
        - "--http-requests=post:/hotel/aubergines:file:/mittens/req_1.json"
        - "--http-headers=content-type:application/json"
        volumeMounts:
          - name: mittens-config
            mountPath: /mittens
    volumes:
        - name: mittens-config
          configMap:
            name: mittens-config
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: mittens-config
data:
  req_1.json: |
    {
      "name": "foo",
      "age": 123
    }

```

### gRPC health checks on Kubernetes

Kubernetes does not natively support gRPC health checks.

This leaves you with a couple of options which are documented [here](https://kubernetes.io/blog/2018/10/01/health-checking-grpc-servers-on-kubernetes/).

## gRPC Server Reflection is needed

Mittens uses `grpcurl` to call the gRPC target server, which needs gRPC reflection to be enabled on the server.

gRPC Server Reflection assists clients in runtime construction of requests without having stub information precompiled into the client.
For more info see [here](https://github.com/grpc/grpc/blob/master/doc/server-reflection.md).

## Note about warm-up duration

`-max-duration-seconds` includes the time needed for your application to start.
Let's say that your application takes 30 seconds to start (ie, for _/ready_ to start returning 200).
What happens is that after these initial 30 seconds, mittens will start but it will only run for 60 seconds. This is because we already spent 30 seconds waiting for the app to start.

If the application is not ready after 90 seconds, we skip the warmup routine.
