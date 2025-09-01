<!-- src/frontend/src/views/agents/AgentsList.vue -->
<template>
  <div class="agents-page">
    <div class="card">
      <header class="head">
        <div class="title">
          <h1>MikroTik Agent Management</h1>
          <p class="muted">Monitor, associate, and manage your deployed routers.</p>
        </div>

        <RouterLink class="btn primary" to="/agents/register">
          <i class="fas fa-plus-circle"></i> Pre-register an Agent
        </RouterLink>
      </header>

      <!-- Toolbar -->
      <div class="toolbar">
        <div class="search-wrap">
          <i class="fas fa-search"></i>
          <input
            v-model="search"
            class="search"
            placeholder="Search by MAC, organization, or site"
            @keydown.enter.prevent
          />
        </div>

        <div class="filters">
          <button
            class="pill"
            :class="{ active: statusFilter === 'all' }"
            @click="statusFilter = 'all'"
            title="All"
          >
            <span class="dot gray"></span> All
          </button>
          <button
            class="pill"
            :class="{ active: statusFilter === 'Associated' }"
            @click="statusFilter = 'Associated'"
            title="Associated"
          >
            <span class="dot green"></span> Associated
          </button>
          <button
            class="pill"
            :class="{ active: statusFilter === 'Unassociated' }"
            @click="statusFilter = 'Unassociated'"
            title="Unassociated"
          >
            <span class="dot amber"></span> Unassociated
          </button>
          <button
            class="pill"
            :class="{ active: statusFilter === 'Deactivated' }"
            @click="statusFilter = 'Deactivated'"
            title="Deactivated"
          >
            <span class="dot red"></span> Deactivated
          </button>
        </div>

        <div class="actions">
          <button class="btn ghost" :disabled="!hasSelection" @click="openAssociate">
            <i class="fas fa-link"></i> Associate
          </button>
          <button class="btn ghost" :disabled="!hasSelection" @click="activateSelected">
            <i class="fas fa-toggle-on"></i> Activate
          </button>
          <button class="btn ghost" :disabled="!hasSelection" @click="deactivateSelected">
            <i class="fas fa-toggle-off"></i> Deactivate
          </button>
          <button class="btn danger" :disabled="!hasSelection" @click="deleteSelected">
            <i class="fas fa-trash"></i> Delete
          </button>
        </div>
      </div>

      <!-- Table -->
      <div class="table-wrapper">
        <table class="table">
          <thead>
            <tr>
              <th class="sel">
                <input type="checkbox" :checked="allVisibleChecked" @change="toggleAllVisible" />
              </th>
              <th>MAC Address</th>
              <th>Status</th>
              <th>Organization</th>
              <th>Site</th>
              <th>IP (WG)</th>
              <th class="right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in filtered" :key="row.mac">
              <td class="sel">
                <input type="checkbox" :checked="selected.has(row.mac)" @change="toggleSel(row.mac)" />
              </td>
              <td class="mono">{{ row.mac }}</td>
              <td>
                <span :class="['badge', badgeClass(row.status)]">{{ row.status }}</span>
              </td>
              <td>—<!-- non renvoyé par l'API /list --> </td>
              <td>{{ row.site_id || '—' }}</td>
              <td class="mono">{{ row.ip || '—' }}</td>
              <td class="right">
                <button class="btn xs ghost" @click="testPing(row.mac)">
                  <i class="fas fa-satellite-dish"></i> Test
                </button>
                <button
                  v-if="row.status === 'Deactivated'"
                  class="btn xs ghost"
                  @click="activateOne(row.mac)"
                >
                  <i class="fas fa-toggle-on"></i>
                </button>
                <button
                  v-else
                  class="btn xs ghost"
                  @click="deactivateOne(row.mac)"
                >
                  <i class="fas fa-toggle-off"></i>
                </button>
                <button class="btn xs danger" @click="deleteOne(row.mac)">
                  <i class="fas fa-trash"></i>
                </button>
              </td>
            </tr>
            <tr v-if="!loading && filtered.length === 0">
              <td colspan="7" class="empty">No agents match your filters.</td>
            </tr>
            <tr v-if="loading">
              <td colspan="7" class="empty">Loading…</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Messages -->
      <p v-if="message" class="message" :class="ok ? 'success' : 'error'">
        {{ message }}
      </p>
    </div>

    <!-- Associate dialog (ultra simple) -->
    <div v-if="showAssociate" class="modal-backdrop" @click.self="closeAssociate">
      <div class="modal">
        <h3><i class="fas fa-link"></i> Associate selected agents</h3>
        <label for="siteId">Site ID</label>
        <input id="siteId" v-model.trim="associateSiteId" placeholder="UUID of the site" />
        <div class="modal-actions">
          <button class="btn ghost" @click="closeAssociate">Cancel</button>
          <button class="btn primary" :disabled="!associateSiteId" @click="confirmAssociate">
            Associate
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { listAgents, enableAgent, disableAgent, deleteAgent, associateAgent, testAgent } from '@/appCore.js'

/* ---------- State ---------- */
const agents = ref([])           // Données des agents du backend
const loading = ref(false)       // État de chargement
const message = ref('')          // Message de statut
const ok = ref(false)            // État du message (succès/erreur)
const error = ref(null)          // Erreur de chargement

const selected = ref(new Set())  // Éléments sélectionnés
const search = ref('')           // Terme de recherche
const statusFilter = ref('all')  // Filtre par statut: 'all' | 'associated' | 'unassociated' | 'deactivated'

/* ---------- Fetch list ---------- */
const load = async () => {
  loading.value = true
  error.value = null
  try {
    agents.value = await listAgents()
  } catch (err) {
    console.error('Failed to load agents:', err)
    error.value = 'Failed to load agents. Please try again later.'
    message.value = error.value
    ok.value = false
  } finally {
    loading.value = false
  }
}
onMounted(load)

/* ---------- Filters / selection ---------- */
const filtered = computed(() => {
  const q = search.value.trim().toLowerCase()
  return agents.value.filter(agent => {
    // Filtre par recherche
    const matchesSearch = !q || 
      agent.mac.toLowerCase().includes(q) ||
      (agent.organization && agent.organization.toLowerCase().includes(q)) ||
      (agent.site_id && agent.site_id.toLowerCase().includes(q))
    
    // Filtre par statut
    const matchesStatus = statusFilter.value === 'all' || 
      agent.status.toLowerCase() === statusFilter.value.toLowerCase()
    
    return matchesSearch && matchesStatus
    const hitSearch =
      !q ||
      r.mac?.toLowerCase().includes(q) ||
      r.site_id?.toLowerCase().includes(q)
    return hitStatus && hitSearch
  })
})
const hasSelection = computed(() => selected.value.size > 0)
const allVisibleChecked = computed(() =>
  filtered.value.length > 0 &&
  filtered.value.every(r => selected.value.has(r.mac))
)
function toggleSel(mac) {
  const s = new Set(selected.value)
  if (s.has(mac)) s.delete(mac)
  else s.add(mac)
  selected.value = s
}
function toggleAllVisible(e) {
  const s = new Set(selected.value)
  if (e.target.checked) {
    filtered.value.forEach(r => s.add(r.mac))
  } else {
    filtered.value.forEach(r => s.delete(r.mac))
  }
  selected.value = s
}

/* ---------- Helpers ---------- */
const badgeClass = (status) => {
  if (!status) return 'secondary'
  
  switch (status.toLowerCase()) {
    case 'associated':
      return 'success'
    case 'unassociated':
      return 'warning'
    case 'deactivated':
      return 'danger'
    default:
      return 'secondary'
  }
}
async function callJson(url, method, body) {
  const token = localStorage.getItem('token') || ''
  const res = await fetch(url, {
    method,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: 'Bearer ' + token } : {})
    },
    body: body ? JSON.stringify(body) : undefined
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) throw new Error(data?.error || data?.message || `${method} ${url} failed`)
  return data
}

/* ---------- Row actions ---------- */
const testPing = async (mac) => {
  try {
    const data = await testAgent(mac)
    message.value = data.message || 'Ping successful!'
    ok.value = true
    // Recharger la liste après un test réussi
    await load()
  } catch (err) {
    message.value = err.message || 'Ping failed'
    ok.value = false
  }
  setTimeout(() => { message.value = '' }, 3000)
}
const activateOne = async (mac) => {
  try {
    await enableAgent(mac)
    ok.value = true
    message.value = `Activated ${mac}`
    await load()
  } catch (e) {
    ok.value = false
    message.value = `Activate failed for ${mac}: ${e.message}`
  }
}
const deactivateOne = async (mac) => {
  try {
    await disableAgent(mac)
    message.value = 'Agent deactivated!'
    ok.value = true
    await load()
  } catch (err) {
    message.value = err.message || 'Deactivation failed'
    ok.value = false
  }
  setTimeout(() => { message.value = '' }, 3000)
}
const deleteOne = async (mac) => {
  if (!confirm('Are you sure you want to delete this agent?')) return
  try {
    await deleteAgent(mac)
    message.value = 'Agent deleted!'
    ok.value = true
    await load()
  } catch (err) {
    message.value = err.message || 'Deletion failed'
    ok.value = false
  }
  setTimeout(() => { message.value = '' }, 3000)
}

/* ---------- Bulk actions ---------- */
const activateSelected = () => bulk(selected.value, enableAgent, 'activate')
const deactivateSelected = () => bulk(selected.value, disableAgent, 'deactivate')
const deleteSelected = () => {
  if (!confirm(`Delete ${selected.value.size} selected agents?`)) return
  bulk(selected.value, deleteAgent, 'delete')
}
const bulk = async (set, fn, verb) => {
  if (set.size === 0) return
  
  const results = { success: 0, error: 0 }
  const promises = Array.from(set).map(async mac => {
    try {
      await fn(mac)
      results.success++
    } catch (err) {
      console.error(`Failed to ${verb} ${mac}:`, err)
      results.error++
    }
  })
  
  await Promise.all(promises)
  message.value = `${verb} completed: ${results.success} succeeded, ${results.error} failed`
  ok.value = results.error === 0
  await load()
  setTimeout(() => { message.value = '' }, 3000)
}

/* ---------- Associate (simple modal) ---------- */
const showAssociate = ref(false)
const associateSiteId = ref('')
const openAssociate = () => {
  if (selected.value.size === 0) return
  associateSiteId.value = ''
  showAssociate.value = true
}

const closeAssociate = () => {
  showAssociate.value = false
}
const confirmAssociate = async () => {
  const siteId = associateSiteId.value.trim()
  if (!siteId) {
    message.value = 'Please enter a site ID'
    ok.value = false
    setTimeout(() => { message.value = '' }, 3000)
    return
  }

  const macs = Array.from(selected.value)
  const results = { success: 0, error: 0 }
  
  for (const mac of macs) {
    try {
      await associateAgent(mac, siteId)
      results.success++
    } catch (err) {
      console.error(`Failed to associate ${mac}:`, err)
      results.error++
    }
  }
  
  message.value = `Associated ${results.success} agents, ${results.error} failed`
  ok.value = results.error === 0
  closeAssociate()
  await load()
  setTimeout(() => { message.value = '' }, 3000)
}
</script>

<style scoped>
@import url('https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css');

.agents-page { padding: 16px; max-width: 1100px; margin: 0 auto; }
.card {
  background: var(--panel-grey);
  border: 1px solid rgba(255,255,255,.06);
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0,0,0,.35), inset 0 1px 0 rgba(255,255,255,.03);
  padding: 16px;
}

.head { display:flex; align-items:center; justify-content:space-between; gap:12px; margin-bottom:10px; }
.title h1 { margin:0; font-size:20px; }
.muted { color: var(--text-secondary); font-size: 13px; }

.toolbar {
  display:flex; gap:10px; align-items:center; flex-wrap: wrap;
  padding: 8px 0 12px;
  border-top: 1px solid var(--divider-grey);
  border-bottom: 1px solid var(--divider-grey);
}

.search-wrap {
  display:flex; align-items:center; gap:8px;
  padding: 8px 10px; border:1px solid var(--divider-grey);
  border-radius: 10px; background: rgba(31,41,55,.35);
  min-width: 260px;
}
.search-wrap i { opacity:.8; }
.search {
  border:none; outline:none; background: transparent; color: var(--text-primary); width: 260px;
}

.filters { display:flex; gap:6px; flex-wrap: wrap; }
.pill {
  display:inline-flex; align-items:center; gap:8px;
  padding:8px 10px; border-radius: 999px;
  background: rgba(255,255,255,.04); border: 1px solid rgba(255,255,255,.08);
  color: var(--text-primary); cursor: pointer;
  font-size: 12px;
}
.pill.active { border-color: var(--primary-accent); background: rgba(0,194,194,.10); }
.dot { width:8px; height:8px; border-radius:50%; display:inline-block; }
.dot.green { background:#22c55e; } .dot.red { background:#ef4444; } .dot.amber{ background:#f59e0b; } .dot.gray{ background:#6b7280; }

.actions { margin-left:auto; display:flex; gap:8px; }

.btn { display:inline-flex; align-items:center; gap:8px; border-radius:8px; padding:10px 12px; cursor:pointer; border:1px solid transparent; }
.btn.primary { background: var(--primary-accent); color:#0e111a; }
.btn.ghost { background: rgba(255,255,255,.06); border:1px solid rgba(255,255,255,.08); color: var(--text-primary); }
.btn.danger { background: rgba(239,68,68,.12); color:#ef4444; border:1px solid rgba(239,68,68,.35); }
.btn:disabled { opacity:.6; cursor: not-allowed; }
.btn.xs { padding:6px 8px; font-size: 12px; }

.table-wrapper { overflow:auto; margin-top: 12px; border-radius: 12px; border:1px solid var(--divider-grey); }
.table { width: 100%; border-collapse: collapse; min-width: 760px; }
th, td { padding: 10px 12px; text-align: left; border-bottom: 1px solid var(--divider-grey); }
th.sel, td.sel { width: 36px; }
td.right, th.right { text-align: right; }
.mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 12px; }

.badge {
  padding: 6px 10px; border-radius: 999px; font-size: 12px; font-weight: 700;
}
.badge.green { background: rgba(34,197,94,.15); color: #22c55e; }
.badge.red { background: rgba(239,68,68,.15); color: #ef4444; }
.badge.amber { background: rgba(245,158,11,.15); color: #f59e0b; }
.badge.gray { background: rgba(107,114,128,.15); color: #9ca3af; }

.empty { text-align:center; color: var(--text-secondary); }

.message { margin-top: 10px; padding: 10px 12px; border-radius: 8px; font-size: 14px; }
.message.success { background: rgba(34,197,94,.1); color: #22c55e; }
.message.error { background: rgba(239,68,68,.1); color: #ef4444; }

/* Associate modal */
.modal-backdrop {
  position: fixed; inset:0; background: rgba(0,0,0,.45);
  display:grid; place-items:center; z-index: 40;
}
.modal {
  width: 100%; max-width: 440px; background: var(--panel-grey);
  border: 1px solid rgba(255,255,255,.06); border-radius: 12px; padding: 16px;
}
.modal h3 { margin: 0 0 10px; }
.modal input {
  width: 100%;
  background: var(--panel-grey);
  border: 1px solid var(--divider-grey);
  border-radius: 8px;
  padding: 10px 12px;
  color: var(--text-primary);
}
.modal-actions { display:flex; justify-content:flex-end; gap:8px; margin-top: 12px; }

/* Light-theme softening */
[data-theme='light'] .card,
[data-theme='light'] .modal {
  background: #ffffff;
  border-color: var(--panel-border);
  box-shadow: var(--panel-shadow);
}
[data-theme='light'] .search-wrap { background: #fff; border-color: #e5e7eb; }
[data-theme='light'] .table-wrapper { border-color: #e5e7eb; }
[data-theme='light'] .badge.green { background: rgba(34,197,94,.12); color:#16a34a; }
[data-theme='light'] .badge.red   { background: rgba(239,68,68,.12); color:#dc2626; }
[data-theme='light'] .badge.amber { background: rgba(245,158,11,.12); color:#d97706; }
</style>
