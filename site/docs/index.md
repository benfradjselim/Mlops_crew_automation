# Ruptura

**The Predictive Action Layer for Cloud-Native Infrastructure.**

Ruptura detects infrastructure ruptures before they cause outages — and acts on them automatically via Kubernetes, webhooks, and alerting integrations. A single Go binary, no external database required.

---

## Why Ruptura?

| Traditional Observability | Ruptura |
|--------------------------|-----------|
| Threshold alerts fire *after* the fact | Rupture Index™ detects divergence **hours early** |
| Rules defined per metric | Adaptive ensemble learns your baseline automatically |
| Manual incident response | Tier-1 actions (scale, restart, rollback) with safety gates |
| 5+ tools: Prom + Grafana + AM + Loki + PD | **One binary**, one `kubectl apply` |
| No reasoning about *why* | Full XAI trace for every prediction |

---

## Core Concepts

### Rupture Index™

```
R(t) = |α_burst(t)| / max(|α_stable(t)|, ε)
```

| R Range | State | Ruptura Action |
|---------|-------|-------------|
| < 1.5 | Stable / Elevated | None |
| 1.5 – 3.0 | Warning | Tier-3 (human) |
| 3.0 – 5.0 | Critical | Tier-2 (suggested) |
| ≥ 5.0 | Emergency | Tier-1 (automated) |

### 8 Composite Signals

`stress` · `fatigue` · `pressure` · `contagion` · `resilience` · `entropy` · `sentiment` · `healthscore`

Each maps raw metrics to a single interpretable 0–1 index with published formulas.

---

## Quick Start

=== "Kubernetes"

    ```bash
    git clone https://github.com/benfradjselim/ruptura.git
    cd ruptura
    docker build -t ruptura:6.2.1 .
    kubectl apply -f deploy/
    kubectl port-forward svc/ruptura 8080:8080 -n ruptura-system
    curl http://localhost:8080/api/v2/health
    ```

=== "Docker"

    ```bash
    docker run -d \
      -p 8080:8080 \
      -p 4317:4317 \
      -v ruptura-data:/var/lib/ruptura \
      -e RUPTURA_API_KEY=$(openssl rand -hex 32) \
      ruptura:6.2.1

    curl http://localhost:8080/api/v2/health
    ```

=== "Helm"

    ```bash
    helm install ruptura ./helm \
      --namespace ruptura-system \
      --create-namespace \
      --set apiKey=$(openssl rand -hex 32)
    ```

---

## Current Release

**v6.2.1** — FusedR in API · anomaly REST endpoints (`/api/v2/anomalies`) · WorkloadRef-native pipeline · stable engine

[Full changelog →](community/roadmap.md) · [Getting Started →](getting-started/installation.md) · [API Reference →](api/reference.md)

 


