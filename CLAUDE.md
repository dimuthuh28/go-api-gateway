# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go run ./cmd              # run the gateway (needs Redis on :6379 and Kafka on :9092 reachable, or it will hang/error on startup)
go build ./...             # compile everything
go vet ./...                # static analysis
go mod tidy                 # sync go.mod/go.sum
go test ./...                # run tests (none exist yet in this repo)
```

There is no Makefile, Dockerfile, or CI config in this repo — the above are the only tooling entry points.

## Architecture

This is an early-stage API gateway skeleton. `cmd/main.go` is the composition root: it constructs each dependency (Redis client, load balancer, job queue, Kafka producer, metrics server, Fiber app) and wires them together with hardcoded config (addresses, ports, backend list, limits) — there is no config file or env var loading yet.

Package responsibilities under `internal/`:

- **server** (`server.go`) — builds the Fiber app, registers middleware, and defines routes. Currently has a single example route (`/api/service1`) that calls the load balancer and returns a string; it does not actually proxy the HTTP request to the chosen backend.
- **middleware** (`redis_middleware.go`) — `RedisMiddleware` provides two Fiber handlers backed by the same Redis client: `RateLimit` (per-IP counter via `INCR`, no expiry set) and `Cache` (keyed by request path, ignores query params/method/body). Both are applied globally in `server.NewServer`.
- **loadbalancer** (`lb.go`) — `RoundRobin` cycles through a static backend list passed in at construction; no health checking or dynamic backend registration.
- **jobs** (`job.go`, `worker.go`) — a channel-based worker pool (`StartWorkers` spawns N goroutines consuming a `chan Job`). Not currently connected to any HTTP route; `ProcessJob` just prints.
- **kafka** (`producer.go`) — thin wrapper around `segmentio/kafka-go`'s writer. Only used in `main.go` to publish a single startup message; not wired into the request path.
- **metrics** (`metrics.go`) — starts a second HTTP server (`:9000`) exposing the default `promhttp` handler. No gateway-specific metrics (request counts, latencies, etc.) are registered — only Go/process default collectors.

Because most of these pieces (job queue, Kafka, metrics) are started but not yet connected to the gateway's actual request flow, treat them as scaffolding rather than fully integrated features when extending the code — check whether a route/middleware actually invokes a given component before assuming it's live.

Note: `go.mod` requires `github.com/lib/pq` (Postgres driver) but no code in the repo uses it yet.
