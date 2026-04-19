# OHE — Claude Code Workspace

## Project
OHE (Observability Holistic Engine) — Go-based MLOps monitoring platform.
Module: `github.com/benfradjselim/ohe`
Working directory: `workdir/` (all Go code lives here)

## Current State
- Released: v5.0.0 (tag pushed, branch v5.1.0 is active dev)
- All 18 packages pass ≥60% coverage gate
- CI: `.github/workflows/ci.yml` — build + vet + race test + coverage gate + lint

## Key Commands
```bash
# From workdir/
go build ./...                          # build check
go test -race -timeout=120s ./...       # full suite with race detector
go test -cover -timeout=120s ./...      # coverage summary
go vet ./...                            # static analysis
golangci-lint run --timeout=5m          # lint (same as CI)
go test -run TestFoo ./internal/pkg/... # single test
```

## Architecture
- `cmd/ohe/` — main binary entry point
- `internal/api/` — HTTP handlers, middleware, router (Gin)
- `internal/storage/` — Badger KV, OrgStore (multi-tenant isolation `o:{orgID}:key`)
- `internal/predictor/` — CAILR (CA+ILR dual-scale), ARIMA, HoltWinters, MAD, Ensemble
- `internal/analyzer/` — dissipative fatigue topology (ILR clusters)
- `internal/orchestrator/` — wires everything; Config struct with yaml tags
- `internal/grpcserver/` — agent gRPC ingest (ohe.v1.AgentService/Ingest)
- `internal/receiver/` — DogStatsD UDP :8125, OTLP HTTP/gRPC
- `internal/billing/` — UsageEvent ring buffer + webhook flush
- `internal/correlator/` — metric correlation engine
- `internal/eventbus/` — pub/sub event bus
- `internal/plugin/` — plugin sandbox
- `internal/vault/` — Vault integration

## Config Defaults (v5.0.0)
```yaml
predictor:
  stable_window: 60m
  burst_window: 5m
  rupture_threshold: 3.0
fatigue:
  r_threshold: 0.3
  lambda: 0.05
```

## Remaining Roadmap (v5.1.0)
- **#13 Go SDK** — typed client wrapping the REST API, publish to pkg.go.dev
- **#13 Python SDK** — `pip install ohe-sdk`, mirrors Go SDK
- Coverage stretch goals: api (61% → 70%), orchestrator (64% → 75%)

## Storage Key Schema
```
o:{orgID}:m:{host}:{metric}:{ts}   metrics
o:{orgID}:k:{host}:{kpi}:{ts}      KPIs
o:{orgID}:a:{id}                   alerts
o:{orgID}:d:{id}                   dashboards
o:{orgID}:ds:{id}                  datasources
o:{orgID}:nc:{id}                  notification channels
o:{orgID}:slo:{id}                 SLOs
o:{orgID}:ak:{id}                  API keys
o:{orgID}:l:{service}:{ts}         logs
o:{orgID}:sp:{traceID}:{spanID}    spans
```

## Conventions
- Errors: always `fmt.Errorf("context: %w", err)`
- Interfaces: accept interfaces (StorageBackend, MetricSink), return structs
- Tests: table-driven, `_test` package suffix for black-box, `-race` always
- No comments unless the WHY is non-obvious
- Version string in `internal/api/handlers.go` const `version`
