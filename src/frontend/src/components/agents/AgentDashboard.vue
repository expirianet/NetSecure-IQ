<!-- /src/frontend/src/components/agents/AgentDashboard.vue -->
<template>
  <div class="agents-page">
    <!-- Background particles -->
    <BackgroundParticles />

    <div class="agents-wrapper">
      <div class="agents-container">
        <div class="agents-card">
          <!-- Header -->
          <div class="header-row">
            <div class="titles">
              <h2>MikroTik Agent Management</h2>
              <p class="subtitle">Monitor, associate, enable/disable and manage your routers.</p>
            </div>

            <div class="header-actions">
              <RouterLink to="/agents/register" class="btn primary">
                <i class="fas fa-microchip"></i>
                <span>Pre-register an Agent</span>
              </RouterLink>
              <button class="btn ghost" @click="loadAgents" :disabled="loading">
                <i class="fas" :class="loading ? 'fa-spinner fa-spin' : 'fa-rotate-right'"></i>
                <span>Refresh</span>
              </button>
            </div>
          </div>

          <!-- Toolbar -->
          <div class="toolbar">
            <div class="search">
              <i class="fas fa-search"></i>
              <input
                v-model.trim="query"
                type="text"
                placeholder="Search by MAC, IP, status or site_id"
                aria-label="Search"
              />
            </div>

            <div class="actions">
              <button class="btn ghost" :disabled="!selectedAgent || acting" @click="onAssociate">
                <i class="fas fa-link"></i><span>Associate</span>
              </button>
              <button class="btn ghost" :disabled="!selectedAgent || acting" @click="onEnable">
                <i class="fas fa-toggle-on"></i><span>Enable</span>
              </button>
              <button class="btn ghost" :disabled="!selectedAgent || acting" @click="onDeactivate">
                <i class="fas fa-ban"></i><span>Deactivate</span>
              </button>
              <button class="btn ghost" :disabled="!selectedAgent || acting" @click="onTest">
                <i class="fas fa-satellite-dish"></i><span>Test</span>
              </button>
              <button class="btn ghost danger" :disabled="!selectedAgent || acting" @click="onDelete">
                <i class="fas fa-trash"></i><span>Delete</span>
              </button>
            </div>

            <div class="legend">
              <span class="chip green"><span class="dot"></span>Associated</span>
              <span class="chip orange"><span class="dot"></span>Unassociated</span>
              <span class="chip red"><span class="dot"></span>Deactivated</span>
            </div>
          </div>

          <!-- Table -->
          <div class="table-wrapper">
            <table class="agents-table" role="table" aria-label="Agent list">
              <thead>
                <tr>
                  <th style="width:72px">Select</th>
                  <th>MAC Address</th>
                  <th>IP</th>
                  <th style="width:170px">Status</th>
                  <th>Site ID</th>
                </tr>
              </thead>

              <tbody>
                <tr
                  v-for="(a, idx) in filteredAgents"
                  :key="a.mac"
                  :class="['row', rowTint(normalizedStatus(a.status)), { selected: idx === selectedIndex } ]"
                  @click="selectAgent(a, idx)"
                >
                  <td>
                    <span class="radio" :class="{ checked: idx === selectedIndex }" aria-hidden="true"></span>
                    <span class="sr-only">Select {{ a.mac }}</span>
                  </td>
                  <td class="mono">{{ a.mac }}</td>
                  <td class="mono">{{ a.ip || '—' }}</td>
                  <td>
                    <span :class="['state', stateClass(normalizedStatus(a.status))]">{{ a.status }}</span>
                  </td>
                  <td class="mono">{{ a.site_id || '—' }}</td>
                </tr>

                <tr v-if="!filteredAgents.length && !loading">
                  <td colspan="5" class="empty">No agents match your search.</td>
                </tr>
                <tr v-if="loading">
                  <td colspan="5" class="empty"><i class="fas fa-spinner fa-spin"></i> Loading…</td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Toast -->
          <p v-if="toast" class="toast" :class="toastType">{{ toast }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import BackgroundParticles from '@/components/common/BackgroundParticles.vue'
import { ref, computed, onMounted } from 'vue'

// ===== API helper (same-origin first, fallback to http://localhost:8000) =====
async function apiFetch(path, opts = {}) {
  const rel = path.startsWith('/') ? path : `/${path}`
  const headers = { 'Content-Type': 'application/json', ...(opts.headers || {}) }
  const token = localStorage.getItem('token')
  if (token) headers.Authorization = headers.Authorization || `Bearer ${token}`

  // Try same-origin
  let res
  try {
    res = await fetch(rel, { ...opts, headers })
    if (res.ok || res.status === 204) return res
  } catch (err) {
    /* eslint-disable-next-line no-console */
    console.debug('Same-origin fetch failed; will fallback to localhost:8000', err)
  }
  // Fallback to localhost:8000
  const url = `http://localhost:8000${rel}`
  res = await fetch(url, { ...opts, headers })
  return res
}

// ===== State =====
const agents = ref([])
const loading = ref(false)
const acting = ref(false)
const query = ref('')
const selectedAgent = ref(null)
const selectedIndex = ref(-1)

// ===== Computed =====
const filteredAgents = computed(() => {
  const q = query.value.toLowerCase().trim()
  if (!q) return agents.value
  return agents.value.filter(a =>
    a.mac?.toLowerCase().includes(q) ||
    a.ip?.toLowerCase().includes(q) ||
    a.site_id?.toLowerCase().includes(q) ||
    a.status?.toLowerCase().includes(q)
  )
})

// ===== Lifecycle =====
onMounted(loadAgents)

async function loadAgents() {
  loading.value = true
  try {
    const res = await apiFetch('/api/mikrotik/list', { method: 'GET' })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    agents.value = await res.json() // [{ mac, ip, status, site_id }]
  } catch (e) {
    console.error('Failed to load agents:', e)
    agents.value = []
    showToast('Failed to load agents', 'error')
  } finally {
    loading.value = false
  }
}

function selectAgent(a, idx) {
  selectedAgent.value = a
  selectedIndex.value = idx
}

// ===== Actions =====
const toast = ref('')
const toastType = ref('success')
const showToast = (msg, type = 'success') => {
  toast.value = msg
  toastType.value = type
  setTimeout(() => (toast.value = ''), 2400)
}

async function onAssociate() {
  if (!selectedAgent.value) return
  const siteId = window.prompt('Enter Site ID (UUID) to associate this router to:')
  if (!siteId) return
  acting.value = true
  try {
    const res = await apiFetch('/api/mikrotik/associate', {
      method: 'POST',
      body: JSON.stringify({ mac: selectedAgent.value.mac, site_id: siteId })
    })
    if (!res.ok) throw new Error(await safeErr(res))
    selectedAgent.value.site_id = siteId
    selectedAgent.value.status = deriveStatus(selectedAgent.value)
    showToast('Agent associated.')
  } catch (e) {
    showToast(e.message || 'Association failed', 'error')
  } finally {
    acting.value = false
  }
}

async function onEnable() {
  if (!selectedAgent.value) return
  acting.value = true
  try {
    const res = await apiFetch('/api/mikrotik/enable', {
      method: 'POST',
      body: JSON.stringify({ mac: selectedAgent.value.mac })
    })
    if (res.status !== 204) throw new Error(await safeErr(res))
    selectedAgent.value.status = deriveStatus(selectedAgent.value)
    showToast('Agent enabled.')
  } catch (e) {
    showToast(e.message || 'Enable failed', 'error')
  } finally {
    acting.value = false
  }
}

async function onDeactivate() {
  if (!selectedAgent.value) return
  acting.value = true
  try {
    const res = await apiFetch('/api/mikrotik/disable', {
      method: 'POST',
      body: JSON.stringify({ mac: selectedAgent.value.mac })
    })
    if (res.status !== 204) throw new Error(await safeErr(res))
    selectedAgent.value.status = 'Deactivated'
    showToast('Agent deactivated.', 'error')
  } catch (e) {
    showToast(e.message || 'Deactivation failed', 'error')
  } finally {
    acting.value = false
  }
}

async function onDelete() {
  if (!selectedAgent.value) return
  if (!confirm(`Delete router ${selectedAgent.value.mac}? This cannot be undone.`)) return
  acting.value = true
  try {
    const res = await apiFetch('/api/mikrotik', {
      method: 'DELETE',
      body: JSON.stringify({ mac: selectedAgent.value.mac })
    })
    if (res.status !== 204) throw new Error(await safeErr(res))
    agents.value = agents.value.filter((_, i) => i !== selectedIndex.value)
    selectedAgent.value = null
    selectedIndex.value = -1
    showToast('Agent deleted.', 'error')
  } catch (e) {
    showToast(e.message || 'Delete failed', 'error')
  } finally {
    acting.value = false
  }
}

async function onTest() {
  if (!selectedAgent.value) return
  acting.value = true
  try {
    const res = await apiFetch('/api/mikrotik/test', {
      method: 'POST',
      body: JSON.stringify({ mac: selectedAgent.value.mac })
    })
    if (!res.ok) throw new Error(await safeErr(res))
    const data = await res.json()
    showToast(data.ok ? '✅ Reachable over WG' : '❌ Unreachable', data.ok ? 'success' : 'error')
  } catch (e) {
    showToast(e.message || 'Test failed', 'error')
  } finally {
    acting.value = false
  }
}

// ===== Display helpers =====
async function safeErr(res) {
  try {
    const t = await res.text()
    return t || `HTTP ${res.status}`
  } catch {
    return `HTTP ${res.status}`
  }
}

function normalizedStatus(status) {
  const s = String(status || '').toLowerCase()
  if (s.startsWith('assoc')) return 'associated'
  if (s.startsWith('unassoc')) return 'unassociated'
  return 'deactivated'
}
function deriveStatus(row) {
  if (normalizedStatus(row.status) === 'deactivated') return 'Deactivated'
  return row.site_id ? 'Associated' : 'Unassociated'
}

function stateClass(status) {
  return status === 'associated'
    ? 'green'
    : status === 'unassociated'
    ? 'orange'
    : 'red'
}
function rowTint(status) {
  return status === 'associated'
    ? 'tint-green'
    : status === 'unassociated'
    ? 'tint-orange'
    : 'tint-red'
}
</script>

<style scoped>
@import url('https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css');

:root {
  --bg-dark: #0e111a;
  --panel-grey: #1a1d26;
  --divider-grey: #2a2d36;
  --text-primary: #f5f7fa;
  --text-secondary: #9ca3af;
  --primary-accent: #00c2c2;
  --primary-hover: #00a7a7;
  --danger: #ef4444;
  --success: #22c55e;
  --radius: 12px;
}

.agents-page { position: relative; min-height: 100vh; overflow: hidden; }
.agents-wrapper { position: relative; z-index: 10; display: flex; justify-content: center; padding: 32px; }
.agents-container { width: 100%; max-width: 1120px; }
.agents-card {
  background: var(--panel-grey);
  border-radius: 16px;
  padding: 20px;
  border: 1px solid rgba(255,255,255,.05);
  box-shadow: 0 0 40px rgba(0,194,194,.05);
}

.header-row { display: grid; grid-template-columns: 1fr auto; align-items: center; gap: 16px; }
.titles h2 { margin: 0; font-size: 20px; color: var(--text-primary); }
.subtitle { margin: 4px 0 0; color: var(--text-secondary); }
.header-actions { display: flex; gap: 8px; }

.btn {
  display: inline-flex; align-items: center; gap: 8px;
  border-radius: 10px; padding: 10px 14px; font-weight: 600;
  border: 1px solid rgba(255,255,255,.08); background: rgba(255,255,255,.06);
  color: var(--text-primary); cursor: pointer; transition: .15s;
}
.btn:hover { background: rgba(255,255,255,.10); }
.btn.primary {
  background: var(--primary-accent); color: #0e111a; border-color: transparent;
}
.btn.primary:hover { background: var(--primary-hover); color: #fff; }
.btn.ghost.danger { border-color: rgba(239,68,68,.35); color: #fca5a5; }
.btn:disabled { opacity: .5; cursor: not-allowed; }

.toolbar {
  margin-top: 16px; display: grid; gap: 12px;
  grid-template-columns: minmax(240px, 1fr) auto auto;
  align-items: center;
}
.search {
  display: flex; align-items: center; gap: 10px;
  background: rgba(255,255,255,.06); border: 1px solid rgba(255,255,255,.08);
  padding: 10px 12px; border-radius: 10px;
}
.search input {
  border: none; outline: none; background: transparent; color: var(--text-primary); width: 100%;
}
.actions { display: flex; gap: 8px; flex-wrap: wrap; }
.legend { display: flex; gap: 8px; justify-self: end; }

.chip {
  display: inline-flex; align-items: center; gap: 8px;
  padding: 6px 10px; border-radius: 999px; font-size: 12px;
  border: 1px solid rgba(255,255,255,.08); background: rgba(255,255,255,.06);
}
.chip .dot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; }
.chip.green  .dot { background: #16a34a; }
.chip.orange .dot { background: #f59e0b; }
.chip.red    .dot { background: #ef4444; }

.table-wrapper { margin-top: 14px; overflow: auto; border-radius: 12px; border: 1px solid rgba(255,255,255,.06); }
.agents-table { width: 100%; border-collapse: separate; border-spacing: 0; font-size: 14px; }
.agents-table thead th {
  text-align: left; padding: 12px 14px; background: rgba(255,255,255,.04);
  border-bottom: 1px solid rgba(255,255,255,.08); color: var(--text-secondary);
}
.agents-table tbody td { padding: 14px; border-bottom: 1px solid rgba(255,255,255,.05); }
.mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; letter-spacing: .2px; }

.row { transition: background .15s, outline .15s; cursor: pointer; }
.row:hover { outline: 1px solid rgba(255,255,255,.08); }
.selected { outline: 2px solid rgba(0,194,194,.35); }

.tint-green  { background: rgba(22,163,74,.06); }
.tint-orange { background: rgba(245,158,11,.06); }
.tint-red    { background: rgba(239,68,68,.06); }

.radio { display:inline-block; width:12px; height:12px; border-radius:50%;
  border:2px solid rgba(255,255,255,.35); vertical-align: middle; margin-right: 6px; }
.radio.checked { background: var(--primary-accent); border-color: var(--primary-accent); }

.state {
  display: inline-block; padding: 6px 10px; border-radius: 999px; font-weight: 600; font-size: 12px;
}
.state.green  { background: rgba(22,163,74,.15); color:#22c55e; }
.state.orange { background: rgba(245,158,11,.15); color:#fbbf24; }
.state.red    { background: rgba(239,68,68,.15); color:#f87171; }

.empty { text-align: center; color: var(--text-secondary); padding: 20px; }

.toast {
  margin-top: 12px; text-align:center; padding: 10px 12px; border-radius: 8px; font-size: 14px;
}
.toast.success { background: rgba(34,197,94,.12); color: var(--success); }
.toast.error   { background: rgba(239,68,68,.12); color: var(--danger); }

@media (max-width: 980px) {
  .toolbar { grid-template-columns: 1fr; }
  .legend { justify-self: start; }
}
</style>
