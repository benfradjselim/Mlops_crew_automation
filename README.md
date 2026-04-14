# Observability Holistic Engine (OHE) v4.0.0

> **"Prevention is better than cure"**

OHE treats your infrastructure as a **living organism**. Instead of reacting to CPU spikes, it tracks six holistic KPIs — Stress, Fatigue, Mood, Pressure, Humidity, Contagion — and predicts failures before they happen.

| Monitoring tool | What it tells you |
|---|---|
| Classic (Prometheus, Grafana) | CPU is at 85% |
| APM (Datadog, Dynatrace) | Service A responded slowly |
| **OHE v4** | **Storm in ~2 hours · fatigue building · contagion spreading to workers** |

---

## Table of Contents

1. [Architecture](#architecture)
2. [Quick Start](#quick-start)
3. [Configuration](#configuration)
4. [Holistic KPIs](#holistic-kpis)
5. [Prediction Engine (ILR)](#prediction-engine-ilr)
6. [Alerting](#alerting)
7. [REST API](#rest-api)
8. [WebSocket Live Feed](#websocket-live-feed)
9. [Embedded UI](#embedded-ui)
10. [Security Model](#security-model)
11. [Kubernetes Deployment](#kubernetes-deployment)
12. [K8s Operator](#k8s-operator)
13. [Project Structure](#project-structure)
14. [Test Coverage](#test-coverage)
15. [Roadmap](#roadmap)

---

## Architecture

```
┌───────────────────────────────────────────────────────────────┐
│                    OHE — Single Binary                        │
├──────────────┬────────────────────────────────────────────────┤
│  Collector   │  /proc · cgroups · container stats · logs      │
│      ↓       │                                                │
│  Processor   │  Normalize · Circular buffer · Downsample      │
│      ↓       │                                                │
│  Analyzer    │  Stress · Fatigue · Mood · Pressure ·          │
│              │  Humidity · Contagion                          │
│      ↓       │                                                │
│  Predictor   │  ILR + Dynamic thresholds + Anomaly detection  │
│      ↓       │                                                │
│  Alerter     │  Rule engine · Ack · Silence · Lifecycle       │
│      ↓       │                                                │
│  Storage     │  Badger v3 embedded TSDB (no external DB)      │
│      ↓       │                                                │
│  REST API    │  :8080  +  WebSocket /api/v1/ws                │
│      ↓       │                                                │
│  Svelte UI   │  Served from ./web (built-in)                  │
└──────────────┴────────────────────────────────────────────────┘
```

### Two Operating Modes

| Mode | Role | Default Port |
|---|---|---|
| `central` | Collects locally, persists to Badger, serves REST API + UI | 8080 |
| `agent` | Collects metrics on a remote node, pushes to central every 15 s | 8081 |

Agents talk to central via `POST /api/v1/ingest`. No message broker. No sidecar. No external dependencies.

---

## Quick Start

### Build from Source

```bash
git clone https://github.com/benfradjselim/Mlops_crew_automation
cd Mlops_crew_automation/workdir

# Build the UI first (outputs to ./web/)
cd ui && npm ci && npm run build && cd ..

# Build the Go binary (embeds ./web at runtime)
go build -ldflags="-s -w" -o ohe ./cmd/ohe/
```

### Run

```bash
# Central node — API + UI + local collection + storage
./ohe central --port 8080 --storage /var/lib/ohe/data

# Agent node — collect and push to central
./ohe agent --central-url http://central:8080 --interval 15s
```

### Docker

```bash
# Build the image (3-stage: Node → Go → distroless)
docker build -t ohe:4.0.0 workdir/

# Run
docker run -p 8080:8080 -v ohe-data:/var/lib/ohe/data ohe:4.0.0
```

### First-Boot Admin Account

On the very first start with an empty database, OHE auto-generates a random 32-character admin password and prints it to stdout:

```
╔══════════════════════════════════════════════════╗
║  FIRST BOOT — admin credentials generated         ║
║  Username : admin                                  ║
║  Password : <random 32-char hex>                  ║
║  Change this password immediately after login!     ║
╚══════════════════════════════════════════════════╝
```

Alternatively, call the setup endpoint once before starting anything else:

```bash
curl -s -X POST http://localhost:8080/api/v1/auth/setup \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"ChangeMe123!"}'
```

The endpoint returns `409 Conflict` if any user already exists — safe to call during provisioning.

---

## Configuration

Both modes accept a YAML config file (`--config`) and CLI flag overrides.

```yaml
# configs/central.yaml
mode: central
host: ""                     # defaults to os.Hostname()
port: 8080
storage_path: /var/lib/ohe/central
collect_interval: 15s
buffer_size: 10000
auth_enabled: false          # set true + jwt_secret for production
jwt_secret: "change-me-in-production"
allowed_origins: []          # empty = wildcard * (dev only)
# allowed_origins:
#   - https://dashboard.example.com
```

```yaml
# configs/agent.yaml
mode: agent
host: ""
port: 8081
central_url: http://localhost:8080
collect_interval: 15s
```

### CLI Flags

```
ohe central --port 8080 --storage /var/lib/ohe --auth --jwt-secret "$SECRET"
ohe agent   --central-url http://central:8080 --interval 15s
```

| Flag | Default | Description |
|---|---|---|
| `--config` | — | Path to YAML config file |
| `--port` | 8080 | HTTP listen port |
| `--storage` | `/var/lib/ohe/data` | Badger data directory |
| `--central-url` | `http://localhost:8080` | Central URL (agent mode) |
| `--interval` | `15s` | Metric collection interval |
| `--auth` | false | Enable JWT authentication |
| `--jwt-secret` | — | JWT signing secret (required with `--auth`) |

---

## Holistic KPIs

OHE computes six composite KPIs every collection cycle, treating infrastructure state like a physiological organism.

| KPI | Formula | States (low → high) |
|---|---|---|
| **Stress** | `0.30·CPU + 0.20·RAM + 0.20·Load + 0.20·Errors + 0.10·Timeouts` | calm · nervous · stressed · panic |
| **Fatigue** | `∫(Stress − RecoveryRate) dt` — time-integrated stress accumulation | rested · tired · exhausted · burnout |
| **Mood** | `(Uptime × Throughput) / (Errors × Timeouts × Restarts + ε)` | depressed · sad · neutral · content · happy |
| **Pressure** | `dStress/dt + ∫Errors dt` — rate of stress change plus error integral | stable · rising · storm_approaching · improving |
| **Humidity** | `(Errors × Timeouts) / (Throughput + ε)` | dry · humid · very_humid · storm |
| **Contagion** | `Errors × CPU_load` | low · moderate · epidemic · pandemic |

All KPI values are normalised to `[0, 1]`. Each has named states so alerts read in plain English rather than raw numbers.

### Storage TTLs

| Data type | Retention |
|---|---|
| Raw metrics | 7 days |
| KPIs | 7 days |
| Alerts | 90 days |
| Logs | 30 days |
| Predictions | 30 days |
| Dashboards / Users / Datasources | No expiry |

---

## Prediction Engine (ILR)

OHE uses **Incremental Linear Regression** — written from scratch in pure Go with zero dependencies.

### How It Works

Each metric gets its own online ILR model that updates on every new sample using Welford's method (numerically stable, O(1) space):

```
α (slope)  = Σ(x·y) - n·x̄·ȳ  /  Σ(x²) - n·x̄²
β (intercept) = ȳ - α·x̄
predicted(t + horizon) = α · (current_x + horizon_steps) + β
```

A `BatchILR` variant re-trains on a fixed sliding window of the most recent N points — providing more stability on volatile metrics.

### Dynamic Thresholds

Instead of static alert thresholds, OHE computes per-metric dynamic bounds from a rolling window of recent values:

```
upper_bound = mean + sigma_multiplier × stddev
```

This means a CPU that normally runs at 90% won't fire alerts, but one that jumps from 30% to 85% will.

### Anomaly Detector

A per-metric Welford-based Z-score detector flags values that deviate more than N standard deviations from the rolling mean. Separate from threshold alerting — catches sudden spikes and drops.

### Storm Detector

Tracks load pressure over a configurable time window. Detects when sustained high pressure predicts an imminent load storm, returning an ETA in hours.

---

## Alerting

The alerter holds a set of named rules. Each rule defines a metric, a threshold, a severity, and a hold duration before firing.

```go
// Rule evaluated every collection cycle
type Rule struct {
    ID          string
    Metric      string
    Threshold   float64
    Operator    string   // ">" | "<" | ">=" | "<="
    Severity    string   // "critical" | "warning" | "info"
    HoldSeconds int      // fire only after condition holds this long
    Description string
}
```

**Alert lifecycle:** `firing` → `acknowledged` / `silenced` → deleted after 90 days.

O(1) resolution via secondary `ruleHostIdx` index. Stale fired entries evicted after 5 minutes of inactivity. Overflow protection via `dropped` counter.

**Built-in predictive alerts:**

| Trigger | Alert |
|---|---|
| Pressure > 0.7 sustained 10 min | Storm approaching — scale up now |
| Fatigue > 0.8 | Burnout imminent — schedule restart |
| Contagion > 0.6 | Epidemic spreading — isolate services |
| Humidity > 0.5 | Error storm — activate circuit breaker |

---

## REST API

All responses follow this envelope:

```json
{
  "success": true,
  "data": {},
  "timestamp": "2026-04-14T10:00:00Z"
}
```

Errors use `"success": false` with an `"error": {"code": "...", "message": "..."}` field.

### Health & System

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/api/v1/health` | — | Full health check with storage status |
| GET | `/api/v1/health/live` | — | Liveness probe — always 200 while process is alive |
| GET | `/api/v1/health/ready` | — | Readiness probe — 200 when ready, 503 otherwise |
| GET | `/api/v1/config` | viewer | Runtime configuration |
| POST | `/api/v1/reload` | admin | Reload configuration |

### Metrics

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/api/v1/metrics` | viewer | All current normalized metric values |
| GET | `/api/v1/metrics/{name}?host=&from=-1h` | viewer | Metric time series |
| GET | `/api/v1/metrics/{name}/aggregate` | viewer | avg / min / max / p95 / p99 |
| POST | `/api/v1/query` | viewer | Query by metric name + time range |
| POST | `/api/v1/ingest` | operator | Agent push (MetricBatch) |

### KPIs & Predictions

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/api/v1/kpis?host=` | viewer | Current KPI snapshot for a host |
| GET | `/api/v1/kpis/{name}?host=` | viewer | Single KPI time series |
| GET | `/api/v1/kpis/{name}/predict` | viewer | ILR prediction for a KPI |
| GET | `/api/v1/predict?host=&metric=&horizon=120` | viewer | All metric predictions |

### Alerts

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/api/v1/alerts` | viewer | List all active alerts |
| GET | `/api/v1/alerts/{id}` | viewer | Get alert by ID |
| POST | `/api/v1/alerts/{id}/acknowledge` | viewer | Acknowledge |
| POST | `/api/v1/alerts/{id}/silence` | viewer | Silence |
| DELETE | `/api/v1/alerts/{id}` | operator | Delete |

### Dashboards & Templates

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/api/v1/dashboards` | viewer | List dashboards |
| POST | `/api/v1/dashboards` | operator | Create dashboard |
| GET | `/api/v1/dashboards/{id}` | viewer | Get dashboard |
| PUT | `/api/v1/dashboards/{id}` | operator | Update dashboard |
| DELETE | `/api/v1/dashboards/{id}` | operator | Delete dashboard |
| GET | `/api/v1/dashboards/{id}/export` | viewer | Export as JSON (download) |
| POST | `/api/v1/dashboards/import` | operator | Import from JSON |
| GET | `/api/v1/templates` | viewer | List built-in dashboard templates |
| GET | `/api/v1/templates/{id}` | viewer | Get template details |
| POST | `/api/v1/templates/{id}/apply` | operator | Instantiate template as a new dashboard |

### Data Sources

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/api/v1/datasources` | viewer | List data sources |
| POST | `/api/v1/datasources` | operator | Create (SSRF-validated) |
| GET | `/api/v1/datasources/{id}` | viewer | Get data source |
| PUT | `/api/v1/datasources/{id}` | operator | Update |
| DELETE | `/api/v1/datasources/{id}` | operator | Delete |
| POST | `/api/v1/datasources/{id}/test` | operator | Test connectivity |

Data source URLs are validated at create/update time against an SSRF allowlist. Loopback, private, and link-local addresses are blocked. AWS/GCP metadata endpoints (`169.254.169.254`, `metadata.google.internal`) are explicitly denied.

### Auth & Users

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/api/v1/auth/setup` | — | First-run admin creation (one-time, 409 after) |
| POST | `/api/v1/auth/login` | — | Login → JWT (24h) |
| POST | `/api/v1/auth/refresh` | viewer | Refresh token |
| POST | `/api/v1/auth/logout` | viewer | Logout (stateless acknowledgement) |
| GET | `/api/v1/auth/users` | admin | List users |
| POST | `/api/v1/auth/users` | admin | Create user |
| GET | `/api/v1/auth/users/{id}` | admin | Get user |
| DELETE | `/api/v1/auth/users/{id}` | admin | Delete user |

### WebSocket

| Path | Auth | Description |
|---|---|---|
| `WS /api/v1/ws` | viewer | Live stream of KPI updates, alerts, and raw metric values |

---

## WebSocket Live Feed

Connect with a valid Bearer token as a query parameter:

```js
const ws = new WebSocket(`ws://localhost:8080/api/v1/ws?token=${jwt}`)
ws.onmessage = (ev) => {
  const msg = JSON.parse(ev.data)
  // msg.type: "kpi" | "alert" | "metric"
  // msg.payload: { stress: 0.42, fatigue: 0.31, ... } | AlertObject | MetricObject
}
```

The hub broadcasts to all connected clients on each collection cycle. CORS origin enforcement mirrors the server's `allowed_origins` list.

---

## Embedded UI

The Svelte SPA is compiled to `workdir/web/` and served by the Go binary's file server. No separate web server needed.

```
ui/src/
├── App.svelte              Sidebar navigation + page routing
├── app.css                 Global dark-theme reset
├── lib/
│   ├── api.js              Fetch wrapper — JWT Bearer, all endpoints
│   ├── store.js            Svelte stores — token, page, KPI state
│   ├── KpiGauge.svelte     SVG half-circle gauge with severity colouring
│   └── Sparkline.svelte    Inline polyline chart for metric history
└── pages/
    ├── Dashboard.svelte    6 KPI gauges · live sparklines · alerts table
    ├── Alerts.svelte       Full alert list with ack / silence / delete
    ├── Dashboards.svelte   Create / list / delete · template apply
    ├── Login.svelte        Login + first-run setup detection
    └── Settings.svelte     User management (create / delete · roles)
```

**Features:**
- Dark theme (slate-900 colour palette)
- WebSocket live feed with automatic 3-second reconnect
- Login page auto-detects first-run mode and shows account creation form
- JWT stored in `localStorage`, cleared on logout
- KPI gauges colour-code by severity (green → yellow → orange → red)

**Build the UI:**

```bash
cd workdir/ui
npm ci
npm run build      # outputs to workdir/web/
```

---

## Security Model

### Authentication

JWT HS256 tokens, 24-hour expiry. Enabled with `--auth --jwt-secret <secret>`.

The `jwt_secret` must be changed from the default value before `--auth` is accepted — the binary refuses to start otherwise.

### RBAC

Three roles, checked per-route by `RequireRole` middleware:

| Role | Capabilities |
|---|---|
| `viewer` | Read-only access to metrics, KPIs, alerts, dashboards |
| `operator` | viewer + create/update/delete datasources, dashboards, ingest |
| `admin` | operator + user management, reload config |

When auth is disabled (development), all roles are permitted — the middleware passes through without a token check.

### Rate Limiting

Login endpoint: token-bucket, 5 attempts per minute per IP. Stale buckets evicted after 10 minutes of inactivity.

### SSRF Protection

Data source URLs are resolved and validated at storage time. Blocked ranges:

- Loopback (`127.0.0.0/8`, `::1`)
- RFC 1918 private ranges (`10/8`, `172.16/12`, `192.168/16`)
- Link-local (`169.254/16`, `fe80::/10`)
- AWS/GCP metadata endpoints

### Other Headers

Every response carries `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, `Referrer-Policy: strict-origin-when-cross-origin`, and `Content-Security-Policy: default-src 'self'`.

---

## Kubernetes Deployment

### Prerequisites

```bash
kubectl apply -f workdir/deploy/crd/oheclusters.yaml
```

### Deploy with Kustomize

```bash
# Edit workdir/deploy/kustomization.yaml to set your JWT secret and image tag
kubectl apply -k workdir/deploy/
```

This creates:
- Namespace `ohe-system`
- `ServiceAccount` + `ClusterRole` / `ClusterRoleBinding` for both central and agent
- `ConfigMap` with `config.yaml`
- `Secret` `ohe-secrets` with the JWT signing key
- `PersistentVolumeClaim` (10 Gi, `ReadWriteOnce`) for Badger data
- `Deployment` `ohe-central` (1 replica, `Recreate` strategy — Badger is single-writer)
- `Service` `ohe-central` (ClusterIP :80 → :8080)
- `DaemonSet` `ohe-agent` (one pod per node, mounts `/proc` and `/sys`)

### Probe Configuration

The central Deployment uses the OHE health endpoints directly:

```yaml
livenessProbe:
  httpGet:
    path: /api/v1/health/live   # always 200
    port: http
  initialDelaySeconds: 5
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /api/v1/health/ready  # 503 until engine is initialised
    port: http
  initialDelaySeconds: 3
  periodSeconds: 5
```

### Resource Requests

| Component | CPU request | Memory request | CPU limit | Memory limit |
|---|---|---|---|---|
| central | 100m | 128 Mi | 500m | 512 Mi |
| agent | 50m | 64 Mi | 200m | 128 Mi |

---

## K8s Operator

The operator lives in `workdir/operator/` — a **standalone Go module** with **zero external dependencies** (no `controller-runtime`, no `client-go`). It talks directly to the K8s API server over HTTPS using the pod's service account token.

### OHECluster CRD

```yaml
apiVersion: ohe.io/v1alpha1
kind: OHECluster
metadata:
  name: production
  namespace: ohe-system
spec:
  mode: central          # "central" or "agent"
  replicas: 1            # central must be 1 (Badger single-writer)
  image: ghcr.io/benfradjselim/ohe:4.0.0
  storageSize: 20Gi
  authEnabled: true
  resources:
    requests: { cpu: "200m", memory: "256Mi" }
    limits:   { cpu: "1",    memory: "1Gi" }
```

```yaml
# Status written back by the operator
status:
  phase: Running
  readyReplicas: 1
  availableReplicas: 1
  lastReconcileTime: "2026-04-14T10:00:00Z"
  observedGeneration: 3
```

### Reconcile Loop

```
every 30s:
  LIST /apis/ohe.io/v1alpha1/oheclusters
  for each OHECluster:
    1. Build desired Deployment spec (central or agent mode)
    2. Server-side PATCH apply → idempotent create-or-update
    3. Read Deployment .status.readyReplicas
    4. PATCH OHECluster .status with phase + replica counts
```

### Build and Deploy the Operator

```bash
cd workdir/operator
go build -o ohe-operator ./

# Apply CRD + RBAC first, then deploy the operator pod
kubectl apply -f ../deploy/crd/oheclusters.yaml
# Deploy the operator binary as a Deployment (RBAC needed for CRD + Deployment access)
```

---

## Project Structure

```
Mlops_crew_automation/
└── workdir/
    ├── cmd/ohe/                Main entry point — agent / central / version / help
    ├── internal/
    │   ├── collector/          System metrics via /proc, cgroups, containers, logs
    │   ├── processor/          Normalise · Circular buffer · Per-host aggregation
    │   ├── analyzer/           6 holistic KPI formulas + state machines
    │   ├── predictor/          ILR · BatchILR · DynamicThreshold · AnomalyDetector · StormDetector
    │   ├── alerter/            Rule engine · O(1) resolution · Ack/silence lifecycle
    │   ├── storage/            Badger v3 TSDB wrapper — metrics, KPIs, alerts, users, dashboards
    │   ├── api/                REST handlers · JWT middleware · RBAC · rate limiter · WebSocket hub
    │   └── orchestrator/       Engine wiring · goroutine lifecycle · graceful shutdown · first-boot seed
    ├── pkg/
    │   ├── models/             Shared data types (User, Alert, Dashboard, Metric, KPI, …)
    │   └── utils/              Math (stddev, percentiles, trapezoid) · CircularBuffer · ID generator
    ├── ui/                     Vite + Svelte SPA source
    ├── web/                    Compiled UI — served by Go file server at runtime
    ├── deploy/                 Kubernetes manifests (kustomize)
    │   ├── crd/                OHECluster CRD v1alpha1
    │   ├── rbac.yaml           Namespace · ServiceAccounts · ClusterRole
    │   ├── configmap.yaml      OHE config.yaml
    │   ├── central-deployment.yaml  PVC + Deployment + Service
    │   ├── agent-daemonset.yaml     DaemonSet with hostPID + /proc mounts
    │   └── kustomization.yaml
    ├── operator/               K8s operator (separate Go module, zero external deps)
    │   ├── main.go             Poll loop · signal handling
    │   ├── controller.go       Reconcile logic · central/agent mode
    │   ├── k8s.go              Raw HTTPS K8s API client (service account auth)
    │   └── types.go            OHECluster CR types · Deployment types
    ├── configs/
    │   ├── central.yaml
    │   └── agent.yaml
    ├── Dockerfile              3-stage: Node (UI) → Go (binary) → distroless runtime
    └── go.mod                  Module: github.com/benfradjselim/ohe  Go 1.18
```

---

## Test Coverage

```
Package                          Coverage
─────────────────────────────────────────
internal/predictor               99.1 %
internal/analyzer                99.0 %
pkg/utils                        97.7 %
internal/alerter                 95.1 %
internal/storage                 91.8 %
internal/processor               89.7 %
internal/api                     72.1 %
internal/collector               74.8 %
─────────────────────────────────────────
Overall                          81.0 %
```

Run the suite:

```bash
cd workdir
go test ./... -v -timeout 180s
go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
```

---

## Roadmap

| Phase | Goal | Status |
|---|---|---|
| 1 | Collector + Processor + Badger storage + 6 holistic KPIs | ✅ Done |
| 2 | ILR predictions + BatchILR + Dynamic thresholds + Storm detector + Alerter | ✅ Done |
| 3 | Full REST API (40+ endpoints) + WebSocket hub | ✅ Done |
| 3.1 | JWT auth + RBAC (viewer/operator/admin) + CORS allowlist + rate limiter + SSRF protection | ✅ Done |
| 3.2 | First-boot admin seed + dashboard templates + setup endpoint | ✅ Done |
| 4 | Svelte UI — KPI gauges · alerts · dashboards · user management · WebSocket live feed | ✅ Done |
| 5 | Liveness/readiness probes + Dockerfile + K8s manifests + CRD + operator | ✅ Done |
| 6 | Distributed tracing + multi-cluster federation | Planned |

---

## Dependencies

| Dependency | Version | Purpose |
|---|---|---|
| `github.com/dgraph-io/badger/v3` | v3.2103.5 | Embedded time-series key-value store |
| `github.com/golang-jwt/jwt/v5` | v5.2.1 | JWT generation and validation |
| `github.com/gorilla/mux` | v1.8.1 | HTTP router with path variables |
| `github.com/gorilla/websocket` | v1.5.1 | WebSocket upgrade + hub |
| `golang.org/x/crypto` | v0.22.0 | bcrypt password hashing (cost 12) |
| `golang.org/x/sys` | v0.22.0 | Linux `/proc` syscall wrappers |
| `gopkg.in/yaml.v3` | v3.0.1 | Config file parsing |

The operator module (`workdir/operator/`) has **no external dependencies** — only the Go standard library.

---

**Author:** Selim Benfradj · **License:** MIT · **Version:** 4.0.0
