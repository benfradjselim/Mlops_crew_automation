# Ruptura

**The Predictive Action Layer for Cloud-Native Infrastructure.**

Ruptura detects workload ruptures before they cause outages — and acts on them automatically via Kubernetes, webhooks, and alerting integrations. A single Go binary, no external database required.

---

## Why Ruptura?

| Traditional Observability | Ruptura |
|--------------------------|---------|
| Threshold alerts fire *after* the fact | Fused Rupture Index™ detects divergence **hours early** |
| Global thresholds — batch jobs always "stressed" | **Adaptive per-workload baselines** after 24 h observation |
| "host-123 CPU 78%" — what does it mean? | "payment-api is exhausted — 72 h fatigue accumulation, cascade from payment-db" |
| Manual incident response | Tier-1 actions (scale, restart, rollback) with safety gates |
| 5+ tools: Prom + Grafana + AM + Loki + PD | **One binary**, one `helm install` |
| Numbers, no reasoning | **Narrative explain** — structured English causal chain |

---

## Core Concepts

### Fused Rupture Index™

Ruptura fuses three independent signal sources — raw metrics, OTLP logs, and OTLP trace spans — into a single rupture index per Kubernetes workload:

```
FusedR = f(metricR, logR, traceR)

  metricR = |α_burst| / max(|α_stable|, ε)   CA-ILR dual-scale slope ratio
  logR    = burst_rate / log_baseline          fires when error/warn > 3σ
  traceR  = span_error_rate × P99_deviation    from OTLP trace spans
```

FusedR requires at least **two** sources — a single noisy signal cannot push a workload to "critical."

| FusedR | State | Default action |
|--------|-------|---------------|
| < 1.5 | Stable / Elevated | None |
| 1.5 – 3.0 | Warning | Tier-3 (human alert) |
| 3.0 – 5.0 | Critical | Tier-2 (suggested action) |
| ≥ 5.0 | Emergency | Tier-1 (automated action) |

### 10 Composite KPI Signals

Every workload gets 10 auditable signals computed from raw telemetry with published formulas:

`stress` · `fatigue` · `mood` · `pressure` · `humidity` · `contagion` · `resilience` · `entropy` · `velocity` · `health_score`

`health_score` (0–100) is an additive-penalty composite. No black boxes — every coefficient is a versioned release artifact.

### WorkloadRef — Kubernetes-Native Treatment Unit

Ruptura groups all signals by **Kubernetes workload** (`namespace/kind/name`), not by host. Multiple pods from the same Deployment are merged into a single health view. OTLP resource attributes (`k8s.deployment.name`, `k8s.namespace.name`, etc.) are extracted automatically.

---

## Quick Start

=== "Helm (recommended)"

    ```bash
    helm install ruptura oci://ghcr.io/benfradjselim/charts/ruptura \
      --namespace ruptura-system \
      --create-namespace \
      --set apiKey=$(openssl rand -hex 32)

    kubectl port-forward svc/ruptura 8080:80 -n ruptura-system
    curl http://localhost:8080/api/v2/health
    ```

=== "Docker"

    ```bash
    docker run -d \
      --name ruptura \
      -p 8080:8080 \
      -p 4317:4317 \
      -v ruptura-data:/var/lib/ruptura/data \
      -e RUPTURA_API_KEY=$(openssl rand -hex 32) \
      ghcr.io/benfradjselim/ruptura:6.2.2

    curl http://localhost:8080/api/v2/health
    ```

=== "kubectl (inline)"

    ```bash
    # 1 — Namespace + RBAC
    kubectl apply -f - <<'EOF'
    apiVersion: v1
    kind: Namespace
    metadata:
      name: ruptura-system
    ---
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: ruptura
      namespace: ruptura-system
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRole
    metadata:
      name: ruptura
    rules:
      - apiGroups: ["apps"]
        resources: ["deployments","statefulsets","daemonsets","replicasets"]
        verbs: ["get","list","watch"]
      - apiGroups: [""]
        resources: ["pods","nodes","namespaces","services"]
        verbs: ["get","list","watch"]
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRoleBinding
    metadata:
      name: ruptura
    roleRef:
      apiGroup: rbac.authorization.k8s.io
      kind: ClusterRole
      name: ruptura
    subjects:
      - kind: ServiceAccount
        name: ruptura
        namespace: ruptura-system
    EOF

    # 2 — API key secret
    kubectl create secret generic ruptura-secrets \
      -n ruptura-system \
      --from-literal=api-key=$(openssl rand -hex 32)

    # 3 — Storage + Deployment + Services (see Installation for full YAML)
    kubectl port-forward svc/ruptura 8080:80 -n ruptura-system
    curl http://localhost:8080/api/v2/health
    ```

---

## Current Release

**v6.2.2** — all v6.x engineering gaps resolved. Production-ready for Kubernetes evaluation.

- WorkloadRef-native pipeline (`namespace/kind/workload`, not host)
- Adaptive per-workload baselines — no false alarms from batch jobs
- Narrative explain at `/api/v2/explain/{id}/narrative`
- Topology-based contagion from real trace service edges (OTLP)
- Maintenance windows via `/api/v2/suppressions`
- Anomaly REST endpoints at `/api/v2/anomalies`
- Fused Rupture Index (metricR + logR + traceR) in every rupture response
- 37 packages pass `go test -race ./...`

[Full changelog →](community/roadmap.md) · [Getting Started →](getting-started/installation.md) · [API Reference →](api/reference.md)
