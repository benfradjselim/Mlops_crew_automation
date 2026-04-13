# Architecture Technique OHE v4.0.0

## 1. Vue d'ensemble

```

┌─────────────────────────────────────────────────────────────────────────────┐
│                         OBSERVABILITY HOLISTIC ENGINE                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         CLI (observability)                          │   │
│  │   observability agent --config=/etc/agent.yaml                      │   │
│  │   observability central --config=/etc/central.yaml                  │   │
│  │   observability version                                              │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ↓                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      CORE LAYER (goroutines)                         │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐  │   │
│  │  │Collector │→│Processor │→│ Analyzer │→│Predictor │→│ Alerter  │  │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ↓                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      STORAGE LAYER                                   │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │                      Badger (embedded)                       │    │   │
│  │  │  - metrics (7d TTL)  - logs (30d TTL)  - predictions (30d)  │    │   │
│  │  │  - alerts (90d)      - dashboards     - users               │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ↓                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      API LAYER (HTTP/2)                              │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │  REST API (port 8080) + WebSocket (port 8081)               │    │   │
│  │  │  - /api/v1/metrics  - /api/v1/kpis  - /api/v1/predict       │    │   │
│  │  │  - /api/v1/alerts   - /api/v1/dashboards - /api/v1/datasources│   │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ↓                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      UI LAYER (Svelte)                               │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │  Dashboard  │  Metrics  │  KPIs  │  Logs  │  Alerts  │ Admin │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

```

## 2. Modes de fonctionnement

| Mode | Rôle | Déploiement | Port |
|------|------|-------------|------|
| **agent** | Collecte locale, push vers central | DaemonSet sur chaque node | 8081 (internal) |
| **central** | Agrégation, API, UI | Deployment unique | 8080 (exposed) |

## 3. Communication

- agent → central: gRPC streaming des métriques, compression Snappy, batch toutes les 5 secondes
- central → agent: Health check (port 8081), configuration push

## 4. Stack technique

| Composant | Choix | Raison |
|-----------|-------|--------|
| Langage | Go 1.22+ | Binaire unique, performance |
| Storage | Badger | Embarcable, rapide, zero config |
| UI | Svelte + ECharts | Léger, réactif, compilé |
| API | net/http + gorilla/mux | Standard, pas de dépendance |
| Auth | JWT + RBAC | Simple, standard |
| Communication | gRPC | Performance, streaming |
| Compression | Snappy | Rapide, efficace |

## 5. Structure des données

### Metrics (TSDB Badger)

Key: m:{metric_name}:{timestamp}
Value: {value}
Exemple: m:cpu_usage:1743552000 = 0.45

### Logs

Key: l:{timestamp}:{pod_name}
Value: {log_message}
Exemple: l:1743552000:collector = {"level":"INFO","message":"collecting metrics"}

### KPIs

Key: k:{kpi_name}:{timestamp}
Value: {computed_value}
Exemple: k:stress:1743552000 = 0.62

### Dashboards

Key: d:{dashboard_id}
Value: {dashboard_json}
Exemple: d:system_overview = {"name":"System","widgets":[...]}

### Alerts

Key: a:{alert_id}:{timestamp}
Value: {alert_json}
Exemple: a:cpu_high:1743552000 = {"severity":"warning","status":"active"}

