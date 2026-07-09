# API Gateway (janus)

A lightweight Go-based API gateway with built-in route proxying, service discovery via config, distributed rate limiting, metrics, tracing, and health monitoring.

## Overview

This repository contains:

- `cmd/janus`: API gateway entrypoint
- `internal/router`: HTTP router and route registration
- `internal/proxy`: reverse proxy, load balancing, circuit breaker, and health checks
- `internal/middleware`: logging, CORS, request IDs, rate limiting, auth, and tracing
- `internal/config`: YAML + env configuration loader
- `internal/redis`: Redis client integration
- `services/users`: sample downstream service
- `services/orders`: sample downstream service
- `docker-compose.yml`: Redis, Prometheus, Grafana stack
- `prometheus/prometheus.yml`: Prometheus scrape configuration
- `grafana/provisioning`: Grafana provisioning files
- `load-test.js`: k6 load test script

## Features

- Route-based proxying to backend services
- Config-driven service instances
- Redis-backed rate limiting
- OpenTelemetry tracing
- Prometheus metrics endpoint
- Graceful shutdown
- Health checks

## Configuration

### `configs/config.yaml`

Example:

```yaml
server:
  port: "8080"

routes:
  - path: /users
    service: users
  - path: /orders
    service: orders

services:
  users:
    instances:
      - http://localhost:8081
      - http://localhost:8082
  orders:
    instances:
      - http://localhost:8083
      - http://localhost:8084
```

### Environment variables

The gateway reads `.env` values for sensitive configuration.

- `POSTGRES_CONN` - optional database connection string
- `REDIS_URL` - Redis connection URL (used for rate limiting)
- `JWT_SECRET` - JWT secret for token validation
- `PORT` - fallback port if not set in YAML

## Running locally

### Start dependencies

Redis is required for rate limiting and metrics.

```bash
docker compose up -d redis
```

### Run sample services

Open two terminals and run:

```bash
cd services/users
go run main.go
```

```bash
cd services/orders
go run main.go
```

> Note: The sample `orders` service currently listens on `:8082` and the sample `users` service listens on `:8081`.

### Run the gateway

```bash
cd cmd/janus
go run main.go
```

The gateway starts on the configured port (`8080` by default).

### Verify

```bash
curl http://localhost:8080/health
curl http://localhost:8080/users
curl http://localhost:8080/orders
```

## Docker Compose stack

To start the observability stack:

```bash
docker compose up -d prometheus grafana
```

Prometheus: http://localhost:9090
Grafana: http://localhost:3000

## Metrics

The gateway exposes Prometheus metrics at:

```bash
http://localhost:8080/metrics
```

## Load testing

Use the included `load-test.js` with k6:

```bash
k6 run load-test.js
```

## Extending routes and services

- Add a new route to `configs/config.yaml` under `routes`
- Add a matching service configuration under `services`
- Each service may define one or more backend instances

## Notes

- The gateway supports JWT auth middleware though it is currently commented out in `internal/router/router.go`
- Update `configs/config.yaml` and `.env` for production settings
- The sample services are minimal and intended for local testing only
