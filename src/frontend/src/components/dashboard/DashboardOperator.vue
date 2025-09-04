<!-- src/frontend/src/views/dashboard/DashboardOperator.vue -->
<template>
  <div class="opdash-page">
    <BackgroundParticles />

    <div class="wrapper">
      <div class="container">
        <!-- ===== Header ===== -->
        <header class="opdash-header card">
          <div class="title">
            <h1>Operator Control</h1>
            <p class="muted">Overview of your organization and quick access to tools</p>
          </div>

          <div class="status-pill" v-if="hasOrganization">
            <span class="dot"></span>
            <span>Operator signed in • org:</span>
            <code class="mono">{{ orgId }}</code>
          </div>
        </header>

        <!-- ===== Top row (Gauge + Organization) ===== -->
        <section class="top-grid">
          <!-- Gauge / KPIs -->
          <div class="card kpi-card">
            <div class="gauge-wrap">
              <div class="gauge">
                <div class="gauge-arc" :style="{ '--pct': overallPct }"></div>
                <div class="gauge-center">
                  <div class="big">{{ totalRouters }}</div>
                  <div class="label">Routers</div>
                </div>
              </div>
            </div>

            <div class="legend">
              <div class="lg"><span class="b b-green"></span> Online <b>{{ onlineCount }}</b></div>
              <div class="lg"><span class="b b-red"></span> Offline <b>{{ offlineCount }}</b></div>
              <div class="lg"><span class="b b-blue"></span> Associated <b>{{ associatedCount }}</b></div>
              <div class="lg"><span class="b b-orange"></span> Unassociated <b>{{ unassociatedCount }}</b></div>
            </div>
          </div>

          <!-- Organization -->
          <div class="card org-card">
            <div class="row-head">
              <h3 class="section-title"><i class="fas fa-building"></i> Organization</h3>
              <div class="actions">
                <RouterLink class="btn ghost sm" to="/organization">Open</RouterLink>
                <RouterLink class="btn primary sm" to="/organization/edit">Edit</RouterLink>
              </div>
            </div>

            <div v-if="orgLoading" class="empty">Loading organization…</div>

            <div v-else class="org-grid" :class="{ dim: !org }">
              <div class="row"><label>Name</label><span>{{ org?.name ?? '—' }}</span></div>
              <div class="row"><label>VAT</label><span>{{ org?.vat_number ?? '—' }}</span></div>
              <div class="row wide"><label>Address</label><span>{{ org?.address ?? '—' }}</span></div>
              <div class="row"><label>Email</label><span>{{ org?.contact_email ?? '—' }}</span></div>
              <div class="row"><label>Phone</label><span>{{ org?.contact_phone ?? '—' }}</span></div>
            </div>

            <p v-if="orgMsg" class="hint">{{ orgMsg }}</p>
          </div>
        </section>

        <!-- ===== Widgets grid ===== -->
        <section class="grid">
          <!-- Router status (preview) -->
          <article class="card tile">
            <div class="row-head">
              <h3 class="section-title"><i class="fas fa-network-wired"></i> Router Status</h3>
              <div class="actions">
                <button class="btn ghost sm" @click="loadRouters" :disabled="loading">
                  <i class="fas fa-rotate"></i>
                  <span>{{ loading ? 'Loading…' : 'Refresh' }}</span>
                </button>
                <RouterLink class="btn primary sm" to="/routertable">Open</RouterLink>
              </div>
            </div>

            <div class="table-wrapper">
              <table v-if="routers.length" class="routers-table" aria-label="Router status">
                <thead>
                  <tr>
                    <th>MAC</th>
                    <th>VPN IP</th>
                    <th>Status</th>
                    <th>Last</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="r in routersPreview" :key="r.mac">
                    <td class="mono ellipsis">{{ r.mac }}</td>
                    <td class="mono ellipsis">{{ r.ip || '—' }}</td>
                    <td>
                      <span :class="['state', stateClass(r)]" :title="statusTitle(r)">
                        {{ statusLabel(r) }}
                      </span>
                    </td>
                    <td class="ellipsis" :title="r.lastTime ? new Date(r.lastTime).toLocaleString() : ''">
                      {{ r.lastTime ? formatRelativeTime(r.lastTime) : '—' }}
                    </td>
                  </tr>
                </tbody>
              </table>

              <div v-else class="empty small">No routers found.</div>
            </div>

            <p v-if="error" class="toast error">{{ error }}</p>
          </article>

          <!-- Agents widget -->
          <article class="card tile">
            <div class="row-head">
              <h3 class="section-title"><i class="fas fa-robot"></i> Agents</h3>
              <RouterLink class="btn primary sm" to="/agents">Open</RouterLink>
            </div>

            <ul class="stats-list">
              <li><span>Total agents</span><b>{{ totalRouters }}</b></li>
              <li><span>Active (online)</span><b class="ok">{{ onlineCount }}</b></li>
              <li><span>Pending/Deactivated</span><b class="warn">{{ deactivatedCount }}</b></li>
            </ul>

            <div class="actions spaced">
              <RouterLink class="btn ghost sm" to="/agents/register">
                <i class="fas fa-wand-magic-sparkles"></i>
                <span>Pre-register (WG & script)</span>
              </RouterLink>
            </div>
          </article>

          <!-- Add User -->
          <article class="card tile">
            <div class="row-head">
              <h3 class="section-title"><i class="fas fa-user-plus"></i> Add User</h3>
            </div>
            <p class="muted small">
              Create an end-user for your organization. Role is fixed to <code>User</code>.
            </p>
            <div class="actions spaced">
              <RouterLink class="btn primary sm" to="/adduser">Create user</RouterLink>
            </div>
            <div class="readonly-pill" v-if="!hasOrganization" title="Operator must have an organization">
              <span class="dot orange"></span>
              <span>Link an organization to enable user creation</span>
            </div>
          </article>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup>
/**
 * DashboardOperator.vue
 * Endpoints:
 *  - GET  /api/mikrotik/list         -> [{ mac, ip(vpn_internal_ip), status }]
 *  - GET  /api/data/routers          -> [{ mac, value('online'|'offline'|1|0), time }] (JWT)
 *  - GET  /api/complete-organization -> { name, vat_number, address, ... } (JWT)
 */
import { ref, computed, onMounted } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import BackgroundParticles from '@/components/common/BackgroundParticles.vue'
import { API } from '@/appCore.js'

const router = useRouter()
function go (path) { router.push(path) }

/* ---- Organization ---- */
const token  = localStorage.getItem('token') || ''
const orgId  = (localStorage.getItem('organization_id') || '').toString()
const hasOrganization = computed(() => !!orgId)

const org = ref(null)
const orgLoading = ref(false)
const orgMsg = ref('')

function readLocalOrg () {
  try { return JSON.parse(localStorage.getItem('organization_profile') || 'null') }
  catch { return null }
}
async function loadOrganization () {
  orgLoading.value = true
  orgMsg.value = ''
  try {
    if (!token) {
      org.value = readLocalOrg()
      orgMsg.value = 'Not authenticated. Showing local data if available.'
      return
    }
    const res = await fetch(`${API}/api/complete-organization`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    const raw = await res.text()
    let data = {}
    try { data = raw ? JSON.parse(raw) : {} } catch { data = {} }

    if (res.ok && (data?.name || data?.vat_number || data?.address)) {
      org.value = { ...(readLocalOrg() || {}), ...data }
      localStorage.setItem('organization_profile', JSON.stringify(org.value))
    } else {
      org.value = readLocalOrg()
      orgMsg.value = 'No server data yet. Complete the organization form.'
    }
  } catch {
    org.value = readLocalOrg()
    orgMsg.value = 'Server unreachable. Showing local data if any.'
  } finally {
    orgLoading.value = false
  }
}

/* ---- Routers / Agents ---- */
const routers = ref([]) // [{ mac, ip, assoc, online, lastTime }]
const loading = ref(false)
const error = ref('')

function normalizeMAC (mac) { return (mac || '').trim().toUpperCase() }
function formatRelativeTime (dateString) {
  try {
    const date = new Date(dateString)
    const now = new Date()
    const diffSec = Math.floor((now - date) / 1000)
    if (diffSec < 60) return 'Just now'
    const diffMin = Math.floor(diffSec / 60)
    if (diffMin < 60) return `${diffMin}m ago`
    const diffH = Math.floor(diffMin / 60)
    if (diffH < 24) return `${diffH}h ago`
    const diffD = Math.floor(diffH / 24)
    if (diffD < 30) return `${diffD}d ago`
    return date.toLocaleDateString()
  } catch { return '—' }
}
function stateClass (r) {
  if (r.online === true) return 'green'
  if (r.online === false) return 'red'
  if (r.assoc === 'Deactivated') return 'orange'
  if (r.assoc === 'Associated') return 'blue'
  return 'gray'
}
function statusLabel (r) {
  if (r.online === true) return 'Online'
  if (r.online === false) return 'Offline'
  return r.assoc || 'Unknown'
}
function statusTitle (r) {
  if (r.online === true) return 'Online (last ping)'
  if (r.online === false) return 'Offline (last ping)'
  return `Association: ${r.assoc || 'Unknown'}`
}

async function fetchList () {
  const res = await fetch(`${API}/api/mikrotik/list`)
  if (!res.ok) throw new Error(`GET /mikrotik/list → ${res.status}`)
  const data = await res.json()
  if (!Array.isArray(data)) throw new Error('Invalid /mikrotik/list payload')
  return data.map(row => {
    const mac = normalizeMAC(row.mac || row.MAC)
    const assoc = (row.status || '').trim() || 'Unknown'
    const ip = row.ip || row.vpn_internal_ip || ''
    return { mac, ip, assoc, online: null, lastTime: null }
  })
}

async function fetchInfluxMap () {
  if (!token) return {}
  const res = await fetch(`${API}/api/data/routers`, {
    headers: { Authorization: `Bearer ${token}` }
  })
  if (!res.ok) return {}
  const rows = await res.json()
  const map = {}
  for (const r of rows || []) {
    const mac  = normalizeMAC(r.mac || r.MAC)
    const val  = String(r.value ?? r.Value ?? '').toLowerCase()
    const time = r.time || r.Time || null
    if (!mac) continue
    let online = null
    if (val === 'online' || val === '1') online = true
    else if (val === 'offline' || val === '0') online = false
    map[mac] = { online, time }
  }
  return map
}

async function loadRouters () {
  loading.value = true
  error.value = ''
  try {
    const [list, influx] = await Promise.allSettled([fetchList(), fetchInfluxMap()])
    let base = []
    if (list.status === 'fulfilled') base = list.value
    else throw list.reason
    const enrich = influx.status === 'fulfilled' ? influx.value : {}
    for (const r of base) {
      const enr = enrich[r.mac]
      if (enr) {
        r.online = enr.online
        r.lastTime = enr.time
      }
    }
    const orderAssoc = { Associated: 0, Unassociated: 1, Deactivated: 2, Unknown: 3 }
    routers.value = base.sort((a, b) => {
      const aScore = a.online === true ? 0 : a.online === false ? 2 : 1
      const bScore = b.online === true ? 0 : b.online === false ? 2 : 1
      if (aScore !== bScore) return aScore - bScore
      const aa = orderAssoc[a.assoc] ?? 9
      const bb = orderAssoc[b.assoc] ?? 9
      if (aa !== bb) return aa - bb
      const ta = a.lastTime ? new Date(a.lastTime).getTime() : 0
      const tb = b.lastTime ? new Date(b.lastTime).getTime() : 0
      return tb - ta
    })
  } catch (e) {
    // eslint-disable-next-line no-console
    console.error('[dashboard-operator] routers load failed:', e)
    error.value = 'Failed to load router data.'
  } finally {
    loading.value = false
  }
}

/* ---- Derived ---- */
const totalRouters      = computed(() => routers.value.length)
const onlineCount       = computed(() => routers.value.filter(r => r.online === true).length)
const offlineCount      = computed(() => routers.value.filter(r => r.online === false).length)
const associatedCount   = computed(() => routers.value.filter(r => r.assoc === 'Associated').length)
const unassociatedCount = computed(() => routers.value.filter(r => r.assoc === 'Unassociated').length)
const deactivatedCount  = computed(() => routers.value.filter(r => r.assoc === 'Deactivated').length)
const overallPct        = computed(() => {
  const t = totalRouters.value || 1
  return Math.round((onlineCount.value / t) * 100)
})
const routersPreview    = computed(() => routers.value.slice(0, 6))

/* ---- Init ---- */
onMounted(async () => {
  document.title = 'NetSecure-IQ – Operator Dashboard'
  await Promise.all([loadOrganization(), loadRouters()])
})
</script>

<style scoped>
@import url('https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css');

:root {
  --bg-dark:#0e111a; --panel-grey:#1a1d26; --divider-grey:#2a2d36;
  --text-primary:#f5f7fa; --text-secondary:#9ca3af;
  --primary-accent:#00c2c2; --primary-hover:#00a7a7;
  --danger:#ef4444; --success:#22c55e; --info:#38bdf8; --muted:#6b7280;
}

/* Page layout */
.opdash-page{ position:relative; min-height:100vh; overflow:hidden; }
.wrapper{ position:relative; z-index:10; display:flex; justify-content:center; padding:28px; }
.container{ width:100%; max-width:1220px; display:flex; flex-direction:column; gap:22px; }

/* Card */
.card{
  background: var(--panel-grey);
  border-radius: 16px;
  padding: 20px;
  border: 1px solid rgba(255,255,255,.05);
  box-shadow: 0 0 40px rgba(0,194,194,.05);
}

/* Header */
.opdash-header{ display:flex; align-items:flex-start; justify-content:space-between; gap:16px; }
.title h1{ margin:0; font-size:20px; color:var(--text-primary); }
.muted{ color:var(--text-secondary); }
.status-pill{
  display:inline-flex; align-items:center; gap:8px; margin-top:4px;
  padding:8px 12px; border-radius:999px;
  background:rgba(255,255,255,.06); border:1px solid rgba(255,255,255,.08);
  color:var(--text-primary);
}
.status-pill .dot{ width:8px; height:8px; border-radius:50%; background:var(--primary-accent); }
.mono{ font-family: ui-monospace, SFMono-Regular, Menlo, monospace; letter-spacing:.2px; }

/* Top grid (more air) */
.top-grid{ display:grid; grid-template-columns: 440px 1fr; gap:20px; }

/* Gauge */
.kpi-card{ display:grid; gap:14px; }
.gauge-wrap{ display:flex; align-items:center; justify-content:center; padding:6px 4px; }
.gauge{
  position:relative; width:100%; max-width:380px; aspect-ratio: 1.85 / 1;
  border-radius: 14px; border:1px solid rgba(255,255,255,.06);
  background: rgba(255,255,255,.04);
}
.gauge-arc{
  --pct: 0;
  position:absolute; inset:0;
  background:
    conic-gradient(from 270deg,
      var(--success) calc(var(--pct) * 1%),
      rgba(255,255,255,.06) 0);
  mask: radial-gradient(closest-side, transparent 62%, #000 63%), linear-gradient(#000,#000);
  -webkit-mask: radial-gradient(closest-side, transparent 62%, #000 63%), linear-gradient(#000,#000);
  border-radius:14px;
}
.gauge-center{ position:absolute; inset:0; display:flex; flex-direction:column; align-items:center; justify-content:center; }
.gauge-center .big{ font-size:32px; font-weight:800; color:var(--text-primary); line-height:1; }
.gauge-center .label{ font-size:12px; color:var(--text-secondary); }

.legend{ display:grid; grid-template-columns: repeat(2,minmax(0,1fr)); gap:10px; }
.lg{ display:flex; align-items:center; gap:8px; font-size:13px; color:var(--text-secondary); }
.b{ width:10px; height:10px; border-radius:3px; display:inline-block; }
.b-green{ background:#22c55e; }
.b-red{ background:#ef4444; }
.b-blue{ background:#93c5fd; }
.b-orange{ background:#fbbf24; }

/* Org card */
.org-card .row-head{ display:flex; align-items:center; justify-content:space-between; gap:10px; margin-bottom:10px; }
.section-title{ color:var(--primary-accent); margin:0; font-size:16px; display:flex; align-items:center; gap:8px; }
.actions{ display:flex; gap:10px; flex-wrap:wrap; }
.actions.spaced{ margin-top:10px; } /* add breathing room below content */
.org-grid{ display:grid; grid-template-columns: repeat(auto-fit, minmax(220px,1fr)); gap:12px 16px; }
.org-grid.dim{ opacity:.85; }
.row{ display:flex; flex-direction:column; gap:6px; }
.row.wide{ grid-column: 1 / -1; }
label{ font-size:12px; color:var(--text-secondary); }
span{ font-size:14px; color:var(--text-primary); }
.hint{ margin-top:6px; color:var(--text-secondary); font-size:12px; }

/* Widgets grid (more gap) */
.grid{ display:grid; grid-template-columns: repeat(auto-fit, minmax(340px, 1fr)); gap:22px; }
.tile .row-head{ display:flex; align-items:center; justify-content:space-between; gap:12px; margin-bottom:10px; }

/* Router table — compact, constrained, and non-overflowing */
.table-wrapper{
  border-radius:12px; border:1px solid rgba(255,255,255,.06);
  overflow-y:auto; overflow-x:hidden; max-height:240px; /* keep inside the tile */
  padding:2px;
}
.routers-table{
  width:100%; border-collapse:separate; border-spacing:0; font-size:12.5px;
  table-layout: fixed; /* makes ellipsis work */
}
.routers-table thead th{
  text-align:left; padding:10px 12px; background:rgba(255,255,255,.04);
  border-bottom:1px solid rgba(255,255,255,.08); color:var(--text-secondary); white-space:nowrap;
}
.routers-table tbody td{ padding:10px 12px; border-bottom:1px solid rgba(255,255,255,.05); color:var(--text-primary); }
.routers-table th:nth-child(1), .routers-table td:nth-child(1){ width:42%; }
.routers-table th:nth-child(2), .routers-table td:nth-child(2){ width:24%; }
.routers-table th:nth-child(3), .routers-table td:nth-child(3){ width:18%; }
.routers-table th:nth-child(4), .routers-table td:nth-child(4){ width:16%; }

.ellipsis{ overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.empty.small{ text-align:center; color:var(--text-secondary); padding:12px; }

.state{ display:inline-block; padding:4px 8px; border-radius:999px; font-weight:700; font-size:11px; }
.state.green  { background:rgba(22,163,74,.15); color:#22c55e; }
.state.red    { background:rgba(239,68,68,.15); color:#f87171; }
.state.blue   { background:rgba(59,130,246,.15); color:#93c5fd; }
.state.orange { background:rgba(245,158,11,.15); color:#fbbf24; }
.state.gray   { background:rgba(156,163,175,.15); color:#d1d5db; }

/* Stats */
.stats-list{ list-style:none; padding:0; margin:0; display:grid; gap:8px; }
.stats-list li{ display:flex; align-items:center; justify-content:space-between; }
.stats-list span{ color:var(--text-secondary); }
.stats-list b{ color:var(--text-primary); }
.stats-list b.ok{ color:#22c55e; }
.stats-list b.warn{ color:#f59e0b; }

/* Buttons (smaller + airy) */
.btn{
  display:inline-flex; align-items:center; gap:8px;
  border-radius:10px; padding:9px 12px; font-weight:700;
  border:1px solid rgba(255,255,255,.08); background:rgba(255,255,255,.06);
  color:var(--text-primary); cursor:pointer; transition:.15s;
}
.btn:hover{ background:rgba(255,255,255,.10); }
.btn.primary{ background:var(--primary-accent); color:#0e111a; border-color:transparent; }
.btn.primary:hover{ background:var(--primary-hover); color:#fff; }
.btn.ghost{ background:rgba(255,255,255,.06); }
.btn.sm{ padding:6px 10px; font-size:12px; border-radius:8px; }

/* Pills / messages */
.readonly-pill{
  margin-top:10px; display:inline-flex; align-items:center; gap:8px;
  padding:8px 10px; border-radius:999px;
  background:rgba(255,255,255,.04); border:1px solid rgba(255,255,255,.06);
  font-size:12px; color:var(--text-primary);
}
.readonly-pill .dot{ width:8px; height:8px; border-radius:50%; background:#aaa; }
.readonly-pill .dot.orange{ background:#f59e0b; }

.toast{ margin-top:12px; text-align:center; padding:10px 12px; border-radius:8px; font-size:14px; }
.toast.error{ background:rgba(239,68,68,.12); color:var(--danger); }

@media (max-width: 960px){
  .top-grid{ grid-template-columns: 1fr; }
}
</style>
