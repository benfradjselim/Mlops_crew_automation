# OHE — Observability Holistic Engine

<p align="center">
  <img src="https://img.shields.io/badge/version-4.1.0-blue?style=for-the-badge" alt="Version">
  <img src="https://img.shields.io/badge/alert%20delivery-slack%20%7C%20pagerduty%20%7C%20webhook-brightgreen?style=for-the-badge" alt="Alert Delivery">
  <img src="https://img.shields.io/badge/go-1.22+-00ADD8?style=for-the-badge&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/kubernetes-native-326CE5?style=for-the-badge&logo=kubernetes" alt="Kubernetes">
  <img src="https://img.shields.io/badge/single%20binary-no%20deps-success?style=for-the-badge" alt="Single Binary">
  <img src="https://img.shields.io/badge/license-MIT-green?style=for-the-badge" alt="License">
</p>

<p align="center">
  <strong>Stop reacting. Start predicting.</strong><br>
  OHE is a self-hosted, Kubernetes-native observability platform with ML-powered predictive alerting.<br>
  One binary replaces Grafana + Prometheus + Alertmanager + Loki — no external databases required.
</p>

---

## v4.1.0 — What's New

### End-to-End Alert Delivery
Alerts now reach your team automatically — no more checking dashboards manually.

- **Slack** — post alerts to any channel with severity-based filtering
- **PagerDuty** — create incidents with routing key support
- **Webhook** — POST to any HTTP endpoint with custom headers
- **Per-severity routing** — a channel can subscribe to `critical` only, or `warning+` — your call
- **Live test** — `POST /api/v1/notifications/{id}/test` fires a real payload immediately
- **Full lifecycle** — active → acknowledged → silenced → resolved, delivered at each transition

### Kubernetes DaemonSet Agent — Real Metrics
- Agent DaemonSet now ships real host metrics (CPU, memory, disk, network) from every node
- Same binary runs as either central server or lightweight agent (`OHE_MODE=agent`)
- `OHE_CENTRAL_URL` fallback for environments where cluster DNS is unavailable

### SSRF Hardening for Datasources
- `*.svc.cluster.local` and `*.svc` hostnames are now allowed in datasource URLs (safe in-cluster)
- `OHE_TRUSTED_DATASOURCE_HOSTS` — explicit allowlist for Prometheus ClusterIPs
- All external IPs still blocked by default

### Prediction Chart Fix
- Forecast dashboards now show future timestamps on the X-axis
- "Next 1h" → labels from `now` to `now+1h`, not historical data

---

## Full Feature Set (v4.0.0 + v4.1.0)

| Feature | Status |
|---------|--------|
| System metrics — CPU, memory, disk, network, load, uptime | ✅ |
| 10 Holistic KPIs — stress, fatigue, mood, pressure, humidity, contagion, resilience, entropy, velocity, health_score | ✅ |
| ML forecasting — exponential smoothing, predicts exhaustion hours ahead | ✅ |
| REST API + WebSocket live feed | ✅ |
| Svelte UI embedded in binary (no CDN, no Node.js at runtime) | ✅ |
| 14 built-in dashboard templates | ✅ |
| Predictive dashboards — Next 1h / 6h / 24h mode | ✅ |
| Dashboard tabs | ✅ |
| Alert rules engine — threshold and KPI-based | ✅ |
| **Alert delivery — Slack, PagerDuty, webhook** | ✅ **NEW** |
| **Kubernetes DaemonSet agent — real host metrics** | ✅ **NEW** |
| **SSRF hardening for datasource URLs** | ✅ **NEW** |
| Prometheus `/metrics` exposition | ✅ |
| OTLP / Loki / Elasticsearch / Datadog / DogStatsD ingestion | ✅ |
| Log collection + viewer | ✅ |
| Distributed traces viewer | ✅ |
| JWT auth + user management | ✅ |
| BadgerDB embedded storage — no external database | ✅ |
| Kubernetes operator + PVC | ✅ |

---

## Quick Start

### Kubernetes

```bash
git clone https://github.com/benfradjselim/Mlops_crew_automation.git
cd Mlops_crew_automation/workdir

cd ui && npm install && npm run build && cd ..
docker build -t ohe:latest .

kubectl apply -f deploy/pvc.yaml
kubectl apply -f deploy/secrets.yaml
kubectl apply -f deploy/configmap.yaml
kubectl apply -f deploy/rbac.yaml
kubectl apply -f deploy/central-deployment.yaml
kubectl apply -f deploy/agent-daemonset.yaml

# Get credentials
kubectl logs -n ohe-system deploy/ohe-central | grep -A4 "FIRST BOOT"

# Open UI
kubectl port-forward svc/ohe-central 8080:80 -n ohe-system
```

### Docker Compose

```yaml
version: "3.9"
services:
  ohe:
    image: ohe:latest
    ports: ["8080:8080"]
    volumes: ["ohe-data:/var/lib/ohe/data"]
    environment:
      OHE_ADMIN_PASSWORD: changeme
volumes:
  ohe-data:
```

---

## Configure Alert Channels

```bash
# Slack
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"ops-slack","type":"slack","url":"https://hooks.slack.com/...","severities":["critical","warning"]}'

# PagerDuty
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"pagerduty","type":"pagerduty","url":"https://events.pagerduty.com/v2/enqueue","headers":{"X-Routing-Key":"<key>"},"severities":["critical"]}'

# Test any channel
curl -X POST http://localhost:8080/api/v1/notifications/{id}/test \
  -H "Authorization: Bearer $TOKEN"
```

---

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `OHE_ADMIN_PASSWORD` | auto-generated | Override admin password |
| `OHE_AUTH_ENABLED` | `false` | Enforce JWT authentication |
| `OHE_JWT_SECRET` | `change-me-in-production` | JWT signing secret |
| `OHE_TRUSTED_DATASOURCE_HOSTS` | — | Comma-separated IPs/hostnames bypassing SSRF check |
| `OHE_PORT` | `8080` | HTTP listen port |
| `OHE_STORAGE_PATH` | `/var/lib/ohe/data` | BadgerDB data directory |
| `OHE_COLLECT_INTERVAL` | `15s` | Metric collection interval |
| `OHE_DOGSTATSD_ADDR` | `:8125` | DogStatsD UDP listener |

---

## What Comes Next

| Version | Focus |
|---------|-------|
| v4.2.0 | PromQL passthrough proxy + Organisations multi-tenancy |
| v4.3.0 | 3-tier long-term retention (400 days) + full RBAC |
| v4.4.0 | SLO engine, widget resize, 20 dashboard templates |

---

## License

MIT — free to use, modify, and deploy.
