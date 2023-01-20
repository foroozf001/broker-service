# Broker Service

## Overview

The broker service provides an endpoint for storing logs in a mongodb back-end via the logger service. It accepts HTTP REST on port 8080 and uses gRPC on the back-end. In addition to this, it will expose custom metrics on an endpoint.

## Getting started

- Spin up the stack by executing `make up`
- Create artificial load using `hey` command-line tool

```bash
hey -z 5m -q 5 -m POST -H "Content-Type: application/json" -d '{"name": "example_log", "data": "this is an example log"}' http://localhost:8080/grpc/log
```

- View documents via mongo-express client at http://localhost:8081
- See custom prometheus metrics by calling the `/api/metrics` endpoint
  - `http_requests_total`
  - `response_status`
  - `http_response_time_seconds`

```bash
curl localhost:8080/grpc/metrics 
```
