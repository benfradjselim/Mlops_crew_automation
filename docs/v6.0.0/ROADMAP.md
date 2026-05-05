# ROADMAP.md — Kairo Core Development Phases

Document ID: KC-ROAD-001
Date: April 2026
Status: Living document — v6.1.0 RELEASED 2026-04-27
Produced by: Orchestrator (Claude Code)

---

## 1. Dependency DAG

```
Phase -1: AUDIT + MIGRATION + SPECS  (Orchestrator — DONE)
    |
Phase 0:  Cleanup + Rename + Structure (Orchestrator)
    |
Phase 1:  Governance + CI/CD           (Orchestrator — IN PROGRESS)
    |
Phase 2a: ALPHA — Core Engine
    |          (predictor -> pipeline/metrics, pkg/rupture)
    |
    +——————————————————+
    |                  |
Phase 2b: BRAVO     Phase 2c: CHARLIE   (parallel after ALPHA green)
  Pipelines           Fusion + Composites
  (ingest,            (fusion, composites,
  pipeline/logs,      pkg/composites)
  pipeline/traces)
    |                  |
    +——————————————————+
               |
          Phase 3: DELTA
          Actions + Explainability
          (actions/*, explain)
               |
          Phase 4: ECHO
          API + Context + Telemetry + Storage
               |
          Phase 5: FOXTROT
          cmd/kairo-core + SDK + Final integration
               |
          Phase 6: Release
          Tag v6.0.0 + ghcr.io + Helm
```

**Parallelism rules:**
- BRAVO and CHARLIE run in parallel — they share no internal dependencies
- DELTA starts only when ALPHA + BRAVO + CHARLIE are all CI-green
- ECHO starts only when DELTA is CI-green
- FOXTROT starts only when ECHO is CI-green

---

## 2. Phase Details

### Phase -1 — Audit ✅ COMPLETE
**Owner:** Orchestrator
**Outputs:** `docs/v6.0.0/AUDIT.md`, `docs/v6.0.0/MIGRATION.md`
**Exit criteria:** All packages categorized (RÉUTILISER / RÉÉCRIRE / JETER / NEW)

### Phase 0 — Cleanup & Structure
**Owner:** Orchestrator
**Branch:** `v6_main`
**Tasks:**
1. Create branch `v6_main` from current HEAD
2. Update `go.mod` module path to `github.com/benfradjselim/kairo-core`
3. Rename `cmd/ohe/` → `cmd/kairo-core/`
4. Delete JETER packages: `internal/collector/`, `internal/billing/`, `internal/cost/`
5. Delete JETER files: `handlers_audit.go`, `handlers_rbac.go`, `handlers_proxy.go`
6. Create empty skeleton dirs for all NEW packages (SPECS.md §2)
7. Global rename: OHE→Kairo, ohe→kairo, version string→"6.0.0"
8. Verify `go build ./...` passes

**Exit criteria:** `go build ./...` green on `v6_main`

### Phase 0.5 — SPECS Extraction ✅ COMPLETE
**Owner:** Orchestrator
**Output:** `docs/v6.0.0/SPECS.md`

### Phase 1 — Governance + CI/CD ✅ COMPLETE
**Owner:** Orchestrator
**Outputs:**
- `docs/v6.0.0/ROADMAP.md` (this file)
- `docs/v6.0.0/AGENTS.md`
- `docs/v6.0.0/TRACEABILITY.md`
- `docs/v6.0.0/DEV-GUIDE.md`
- `.github/workflows/ci.yml`
- `Dockerfile`
- `.golangci.yml`
- `deploy/helm/kairo-core/`
- `CODEOWNERS`

**Exit criteria:** CI pipeline runs on `v6_main` skeleton (build + vet pass)

### Phase 2a — ALPHA: Core Engine
**Owner:** Orchestrator (light — heritage code reuse)
**Branch:** `v6_alpha`
**Packages:**
- `internal/pipeline/metrics/` — from `internal/predictor/` + `internal/processor/`
- `pkg/rupture/` — public Rupture Index + TTF formulas

**Exit criteria:**
- CI green on `v6_alpha`
- `internal/pipeline/metrics` coverage >= 80%
- `pkg/rupture` coverage >= 85%

### Phase 2b — BRAVO: Signal Pipelines
**Owner:** OpenCode
**Branch:** `v6_bravo`
**Packages:**
- `internal/ingest/` — Prom remote_write + OTLP + DogStatsD + gRPC (merged from receiver + grpcserver)
- `internal/pipeline/logs/` — 4 extractors: error_rate, keyword_counter, burst_detector, novelty_score
- `internal/pipeline/traces/` — 4 analyzers: latency_propagation, bottleneck_score, error_cascade, fanout_pressure

**Exit criteria:**
- CI green on `v6_bravo`
- All three packages coverage >= 80%

### Phase 2c — CHARLIE: Fusion + Composites
**Owner:** OpenCode
**Branch:** `v6_charlie`
**Packages:**
- `internal/fusion/` — weighted Bayesian signal fusion (WP gap: use weighted average as pragmatic default)
- `internal/composites/` — all 8 composite signals: Stress, Fatigue, Pressure, Contagion, Resilience, Entropy, Sentiment, HealthScore
- `pkg/composites/` — public composite formulas (importable)

**Exit criteria:**
- CI green on `v6_charlie`
- All three packages coverage >= 80%

### Phase 3 — DELTA: Actions + Explainability
**Owner:** OpenCode
**Branch:** `v6_delta`
**Packages:**
- `internal/actions/engine/` — rule evaluation, tier determination (Tier 1/2/3)
- `internal/actions/providers/` — Kubernetes, Webhook, Alertmanager, PagerDuty
- `internal/actions/arbitration/` — conflict detection, deduplication
- `internal/actions/safety/` — rate limiting, cooldown, rollback, emergency stop
- `internal/explain/` — XAI trace, formula audit, pipeline debug

**Exit criteria:**
- CI green on `v6_delta`
- `internal/actions/*` coverage >= 75%
- `internal/explain` coverage >= 75%

### Phase 4 — ECHO: API + Integrations
**Owner:** OpenCode
**Branch:** `v6_echo`
**Packages:**
- `internal/api/` — full v2 API; v1 behind `--compat-ohe-v5`; wire all handlers
- `internal/context/` — 4-layer context awareness
- `internal/telemetry/` — self-monitoring, `/metrics`, health endpoint
- `internal/storage/` — strip org isolation; single-tenant key schema

**Exit criteria:**
- CI green on `v6_echo`
- `internal/api` coverage >= 70%
- `internal/context`, `internal/telemetry` coverage >= 75%
- `internal/storage` coverage >= 70%

### Phase 5 — FOXTROT: Integration & Operability
**Owner:** Orchestrator
**Branch:** `v6_foxtrot`
**Packages:**
- `cmd/kairo-core/` — main binary; wire all packages; `kairo.yaml` parsing
- `internal/vault/` — keep as-is
- `internal/eventbus/` — extend for rupture events
- `sdk/go/` — rename + update to v2 API
- `sdk/python/` — rename + update to v2 API

**Exit criteria:**
- CI green on `v6_foxtrot`
- Binary size <= 25 MB (`go build -o kairo-core && du -sh kairo-core`)
- Total coverage >= 70%
- Smoke test: start binary, POST /api/v2/write, GET /api/v2/health returns 200

### Phase 6 — Release ✅ COMPLETE (2026-04-25)
**Owner:** Orchestrator
**Tasks:**
1. Merge all PRs: v6_alpha → v6_main, then BRAVO+CHARLIE (parallel), then DELTA, ECHO, FOXTROT
2. Final coverage check: total >= 70%
3. `git tag v6.0.0 && git push origin v6.0.0`
4. CI Stage 4-5-6 auto-triggers: Docker buildx + push ghcr.io + cosign signature
5. Helm chart package and publish

**Exit criteria:** Image live on `ghcr.io/benfradjselim/kairo-core:v6.0.0`; Helm chart installable

---

## v6.1.0 ✅ RELEASED 2026-04-27

| Item | Spec | Agent | Branch | PR | Coverage |
|------|------|-------|--------|----|---------|
| Real gRPC ingest server | §23 | GOLF | v6.1_golf | #8 | 83.2% |
| NATS/Kafka eventbus (JetStream + franz-go) | §24 | HOTEL | v6.1_hotel | #9 | 88.0% |
| Adaptive ensemble weighting (online MAE) | §25 | INDIA | v6.1_india | #10 | 89.2% |
| Kubernetes operator (KairoInstance CRD) | §26 | JULIET | v6.1_juliet | #11 | 85.1% |
| Go SDK kairo-client-go (full v2 coverage) | — | Orchestrator | v6.1 | direct | — |

---

## v6.2.x — ✅ RELEASED 2026-04-30

Theme: **Fused Rupture Index + workload-level signals**

| Item | Detail |
|------|--------|
| Fused Rupture Index (metricR · logR · traceR) | Three-source fusion; requires ≥ 2 sources to reach critical |
| Workload-level signals (WorkloadRef) | Pods merged per `namespace/kind/name`; adaptive per-workload baselines |
| Narrative explain | `GET /api/v2/explain/{id}/narrative` — structured English causal chain |
| Topology contagion | Real service edges from trace spans with proxy fallback |
| Maintenance windows | `POST/GET/DELETE /api/v2/suppressions` |
| Anomaly REST endpoints (v6.2.2) | `GET /api/v2/anomalies`, `GET /api/v2/anomalies/{host}` |

---

## v6.3.0 — ✅ RELEASED 2026-05-04

Theme: **Calibration + forecasting + simulation**

| Item | Detail |
|------|--------|
| Calibration warm-up | `calibrating` state for first 96 observations; `calibration_progress` + `calibration_eta_minutes` in every snapshot |
| HealthScore trend forecast | OLS regression → `in_15min`, `in_30min`, `critical_eta_minutes` |
| ruptura-sim | Companion binary; 4 patterns: `memory-leak`, `cascade-failure`, `traffic-surge`, `slow-burn` via `POST /api/v2/sim/inject` |

---

## v6.4.0 — ✅ RELEASED 2026-05-05

Theme: **Rupture fingerprinting + business signals**

| Item | Detail |
|------|--------|
| Rupture fingerprinting | 11-dimensional KPI vector per confirmed rupture; cosine similarity ≥ 0.85 → `pattern_match` in response |
| Business signal layer | `slo_burn_velocity`, `blast_radius`, `recovery_debt` in every snapshot's `business` block |

---

## v6.5.0 — ✅ RELEASED 2026-05-05

Theme: **Edition gate**

| Item | Detail |
|------|--------|
| Edition gate | `RUPTURA_EDITION=community|autopilot`; community blocks `POST .../approve` with 402; autopilot enables full execution |

---

## v6.6.0 — ✅ RELEASED 2026-05-05

Theme: **Per-workload signal weight tuning**

| Item | Detail |
|------|--------|
| Per-workload signal weights | `POST/GET /api/v2/config/weights`; selector syntax: exact, `ns/*` prefix, `*` wildcard; `RUPTURA_WORKLOAD_WEIGHTS` env var for K8s bootstrap |

---

## v7.0.0 — PLANNED (Q3 2026)

| Feature | Detail |
|---------|--------|
| `ruptura-ctl` CLI | `ruptura-ctl status`, `ruptura-ctl explain <id>`, `ruptura-ctl suppress <workload> 30m` |
| Web dashboard v2 | Embedded Svelte UI: FusedR heatmap, signal timelines, action log, narrative explain panel |
| Multi-tenant opt-in | X-Org-ID header → namespace filter on all queries; per-org storage namespacing |

---

## 3. Exit Criteria Summary Table

| Phase | Hard Gate | Coverage Gate | Result |
|-------|-----------|--------------|--------|
| 0 | `go build ./...` green | N/A | ✅ DONE |
| 1 | CI pipeline runs (build+vet) | N/A | ✅ DONE |
| 2a ALPHA | CI green | pipeline/metrics >= 80%, pkg/rupture >= 85% | ✅ 89.2% — MERGED PR#1 |
| 2b BRAVO | CI green | ingest >= 80%, pipeline/logs >= 80%, pipeline/traces >= 80% | ✅ 85-86% — MERGED PR#4 |
| 2c CHARLIE | CI green | fusion >= 80%, composites >= 80%, pkg/composites >= 85% | ✅ 85-93% — MERGED PR#3 |
| 3 DELTA | CI green | actions/* >= 75%, explain >= 75% | ✅ 83-100% — MERGED PR#5 |
| 4 ECHO | CI green | api >= 70%, context >= 75%, storage >= 70% | ✅ 72-95% — MERGED PR#6 |
| 5 FOXTROT | CI green + binary <= 25MB | total >= 70% | ✅ 63-88%* — MERGED PR#7 |
| 6 Release | image on ghcr.io | total >= 70% | ✅ TAGGED v6.0.0 — 2026-04-25 |

*cmd/kairo-core à 63% — sous le seuil de 70%, correction planifiée en v6.1.

---

## 4. Branch Strategy

| Branch | Owner | PR Target |
|--------|-------|-----------|
| `v6_main` | Orchestrator | base branch |
| `v6_alpha` | Orchestrator | → v6_main |
| `v6_bravo` | OpenCode | → v6_main |
| `v6_charlie` | OpenCode | → v6_main |
| `v6_delta` | OpenCode | → v6_main |
| `v6_echo` | OpenCode | → v6_main |
| `v6_foxtrot` | Orchestrator | → v6_main |

Merge order enforced by CI required checks:
ALPHA → (BRAVO || CHARLIE) → DELTA → ECHO → FOXTROT → tag v6.0.0

---

Produced: 2026-04-24
Last updated: 2026-05-05 — v6.6.0 released; v7.0.0 planned
