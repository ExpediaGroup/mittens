---
id: introduction
title: Introduction
---

<h1 align="left">
  <img width="300" alt="Mittens" src="../assets/mittens_logo.svg">
</h1>

Mittens is a warm-up routine for http applications over REST and gRPC.

## The Problem

When an app starts (as part of a deploy, redeploy, or restart) the very first requests are expected to be slow. 
Although there are many reasons that contribute to this and to various extents (e.g. class loading, caches, SSL handshake), having a warmup routine that sends dummy requests before the app receives any traffic is a common practice to reduce the impact of the so called "cold-start" issue.

A warm-up routine can be particularly useful in the following cases:
- For apps which are deployed on Kubernetes where pods are created/destroyed quite often and as a result the app (re)starts frequently.
- For apps exposing gRPC endpoints, since the first gRPC requests are expected to face significantly high latency until the gRPC channel is warm.

Although many tools exist for making either HTTP/REST or gRPC calls none of these supports both REST and gRPC requests, and it’s not trivial to run these in Kubernetes or as linked containers.

## The Solution

Mittens is a simple tool that can be used as a warm-up routine for http applications over REST and gRPC.

Its main features are summarised below:
- Sends requests continuously for X seconds
- Supports REST and gRPC
- Supports HTTP and gRPC headers
- Supports concurrent requests
- Provides files or/and endpoints that can be used as liveness/readiness probes in Kubernetes

Mittens can run as a standalone [command-line tool](https://expediagroup.github.io/mittens/docs/installation/running#run-as-a-cmd-application), as a [linked Docker container](https://expediagroup.github.io/mittens/docs/installation/running#run-as-a-linked-docker-container), or even as a [sidecar in Kubernetes](https://expediagroup.github.io/mittens/docs/installation/running#run-as-a-sidecar-on-kubernetes).
