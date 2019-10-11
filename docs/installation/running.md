---
id: running
title: How to Run
---

## Run as a cmd application

You can run the binary executable as follows:

    ./mittens -readinessPath=/ready -durationSeconds=60 -httpTimeoutSeconds=15 -concurrency=3 -httpHeader=X-Forwarded-Proto=https -warmupRequest=http:get:/hotel/potatoes -warmupRequest=http:get:/hotel/tomatoes -warmupRequest="http:post:/hotel/aubergines:{\"foo\":\"bar\"}" -warmupRequest="grpc:service/method:{\"foo\":\"bar\"}" -requestDelayMilliseconds=10

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
        command: "-targetHost=app -readinessPath=/ready -durationSeconds=60 -httpTimeoutSeconds=15 -concurrency=3 -httpHeader=X-Forwarded-Proto=https -warmupRequest=http:get:/hotel/potatoes -warmupRequest=http:get:/hotel/tomatoes -warmupRequest="http:post:/hotel/aubergines:{\"foo\":\"bar\"}" -warmupRequest="grpc:service/method:{\"foo\": \"bar\"}"
    
## Run as a sidecar on Kubernetes

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
        - "-httpHeader"
        - "X-Forwarded-Proto=https"
```
