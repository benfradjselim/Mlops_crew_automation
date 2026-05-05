# Ruptura

<p align="center">
  <img src="https://img.shields.io/badge/version-6.6.0-0069ff?style=for-the-badge" alt="v6.6.0">
  <img src="https://img.shields.io/badge/go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go 1.21+">
  <img src="https://img.shields.io/badge/license-Apache%202.0-green?style=for-the-badge" alt="Apache 2.0">
  <img src="https://img.shields.io/badge/kubernetes-native-326CE5?style=for-the-badge&logo=kubernetes&logoColor=white" alt="Kubernetes Native">
  <img src="https://img.shields.io/badge/tests-passing-brightgreen?style=for-the-badge" alt="Tests Passing">
</p>

<p align="center">
  <b>The Predictive Action Layer for Cloud-Native Infrastructure.</b><br>
  Ruptura detects workload ruptures before they cause outages — and acts on them automatically.
</p>

---

## What Ruptura Does

Traditional observability tells you what broke. Ruptura tells you **what is about to break** — and triggers the right action before users feel it.

| Traditional Observability | Ruptura |
|--------------------------|-----------|
| Threshold alerts fire after the fact | Fused Rupture Index™ detects divergence **hours early** |
| Global thresholds that alert on batch jobs | **Adaptive per-workload baselines** — no false alarms |
| You define rules per metric | Signals are relative to each workload's normal behavior |
| Manual incident response | Tier-1 actions (scale, restart, rollback) fire automatically |
| 5+ tools: Prom + Grafana + AM + Loki + PD | **One binary**, one `kubectl apply` |
| "CPU 78%" — what does that mean? | **Narrative explain**: "payment-api has been accumulating fatigue for 72h — cascade from payment-db" |

---

## Core Concepts

### Fused Rupture Index™

Ruptura combines metric, log, and trace signals into a single rupture index per workload:

```
FusedR = f(metricR, logR, traceR)

  metricR  = CA-ILR dual-scale slope ratio on raw metrics
  logR     = burst detector ratio when error/warn rate > 3σ baseline
  traceR   = span error rate + P99 latency vs rolling baseline
```

| FusedR | State | Action |
|--------|-------|--------|
| < 1.0 | Stable | None |
| 1.0–1.5 | Elevated | None |
| 1.5–3.0 | Warning | Tier-3 (human alert) |
| 3.0–5.0 | Critical | Tier-2 (suggested action) |
| ≥ 5.0 | Emergency | Tier-1 (automated action) |

### 10 Composite KPI Signals

`stress` · `fatigue` · `mood` · `pressure` · `humidity` · `contagion` · `resilience` · `entropy` · `velocity` · `health_score`

Each maps multiple raw metrics to a single interpretable 0–1 index with published formulas. `health_score` is a 0–100 composite of the primary signals.

### Adaptive Per-Workload Baselines

After 24h of observation per workload, all thresholds become relative to that workload's own Welford baseline. A batch job that always runs at 90% CPU is never "stressed." An API that normally runs at 10% and spikes to 40% is flagged.

### Calibration Warm-Up

During the first 24h (`status: "calibrating"`), signals are computed and stored but rupture predictions and actions are suppressed — the baseline isn't ready yet. Every snapshot carries `calibration_progress` (0–100) and `calibration_eta_minutes` so you always know where you stand. Use `ruptura-sim` to demo immediately without waiting.

### HealthScore Trend Forecast

Once calibrated, every snapshot includes a linear projection of where HealthScore is heading:

```json
"health_forecast": { "trend": "degrading", "in_15min": 51.2, "in_30min": 38.7, "critical_eta_minutes": 28 }
```

Turns "your score is 54" into "you have 28 minutes."

### Rupture Fingerprinting

At every confirmed rupture (FusedR ≥ 3.0), Ruptura captures an 11-dimensional KPI vector. On subsequent queries, the current state is compared against all past fingerprints using cosine similarity. A match ≥ 0.85 surfaces as `pattern_match` in the response — with the matched rupture ID, similarity score, and resolution note from last time.

### Business Signals

Three business-layer signals are included in every snapshot:

| Signal | Meaning |
|--------|---------|
| `slo_burn_velocity` | `current_error_rate / allowed_error_rate` — > 1.0 means burning error budget too fast |
| `blast_radius` | Downstream services that depend on this workload (from trace topology) |
| `recovery_debt` | Near-miss count (FusedR 2–3, recovered without rupturing) in the last 7 days |

### Narrative Explain

```
GET /api/v2/explain/{id}/narrative
```

Returns a structured English narrative — not a JSON of numbers:

> "payment-api has been accumulating fatigue for 72h (fatigue 0.81, burnout threshold 0.80). Contagion wave from payment-db propagated via the payment-api→payment-db edge and pushed FusedR from 1.8 to 4.2 in 18 minutes. This is a cascade rupture, not an isolated spike. Recommended action: scale payment-api by 2 replicas."

### Topology-Based Contagion

When OTLP trace spans are ingested, Ruptura builds a real service dependency graph. Contagion is computed from actual edge error rates weighted by call volume — not a `cpu × errors` proxy.

### Edition Gate

`RUPTURA_EDITION=community` (default) — action recommendations are visible read-only; `POST .../approve` returns 402. Set `RUPTURA_EDITION=autopilot` to enable full Tier-1 auto-execution and manual approval.

### Per-Workload Signal Weights

Override HealthScore weights per namespace or workload via `POST /api/v2/config/weights` at runtime, or via `RUPTURA_WORKLOAD_WEIGHTS` (JSON) / `workloadWeights:` in Helm values. A latency-sensitive API and a batch job should not share the same `stress` weight.

---

## Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                         ruptura                              │
│                                                              │
│  Ingest ──► Metric/Log/Trace pipelines ──► Fusion            │
│     │              │                          │              │
│  Prom rw        Analyzer                  FusedR             │
│  OTLP HTTP      (10 KPI signals)              │              │
│  DogStatsD      Adaptive baselines        ActionEngine       │
│     │           Calibration warm-up      (K8s/Webhook/PD)    │
│     │           Topology contagion        Edition gate       │
│     │           Narrative explain             │              │
│  Correlator     Fingerprinting           NarrativeExplain    │
│  (BurstDetector,Business signals              │              │
│   TopologyBuilder) HealthScore forecast  REST API v2         │
│                                          WorkloadRef routes  │
│              BadgerDB embedded storage (single binary)       │
└──────────────────────────────────────────────────────────────┘
```

**Single binary.** No external database. No sidecars. One `kubectl apply` or `helm install`.

---

## Quick Start

### Kubernetes (recommended)

```bash
# Option 1: Helm (recommended)
helm install ruptura ../helm \
  --namespace ruptura-system \
  --create-namespace \
  --set apiKey=$(openssl rand -hex 32)

# Option 2: kustomize manifests
kubectl apply -f deploy/

# Port-forward the API
kubectl port-forward svc/ruptura 8080:80 -n ruptura-system

# Verify
curl http://localhost:8080/api/v2/health
```

### Docker

```bash
docker run -d \
  --name ruptura \
  -p 8080:8080 \
  -p 4317:4317 \
  -v ruptura-data:/var/lib/ruptura/data \
  -e RUPTURA_API_KEY=$(openssl rand -hex 32) \
  ghcr.io/benfradjselim/ruptura:6.6.0
```

| Port | Purpose |
|------|---------|
| 8080 | REST API v2 + Prometheus metrics scrape |
| 4317 | OTLP ingest (metrics, logs, traces) |

### Build from Source

```bash
git clone https://github.com/benfradjselim/ruptura.git
cd ruptura/workdir
go build -o ruptura ./cmd/ruptura
./ruptura --port=8080 --otlp-port=4317 --storage=/tmp/ruptura-data
```

---

## Sending Metrics

### Prometheus remote_write

```yaml
# prometheus.yml
remote_write:
  - url: http://ruptura:8080/api/v2/write
    authorization:
      credentials: <your-api-key>
```

### OTLP (OpenTelemetry)

```yaml
# otel-collector-config.yaml
exporters:
  otlphttp:
    endpoint: http://ruptura:4317
    headers:
      Authorization: "Bearer <your-api-key>"
```

Ruptura extracts `k8s.namespace.name`, `k8s.deployment.name`, `k8s.statefulset.name` from OTLP resource attributes and groups signals by **Kubernetes workload** — not by node/host. Multiple pods from the same Deployment are merged into a single workload health view.

---

## API

All endpoints at `/api/v2/`. Auth via `Authorization: Bearer <api-key>`.

```
# Health & readiness
GET  /api/v2/health
GET  /api/v2/ready
GET  /api/v2/metrics              Prometheus self-metrics (scrape endpoint)

# Ingest
POST /api/v2/write                Prometheus remote_write
# OTLP → send to port 4317 (separate OTLP HTTP server)

# Rupture index
GET  /api/v2/ruptures                              all workloads
GET  /api/v2/rupture/{namespace}/{workload}        WorkloadRef (primary)
GET  /api/v2/rupture/{host}                        legacy host-based

# KPI signals (stress/fatigue/mood/pressure/humidity/contagion/resilience/entropy/velocity/health_score)
GET  /api/v2/kpi/{signal}/{namespace}/{workload}   WorkloadRef
GET  /api/v2/kpi/{signal}/{host}                   legacy

# Forecast
POST /api/v2/forecast
GET  /api/v2/forecast/{metric}/{namespace}/{workload}

# Anomalies
GET  /api/v2/anomalies
GET  /api/v2/anomalies/{host}

# Actions (approve/reject T2; T1 auto in autopilot edition)
GET  /api/v2/actions
POST /api/v2/actions/{id}/approve      ← 402 in community edition
POST /api/v2/actions/{id}/reject
POST /api/v2/actions/emergency-stop

# Maintenance windows (suppress alerts during deploys)
POST /api/v2/suppressions         { workload, start, end, [signals] }
GET  /api/v2/suppressions
DELETE /api/v2/suppressions/{id}

# Simulation (demo without real incidents)
POST /api/v2/sim/inject           { pattern, host, duration_minutes }

# Signal weight configuration (per-workload HealthScore tuning)
GET  /api/v2/config/weights
POST /api/v2/config/weights       [{ selector, stress, fatigue, ... }]

# Explainability
GET  /api/v2/explain/{id}
GET  /api/v2/explain/{id}/formula
GET  /api/v2/explain/{id}/pipeline
GET  /api/v2/explain/{id}/narrative    ← human-readable explanation
```

---

## Configuration

Environment variables (no config file required for basic use):

| Variable | Default | Description |
|----------|---------|-------------|
| `RUPTURA_API_KEY` | _(empty, auth disabled)_ | Bearer token for all API requests |
| `RUPTURA_EDITION` | `community` | `community` (read-only actions) or `autopilot` (full execution) |
| `RUPTURA_WORKLOAD_WEIGHTS` | _(empty)_ | JSON array of `SignalWeights` for per-workload HealthScore tuning |
| `RUPTURA_INGEST_RPS` | `1000` | Token-bucket rate limit on ingest |
| `RUPTURA_LOG_LEVEL` | `info` | Log verbosity (debug/info/warn/error) |

CLI flags: `--port=8080 --otlp-port=4317 --storage=/var/lib/ruptura/data`

---

## Prometheus Self-Metrics

Scrape at `GET /api/v2/metrics`.

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `ruptura_kpi` | gauge | `namespace,kind,workload,signal` | All 10 KPI signals per workload |
| `rpt_rupture_index` | gauge | `host,metric,severity` | Per-metric rupture index |
| `rpt_time_to_failure_seconds` | gauge | `host,metric` | Estimated time to failure |
| `rpt_kpi_healthscore` | gauge | `host` | HealthScore (legacy host label) |
| `rpt_actions_total` | counter | `type,tier,outcome` | Actions fired by tier |
| `rpt_ingest_samples_total` | counter | `source` | Ingested samples by source |
| `rpt_memory_bytes` | gauge | — | Memory usage |
| `rpt_uptime_seconds` | gauge | — | Uptime |

For Grafana: import `deploy/grafana/dashboards/ruptura_overview.json` or use the provisioning config in `deploy/grafana/provisioning.yaml`.

---

## Client library (Go)

The embeddable Go client is in `pkg/client`:

```go
import "github.com/benfradjselim/ruptura/pkg/client"

c := client.New("http://ruptura:8080", client.WithAPIKey("your-key"))
rupture, _ := c.RuptureIndex(ctx, "default", "payment-api")
```

For REST-only access from any language, use the API directly with a Bearer token — see [API reference](../docs/v6.0.0/SPECS.md).

---

## Kubernetes Operator

```yaml
apiVersion: ruptura.io/v1alpha1
kind: RupturaInstance
metadata:
  name: production
  namespace: ruptura-system
spec:
  image: ghcr.io/benfradjselim/ruptura:6.6.0
  port: 8080
  storageSize: 20Gi
  apiKey:
    secretRef: ruptura-api-key
```

The operator reconciles Deployment + Service + PVC per `RupturaInstance`. See `ohe/operator/`.

---

## Development

```bash
go build ./...
go test -race -timeout=120s ./...
go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | grep total
helm lint helm/
```

---

## Changelog

### v6.6.0 — 2026-05-05
- **Per-workload signal weight tuning**: `POST/GET /api/v2/config/weights` for runtime override. `RUPTURA_WORKLOAD_WEIGHTS` JSON env var for Helm bootstrap. Selector syntax: exact, `ns/*`, or `*`. Weights auto-normalised to 1.0.

### v6.5.0 — 2026-05-05
- **Edition gate**: `RUPTURA_EDITION=community|autopilot`. `POST .../approve` returns 402 in community — recommendations stay read-only. Full action execution in autopilot.

### v6.4.0 — 2026-05-05
- **Rupture fingerprinting**: 11-dimensional KPI vector per confirmed rupture, cosine similarity ≥ 0.85 → `pattern_match` in every rupture response.
- **Business signal layer**: `slo_burn_velocity`, `blast_radius`, `recovery_debt` in every snapshot's `business` block.

### v6.3.0 — 2026-05-04
- **Calibration warm-up**: `status`, `calibration_progress`, `calibration_eta_minutes` in every snapshot.
- **HealthScore trend forecast**: `health_forecast` block — OLS slope → `in_15min`, `in_30min`, `critical_eta_minutes`.
- **`ruptura-sim`**: four simulation patterns via `POST /api/v2/sim/inject` for local demos.

### v6.2.2 — 2026-04-30
- GAP-04 closed: anomaly REST endpoints (`GET /api/v2/anomalies`, `/api/v2/anomalies/{host}`)
- Dead duplicate `internal/predictor/anomaly_engine.go` removed
- Release workflow fixed; docs updated with correct API key env var

### v6.2.1 — 2026-04-30
- `FusedRuptureIndex` exposed in API response + integration test coverage
- Grafana dashboard corrected: queries `ruptura_kpi{signal="fused_rupture_index"}`, 6 panels, workload variable

### v6.2.0 — 2026-04-30
- **WorkloadRef treatment unit**: OTLP extracts `k8s.namespace.name` / `k8s.deployment.name` / etc. Multiple pods from one Deployment merged into one workload view.
- **Adaptive per-workload baselines**: After 96 observations (~24h), thresholds become relative z-score deviations from Welford baseline.
- **Narrative explain** at `/api/v2/explain/{id}/narrative`
- **Topology-based contagion**: real service edges from trace spans; falls back to proxy when no edges exist
- **Maintenance windows**: POST/GET/DELETE `/api/v2/suppressions`
- All 10 KPI signals wired end-to-end: stress, fatigue, mood, pressure, humidity, contagion, resilience, entropy, velocity, health_score
- HealthScore formula: additive penalty model (was multiplicative — collapsed too aggressively)
- Fusion: metricR + logR + traceR → FusedR fully wired
- BadgerDB: FlushSnapshots() on SIGTERM (no data loss on graceful shutdown)
- Token-bucket rate limiter on ingest (default 1000 req/s)
- 37 packages pass `go test -race ./...`

### v6.1.0 — 2026-04-27
- Real gRPC ingest server (port 9090)
- NATS/Kafka eventbus (JetStream + franz-go)
- Adaptive ensemble weighting (online MAE-based, 60s update)
- Kubernetes operator (RupturaInstance CRD)

---

## Roadmap

```
v6.2.x ✅  Fused Rupture Index · workload-level signals · adaptive baselines · narrative explain
v6.1.0 ✅  gRPC ingest · NATS/Kafka eventbus · adaptive ensemble · K8s operator
v7.0.0 ⏳  ruptura-ctl CLI · web dashboard v2 · multi-tenant (X-Org-ID)
```

---

## License

Apache 2.0 — see [LICENSE](../LICENSE)
