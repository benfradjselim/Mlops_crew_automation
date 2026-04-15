<script>
  import { onMount, onDestroy } from 'svelte'
  import { api } from '../lib/api.js'

  export let dashboardId
  export let onBack

  let dashboard = null
  let widgetData = {}   // widgetIndex -> { points, current, unit, max, label }
  let loading = true
  let error = ''
  let host = ''
  let refreshTimer = null

  // ── Unit classification ────────────────────────────────────────────────────
  const PCT_METRICS  = new Set(['cpu_percent','memory_percent','disk_percent',
                                 'container_cpu_percent','container_mem_percent'])
  const RATE_METRICS = new Set(['error_rate','timeout_rate'])
  const BPS_METRICS  = new Set(['net_rx_bps','net_tx_bps'])
  const BYTES_METRICS = new Set(['container_net_rx_bytes','container_net_tx_bytes'])
  const MB_METRICS   = new Set(['container_mem_used_mb'])
  const KPI_NAMES    = new Set(['stress','fatigue','mood','pressure','humidity',
                                 'contagion','resilience','entropy','velocity','health_score'])

  // Returns a human-readable {value, unit} pair
  function humanise(metricName, raw) {
    if (metricName === 'uptime_seconds') {
      return { value: fmtUptime(raw), unit: '' }
    }
    if (BPS_METRICS.has(metricName)) {
      return fmtBits(raw)
    }
    if (BYTES_METRICS.has(metricName)) {
      return fmtBytes(raw)
    }
    if (MB_METRICS.has(metricName)) {
      const gb = raw / 1024
      return gb >= 1
        ? { value: gb.toFixed(2), unit: 'GB' }
        : { value: raw.toFixed(1), unit: 'MB' }
    }
    if (PCT_METRICS.has(metricName) || RATE_METRICS.has(metricName)) {
      return { value: raw.toFixed(1), unit: '%' }
    }
    if (metricName === 'request_rate') {
      return { value: raw.toFixed(1), unit: 'req/s' }
    }
    if (metricName === 'load_avg_1' || metricName === 'load_avg_5' || metricName === 'load_avg_15') {
      return { value: raw.toFixed(2), unit: '' }
    }
    if (metricName === 'processes') {
      return { value: Math.round(raw).toString(), unit: 'procs' }
    }
    if (KPI_NAMES.has(metricName)) {
      // KPIs are [0,1] stored; multiply by 100 for display
      const pct = raw * 100
      return { value: pct.toFixed(1), unit: '%' }
    }
    return { value: raw.toFixed(2), unit: '' }
  }

  // For display in labels: kpi comes in already scaled [0,1] from storage
  function kpiHumanise(kpiName, raw) {
    const pct = raw * 100
    if (kpiName === 'health_score') return { value: pct.toFixed(1), unit: '/100' }
    return { value: pct.toFixed(1), unit: '%' }
  }

  function fmtUptime(seconds) {
    const s = Math.round(seconds)
    if (s < 60)  return s + 's'
    if (s < 3600) return Math.floor(s / 60) + 'm ' + (s % 60) + 's'
    const h = Math.floor(s / 3600)
    const m = Math.floor((s % 3600) / 60)
    if (h < 24)  return h + 'h ' + m + 'm'
    const d = Math.floor(h / 24)
    return d + 'd ' + (h % 24) + 'h'
  }

  function fmtBits(bps) {
    if (bps >= 1e9)  return { value: (bps / 1e9).toFixed(2),  unit: 'Gbps' }
    if (bps >= 1e6)  return { value: (bps / 1e6).toFixed(2),  unit: 'Mbps' }
    if (bps >= 1e3)  return { value: (bps / 1e3).toFixed(1),  unit: 'Kbps' }
    return { value: bps.toFixed(0), unit: 'bps' }
  }

  function fmtBytes(b) {
    if (b >= 1073741824) return { value: (b / 1073741824).toFixed(2), unit: 'GB' }
    if (b >= 1048576)    return { value: (b / 1048576).toFixed(2),    unit: 'MB' }
    if (b >= 1024)       return { value: (b / 1024).toFixed(1),       unit: 'KB' }
    return { value: b.toFixed(0), unit: 'B' }
  }

  // Format a timestamp ms into a short time label for the X axis
  function fmtTime(ms) {
    const d = new Date(ms)
    const hh = String(d.getHours()).padStart(2, '0')
    const mm = String(d.getMinutes()).padStart(2, '0')
    return hh + ':' + mm
  }

  async function detectHost() {
    try { const r = await api.health(); host = r.data?.host || '' } catch {}
  }

  async function loadWidget(idx, widget) {
    const isKpi    = !!widget.kpi
    const key      = isKpi ? widget.kpi : widget.metric
    if (!key) return

    try {
      let points  = []
      let current = 0
      let humanVal = ''
      let unit = ''

      if (widget.type === 'timeseries') {
        let raw = []
        if (isKpi) {
          // GET /api/v1/kpis/{name}?host=...
          const r = await api.kpi(key, host)
          raw = (r.data?.points || []).map(p => ({
            timestamp: p.timestamp,
            value:     p.value
          }))
        } else {
          const r = await api.metricRange(key, host, '-1h')
          raw = r.data?.points || []
        }

        // Sort by time ascending (fixes unsorted API responses)
        points = raw
          .map(p => ({ t: new Date(p.timestamp).getTime(), v: p.value }))
          .filter(p => !isNaN(p.t) && p.t > 1000000)
          .sort((a, b) => a.t - b.t)

        const lastV = points.length ? points[points.length - 1].v : 0
        const h = isKpi ? kpiHumanise(key, lastV) : humanise(key, lastV)
        humanVal = h.value
        unit     = h.unit
        current  = isKpi ? lastV * 100 : lastV

      } else {
        // gauge / stat — get current snapshot value
        if (isKpi) {
          const r = await api.kpis(host)
          const snap = r.data || {}
          // snap has {stress:{value,state}, fatigue:{value,state}, ...}
          const kpiKey = key === 'health_score' ? 'HealthScore'
                       : key.charAt(0).toUpperCase() + key.slice(1)
          const kpiObj = snap[kpiKey] || snap[key] || {}
          const raw = typeof kpiObj === 'object' ? (kpiObj.value ?? 0) : (kpiObj ?? 0)
          const h = kpiHumanise(key, raw)
          humanVal = h.value
          unit     = h.unit
          current  = raw * 100
        } else {
          const r = await api.metrics(host)
          const raw = (r.data?.metrics || {})[key] ?? 0
          const h = humanise(key, raw)
          humanVal = h.value
          unit     = h.unit
          current  = PCT_METRICS.has(key) || RATE_METRICS.has(key) ? raw : raw
        }
      }

      const max = (PCT_METRICS.has(key) || RATE_METRICS.has(key) || KPI_NAMES.has(key)) ? 100 : undefined

      widgetData = {
        ...widgetData,
        [idx]: { points, current, humanVal, unit, max }
      }
    } catch {}
  }

  async function loadAll() {
    if (!dashboard) return
    await Promise.all(dashboard.widgets.map((w, i) => loadWidget(i, w)))
  }

  onMount(async () => {
    loading = true
    try {
      await detectHost()
      const r = await api.dashboardGet(dashboardId)
      dashboard = r.data
    } catch (e) {
      error = e.message
    } finally {
      loading = false
    }
    await loadAll()
    if (dashboard?.refresh_seconds > 0) {
      refreshTimer = setInterval(loadAll, dashboard.refresh_seconds * 1000)
    }
  })

  onDestroy(() => clearInterval(refreshTimer))

  // ── SVG helpers ────────────────────────────────────────────────────────────
  function buildPath(points, width, height, pad) {
    if (!points || points.length < 2) return ''
    const xs = points.map(p => p.t)
    const ys = points.map(p => p.v)
    const xMin = Math.min(...xs), xMax = Math.max(...xs)
    const yMin = 0, yMax = Math.max(...ys, 0.001)
    const sx = v => pad + (v - xMin) / (xMax - xMin || 1) * (width - pad * 2)
    const sy = v => height - pad - (v - yMin) / (yMax - yMin || 1) * (height - pad * 2)
    return points.map((p, i) => (i === 0 ? 'M' : 'L') + sx(p.t).toFixed(1) + ',' + sy(p.v).toFixed(1)).join(' ')
  }

  function buildArea(points, width, height, pad) {
    const path = buildPath(points, width, height, pad)
    if (!path) return ''
    const xs = points.map(p => p.t)
    const xMin = Math.min(...xs), xMax = Math.max(...xs)
    const sx = v => pad + (v - xMin) / (xMax - xMin || 1) * (width - pad * 2)
    const bY = height - pad
    const lastX = sx(xs[xs.length - 1]).toFixed(1)
    const firstX = sx(xs[0]).toFixed(1)
    return path + ` L${lastX},${bY} L${firstX},${bY} Z`
  }

  // Build 4 evenly-spaced time tick labels along the X axis
  function buildTimeTicks(points, width, height, pad) {
    if (!points || points.length < 2) return []
    const xs = points.map(p => p.t)
    const xMin = Math.min(...xs), xMax = Math.max(...xs)
    const sx = v => pad + (v - xMin) / (xMax - xMin || 1) * (width - pad * 2)
    const ticks = [0, 0.33, 0.66, 1].map(f => {
      const ts = xMin + f * (xMax - xMin)
      return { x: sx(ts).toFixed(1), label: fmtTime(ts) }
    })
    return ticks
  }

  // Build Y-axis labels (min, mid, max)
  function buildYLabels(points, height, pad) {
    if (!points || points.length < 2) return []
    const ys = points.map(p => p.v)
    const yMax = Math.max(...ys, 0.001)
    const sy = v => height - pad - v / yMax * (height - pad * 2)
    return [
      { y: sy(yMax).toFixed(1), label: yMax >= 1000 ? (yMax/1000).toFixed(1)+'k' : yMax.toFixed(1) },
      { y: sy(yMax / 2).toFixed(1), label: (yMax/2 >= 1000) ? (yMax/2000).toFixed(1)+'k' : (yMax/2).toFixed(1) },
      { y: sy(0).toFixed(1), label: '0' },
    ]
  }

  function arcPath(pct, r, cx, cy) {
    const clamped = Math.min(Math.max(pct, 0), 100)
    const startAngle = -Math.PI * 0.75
    const endAngle = startAngle + (clamped / 100) * Math.PI * 1.5
    const x1 = cx + r * Math.cos(startAngle), y1 = cy + r * Math.sin(startAngle)
    const x2 = cx + r * Math.cos(endAngle),   y2 = cy + r * Math.sin(endAngle)
    const large = endAngle - startAngle > Math.PI ? 1 : 0
    return `M ${x1.toFixed(1)} ${y1.toFixed(1)} A ${r} ${r} 0 ${large} 1 ${x2.toFixed(1)} ${y2.toFixed(1)}`
  }

  function gaugeColor(pct) {
    if (pct >= 80) return '#ef4444'
    if (pct >= 60) return '#f97316'
    if (pct >= 40) return '#eab308'
    return '#22c55e'
  }
</script>

<div class="view">
  <div class="view-header">
    <button class="back-btn" on:click={onBack}>← Boards</button>
    {#if dashboard}
      <h1>{dashboard.name}</h1>
      <span class="refresh-badge">↻ {dashboard.refresh_seconds}s</span>
    {/if}
  </div>

  {#if loading}
    <p class="muted">Loading dashboard…</p>
  {:else if error}
    <p class="err">{error}</p>
  {:else if dashboard}
    <div class="widget-grid">
      {#each dashboard.widgets as widget, idx}
        {@const wd = widgetData[idx]}
        <div class="widget card type-{widget.type}">
          <div class="widget-title">{widget.title}</div>

          {#if widget.type === 'timeseries'}
            <div class="chart-wrap">
              {#if wd?.points?.length > 1}
                <svg viewBox="0 0 340 110" class="chart">
                  <!-- Y-axis grid lines and labels -->
                  {#each buildYLabels(wd.points, 95, 24) as yl}
                    <line x1="36" y1={yl.y} x2="326" y2={yl.y} stroke="#1e293b" stroke-width="1"/>
                    <text x="34" y={+yl.y + 3} fill="#475569" font-size="7" text-anchor="end">{yl.label}</text>
                  {/each}
                  <!-- Axes -->
                  <line x1="36" y1="10" x2="36"  y2="83" stroke="#334155" stroke-width="0.8"/>
                  <line x1="36" y1="83" x2="326" y2="83" stroke="#334155" stroke-width="0.8"/>
                  <!-- Area fill -->
                  <path d={buildArea(wd.points, 340, 95, 36)} fill="#38bdf815" />
                  <!-- Line -->
                  <path d={buildPath(wd.points, 340, 95, 36)} fill="none" stroke="#38bdf8" stroke-width="1.5"/>
                  <!-- X-axis time ticks -->
                  {#each buildTimeTicks(wd.points, 340, 95, 36) as tick}
                    <line x1={tick.x} y1="83" x2={tick.x} y2="86" stroke="#475569" stroke-width="0.8"/>
                    <text x={tick.x} y="95" fill="#475569" font-size="7" text-anchor="middle">{tick.label}</text>
                  {/each}
                  <!-- Latest value badge -->
                  <text x="325" y="12" fill="#94a3b8" font-size="8" text-anchor="end">
                    {wd.humanVal}{wd.unit}
                  </text>
                </svg>
              {:else}
                <div class="no-data">No data yet — waiting for first collection cycle</div>
              {/if}
            </div>

          {:else if widget.type === 'gauge'}
            <div class="gauge-wrap">
              <svg viewBox="0 0 130 95" class="gauge-svg">
                <!-- Background arc -->
                <path d={arcPath(100, 44, 65, 68)} fill="none" stroke="#1e293b" stroke-width="10" stroke-linecap="round"/>
                <!-- Tick marks at 25 / 50 / 75 -->
                {#each [25, 50, 75] as tick}
                  {@const ang = -Math.PI * 0.75 + (tick / 100) * Math.PI * 1.5}
                  <line
                    x1={65 + 36 * Math.cos(ang)}
                    y1={68 + 36 * Math.sin(ang)}
                    x2={65 + 44 * Math.cos(ang)}
                    y2={68 + 44 * Math.sin(ang)}
                    stroke="#334155" stroke-width="1.2"/>
                {/each}
                <!-- Value arc -->
                {#each [wd?.max ? (wd.current / wd.max) * 100 : (wd?.current ?? 0)] as gPct}
                <path d={arcPath(gPct, 44, 65, 68)}
                      fill="none"
                      stroke={gaugeColor(gPct)}
                      stroke-width="10" stroke-linecap="round"/>
                {/each}
                <!-- Center value -->
                <text x="65" y="68" text-anchor="middle" fill="#e2e8f0" font-size="15" font-weight="700">
                  {wd?.humanVal ?? '—'}
                </text>
                <!-- Unit below -->
                <text x="65" y="80" text-anchor="middle" fill="#64748b" font-size="8">
                  {wd?.unit ?? ''}
                </text>
                <!-- Min / Max labels -->
                <text x="24" y="90" fill="#475569" font-size="7">0</text>
                <text x="106" y="90" fill="#475569" font-size="7" text-anchor="end">
                  {wd?.max ?? ''}
                </text>
              </svg>
            </div>

          {:else if widget.type === 'kpi'}
            <div class="kpi-wrap">
              <div class="kpi-val" style="color: {gaugeColor(wd?.current ?? 0)}">
                {wd?.humanVal ?? '—'}
              </div>
              <div class="kpi-unit">{wd?.unit ?? ''}</div>
            </div>

          {:else if widget.type === 'alerts'}
            <div class="stat-val muted-sm">See Alerts tab</div>

          {:else}
            <!-- stat -->
            <div class="kpi-wrap">
              <div class="stat-val">{wd?.humanVal ?? (wd?.current ?? 0).toFixed(2)}</div>
              <div class="kpi-unit">{wd?.unit ?? ''}</div>
            </div>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .view { padding: 0; }
  .view-header { display: flex; align-items: center; gap: 1rem; margin-bottom: 1.25rem; }
  .back-btn { background: #334155; border: none; color: #94a3b8; padding: 0.3rem 0.7rem;
              border-radius: 5px; cursor: pointer; font-size: 0.8rem; }
  .back-btn:hover { color: #e2e8f0; }
  h1 { margin: 0; font-size: 1.1rem; color: #e2e8f0; flex: 1; }
  .refresh-badge { font-size: 0.7rem; color: #475569; }

  .widget-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 1rem;
  }

  .card { background: #1e293b; border: 1px solid #334155; border-radius: 8px;
          padding: 0.75rem 1rem; }

  .widget-title { font-size: 0.72rem; color: #64748b; text-transform: uppercase;
                  letter-spacing: 0.05em; margin-bottom: 0.5rem; }

  /* Timeseries spans 2 columns */
  .type-timeseries { grid-column: span 2; }
  .chart-wrap { width: 100%; }
  .chart { width: 100%; height: auto; display: block; }
  .no-data { color: #475569; font-size: 0.8rem; padding: 1rem 0; text-align: center; }

  /* Gauge */
  .gauge-wrap { display: flex; justify-content: center; }
  .gauge-svg { width: 130px; height: 95px; }

  /* KPI / stat big number */
  .kpi-wrap { display: flex; flex-direction: column; align-items: center; padding: 0.5rem 0; }
  .kpi-val  { font-size: 2.2rem; font-weight: 700; line-height: 1; }
  .kpi-unit { font-size: 0.8rem; color: #64748b; margin-top: 0.2rem; }
  .stat-val { font-size: 1.8rem; font-weight: 600; color: #e2e8f0; }
  .muted-sm { color: #475569; font-size: 0.85rem; text-align: center; padding: 0.5rem; }

  .muted { color: #64748b; }
  .err   { color: #f87171; }
</style>
