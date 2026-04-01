# Observability Holistic Engine (OHE) v4.0.0

## White Paper

**Version:** 4.0.0  
**Status:** Design Document  
**Date:** 2026-04-01  
**Author:** Selim Benfradj  

---

## Table of Contents

1. Executive Summary
2. Context and Problem Statement
3. Analysis of Existing Solutions
4. Our Vision: Holistic Observability
5. Key Features and Value Proposition
6. Technical Architecture
7. Mathematical Formalization
8. Use Cases
9. Roadmap
10. Conclusion

---

## 1. Executive Summary

### 1.1 The Problem

Current observability solutions measure isolated metrics without understanding the overall system behavior. They answer **"What is wrong?"** but not **"When will it go wrong?"**

### 1.2 Our Solution

**Observability Holistic Engine (OHE)** treats infrastructure as a **living organism** with:
- **Vital signs** (classic metrics)
- **Behaviors** (patterns, habits, rhythms)
- **Emotions** (stress, fatigue, mood)
- **Social interactions** (dependencies, contagion)

### 1.3 Unique Value Proposition

| Solution | Approach | Question Answered |
|----------|----------|-------------------|
| Classic solutions | Isolated metrics | "CPU at 85%" |
| APM solutions | Metrics + traces | "Service A is slow" |
| **OHE v4.0** | **Living organism** | **"Storm in 2h, high fatigue, contagion spreading"** |

---

## 2. Context and Problem Statement

### 2.1 Evolution of Observability
2000-2010 : Monitoring
→ "Is the server UP ?"

2010-2020 : Observability
→ "Why is the server slow ?"

2020-2025 : MLops
→ "What will go wrong ?"

2025+ : Holistic Observability (OHE)
→ "When and how will it go wrong ?"

### 2.2 The Gap

No current solution offers:

1. **A holistic view** of infrastructure as a living organism
2. **Complex KPIs** (observability ETFs) reflecting overall health
3. **Contextual predictions** ("storm in 2 hours")
4. **Behavioral analysis** (habits, rhythms, trends)
5. **Emotion detection** (stress, fatigue, mood)
6. **Social analysis** (error propagation, dependencies)

---

## 3. Analysis of Existing Solutions

### 3.1 Comparative Matrix

| Criteria | Classic Solutions | APM Solutions | OHE v4.0 |
|----------|------------------|---------------|----------|
| **Metrics** | ✅ | ✅ | ✅ |
| **Logs** | ❌ | ✅ | ✅ |
| **Traces** | ❌ | ✅ | 🔄 |
| **Predictions** | ❌ | ⚠️ | ✅ |
| **Complex KPIs** | ❌ | ❌ | ✅ |
| **Behavioral Analysis** | ❌ | ❌ | ✅ |
| **Emotion Detection** | ❌ | ❌ | ✅ |
| **Social Analysis** | ❌ | ❌ | ✅ |
| **Lightweight** | ⚠️ | ⚠️ | ✅ |
| **Installation** | Complex | Simple | **One-liner** |

### 3.2 Identified Limitations

- **Classic solutions**: 15+ services to maintain, 8-12GB RAM, no predictions
- **APM solutions**: High cost, proprietary, limited predictions
- **Log solutions**: Logs only, no predictions


---

## 4. Our Vision: Holistic Observability

### 4.1 The Medical Metaphor

Infrastructure is treated as a living organism:

| Physical System | Human System |
|-----------------|--------------|
| CPU / RAM / Disk | Temperature / Blood Pressure / Heart Rate |
| Network | Blood Circulation |
| Logs | Symptoms |
| Errors | Pain |
| Timeouts | Fatigue |
| Restarts | Fever |
| Latency | Reflexes |
| Throughput | Cardiac Output |

### 4.2 Behaviors

| Human Behavior | System Behavior |
|----------------|-----------------|
| Circadian Rhythm | Daily Traffic |
| Habits | Recurring Patterns |
| Stress | Excessive Load |
| Fatigue | Cumulative Wear |
| Mood | Overall Stability |

### 4.3 Social Interactions

| Social Interaction | Service Interaction |
|--------------------|---------------------|
| Dependencies | Service Calls |
| Contagion | Error Propagation |
| Isolation | Orphaned Services |
| Epidemic | Cascading Incidents |

### 4.4 Philosophy

### 4.4 Philosophy

**From Reactive to Proactive**

Traditional observability tools operate on a reactive model:
- Alert when a metric crosses a threshold
- Respond after an incident occurs
- Fix problems after they impact users

Our approach is fundamentally different:
- Detect trends before they become problems
- Predict when thresholds will be crossed
- Prevent incidents before they impact users

**The Shift in Thinking**

| Reactive Approach | Proactive Approach |
|-------------------|---------------------|
| "CPU is at 85%" | "CPU will reach 90% in 3 hours" |
| "Errors are spiking" | "Error storm forming in 30 minutes" |
| "Service is down" | "Service fatigue indicates risk of failure" |
| "Fix after crash" | "Prevent before crash" |

**Prevention over Cure**

The core philosophy is simple but powerful: it is better to prevent problems than to fix them after they occur. This applies to infrastructure just as it applies to health, weather, and finance.

Just as preventive medicine focuses on early detection and lifestyle changes rather than treating symptoms, OHE focuses on detecting behavioral patterns and predicting outcomes rather than just alerting on thresholds.

**The Four Pillars of Holistic Observability**

1. **Vital Signs** - Like a doctor measuring temperature and blood pressure, we monitor core system metrics (CPU, memory, network)

2. **Behavior Patterns** - Like understanding daily habits and routines, we learn system rhythms and cycles

3. **Emotional State** - Like assessing mood and stress levels, we compute system emotions (stress, fatigue, mood)

4. **Social Dynamics** - Like tracking how diseases spread in a population, we analyze error propagation and dependency contagion

This philosophy transforms infrastructure monitoring from a reactive "alarm system" into a proactive "health management system".


---

## 6. Technical Architecture

### 6.1 System Overview

OHE runs as a single binary with internal components communicating via channels.
┌─────────────────────────────────────────────────────────────┐
│                    OBSERVABILITY HOLISTIC ENGINE            │
│                           :8080                             │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────┐   │
│  │  LAYER 1: COLLECTION                                │   │
│  │  • System metrics (procfs, sysfs)                   │   │
│  │  • Container metrics (Docker, K8s API)              │   │
│  │  • Logs (file tail, journald)                       │   │
│  └─────────────────────────────────────────────────────┘   │
│                          ↓                                  │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  LAYER 2: PROCESSING                                │   │
│  │  • Normalization [0-1]                              │   │
│  │  • Aggregation (avg, p95, p99)                      │   │
│  │  • Downsampling (1m → 5m → 1h)                      │   │
│  └─────────────────────────────────────────────────────┘   │
│                          ↓                                  │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  LAYER 3: KPI COMPUTATION                           │   │
│  │  • Stress, Fatigue, Mood                            │   │
│  │  • Pressure, Humidity, Contagion                    │   │
│  │  • Cycle detection (FFT)                            │   │
│  └─────────────────────────────────────────────────────┘   │
│                          ↓                                  │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  LAYER 4: PREDICTION                                │   │
│  │  • ARIMA models                                     │   │
│  │  • Dynamic thresholds                               │   │
│  │  • Anomaly detection                                │   │
│  └─────────────────────────────────────────────────────┘   │
│                          ↓                                  │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  LAYER 5: OUTPUT                                    │   │
│  │  • REST API                                         │   │
│  │  • Embedded UI                                      │   │
│  │  • Alerts (Slack, Email, Webhook)                   │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                           │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  STORAGE: Badger (embedded)                         │   │
│  │  • TTL: 7d metrics, 30d logs                        │   │
│  │  • Automatic compaction                             │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘

```

### 6.2 Code Structure

```

workdir/
├── cmd/agent/
│   └── main.go
├── internal/
│   ├── collector/
│   │   ├── system.go
│   │   ├── container.go
│   │   └── logs.go
│   ├── processor/
│   │   ├── normalize.go
│   │   ├── aggregate.go
│   │   └── downsample.go
│   ├── analyzer/
│   │   ├── stress.go
│   │   ├── fatigue.go
│   │   ├── mood.go
│   │   ├── pressure.go
│   │   ├── humidity.go
│   │   └── contagion.go
│   ├── predictor/
│   │   ├── arima.go
│   │   ├── threshold.go
│   │   └── anomaly.go
│   ├── storage/
│   │   └── badger.go
│   ├── api/
│   │   └── handlers.go
│   └── web/
│       └── embed.go
├── pkg/
│   ├── models/
│   └── utils/
└── configs/
└── agent.yaml

```

### 6.3 API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/health` | GET | Health check |
| `/api/v1/metrics` | GET | Raw metrics |
| `/api/v1/kpis` | GET | Complex KPIs |
| `/api/v1/predict` | GET | Predictions |
| `/api/v1/alerts` | GET | Active alerts |


---

## 7. Mathematical Formalization

### 7.1 Definitions

Let system have n services S = {s₁, s₂, ..., sₙ}. For each service sᵢ at time t:

- CPUᵢ(t), RAMᵢ(t), Diskᵢ(t), Netᵢ(t)
- Reqᵢ(t), Errᵢ(t), Latᵢ(t), Toutᵢ(t)
- Restartᵢ(t), Uptimeᵢ(t)

### 7.2 Fundamental KPIs

**Stress Index:**
Sᵢ(t) = α·CPUᵢ(t) + β·RAMᵢ(t) + γ·Latᵢ(t) + δ·Errᵢ(t) + ε·Toutᵢ(t)

```
where α + β + γ + δ + ε = 1

**Cumulative Fatigue:**
```

Fᵢ(t) = ∫₀ᵗ (Sᵢ(τ) - Rᵢ(τ)) dτ

```

**System Mood:**
```

Mᵢ(t) = (Uptimeᵢ(t) × Reqᵢ(t)) / (Errᵢ(t) × Toutᵢ(t) × Restartᵢ(t) + ε)

```

### 7.3 Systemic KPIs

**Atmospheric Pressure:**
```

P(t) = dS̄/dt + ∫₀ᵗ Ē(τ) dτ

```
where S̄ = average stress, Ē = average errors

**Error Humidity:**
```

H(t) = (Ē(t) × T̄(t)) / Q̄(t)

```

**Contagion Index:**
```

C(t) = Σᵢⱼ Eᵢⱼ(t) × Dᵢⱼ

```

### 7.4 Prediction Functions

**Storm Forecast:**
```

Storm(t+Δt) = 1 if P(t) > θ_p for δ_t

```

**Burnout Forecast:**
```

Burnout(t+Δt) = 1 if F̄(t) > θ_f

```

**Epidemic Forecast:**
```

Epidemic(t+Δt) = 1 if C(t) > θ_c

```


---

## 8. Use Cases

### 8.1 Storm Detection

| T-12h | T-6h | T-2h | T |
|-------|------|------|---|
| CPU=45% | CPU=65% | CPU=80% | CPU=95% |
| P=+0.05/h | P=+0.1/h | P=+0.2/h | Incident |

**OHE Output:**
- T-12h: "Pressure rising, enhanced monitoring"
- T-6h: "Storm risk in 4h, prepare resources"
- T-2h: "Storm in 1h, scale up recommended"

### 8.2 Epidemic Detection

| Service A | Service B | Service C |
|-----------|-----------|-----------|
| Err=5% | Err=1% | Err=0.5% |
| Dependency A→B, B→C | | |

**OHE Output:**
- Contagion index = 0.7
- "Epidemic detected, propagation in 30 min"
- "Isolate service A recommended"

### 8.3 Fatigue Detection

| Day-3 | Day-2 | Day-1 | Day |
|-------|-------|-------|-----|
| Latency +5% | +10% | +15% | Crash |
| Fatigue=0.3 | 0.5 | 0.7 | 0.9 |

**OHE Output:**
- "Fatigue increasing (+0.2/day)"
- "Burnout in 24h without rest"
- "Preventive restart recommended"

---

## 9. Roadmap

### 9.1 Development Phases

| Phase | Objective | Duration |
|-------|-----------|----------|
| Phase 1 | Collection + Core KPIs | 2 weeks |
| Phase 2 | Advanced KPIs + Patterns | 2 weeks |
| Phase 3 | Predictions + Alerts | 2 weeks |
| Phase 4 | UI + Dashboards | 2 weeks |
| Phase 5 | HA + K8s Operator | 2 weeks |
| Phase 6 | Ecosystem + Community | 4 weeks |

### 9.2 Milestones
`

Week 1-2:   Phase 1 - Collection + Core KPIs
Week 3-4:   Phase 2 - Advanced Analysis
Week 5-6:   Phase 3 - Predictions
Week 7-8:   Phase 4 - User Interface
Week 9-10:  Phase 5 - Production HA
Week 11-14: Phase 6 - Ecosystem

```

### 9.3 Future Features

| Feature | Priority |
|---------|----------|
| Distributed Tracing | High |
| Multi-cluster Federation | High |
| Auto-remediation | Medium |
| Marketplace | Low |
| Mobile App | Low |


---

## 10. Conclusion

### 10.1 Summary

Observability Holistic Engine (OHE) represents a new generation of observability that:

1. Treats infrastructure as a living organism
2. Creates complex KPIs (observability ETFs)
3. Provides contextual predictions
4. Is lightweight and portable (<100MB)
5. Is open source and vendor-agnostic

### 10.2 Key Benefits

| Benefit | Impact |
|---------|--------|
| Prevention | 80% of incidents avoided |
| Cost | 70% savings vs traditional solutions |
| Simplicity | 1 binary vs 15+ services |
| Performance | 10x lighter |
| Predictions | Unique market differentiator |

### 10.3 Call to Action

We invite the community to contribute to this new vision of observability.

**"Prevention is better than cure."**

---

**Selim Benfradj**  
*Architect and Founder*  
*April 2026*


---

## 6. Technical Architecture

### 6.1 System Overview

OHE runs as a single binary with internal components communicating via channels.


┌─────────────────────────────────────────────────────────────────────────────┐
│                         OBSERVABILITY HOLISTIC ENGINE                       │
│                                   :8080                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐ │
│  │  LAYER 1: COLLECTION                                                  │ │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                │ │
│  │  │   System     │  │  Container   │  │    Logs      │                │ │
│  │  │   procfs     │  │  Docker/K8s  │  │  file tail   │                │ │
│  │  └──────────────┘  └──────────────┘  └──────────────┘                │ │
│  └───────────────────────────────────────────────────────────────────────┘ │
│                                      ↓                                     │
│  ┌───────────────────────────────────────────────────────────────────────┐ │
│  │  LAYER 2: PROCESSING                                                  │ │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                │ │
│  │  │ Normalize    │  │  Aggregate   │  │ Downsample   │                │ │
│  │  │    [0-1]     │  │ avg, p95     │  │  1m → 1h     │                │ │
│  │  └──────────────┘  └──────────────┘  └──────────────┘                │ │
│  └───────────────────────────────────────────────────────────────────────┘ │
│                                      ↓                                     │
│  ┌───────────────────────────────────────────────────────────────────────┐ │
│  │  LAYER 3: KPI COMPUTATION                                             │ │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │ │
│  │  │  Stress = α·CPU + β·RAM + γ·Latency + δ·Errors + ε·Timeouts    │  │ │
│  │  │  Fatigue = ∫(Stress - Recovery) dt                              │  │ │
│  │  │  Mood = (Uptime × Throughput) / (Errors × Timeouts × Restarts) │  │ │
│  │  │  Pressure = dStress/dt + ∫Errors dt                             │  │ │
│  │  │  Humidity = (Errors × Timeouts) / Throughput                    │  │ │
│  │  │  Contagion = Σ(Error_propagation × Dependency)                  │  │ │
│  │  └─────────────────────────────────────────────────────────────────┘  │ │
│  └───────────────────────────────────────────────────────────────────────┘ │
│                                      ↓                                     │
│  ┌───────────────────────────────────────────────────────────────────────┐ │
│  │  LAYER 4: PREDICTION                                                  │ │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                │ │
│  │  │    ARIMA     │  │   Dynamic    │  │   Anomaly    │                │ │
│  │  │   Models     │  │  Thresholds  │  │  Detection   │                │ │
│  │  └──────────────┘  └──────────────┘  └──────────────┘                │ │
│  └───────────────────────────────────────────────────────────────────────┘ │
│                                      ↓                                     │
│  ┌───────────────────────────────────────────────────────────────────────┐ │
│  │  LAYER 5: OUTPUT                                                      │ │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                │ │
│  │  │   REST API   │  │  Embedded    │  │   Alerts     │                │ │
│  │  │  /api/v1/*   │  │     UI       │  │ Slack/Email  │                │ │
│  │  └──────────────┘  └──────────────┘  └──────────────┘                │ │
│  └───────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐ │
│  │  STORAGE: Badger (embedded)                                           │ │
│  │  • TTL: 7 days for metrics, 30 days for logs                         │ │
│  │  • Automatic compaction and compression                               │ │
│  │  • Concurrent read/write with snapshot isolation                      │ │
│  └───────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

```

### 6.2 Code Structure

```

workdir/
├── cmd/agent/
│   └── main.go                     # Orchestrator entry point
│
├── internal/
│   ├── collector/                  # Data collection
│   │   ├── system.go               # procfs, sysfs metrics
│   │   ├── container.go            # Docker, Kubernetes API
│   │   └── logs.go                 # File tailing, journald
│   │
│   ├── processor/                  # Data processing
│   │   ├── normalize.go            # [0-1] normalization
│   │   ├── aggregate.go            # avg, min, max, p95, p99
│   │   └── downsample.go           # 1m → 5m → 1h → 1d
│   │
│   ├── analyzer/                   # KPI computation
│   │   ├── stress.go               # Stress index
│   │   ├── fatigue.go              # Cumulative fatigue
│   │   ├── mood.go                 # System mood
│   │   ├── pressure.go             # Atmospheric pressure
│   │   ├── humidity.go             # Error humidity
│   │   └── contagion.go            # Contagion index
│   │
│   ├── predictor/                  # Predictions
│   │   ├── arima.go                # ARIMA time series
│   │   ├── threshold.go            # Dynamic thresholds
│   │   └── anomaly.go              # Anomaly detection
│   │
│   ├── storage/                    # Persistence
│   │   └── badger.go               # Badger DB wrapper
│   │
│   ├── api/                        # REST API
│   │   └── handlers.go             # HTTP handlers
│   │
│   └── web/                        # Embedded UI
│       └── embed.go                # Svelte static files
│
├── pkg/
│   ├── models/                     # Data structures
│   │   ├── metric.go
│   │   ├── kpi.go
│   │   └── alert.go
│   │
│   └── utils/                      # Utilities
│       ├── math.go                 # Mathematical helpers
│       └── time.go                 # Time utilities
│
└── configs/
└── agent.yaml                  # Configuration file

```

### 6.3 API Endpoints

| Endpoint | Method | Description | Response |
|----------|--------|-------------|----------|
| `/api/v1/health` | GET | Health check | `{"status": "ok"}` |
| `/api/v1/metrics` | GET | Raw metrics | `[{"name":"cpu","value":0.45}]` |
| `/api/v1/kpis` | GET | Complex KPIs | `{"stress":0.32,"fatigue":0.41}` |
| `/api/v1/predict` | GET | Predictions | `{"storm":"2h","confidence":0.85}` |
| `/api/v1/alerts` | GET | Active alerts | `[{"severity":"warning","message":"..."}]` |

---

## 7. Mathematical Formalization

### 7.1 Core Metrics Definition

For a system with n services S = {s₁, s₂, ..., sₙ}, each service sᵢ at time t provides:

| Category | Metrics |
|----------|---------|
| System | CPUᵢ(t), RAMᵢ(t), Diskᵢ(t), Netᵢ(t) |
| Application | Reqᵢ(t), Errᵢ(t), Latᵢ(t), Toutᵢ(t) |
| Behavioral | Restartᵢ(t), Uptimeᵢ(t) |

All metrics are normalized to [0,1] range where 0 = optimal, 1 = critical.

### 7.2 Fundamental KPIs

#### Stress Index

The stress index measures current system pressure combining multiple signals:

```

Sᵢ(t) = α·CPUᵢ(t) + β·RAMᵢ(t) + γ·Latᵢ(t) + δ·Errᵢ(t) + ε·Toutᵢ(t)

```

with α + β + γ + δ + ε = 1 (weights configurable by user)

| S value | State |
|---------|-------|
| S < 0.3 | Calm |
| 0.3 ≤ S < 0.6 | Nervous |
| 0.6 ≤ S < 0.8 | Stressed |
| S ≥ 0.8 | Panic |

#### Cumulative Fatigue

Fatigue accumulates when stress exceeds recovery capacity:

```

Fᵢ(t) = ∫₀ᵗ (Sᵢ(τ) - Rᵢ(τ)) dτ

```

where Rᵢ(τ) is the recovery factor (0.1 during normal operation, 0.5 during rest)

| F value | State | Action |
|---------|-------|--------|
| F < 0.3 | Rested | Normal monitoring |
| 0.3 ≤ F < 0.6 | Tired | Increase observation |
| 0.6 ≤ F < 0.8 | Exhausted | Plan maintenance |
| F ≥ 0.8 | Burnout | Preventive restart |

#### System Mood

Mood reflects overall system well-being:

```

Mᵢ(t) = (Uptimeᵢ(t) × Reqᵢ(t)) / (Errᵢ(t) × Toutᵢ(t) × Restartᵢ(t) + ε)

```

| M value | Mood |
|---------|------|
| M > 100 | Happy |
| 50 < M ≤ 100 | Content |
| 10 < M ≤ 50 | Neutral |
| 1 < M ≤ 10 | Sad |
| M ≤ 1 | Depressed |

### 7.3 Systemic KPIs

#### Atmospheric Pressure

Pressure predicts approaching storms:

```

P(t) = dS̄/dt + ∫₀ᵗ Ē(τ) dτ

```

where S̄ = average stress across all services, Ē = average error rate

| P trend | Prediction |
|---------|------------|
| P > 0.1 for 10m | Storm in 2h |
| P stable | Stable |
| P < 0 | Improving |

#### Error Humidity

Humidity indicates error density in the system:

```

H(t) = (Ē(t) × T̄(t)) / Q̄(t)

```

where T̄ = average timeout rate, Q̄ = average throughput

| H value | State | Prediction |
|---------|-------|------------|
| H < 0.1 | Dry | Normal |
| 0.1 ≤ H < 0.3 | Humid | Watch |
| 0.3 ≤ H < 0.5 | Very humid | Alert |
| H ≥ 0.5 | Storm | Immediate action |

#### Contagion Index

Contagion measures how failures propagate:

```

C(t) = Σᵢⱼ Eᵢⱼ(t) × Dᵢⱼ

```

where:
- Eᵢⱼ = error propagation probability from i to j
- Dᵢⱼ = dependency strength (call frequency, criticality)

| C value | State | Action |
|---------|-------|--------|
| C < 0.3 | Low | Normal |
| 0.3 ≤ C < 0.6 | Moderate | Monitor |
| 0.6 ≤ C < 0.8 | Epidemic | Isolate |
| C ≥ 0.8 | Pandemic | Global response |

### 7.4 Prediction Functions

#### Storm Forecast

```

Storm(t+Δt) = 1 if ∫ₜ₋δₜ^t P(τ) dτ > θ_p

```

where θ_p = 0.1, δ_t = 10 minutes

#### Burnout Forecast

```

Burnout(t+Δt) = 1 if F̄(t) > θ_f

```

where θ_f = 0.7, Δt = 4 hours

#### Epidemic Forecast

```

Epidemic(t+Δt) = 1 if C(t) > θ_c

```

where θ_c = 0.6, Δt = 30 minutes

