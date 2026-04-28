# Kairo Core

<p align="center">
  <img src="https://img.shields.io/badge/version-6.1.1-0069ff?style=for-the-badge" alt="v6.1.1">
  <img src="https://img.shields.io/badge/go-1.18+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go 1.18+">
  <img src="https://img.shields.io/badge/license-Apache%202.0-green?style=for-the-badge" alt="Apache 2.0">
  <img src="https://img.shields.io/badge/kubernetes-native-326CE5?style=for-the-badge&logo=kubernetes&logoColor=white" alt="Kubernetes Native">
  <img src="https://img.shields.io/badge/coverage-70%25+-brightgreen?style=for-the-badge" alt="Coverage">
</p>

<p align="center">
  <b>The Predictive Action Layer for Cloud-Native Infrastructure.</b><br>
  Kairo detects infrastructure ruptures before they cause outages — and acts on them automatically.
</p>

**[Documentation →](https://benfradjselim.github.io/Mlops_crew_automation/)**

---

## What Kairo Does

Traditional observability tells you what broke. Kairo tells you **what is about to break** — and triggers the right action before users feel it.

| Traditional Observability | Kairo Core |
|--------------------------|-----------|
| Threshold alerts fire after the fact | Rupture Index™ detects divergence **hours early** |
| You define rules per metric | Adaptive ensemble learns your baseline automatically |
| Manual incident response | Tier-1 actions (scale, restart, rollback) fire automatically with safety gates |
| 5+ tools: Prom + Grafana + AM + Loki + PD | **One binary**, one `kubectl apply` |
| No reasoning about why | Full XAI trace for every prediction |

---

## Core Concepts

### Rupture Index™

```
R(t) = |α_burst(t)| / max(|α_stable(t)|, ε)

  α_burst  = slope from 5-min CA-ILR tracker  (detects sudden change)
  α_stable = slope from 60-min CA-ILR tracker (tracks baseline)
  ε        = 1e-6 (numerical stability)
```

| R Range | State | Action |
|---------|-------|--------|
| < 1.5 | Stable / Elevated | None |
| 1.5–3.0 | Warning | Tier-3 (human) |
| 3.0–5.0 | Critical | Tier-2 (suggested) |
| ≥ 5.0 | Emergency | Tier-1 (automated) |

### Adaptive Ensemble (v6.1)

Five models (CA-ILR, ARIMA, Holt-Winters, MAD, EWMA) with online MAE-based weight adaptation — weights update every 60s over a 1-hour sliding window.

### 8 Composite Signals

`stress` · `fatigue` · `pressure` · `contagion` · `resilience` · `entropy` · `sentiment` · `healthscore`

Each maps raw telemetry to an interpretable 0–1 index with a published formula. No black boxes.

---

## Architecture

```
┌──────────────────────────────────────────────────────────┐
│                      kairo-core                          │
│                                                          │
│  Ingest ──► Metric/Log/Trace pipelines ──► Fusion        │
│     │              │                         │           │
│  gRPC           Composites              RuptureDetector  │
│  OTLP           (8 signals)                  │           │
│  Prom rw        Adaptive                  Actions        │
│  DogStatsD      Ensemble              (K8s/Webhook/PD)   │
│     │                                        │           │
│  NATS/Kafka eventbus ◄──────────────── XAI Explain       │
│                                              │           │
│              REST API v2 (44 endpoints) ─────┘           │
│              K8s Operator (KairoInstance CRD)            │
└──────────────────────────────────────────────────────────┘
```

Single binary — BadgerDB embedded, no external database required.

---

## Quick Start

### Kubernetes (recommended)

```bash
git clone https://github.com/benfradjselim/Mlops_crew_automation.git
cd Mlops_crew_automation/workdir
docker build -t kairo-core:6.1.1 .
kubectl apply -f deploy/
kubectl port-forward svc/kairo-core 8080:8080 -n kairo-system
curl http://localhost:8080/api/v2/health
```

### Docker

```bash
docker run -d \
  -p 8080:8080 \
  -v kairo-data:/var/lib/kairo \
  -e KAIRO_JWT_SECRET=$(openssl rand -hex 32) \
  kairo-core:6.1.1
curl http://localhost:8080/api/v2/health
```

### Helm

```bash
helm install kairo-core ./workdir/helm \
  --namespace kairo-system \
  --create-namespace \
  --set auth.jwtSecret=$(openssl rand -hex 32)
```

---

## Configuration (`kairo.yaml`)

```yaml
mode: connected          # connected | stateless | shadow

ingest:
  http_port: 8080
  grpc_port: 9090

eventbus:
  driver: none           # none | nats | kafka

ensemble:
  adaptive: false        # true = online MAE-based weight adaptation (v6.1)

predictor:
  rupture_threshold: 3.0

actions:
  execution_mode: shadow  # shadow | suggest | auto
  safety:
    rate_limit_per_hour: 6

auth:
  jwt_secret: ""         # set via KAIRO_JWT_SECRET env var

storage:
  path: /var/lib/kairo
```

---

## SDKs

**Go**
```go
import ohe "github.com/benfradjselim/kairo-core/sdk/go"

c := ohe.New("http://kairo-core:8080", ohe.WithAPIKey("ohe_your_api_key"))
rupture, _ := c.RuptureIndex(ctx, "web-01")
weights, _ := c.EnsembleWeights(ctx, "web-01")  // v6.1
```

**Python**
```python
from kairo import KairoClient

c = KairoClient("http://kairo-core:8080", api_key="ohe_your_api_key")
rupture = c.rupture_index("web-01")
```

---

## Changelog

### v6.1.1 — 2026-04-28
- Documentation site launched at [benfradjselim.github.io/Mlops_crew_automation](https://benfradjselim.github.io/Mlops_crew_automation/)
- All 8 composite signal formulas published
- Bug fixes

### v6.1.0 — 2026-04-27
- **§23** gRPC ingest server (google.golang.org/grpc, 4MB max, back-pressure)
- **§24** NATS/Kafka eventbus — JetStream at-least-once + franz-go exactly-once
- **§25** Adaptive ensemble weighting — online MAE-based, 1-hour sliding window, 60s update
- **§26** Kubernetes operator — KairoInstance CRD, controller-runtime reconcile loop
- Go SDK `kairo-client-go` — full v2 API coverage

### v6.0.0 — 2026-04-25
- Complete clean-room rewrite from OHE v5.1 as `github.com/benfradjselim/kairo-core`
- CA-ILR dual-scale engine, 5-model ensemble, 8 composite signals
- 44-endpoint REST API v2, XAI explainability, single-tenant BadgerDB storage
- Action engine with safety gates, OTLP + Prom remote_write + DogStatsD ingest
- ≥70% test coverage

---

## Development

```bash
cd workdir
go build ./...
go test -race -timeout=120s ./...
go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | grep total
```

---

## License

Apache 2.0 — see [LICENSE](LICENSE)
