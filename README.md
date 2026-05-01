# Shop — E-commerce Backend

Backend for an e-commerce platform built with Go. Event-driven microservice architecture with async communication via Kafka, Redis caching, and full observability stack.

## Architecture

The system consists of three independent services that communicate through Kafka:

```
Client
  └── HTTP REST API (api)
        ├── publishes events → Kafka
        │     ├── catalog_event  → changeDataCapture (CDC logging)
        │     └── cart_event     → cacheInvalidator (Redis invalidation)
        └── Outbox Worker → guaranteed event delivery
```

| Service | Responsibility |
|---|---|
| `api` | Main HTTP service. Auth, catalog, cart, orders, payments |
| `cacheInvalidator` | Listens to Kafka, invalidates Redis cache on data changes |
| `changeDataCapture` | Listens to Kafka, logs all data mutations for audit trail |

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.22+ |
| Database | PostgreSQL 15 |
| Cache | Redis 7.2 |
| Message Broker | Apache Kafka (Confluent) |
| Containerization | Docker / Docker Compose |
| Observability | Grafana + Loki + Promtail |
| CI/CD | GitLab CI |
| Linter | golangci-lint |

## Key Design Decisions

**Outbox Pattern** — events are first written to an `outbox` table in PostgreSQL within the same transaction as the business operation. A background worker then reads from the outbox and publishes to Kafka. This guarantees no event is lost even if Kafka is temporarily unavailable.

**Cache Invalidation via CDC** — instead of invalidating cache directly in handlers, the `api` publishes events to Kafka. The `cacheInvalidator` service consumes these events and clears the relevant Redis keys. This keeps the main service decoupled from cache logic.

**Graceful Shutdown** — all three services listen for `SIGTERM`/`SIGINT` via `signal.NotifyContext`. On shutdown, the HTTP server stops accepting new connections and waits up to 10 seconds for in-flight requests to complete before closing DB, Kafka, and Redis connections.

**Envelope format for Kafka messages:**
```json
{
  "type": "product_create",
  "payload": { ... },
  "meta_date": { "id_user": 1, "role_user": "admin" }
}
```

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Go 1.22+

### Run locally

1. Clone the repository:
```bash
git clone https://github.com/patato8984/Shop.git
cd Shop
```

2. Create secrets:
```bash
mkdir secrets
echo "your-jwt-secret" > secrets/jwt_key.txt
echo "your-bank-secret" > secrets/bank_key.txt
```

3. Start all services:
```bash
docker compose up --build
```

4. The API will be available at `http://localhost:8080`

### Observability

| Service | URL |
|---|---|
| Grafana | http://localhost:3000 (admin/admin) |
| Loki | http://localhost:3100 |

Logs from all containers are collected by Promtail and shipped to Loki. Grafana is pre-configured to auto-refresh dashboards every 1 second.

## Project Structure

```
.
├── cmd/
│   ├── api/                  # Main HTTP service entrypoint
│   ├── cacheInvalidator/     # Cache invalidation worker
│   └── changeDataCapture/    # CDC audit logging worker
├── internal/
│   ├── app/                  # DI wiring, config loading
│   ├── infra/
│   │   ├── kafka/            # Producer & consumer wrappers
│   │   ├── postgres/         # DB connection, migrations
│   │   └── redis/            # Redis client init
│   ├── modules/
│   │   └── catalog/          # Product & SKU domain logic
│   └── shared/
│       ├── events/           # Event publisher, Kafka listener
│       └── logger/           # zap logger init
├── configs/                  # YAML configs per service
├── build/
│   ├── loki/                 # Loki config
│   └── promtail/             # Promtail scrape config
├── secrets/                  # Docker secrets (not committed)
├── docker-compose.yml
├── Dockerfile
└── .golangci.yml
```

## CI/CD

GitLab CI pipeline runs on every push:

1. **lint** — `golangci-lint` with `errcheck`, `govet`, `staticcheck`, `gocritic`, `unused`
2. **build** — builds Docker image and pushes to GitLab Container Registry

## API Endpoints

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/v1/webhook/bank` | Bank payment webhook |
| ... | `/api/v1/...` | Auth, catalog, cart, orders (registered via `RegisterAuthRoutes`) |
