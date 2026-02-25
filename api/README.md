# Distributed Configuration Management API

Distributed configuration management system built with Go. Three services work together to manage, distribute, and apply configuration across a distributed system.

## Architecture

```
                                  +-----------------+
                                  |  Administrator  |
                                  +--------+--------+
                                           |
                                    POST /config
                                           |
                                           v
+----------+    POST /register    +--------+--------+    GET /config     +----------+
|  Agent   | ------------------> |    Controller    | <---------------- |  Agent   |
|          | <------------------ |    (SQLite)      |    (ETag poll)    |          |
+----+-----+   agent_id, poll    +---------+--------+                   +----+-----+
     |                                                                       |
     |  POST /config (forward)                          POST /config         |
     +---------------------------+                +-------------------------+
                                 |                |
                                 v                v
                          +------+----------------+------+
                          |         Worker (6002)        |
                          |     HTTP Server + Asynq      |
                          +------+--------------+--------+
                                 |              |
                           GET /hit       GET /hit/:taskId
                                 |              |
                                 v              v
                          +------+--------------+--------+
                          |         Redis                |
                          |  DB 0: results (1h TTL)      |
                          |  DB 1: asynq task queue      |
                          +------------------------------+
```

- **Controller** (port 6001): Central config management + agent registration. Stores configs in SQLite with versioning and ETag support.
- **Agent**: Registers with Controller, polls for config changes with exponential backoff, forwards updates to Worker. No HTTP server - pure client.
- **Worker** (port 6002): Receives config from Agent, stores in memory. `GET /hit` enqueues async task to Redis, returns task ID (202). Background asynq worker executes the HTTP request and stores result in Redis (1h TTL). `GET /hit/:taskId` retrieves the result.

## Tech Stack

- Go 1.26, Gin HTTP framework, SQLite (CGO)
- Redis + [Asynq](https://github.com/hibiken/asynq) for async task processing
- Clean Architecture with CQRS pattern
- Graceful shutdown, exponential backoff with jitter
- Multi-stage Docker builds

## Project Structure

```
domain/                          # Business rules (no dependencies)
  dto/                           # Data transfer objects
  dto/mapper/                    # Entity -> DTO mappers
  entity/                        # Business entities
  repository/                    # Repository interfaces
  request/                       # Request models
  usecases/                      # Usecase interfaces
  valueobject/                   # Value objects (status constants)

internal/                        # Application layer
  config/                        # Config loading (controller, worker, agent, redis, db)
  delivery/http/controller/      # Controller HTTP handlers + router
  delivery/http/worker/          # Worker HTTP handlers + router
  middleware/                    # Auth + logging middleware
  repository/commands/           # CQRS write implementations (SQLite)
  repository/queries/            # CQRS read implementations (SQLite)
  repository/memory/             # In-memory config store (Agent/Worker)
  usecases/controller/           # Controller command + query usecases
  usecases/worker/               # Worker command + query usecases
  usecases/agent/                # Agent command + query usecases

pkg/                             # Shared packages
  backoff/                       # Exponential backoff with jitter
  cache/                         # Redis client wrapper
  controller/                    # Controller HTTP client
  hit/queue/                     # Asynq task queue (client, processor, result store)
  httpclient/                    # Generic HTTP client wrapper
  response/                      # Standardized API response
  shutdown/                      # Graceful shutdown handler
  worker/                        # Worker HTTP client

cmd/controller/                  # Controller entrypoint
cmd/agent/                       # Agent entrypoint
cmd/worker/                      # Worker entrypoint
deployments/                     # Dockerfiles + docker-compose
test/                            # Table-driven tests (blackbox)
```

## Quick Start

```bash
cp .env.example .env

# Start Redis (skip if already running)
docker run -d -p 6379:6379 redis:7-alpine

# Run all services
make run-all

# Or run individually (3 terminals)
make run-controller    # Terminal 1
make run-worker        # Terminal 2
make run-agent         # Terminal 3
```

## API Endpoints

### Controller (port 6001)

All endpoints require `X-API-Key` header.

| Method | Path             | Description                     |
|--------|------------------|---------------------------------|
| POST   | /register        | Register an agent               |
| POST   | /config          | Create/update config            |
| GET    | /config          | Get latest config (supports ETag) |
| GET    | /config/:version | Get config by version           |

### Worker (port 6002)

| Method | Path           | Description                              |
|--------|----------------|------------------------------------------|
| POST   | /config        | Receive config from agent                |
| GET    | /config        | Get current config                       |
| GET    | /hit           | Enqueue async hit (returns 202 + task_id)|
| GET    | /hit/:taskId   | Get hit result by task ID                |

## Async Hit Flow

```
User: GET /hit
  -> Worker validates config (url must exist)
  -> Generates task UUID
  -> Enqueues to Redis via asynq (worker:hit:execute)
  -> Returns 202 Accepted: {task_id, status: "queued"}

Background (asynq worker):
  -> Picks up task from Redis queue
  -> Executes HTTP GET to configured URL
  -> Stores result in Redis (worker:hit:result:{task_id}, TTL: 1 hour)

User: GET /hit/:taskId
  -> Reads result from Redis
  -> Returns: {status: "completed", status_code: 200, body: "..."}
  -> Or: {status: "pending"} if still processing
  -> Or: {status: "failed", error: "..."} if execution failed
```

### Task Statuses

| Status      | Description                                    |
|-------------|------------------------------------------------|
| `queued`    | Task enqueued to Redis, not yet picked up      |
| `pending`   | Task not found in result store (still processing) |
| `completed` | HTTP request executed successfully             |
| `failed`    | HTTP request failed (error message stored)     |

## Config Format

```json
{
  "data": {
    "url": "https://httpbin.org/get"
  },
  "poll_interval_seconds": 30
}
```

## Build & Test

```bash
make build       # Build all binaries
make test        # Run tests with race detector
make tidy        # go mod tidy
make fix         # go fix ./...
make clean       # Remove binaries
```

## Docker

```bash
# Controller (standalone)
make docker-controller

# Agent + Worker + Redis (combined)
make docker-agent-worker
```

## Environment Variables

| Variable                | Default             | Description                    |
|-------------------------|---------------------|--------------------------------|
| `CONTROLLER_PORT`       | `6001`              | Controller HTTP port           |
| `CONTROLLER_DB_PATH`    | `controller.db`     | SQLite database path           |
| `API_KEY`               | `default-api-key`   | API authentication key         |
| `AGENT_HOSTNAME`        | `agent-01`          | Agent hostname for registration|
| `AGENT_IP`              | `127.0.0.1`         | Agent IP address               |
| `AGENT_PORT`            | `8081`              | Agent port                     |
| `CONTROLLER_URL`        | `http://localhost:6001` | Controller URL for agent   |
| `WORKER_URL`            | `http://localhost:6002` | Worker URL for agent       |
| `POLL_INTERVAL_SECONDS` | `30`                | Agent polling interval         |
| `REQUEST_TIMEOUT_SECONDS`| `10`               | HTTP request timeout           |
| `WORKER_PORT`           | `6002`              | Worker HTTP port               |
| `REDIS_HOST`            | `localhost`          | Redis host                     |
| `REDIS_PORT`            | `6379`              | Redis port                     |
| `REDIS_DB`              | `0`                 | Redis DB for result storage    |
| `ASYNQ_DB`              | `1`                 | Redis DB for asynq queue       |
| `WORKER_RETRY_MAX`      | `3`                 | Max retry for failed tasks     |

## Worker Log Prefixes

| Prefix        | Source                     | Description                           |
|---------------|----------------------------|---------------------------------------|
| `[init]`      | `cmd/worker/main.go`       | Config loading, redis connection      |
| `[http]`      | `cmd/worker/main.go`       | HTTP server lifecycle                 |
| `[asynq]`     | `cmd/worker/main.go`       | Asynq background worker lifecycle     |
| `[queue]`     | `pkg/hit/queue/client.go`  | Queue client connection               |
| `[config]`    | `usecases/worker/`         | Config received from agent            |
| `[enqueue]`   | `usecases/worker/`         | Task creation and enqueue             |
| `[processor]` | `pkg/hit/queue/processor.go`| Task pickup, execution, result save  |
| `[shutdown]`  | `cmd/worker/main.go`       | Graceful shutdown sequence            |
