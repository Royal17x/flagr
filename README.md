# Flagr — Feature Flags as a Service

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8.svg)

Deploy features without deploying code.  
Built for the Russian market — no vendor lock-in, no sanctions risk,
no data leaving the country.

## Why Flagr?

LaunchDarkly is an American company. Unleash Cloud is Norwegian.  
Both can cut off your account tomorrow — and in 2022, many Russian teams
learned this the hard way.

Flagr is different:

- 🇷🇺 **Data stays in Russia** — 152-ФЗ compliant by design
- 🔒 **No geopolitical risk** — self-hosted or deployed to Russian cloud
- 💰 **Honest pricing in rubles** — no currency risk, no per-seat shock
- ⚡ **< 1ms flag evaluation** — Redis-powered, built for production load
- 🛠 **Go SDK** — drop-in integration, no heavy dependencies

## Architecture

> Coming soon — diagram after core services are implemented.

## Tech Stack

| Layer | Technology |
|---|---|
| HTTP API | Go + Chi |
| RPC | gRPC |
| Database | PostgreSQL |
| Cache | Redis |
| Audit log | Kafka |
| Metrics | Prometheus + Grafana |
| Tracing | Jaeger |
| Frontend | React 19 + TypeScript + shadcn/ui |
| Infra | Docker + GitHub Actions |

## Quick Start
```bash
git clone https://github.com/ТВО_НИК/flagr.git
cd flagr/infra
docker-compose up -d
cd ../backend
go run ./cmd/server/
```

Dashboard: http://localhost:3000  
API: http://localhost:8080  
Swagger: http://localhost:8080/swagger/index.html

## Roadmap

- [x] Project setup & infrastructure
- [ ] Core flag CRUD API
- [ ] Flag evaluation engine (<1ms via Redis)
- [ ] Go SDK
- [ ] Web dashboard
- [ ] Audit log via Kafka
- [ ] Percentage rollouts & targeting rules
- [ ] AI agent: auto-rollout recommendations
- [ ] Billing & multi-tenancy

## License

MIT