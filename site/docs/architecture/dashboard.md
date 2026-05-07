# Embedded Dashboard

Ruptura ships a self-contained web dashboard served directly by the Go binary. No Grafana, no external server, no CDN — open a browser and go.

## Accessing the dashboard

```bash
# Local / Docker
open http://localhost:8080

# Kubernetes
kubectl port-forward svc/ruptura 8080:80 -n ruptura-system
open http://localhost:8080
```

The dashboard is served at `/` by the same binary that runs the REST API. No additional configuration is required.

## Air-gap compatibility

All assets are embedded in the binary at compile time via `go:embed`:

| Asset | Source | Purpose |
|-------|--------|---------|
| `index.html` | `internal/ui/static/` | Single-page application shell |
| `vendor/chart.min.js` | bundled locally | Chart rendering (no CDN) |
| `vendor/alpine.min.js` | bundled locally | Reactive UI (no CDN) |
| `ruptura-icon.png` | bundled locally | Brand icon in topbar |

Fonts (`Inter`, `JetBrains Mono`) are loaded non-blocking via `<link rel="preload">` with a `<noscript>` fallback. If Google Fonts is unreachable (air-gapped environment), the system UI font stack is used automatically — the dashboard remains fully functional.

## Layout

```
┌─ topbar (logo · workload selector · edition badge · refresh) ─────────┐
│                                                                         │
│  sidebar          center panel              right panel                 │
│  ─────────        ─────────────             ───────────                 │
│  Workload list    KPI signal grid           Action log                  │
│  FusedR badges    FusedR heatmap            Pending approvals           │
│  Health bars      HealthScore timeline      Emergency stop              │
│                   Health forecast           Narrative explain           │
│                   SLO widget                                            │
└─────────────────────────────────────────────────────────────────────────┘
```

## Panels

### Fused Rupture Index heatmap

Colour-coded cells per workload × time interval. Cell colour maps directly to FusedR thresholds:

| Colour | FusedR | State |
|--------|--------|-------|
| Green | < 1.5 | Stable / Elevated |
| Yellow | 1.5 – 3.0 | Warning |
| Orange | 3.0 – 5.0 | Critical |
| Red | ≥ 5.0 | Emergency |

### Per-workload signal timelines

Sparkline charts for all 10 KPI signals (`stress`, `fatigue`, `mood`, `pressure`, `humidity`, `contagion`, `resilience`, `entropy`, `velocity`, `health_score`) over a rolling 30-minute window. Signals are fetched from `GET /api/v2/kpi/{signal}/{namespace}/{workload}`.

### HealthScore forecast

Displays `health_forecast` from the snapshot: current score, trend direction (`improving` / `stable` / `degrading`), projected values at 15 min and 30 min, and `critical_eta_minutes` when degrading.

### SLO widget

Shows `slo_burn_velocity` (error budget burn rate relative to the SLO contract), `blast_radius` (downstream dependency count), and `recovery_debt` (near-miss count over 7 days).

### Action log

Lists all pending and recently resolved actions from `GET /api/v2/actions`. Each entry shows workload, tier, recommended action, and confidence score. Approve and Reject buttons call `POST /api/v2/actions/{id}/approve` and `POST /api/v2/actions/{id}/reject` respectively.

!!! note "Community edition"
    In `community` edition, the Approve button is visible but calls return `402 Payment Required`. Switch to `autopilot` edition (`RUPTURA_EDITION=autopilot`) to enable full action execution.

### Emergency stop

A dedicated button in the right panel calls `POST /api/v2/actions/emergency-stop` to immediately halt all automated Tier-1 action execution. A toast notification confirms the call.

### Narrative explain panel

Clicking any rupture event in the heatmap or action log opens the narrative panel, which fetches `GET /api/v2/explain/{id}/narrative` and displays the structured English explanation:

> "payment-api has been accumulating fatigue for 72h (fatigue 0.81). A contagion wave from payment-db propagated via the payment-api→payment-db call edge and pushed FusedR from 1.8 to 4.2 in 18 minutes. This is a cascade rupture, not an isolated spike. Recommended action: scale payment-api by 2 replicas."

## Auto-refresh

The dashboard polls the API every **15 seconds** by default. The refresh interval and API base URL are configurable via the topbar settings icon. All fetch calls include the `Authorization: Bearer` header when an API key is configured.

## Authentication

If `RUPTURA_API_KEY` is set, the dashboard prompts for the key on first load and stores it in `sessionStorage`. The key is sent as `Authorization: Bearer <key>` on every API call. Clearing the session (closing the tab) clears the stored key.

## Source

The dashboard source lives at `workdir/internal/ui/static/index.html`. It is a single self-contained HTML file — no build step, no bundler, no Node.js required. Changes take effect on the next binary build (`go build ./cmd/ruptura`).
