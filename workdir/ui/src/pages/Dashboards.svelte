<script>
  import { onMount } from 'svelte'
  import { api } from '../lib/api.js'
  import DashboardView from './DashboardView.svelte'
  import DashboardEdit from './DashboardEdit.svelte'

  let dashboards = [], templates = [], loading = true
  let showNew = false, newName = '', newRefresh = 30, creating = false
  let openId = null
  let editId = null

  // Template apply modal state
  let applyModal = null   // null | { template }
  let applyName  = ''
  let applyMode  = 'current'   // 'current' | 'predicted'
  let applying   = false

  // Active category filter
  let activeCategory = 'All'

  async function load() {
    loading = true
    const [dr, tr] = await Promise.all([
      api.dashboards().catch(() => ({ data: [] })),
      api.templates().catch(() => ({ data: [] })),
    ])
    dashboards = dr.data || []
    templates  = tr.data || []
    loading = false
  }

  async function create() {
    creating = true
    await api.dashboardCreate({ name: newName, refresh_seconds: newRefresh }).catch(() => {})
    showNew = false; newName = ''; creating = false
    load()
  }

  function openApplyModal(t) {
    applyModal = t
    applyName  = ''
    applyMode  = 'current'
  }

  async function confirmApply() {
    if (!applyModal) return
    applying = true
    const name = applyName.trim() || undefined
    await api.templateApply(applyModal.id, name, applyMode).catch(() => {})
    applying = false
    applyModal = null
    load()
  }

  async function del(id) {
    if (!confirm('Delete this dashboard?')) return
    await api.dashboardDelete(id).catch(() => {})
    load()
  }

  onMount(load)

  $: categories = ['All', ...new Set(templates.map(t => t.category).filter(Boolean))]
  $: filtered = activeCategory === 'All'
    ? templates
    : templates.filter(t => t.category === activeCategory)

  const ICONS = {
    'server':        'M20 4H4a2 2 0 00-2 2v3a2 2 0 002 2h16a2 2 0 002-2V6a2 2 0 00-2-2zm0 9H4a2 2 0 00-2 2v3a2 2 0 002 2h16a2 2 0 002-2v-3a2 2 0 00-2-2zM7 7h.01M7 16h.01',
    'activity':      'M22 12h-4l-3 9L9 3l-3 9H2',
    'trending-up':   'M23 6l-9.5 9.5-5-5L1 18M17 6h6v6',
    'box':           'M21 16V8a2 2 0 00-1-1.73l-7-4a2 2 0 00-2 0l-7 4A2 2 0 002 8v8a2 2 0 001 1.73l7 4a2 2 0 002 0l7-4A2 2 0 0021 16z',
    'layers':        'M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5',
    'globe':         'M12 2a10 10 0 100 20A10 10 0 0012 2zM2 12h20M12 2a15.3 15.3 0 010 20M12 2a15.3 15.3 0 000 20',
    'shield':        'M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z',
    'alert-triangle':'M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0zM12 9v4M12 17h.01',
    'bar-chart-2':   'M18 20V10M12 20V4M6 20v-6',
    'eye':           'M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8zM12 9a3 3 0 110 6 3 3 0 010-6z',
    'award':         'M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z',
    'zap':           'M13 2L3 14h9l-1 8 10-12h-9l1-8z',
    'wifi':          'M5 12.55a11 11 0 0114.08 0M1.42 9a16 16 0 0121.16 0M8.53 16.11a6 6 0 016.95 0M12 20h.01',
    'database':      'M12 2C6.48 2 2 4.24 2 7v10c0 2.76 4.48 5 10 5s10-2.24 10-5V7c0-2.76-4.48-5-10-5zM2 12c0 2.76 4.48 5 10 5s10-2.24 10-5M2 7c0 2.76 4.48 5 10 5s10-2.24 10-5',
    'cpu':           'M18 4H6a2 2 0 00-2 2v12a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2zM9 9h6v6H9V9zM9 1v3M15 1v3M9 20v3M15 20v3M20 9h3M20 15h3M1 9h3M1 15h3',
    'percent':       'M19 5L5 19M6.5 6.5a1 1 0 100 2 1 1 0 000-2zM17.5 15.5a1 1 0 100 2 1 1 0 000-2z',
    'git-branch':    'M6 3v12M18 9a3 3 0 100-6 3 3 0 000 6zM6 21a3 3 0 100-6 3 3 0 000 6zM18 9a9 9 0 01-9 9',
  }

  const CATEGORY_COLOR = {
    'Infrastructure': '#0ea5e9',
    'OHE KPIs':       '#a855f7',
    'Kubernetes':     '#3b82f6',
    'SRE':            '#f59e0b',
    'Containers':     '#10b981',
    'Security':       '#ef4444',
    'Applications':   '#6366f1',
    'Prediction':     '#ec4899',
  }
</script>

{#if openId}
  <DashboardView
    dashboardId={openId}
    onBack={() => { openId = null; load() }}
    onEdit={(id) => { openId = null; editId = id }}
  />
{:else if editId}
  <DashboardEdit
    dashboardId={editId}
    onBack={() => { editId = null; load() }}
  />
{:else}

<!-- Apply Template Modal -->
{#if applyModal}
  <div class="modal-backdrop" on:click|self={() => applyModal = null}>
    <div class="modal">
      <div class="modal-header">
        <svg viewBox="0 0 24 24" class="modal-icon" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
          <path d={ICONS[applyModal.icon] || ICONS['activity']}/>
        </svg>
        <div>
          <div class="modal-title">{applyModal.name}</div>
          <div class="modal-desc">{applyModal.description}</div>
        </div>
      </div>

      <div class="modal-field">
        <label class="field-label">Dashboard name (optional)</label>
        <input class="inp" bind:value={applyName} placeholder={applyModal.name} />
      </div>

      <div class="modal-field">
        <label class="field-label">Data mode</label>
        <div class="mode-toggle">
          <button
            class="mode-btn {applyMode === 'current' ? 'active' : ''}"
            on:click={() => applyMode = 'current'}
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mode-icon">
              <path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
            </svg>
            Current
            <span class="mode-hint">Live data · timeseries charts</span>
          </button>
          <button
            class="mode-btn {applyMode === 'predicted' ? 'active' : ''}"
            on:click={() => applyMode = 'predicted'}
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mode-icon">
              <path d="M23 6l-9.5 9.5-5-5L1 18M17 6h6v6"/>
            </svg>
            Predicted
            <span class="mode-hint">ML forecast · trend overlay</span>
          </button>
        </div>
      </div>

      <div class="modal-meta">
        <span class="meta-tag" style="background:{CATEGORY_COLOR[applyModal.category] || '#475569'}22;color:{CATEGORY_COLOR[applyModal.category] || '#94a3b8'}">{applyModal.category}</span>
        <span class="meta-count">{applyModal.widget_count} widgets</span>
        {#each (applyModal.tags || []).slice(0,4) as tag}
          <span class="meta-tag-dim">{tag}</span>
        {/each}
      </div>

      <div class="modal-actions">
        <button class="btn-ghost" on:click={() => applyModal = null}>Cancel</button>
        <button class="btn-primary" on:click={confirmApply} disabled={applying}>
          {applying ? 'Creating…' : 'Apply Template'}
        </button>
      </div>
    </div>
  </div>
{/if}

<div class="page">
  <div class="header">
    <h1>Dashboards</h1>
    <button class="btn" on:click={() => showNew = !showNew}>+ New</button>
  </div>

  {#if showNew}
    <div class="new-form card">
      <input bind:value={newName} placeholder="Dashboard name" class="inp"/>
      <label class="inline-label">Refresh (s):<input type="number" bind:value={newRefresh} min="5" max="3600" class="inp-num"/></label>
      <button class="btn-primary" on:click={create} disabled={!newName || creating}>Create</button>
      <button class="btn-ghost" on:click={() => showNew = false}>Cancel</button>
    </div>
  {/if}

  {#if loading}
    <p class="muted">Loading…</p>
  {:else}

    {#if dashboards.length > 0}
      <section>
        <h2>My Dashboards</h2>
        <div class="grid">
          {#each dashboards as d}
            <div class="dash-card card">
              <div class="dash-name">{d.name}</div>
              <div class="dash-meta">{d.widgets?.length || 0} widgets · refresh {d.refresh_seconds}s</div>
              <div class="dash-actions">
                <button class="btn-sm"        on:click={() => openId = d.id}>View</button>
                <button class="btn-sm edit"   on:click={() => editId = d.id}>Edit</button>
                <button class="btn-sm danger" on:click={() => del(d.id)}>Delete</button>
              </div>
            </div>
          {/each}
        </div>
      </section>
    {:else}
      <p class="muted empty">No dashboards yet — apply a template below to get started.</p>
    {/if}

    {#if templates.length > 0}
      <section class="tmpl-section">
        <div class="tmpl-header">
          <h2>Template Gallery</h2>
          <span class="tmpl-count">{templates.length} templates</span>
        </div>

        <div class="category-pills">
          {#each categories as cat}
            <button
              class="pill {activeCategory === cat ? 'pill-active' : ''}"
              style={activeCategory === cat && cat !== 'All' ? `background:${CATEGORY_COLOR[cat] || '#475569'}22;color:${CATEGORY_COLOR[cat] || '#94a3b8'};border-color:${CATEGORY_COLOR[cat] || '#475569'}55` : ''}
              on:click={() => activeCategory = cat}
            >{cat}</button>
          {/each}
        </div>

        <div class="tmpl-grid">
          {#each filtered as t}
            {@const catColor = CATEGORY_COLOR[t.category] || '#475569'}
            <div class="tmpl-card">
              <div class="tmpl-top">
                <div class="tmpl-icon-wrap" style="background:{catColor}18;color:{catColor}">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"
                       stroke-linecap="round" stroke-linejoin="round" class="tmpl-icon">
                    <path d={ICONS[t.icon] || ICONS['activity']}/>
                  </svg>
                </div>
                <div class="tmpl-info">
                  <div class="tmpl-name">{t.name}</div>
                  <span class="tmpl-cat" style="color:{catColor}">{t.category}</span>
                </div>
              </div>
              <div class="tmpl-desc">{t.description}</div>
              <div class="tmpl-footer">
                <span class="tmpl-wcount">{t.widget_count} widgets</span>
                <div class="tmpl-tags">
                  {#each (t.tags || []).slice(0, 3) as tag}
                    <span class="tag">{tag}</span>
                  {/each}
                </div>
              </div>
              <button class="apply-btn" on:click={() => openApplyModal(t)}>Apply</button>
            </div>
          {/each}
        </div>
      </section>
    {/if}

  {/if}
</div>
{/if}

<style>
  .page { padding: 0; }
  .header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 1rem; }
  h1 { margin: 0; font-size: 1.2rem; color: #e2e8f0; }
  h2 { font-size: 0.72rem; color: #64748b; text-transform: uppercase; letter-spacing: 0.08em; margin: 0 0 0.6rem; }
  section { margin-bottom: 2rem; }

  .card { background: #1e293b; border: 1px solid #334155; border-radius: 8px; padding: 1rem; }
  .new-form { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 1rem; flex-wrap: wrap; }
  .inline-label { display: flex; align-items: center; gap: 4px; font-size: 0.8rem; color: #94a3b8; }
  .inp { background: #0f172a; border: 1px solid #334155; color: #e2e8f0; padding: 0.4rem 0.6rem; border-radius: 5px; font-size: 0.85rem; }
  .inp-num { width: 70px; background: #0f172a; border: 1px solid #334155; color: #e2e8f0; padding: 0.3rem 0.4rem; border-radius: 4px; }
  .grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(210px, 1fr)); gap: 0.75rem; }
  .dash-card { display: flex; flex-direction: column; gap: 0.4rem; }
  .dash-name { font-weight: 600; color: #e2e8f0; font-size: 0.9rem; }
  .dash-meta { font-size: 0.72rem; color: #64748b; }
  .dash-actions { margin-top: auto; display: flex; gap: 0.4rem; flex-wrap: wrap; }

  .btn         { background: #334155; border: none; color: #e2e8f0; padding: 0.35rem 0.75rem; border-radius: 5px; cursor: pointer; font-size: 0.85rem; }
  .btn-primary { background: #0284c7; border: none; color: #fff; padding: 0.4rem 0.9rem; border-radius: 6px; cursor: pointer; font-size: 0.85rem; font-weight: 600; }
  .btn-primary:disabled { opacity: 0.5; cursor: not-allowed; }
  .btn-ghost   { background: transparent; border: 1px solid #334155; color: #94a3b8; padding: 0.4rem 0.9rem; border-radius: 6px; cursor: pointer; font-size: 0.85rem; }
  .btn-sm      { background: #334155; border: none; color: #e2e8f0; padding: 2px 8px; border-radius: 4px; cursor: pointer; font-size: 0.75rem; }
  .btn-sm.edit   { background: #0f3460; color: #38bdf8; }
  .btn-sm.danger { background: #7f1d1d; color: #fca5a5; }

  .tmpl-header { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 0.75rem; }
  .tmpl-count  { font-size: 0.7rem; color: #475569; background: #0f172a; border: 1px solid #334155; border-radius: 10px; padding: 1px 8px; }

  .category-pills { display: flex; flex-wrap: wrap; gap: 0.4rem; margin-bottom: 0.9rem; }
  .pill { background: #1e293b; border: 1px solid #334155; color: #64748b; padding: 3px 10px; border-radius: 20px; cursor: pointer; font-size: 0.72rem; transition: all 0.15s; }
  .pill:hover { border-color: #475569; color: #94a3b8; }
  .pill-active { background: #0f172a; border-color: #475569; color: #e2e8f0; }

  .tmpl-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(255px, 1fr)); gap: 0.75rem; }
  .tmpl-card {
    background: #1e293b; border: 1px solid #334155; border-radius: 10px;
    padding: 0.9rem; display: flex; flex-direction: column; gap: 0.55rem;
    transition: border-color 0.15s, transform 0.1s;
  }
  .tmpl-card:hover { border-color: #475569; transform: translateY(-1px); }

  .tmpl-top { display: flex; align-items: flex-start; gap: 0.65rem; }
  .tmpl-icon-wrap { width: 36px; height: 36px; border-radius: 8px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
  .tmpl-icon { width: 18px; height: 18px; }
  .tmpl-info { display: flex; flex-direction: column; gap: 1px; min-width: 0; }
  .tmpl-name { font-weight: 700; color: #e2e8f0; font-size: 0.87rem; line-height: 1.2; }
  .tmpl-cat  { font-size: 0.67rem; font-weight: 600; letter-spacing: 0.04em; }
  .tmpl-desc { font-size: 0.73rem; color: #64748b; line-height: 1.45; flex: 1; }

  .tmpl-footer { display: flex; align-items: center; gap: 0.5rem; flex-wrap: wrap; }
  .tmpl-wcount { font-size: 0.67rem; color: #475569; }
  .tmpl-tags { display: flex; flex-wrap: wrap; gap: 3px; }
  .tag { background: #0f172a; border: 1px solid #1e3a5f; color: #38bdf8; font-size: 0.6rem; padding: 1px 5px; border-radius: 3px; }

  .apply-btn {
    width: 100%; background: linear-gradient(135deg, #0284c7, #0369a1);
    border: none; color: #fff; padding: 0.38rem 0; border-radius: 6px;
    cursor: pointer; font-size: 0.8rem; font-weight: 600; transition: opacity 0.15s;
    margin-top: auto;
  }
  .apply-btn:hover { opacity: 0.88; }

  /* Modal */
  .modal-backdrop {
    position: fixed; inset: 0; background: rgba(0,0,0,0.65);
    display: flex; align-items: center; justify-content: center; z-index: 200;
  }
  .modal {
    background: #1e293b; border: 1px solid #334155; border-radius: 12px;
    padding: 1.5rem; width: 420px; max-width: 95vw;
    display: flex; flex-direction: column; gap: 1rem;
  }
  .modal-header { display: flex; align-items: flex-start; gap: 0.75rem; }
  .modal-icon { width: 28px; height: 28px; color: #38bdf8; flex-shrink: 0; margin-top: 2px; }
  .modal-title { font-weight: 700; color: #e2e8f0; font-size: 1rem; }
  .modal-desc  { font-size: 0.76rem; color: #64748b; margin-top: 2px; line-height: 1.4; }
  .modal-field { display: flex; flex-direction: column; gap: 0.3rem; }
  .field-label { font-size: 0.7rem; color: #94a3b8; font-weight: 600; text-transform: uppercase; letter-spacing: 0.06em; }
  .modal-field .inp { width: 100%; box-sizing: border-box; }

  .mode-toggle { display: flex; gap: 0.5rem; }
  .mode-btn {
    flex: 1; display: flex; flex-direction: column; align-items: center; gap: 4px;
    padding: 0.65rem 0.4rem; border-radius: 8px;
    background: #0f172a; border: 2px solid #334155; color: #64748b;
    cursor: pointer; font-size: 0.8rem; font-weight: 600; transition: all 0.15s;
  }
  .mode-btn.active { border-color: #0284c7; background: #0284c720; color: #38bdf8; }
  .mode-btn:hover:not(.active) { border-color: #475569; color: #94a3b8; }
  .mode-icon { width: 18px; height: 18px; }
  .mode-hint { font-size: 0.63rem; font-weight: 400; opacity: 0.75; text-align: center; }

  .modal-meta { display: flex; align-items: center; flex-wrap: wrap; gap: 0.4rem; }
  .meta-tag     { font-size: 0.68rem; font-weight: 600; padding: 2px 8px; border-radius: 4px; }
  .meta-tag-dim { font-size: 0.65rem; background: #0f172a; color: #475569; padding: 2px 6px; border-radius: 3px; }
  .meta-count   { font-size: 0.68rem; color: #475569; }
  .modal-actions { display: flex; justify-content: flex-end; gap: 0.5rem; margin-top: 0.25rem; }

  .muted { color: #64748b; font-size: 0.85rem; }
  .empty { margin: 0.25rem 0 1.5rem; }
</style>
