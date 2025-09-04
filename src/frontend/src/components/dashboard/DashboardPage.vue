<!-- src/frontend/src/views/dashboard/DashboardUser.vue -->
<template>
  <div class="userdash-page">
    <BackgroundParticles />

    <div class="wrapper">
      <div class="container">
        <!-- ===== Header ===== -->
        <header class="ud-head card">
          <div class="title">
            <h1>User Dashboard</h1>
            <p class="muted">Quick access to router status</p>
          </div>

          <div v-if="hasOrganization" class="status-pill">
            <span class="dot"></span>
            <span>org:</span>
            <code class="mono">{{ orgId }}</code>
          </div>
        </header>

        <!-- ===== Single widget: Router Status ===== -->
        <section class="grid">
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

            <!-- Small legend with counts -->
            <ul class="counts">
              <li><span class="dot g"></span>Online <b>{{ onlineCount }}</b></li>
              <li><span class="dot r"></span>Offline <b>{{ offlineCount }}</b></li>
              <li><span class="dot b"></span>Associated <b>{{ associatedCount }}</b></li>
              <li><span class="dot o"></span>Unassociated <b>{{ unassociatedCount }}</b></li>
            </ul>

            <!-- Compact preview table -->
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
        </section>
      </div>
    </div>
  </div>
</template>

<script setup>
/**
 * DashboardUser.vue
 * Access limited to RouterTable preview & link.
 * Endpoints used:
 *  - GET  ${API}/api/mikrotik/list
 *  - GET  ${API}/api/data/routers (JWT optional → enrich online/offline)
 */
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import BackgroundParticles from '@/components/common/BackgroundParticles.vue'
import { API } from '@/appCore.js'

/* Org context (display-only) */
const orgId = ref(localStorage.getItem('organization_id') || '')
const hasOrganization = computed(() => !!String(orgId.value).trim())

/* Routers state */
const routers = ref([]) // [{ mac, ip, assoc, online, lastTime }]
const loading = ref(false)
const error = ref('')
const token = localStorage.getItem('token') || ''

/* Helpers */
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

/* Fetchers */
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
      if (enr) { r.online = enr.online; r.lastTime = enr.time }
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
    console.error('[dashboard-user] routers load failed:', e)
    error.value = 'Failed to load router data.'
  } finally {
    loading.value = false
  }
}

/* Derived */
const onlineCount       = computed(() => routers.value.filter(r => r.online === true).length)
const offlineCount      = computed(() => routers.value.filter(r => r.online === false).length)
const associatedCount   = computed(() => routers.value.filter(r => r.assoc === 'Associated').length)
const unassociatedCount = computed(() => routers.value.filter(r => r.assoc === 'Unassociated').length)
const routersPreview    = computed(() => routers.value.slice(0, 6))

onMounted(loadRouters)
</script>

<style scoped>
@import url('https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css');

:root{
  --bg-dark:#0e111a; --panel-grey:#1a1d26; --divider-grey:#2a2d36;
  --text-primary:#f5f7fa; --text-secondary:#9ca3af;
  --primary-accent:#00c2c2; --primary-hover:#00a7a7;
  --danger:#ef4444; --success:#22c55e; --muted:#6b7280;
}

.userdash-page{ position:relative; min-height:100vh; overflow:hidden; }
.wrapper{ position:relative; z-index:10; display:flex; justify-content:center; padding:28px; }
.container{ width:100%; max-width:1200px; display:flex; flex-direction:column; gap:16px; }

.card{
  background: var(--panel-grey);
  border-radius: 16px;
  padding: 20px;
  border: 1px solid rgba(255,255,255,.05);
  box-shadow: 0 0 40px rgba(0,194,194,.05);
}

/* Header */
.ud-head{ display:flex; align-items:flex-start; justify-content:space-between; gap:16px; }
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

/* Grid */
.grid{ display:grid; grid-template-columns: 1fr; gap:20px; }

/* Tile */
.tile .row-head{ display:flex; align-items:center; justify-content:space-between; gap:12px; margin-bottom:10px; }
.section-title{ color:var(--primary-accent); margin:0; font-size:16px; display:flex; align-items:center; gap:8px; }
.actions{ display:flex; gap:10px; flex-wrap:wrap; }

/* Counts legend */
.counts{ display:flex; gap:14px; flex-wrap:wrap; list-style:none; padding:0; margin:0 0 8px; color:var(--text-secondary); font-size:13px; }
.counts .dot{ width:8px; height:8px; border-radius:50%; display:inline-block; margin-right:6px; }
.counts .g{ background:#22c55e; } .counts .r{ background:#ef4444; }
.counts .b{ background:#93c5fd; } .counts .o{ background:#fbbf24; }

/* Table */
.table-wrapper{
  border-radius:12px; border:1px solid rgba(255,255,255,.06);
  overflow-y:auto; overflow-x:hidden; max-height:260px; padding:2px;
}
.routers-table{
  width:100%; border-collapse:separate; border-spacing:0; font-size:12.5px; table-layout: fixed;
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

/* Status badge */
.state{ display:inline-block; padding:4px 8px; border-radius:999px; font-weight:700; font-size:11px; }
.state.green  { background:rgba(22,163,74,.15); color:#22c55e; }
.state.red    { background:rgba(239,68,68,.15); color:#f87171; }
.state.blue   { background:rgba(59,130,246,.15); color:#93c5fd; }
.state.orange { background:rgba(245,158,11,.15); color:#fbbf24; }
.state.gray   { background:rgba(156,163,175,.15); color:#d1d5db; }

/* Buttons */
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

/* Toast */
.toast{ margin-top:12px; text-align:center; padding:10px 12px; border-radius:8px; font-size:14px; }
.toast.error{ background:rgba(239,68,68,.12); color:#f87171; }
</style>
