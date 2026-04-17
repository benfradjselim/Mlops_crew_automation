# OHE — Observability Holistic Engine

<p align="center">
  <img src="https://img.shields.io/badge/version-4.3.0-blue?style=for-the-badge" alt="Version">
  <img src="https://img.shields.io/badge/retention-400%20days-blueviolet?style=for-the-badge" alt="Retention">
  <img src="https://img.shields.io/badge/RBAC-org%20isolation-orange?style=for-the-badge" alt="RBAC">
  <img src="https://img.shields.io/badge/go-1.22+-00ADD8?style=for-the-badge&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/kubernetes-native-326CE5?style=for-the-badge&logo=kubernetes" alt="Kubernetes">
  <img src="https://img.shields.io/badge/license-MIT-green?style=for-the-badge" alt="License">
</p>

<p align="center">
  <strong>Stop reacting. Start predicting.</strong><br>
  OHE is a self-hosted, Kubernetes-native observability platform with ML-powered predictive alerting.<br>
  One binary replaces Grafana + Prometheus + Alertmanager + Loki — no external databases required.
</p>

---

## v4.3.0 — What's New

### 3-Tier Long-Term Retention
Data no longer expires after 7 days. OHE automatically downsample and keeps history for over a year.

| Tier | Key Prefix | Retention | Resolution | Use Case |
|------|-----------|-----------|------------|---------|
| Raw | `m:` / `k:` | 7 days | As-collected (~15s) | Recent debugging |
| 5-min rollup | `r5:` / `kr5:` | 35 days | 5-minute average | Last-month trend |
| 1-hour rollup | `r1h:` / `kr1h:` | **400 days** | 1-hour average | Capacity planning |

- Compaction goroutine runs every **30 minutes** — raw → 5m after 2h, 5m → 1h after 12h
- `GET /api/v1/metrics/{name}/range` — tiered query, **auto-selects the right tier by window size**
  - < 6h window → raw tier
  - 6h – 7d window → 5-min tier
  - > 7d window → 1-hour tier
- `GET /api/v1/retention/stats` — live data point counts per tier
- `POST /api/v1/retention/compact` — trigger on-demand compaction (operator role)

### Full RBAC with Org Isolation
- Users now carry an `org_id` — propagated from JWT into every request context
- `PUT /api/v1/auth/users/{id}/org` — assign a user to an organisation (admin only)
- `GET /api/v1/orgs/{id}/members` — list all members of an org
- `POST /api/v1/orgs/{id}/members` — invite a user directly into an org
- Dashboards and datasources carry `org_id` — foundation for resource-level org isolation

---

## Full Feature Set (v4.0.0 → v4.3.0)

| Feature | Version |
|---------|---------|
| System metrics — CPU, memory, disk, network, load, uptime | v4.0.0 |
| 10 Holistic KPIs — stress, fatigue, mood, pressure, entropy, velocity… | v4.0.0 |
| ML forecasting — exponential smoothing, predicts exhaustion hours ahead | v4.0.0 |
| REST API + WebSocket live feed | v4.0.0 |
| Svelte UI embedded in binary (no CDN, no Node.js at runtime) | v4.0.0 |
| 14 built-in dashboard templates | v4.0.0 |
| Predictive dashboards — Next 1h / 6h / 24h mode | v4.0.0 |
| Alert rules engine — threshold and KPI-based | v4.0.0 |
| Alert delivery — Slack, PagerDuty, webhook | v4.1.0 |
| Kubernetes DaemonSet agent — real host metrics from every node | v4.1.0 |
| SSRF hardening for datasource URLs | v4.1.0 |
| PromQL passthrough proxy | v4.2.0 |
| Query widget — live PromQL results in dashboards | v4.2.0 |
| Organisations API — multi-tenant workspaces | v4.2.0 |
| **3-tier retention — 400-day history, automatic compaction** | **v4.3.0** |
| **Tiered range queries — auto-selects tier by window** | **v4.3.0** |
| **Full RBAC — org_id on users, dashboards, datasources** | **v4.3.0** |
| **Org member management** | **v4.3.0** |
| Prometheus / OTLP / Loki / Elasticsearch / Datadog / DogStatsD ingestion | v4.0.0 |
| JWT auth + user management | v4.0.0 |
| BadgerDB embedded storage — no external database | v4.0.0 |

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

kubectl logs -n ohe-system deploy/ohe-central | grep -A4 "FIRST BOOT"
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

## Retention API Examples

```bash
# Check data point counts per tier
curl http://localhost:8080/api/v1/retention/stats \
  -H "Authorization: Bearer $TOKEN"

# Trigger on-demand compaction
curl -X POST http://localhost:8080/api/v1/retention/compact \
  -H "Authorization: Bearer $TOKEN"

# Query 30-day trend (auto-selects 5-min tier)
curl "http://localhost:8080/api/v1/metrics/cpu/range?from=2024-01-01T00:00:00Z&to=2024-01-31T00:00:00Z" \
  -H "Authorization: Bearer $TOKEN"
```

---

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `OHE_ADMIN_PASSWORD` | auto-generated | Override admin password |
| `OHE_AUTH_ENABLED` | `false` | Enforce JWT authentication |
| `OHE_JWT_SECRET` | `change-me-in-production` | JWT signing secret |
| `OHE_TRUSTED_DATASOURCE_HOSTS` | — | Comma-separated IPs bypassing SSRF check |
| `OHE_PORT` | `8080` | HTTP listen port |
| `OHE_STORAGE_PATH` | `/var/lib/ohe/data` | BadgerDB data directory |
| `OHE_COLLECT_INTERVAL` | `15s` | Metric collection interval |

---

## What Comes Next

| Version | Focus |
|---------|-------|
| v4.4.0 | SLO / Error Budget engine, widget resize, 20 dashboard templates, 300° gauge |

---

## License

MIT — free to use, modify, and deploy.
