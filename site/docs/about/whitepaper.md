# Whitepaper

## The Problem

Current observability solutions split along two failing axes:

- **Open-source stacks** (Prometheus + Grafana + Loki) demand 5+ services, 8 GB+ RAM, and weeks of integration. They answer *"What is broken?"* via static thresholds — never *"When will it break?"*
- **Enterprise SaaS** (Datadog, Dynatrace) provide black-box AI at prohibitive cost with opaque decision logic.

Neither predicts. Neither explains.

## The Kairo Approach

Kairo treats infrastructure as a **living organism** — measuring vital signs, behaviours, stress responses, and social dynamics through auditable composite KPIs.

### Rupture Index™ — the core innovation

```
R(t) = |α_burst(t)| / max(|α_stable(t)|, ε)
```

Two Incremental Linear Regression (ILR) windows run in parallel per metric:

| Window | Captures |
|--------|---------|
| 60 min (`ILR_stable`) | Long-term baseline — what is normal? |
| 5 min (`ILR_burst`) | Short-term acceleration — is something happening *right now*? |

When R > 3, a metric is accelerating 3× faster than its own baseline — a memory leak, cascade failure, or saturation event in progress. Kairo detects this **before** the metric reaches 80% saturation.

### Why not LSTM?

| Model | MAE | RAM | Inference | Efficiency Score |
|-------|-----|-----|-----------|-----------------|
| LSTM | 2.0% | 200+ MB | 500 ms | < 0.0001 |
| ARIMA | 4.1% | 85 MB | 210 ms | 0.0001 |
| **ILR (Kairo)** | **6.2%** | **0.5 MB** | **0.8 ms** | **1,550×** |

ILR trades +2.1% MAE for 170× less RAM and 262× faster inference. **1,550× more efficient than ARIMA** — validated on a Raspberry Pi 4 over 40,320 samples.

### 8 Composite Signals

Kairo fuses raw telemetry into 8 interpretable signals with published formulas:

| Signal | Formula | Example state |
|--------|---------|--------------|
| stress | 0.3·CPU + 0.2·RAM + 0.2·Latency + 0.2·Errors + 0.1·Timeouts | 0.72 → "Stressed" |
| fatigue | max(0, F_prev + (stress − 0.3) − λ) | 0.45 → "Tired" |
| pressure | d(stress̄)/dt + ∫errors | 0.15 → stable |
| contagion | Σ E_ij × D_ij | 0.05 → isolated |
| healthscore | (1−stress)(1−fatigue)(1−pressure)(1−contagion)×100 | 43 → needs attention |

### Adaptive Ensemble (v6.1)

Five models weighted by online MAE over a 1-hour window: CA-ILR, ARIMA, Holt-Winters, MAD, EWMA. Weights update every 60 seconds. No manual tuning. No profile configuration.

## Benchmarks

| Criterion | Prom/Grafana/Loki | Datadog | **Kairo Core** |
|-----------|-------------------|---------|---------------|
| RAM idle | ~450 MB | ~180 MB | **22 MB** |
| Setup time | ~30 min | ~5 min | **< 1 min** |
| Prediction accuracy | ❌ None | ✅ Black-box | **✅ Transparent, 6.2% MAE** |
| False positives (backup spikes) | ❌ Yes | ⚠️ Sometimes | **✅ No (λ dissipation)** |
| Exponential crash detection | ❌ No | ✅ Black-box | **✅ R > 3 (auditable)** |
| Air-gapped ready | ⚠️ Complex | ❌ Impossible | **✅ Native** |
| Efficiency score | 1× | ~0.0001× | **1,550×** |

## Full Whitepaper

The v5.0 whitepaper contains the complete mathematical formalization, all KPI formulas, and the canonical METRICS.md standard:

[Read OHE v5.0 Whitepaper (GitHub) →](https://github.com/benfradjselim/Mlops_crew_automation/blob/v6.1/workdir/docs/v5.0.0/WHITEPAPER-v5.0.0.md)

> "Stop staring at dashboards hoping for the best. Sleep. Kairo watches."
>
> — Selim Benfradj, Architect & Founder
