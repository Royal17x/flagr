# Flagr

**Feature Flags as a Service** — open-source alternative to LaunchDarkly, built for production.

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)
![React](https://img.shields.io/badge/React-18-61DAFB?logo=react)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql)
![Kafka](https://img.shields.io/badge/Kafka-7.5-231F20?logo=apachekafka)

## What is Flagr?

Flagr lets you control feature rollouts without redeploying. Enable or disable features for specific environments instantly through a dashboard or SDK.

**Use cases:**
- A/B testing new features
- Gradual rollouts (10% → 50% → 100%)
- Emergency kill switches
- Environment-specific features (staging vs production)


## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.25, Chi, sqlx, pgx |
| Database | PostgreSQL 16 |
| Cache | Redis 7 (5x latency reduction) |
| Messaging | Apache Kafka (audit log, DLQ) |
| SDK | Go SDK (gRPC + HTTP, local cache) |
| Frontend | React 18, TypeScript, Tailwind CSS |
| Observability | Prometheus, Grafana, Jaeger (OpenTelemetry) |
| Auth | JWT (access 15min + refresh 7d, token rotation) |
| API | REST (Swagger) + gRPC (protobuf) |

## Performance

| Metric | Result |
|--------|--------|
| EvaluateFlag with Redis cache | ~640µs |
| EvaluateFlag without cache | ~3.2ms |
| Cache speedup | **5x** |
| SDK local cache (repeated calls) | ~2µs |
| SDK vs network | **8000x** |

Measured via Go benchmarks on AMD Ryzen 5 7520U.

## Quick Start

**Prerequisites:** Docker, Docker Compose, Go 1.25

```bash
# 1. Clone
git clone https://github.com/Royal17x/flagr
cd flagr

# 2. Configure
cp infra/.env.example infra/.env
# Edit infra/.env — set AUTH_JWT_SECRET (min 32 chars)
# Generate: openssl rand -hex 32

# 3. Start infrastructure
cd infra && docker compose up -d

# 4. Start backend
cd backend
export AUTH_JWT_SECRET=your-secret-here
go run cmd/server/main.go

# 5. Start frontend
cd frontend && npm install && npm run dev
```

Open http://localhost:5173 — register and start creating flags.

## Seed Demo Data

```bash
cd backend
./scripts/seed_demo_data.sh
```

Creates a demo user, 3 feature flags, and an SDK key.

## API

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/api/v1/auth/register` | POST | — | Register (creates org + project + envs) |
| `/api/v1/auth/login` | POST | — | Login |
| `/api/v1/flags` | GET/POST | JWT | List/create flags |
| `/api/v1/flags/{id}/toggle` | POST | JWT | Toggle flag in environment |
| `/api/v1/flags/evaluate` | GET | SDK Key | Evaluate flag (hot path) |
| `/api/v1/sdk-keys` | GET/POST | JWT | Manage SDK keys |
| `/health/live` | GET | — | Liveness probe |
| `/health/ready` | GET | — | Readiness probe (DB + Redis) |
| `/metrics` | GET | — | Prometheus metrics |

Full API docs: http://localhost:8080/swagger/index.html

## gRPC SDK

```go
import flagr "github.com/Royal17x/flagr/sdk"

client, err := flagr.NewClient(
    "http://localhost:8080",
    "sdk-key-xxx",
    flagr.WithCacheTTL(1 * time.Minute),
    flagr.WithDefaultValue(false),
)
defer client.Close()

if client.IsEnabled(ctx, "checkout-v2", projectID, envID) {
    // show new checkout
}
```

First call: ~16ms (network). Repeated calls: ~2µs (local cache).

## Observability

| Service | URL | Description |
|---------|-----|-------------|
| Grafana | http://localhost:3000 | Metrics dashboard |
| Prometheus | http://localhost:9090 | Raw metrics |
| Jaeger | http://localhost:16686 | Distributed traces |
| Kafka UI | http://localhost:8090 | Audit log viewer |
| Swagger | http://localhost:8080/swagger/ | API docs |

**Grafana metrics:**
- HTTP requests/sec and p99 latency
- Cache hit rate
- Flag evaluation by source (cache vs DB)
- Error rate

## Security

- JWT with short-lived access tokens (15min) + rotating refresh tokens (7d)
- Passwords hashed with bcrypt (cost 12)
- SDK keys stored as SHA-256 hashes
- Rate limiting on auth endpoints (10 req/min)
- CORS with origin whitelist
- Security headers (CSP, X-Frame-Options, etc.)
- Input validation on all endpoints

## Project Structure

```
flagr/
├── backend/
│   ├── cmd/server/        # Entry point
│   ├── internal/
│   │   ├── domain/        # Business entities
│   │   ├── repository/    # PostgreSQL queries
│   │   ├── service/       # Business logic
│   │   ├── handler/       # HTTP handlers
│   │   ├── grpc/          # gRPC server
│   │   ├── middleware/     # Auth, rate limit, CORS
│   │   ├── metrics/        # Prometheus
│   │   └── tracing/        # Jaeger/OpenTelemetry
│   ├── pkg/
│   │   ├── kafka/          # Producer, consumer, DLQ
│   │   ├── redis/          # Connection pool
│   │   └── postgres/       # Connection pool
│   ├── migrations/         # SQL migrations
│   └── scripts/            # Seed, test, load
├── sdk/                    # Go SDK (separate module)
├── frontend/               # React dashboard
└── infra/                  # Docker Compose, Prometheus
```

## Testing

```bash
cd backend

# Unit tests (no Docker needed)
go test ./internal/service/... -v

# Integration tests (requires Docker)
go test ./internal/repository/... -v

# Benchmarks
go test ./internal/service/... -bench=BenchmarkEvaluateFlag -benchmem -count=3

# Smoke tests
./scripts/test_api.sh
```

## Roadmap

- [ ] Percentage rollouts (10% of users)
- [ ] Targeting rules (enable for `country=RU`, `plan=pro`)
- [ ] Deploy to Yandex Cloud / VK Cloud
- [ ] CI/CD GitHub Actions (build + test, deploy on push)
- [ ] AI agent for flag rollout recommendations

## Author

Built by [Royal17x](https://github.com/Royal17x) as a portfolio project demonstrating production Go development.
