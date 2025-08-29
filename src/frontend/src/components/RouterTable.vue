<!-- src/frontend/src/components/RouterTable.vue -->
<template>
  <div class="routers-page">
    <BackgroundParticles />

    <div class="wrapper">
      <div class="container">
        <div class="card">
          <!-- Simple header -->
          <div class="header-row">
            <h2>Router Status Table</h2>
            <button class="btn ghost" @click="loadRouters" :disabled="loading">
              <i class="fas fa-rotate"></i>
              <span>{{ loading ? 'Loading…' : 'Refresh' }}</span>
            </button>
          </div>

          <!-- Table -->
          <div class="table-wrapper">
            <table v-if="routers.length" class="routers-table" aria-label="Router status">
              <thead>
                <tr>
                  <th>MAC Address</th>
                  <th>Status</th>
                  <th>Organization</th>
                  <th>Site</th>
                  <th>Last Contact</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="r in routers" :key="r.mac">
                  <td class="mono">{{ r.mac }}</td>
                  <td><span :class="['state', stateClass(r.status)]">{{ label(r.status) }}</span></td>
                  <td>{{ r.organization }}</td>
                  <td>{{ r.site }}</td>
                  <td :title="new Date(r.time).toLocaleString()">
                    {{ formatRelativeTime(r.time) }}
                  </td>
                </tr>
              </tbody>
            </table>

            <div v-else class="empty">No routers found.</div>
          </div>

          <p v-if="error" class="toast error">{{ error }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import BackgroundParticles from '@/components/BackgroundParticles.vue'
import { ref, onMounted } from 'vue'

const routers = ref([])
const loading = ref(false)
const error = ref('')

function normalize(router) {
  // Vérifier si le routeur est en ligne (dernier contact < 5 minutes)
  const lastSeen = new Date(router.last_seen || router.time || new Date().toISOString())
  const minutesSinceLastSeen = (new Date() - lastSeen) / (1000 * 60)
  const isOnline = minutesSinceLastSeen < 5
  
  return {
    mac: router.mac || '—',
    status: isOnline ? 'online' : 'offline',
    time: router.last_seen || router.time || new Date().toISOString(),
    organization: router.organization || '—',
    site: router.site || '—',
    site_id: router.site_id || null
  }
}
function label(s) {
  return s === 'online' ? 'Online' : s === 'offline' ? 'Offline' : 'Unknown'
}
function stateClass(s) {
  return s === 'online' ? 'green' : s === 'offline' ? 'red' : 'orange'
}
function formatDate(iso) {
  try { return new Date(iso).toLocaleString() } catch { return '—' }
}

function formatRelativeTime(dateString) {
  try {
    const date = new Date(dateString)
    const now = new Date()
    const diffInSeconds = Math.floor((now - date) / 1000)
    
    if (diffInSeconds < 60) return 'Just now'
    
    const diffInMinutes = Math.floor(diffInSeconds / 60)
    if (diffInMinutes < 60) return `${diffInMinutes}m ago`
    
    const diffInHours = Math.floor(diffInMinutes / 60)
    if (diffInHours < 24) return `${diffInHours}h ago`
    
    const diffInDays = Math.floor(diffInHours / 24)
    if (diffInDays < 30) return `${diffInDays}d ago`
    
    return date.toLocaleDateString()
  } catch {
    return '—'
  }
}

async function loadRouters() {
  loading.value = true
  error.value = ''
  try {
    const API_BASE_URL = 'http://localhost:8000/api'
    const response = await fetch(`${API_BASE_URL}/mikrotik/list`)
    if (!response.ok) throw new Error('Failed to fetch routers')
    
    const data = await response.json()
    if (!Array.isArray(data)) {
      throw new Error('Invalid response format')
    }
    
    // Trier par statut (en ligne d'abord) puis par date de dernier contact
    routers.value = data
      .map(normalize)
      .sort((a, b) => {
        if (a.status === b.status) {
          return new Date(b.time) - new Date(a.time)
        }
        return a.status === 'online' ? -1 : 1
      })
    
  } catch (e) {
    console.error('Failed to load routers:', e)
    error.value = 'Failed to load router data. Please try again.'
    // Garder les données existantes en cas d'erreur
    if (!routers.value.length) {
      routers.value = [
        { mac: 'E4:8D:8C:AA:01:11', status: 'offline', time: new Date().toISOString() },
      ]
    }
  } finally {
    loading.value = false
  }
}

onMounted(loadRouters)
</script>

<style scoped>
@import url('https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css');

:root {
  --bg-dark:#0e111a; --panel-grey:#1a1d26; --divider-grey:#2a2d36;
  --text-primary:#f5f7fa; --text-secondary:#9ca3af;
  --primary-accent:#00c2c2; --primary-hover:#00a7a7;
  --danger:#ef4444; --success:#22c55e;
}

/* Card layout */
.routers-page { position:relative; min-height:100vh; overflow:hidden; }
.wrapper { position:relative; z-index:10; display:flex; justify-content:center; padding:32px; }
.container { width:100%; max-width:1000px; }
.card {
  background: var(--panel-grey);
  border-radius: 16px;
  padding: 20px;
  border: 1px solid rgba(255,255,255,.05);
  box-shadow: 0 0 40px rgba(0,194,194,.05);
}

/* Header */
.header-row { display:flex; align-items:center; justify-content:space-between; gap:12px; }
.header-row h2 { margin:0; font-size:20px; color:var(--text-primary); }

/* Buttons */
.btn {
  display:inline-flex; align-items:center; gap:8px;
  border-radius:10px; padding:10px 14px; font-weight:600;
  border:1px solid rgba(255,255,255,.08); background:rgba(255,255,255,.06);
  color:var(--text-primary); cursor:pointer; transition:.15s;
}
.btn:hover { background:rgba(255,255,255,.10); }
.btn:disabled { opacity:.5; cursor:not-allowed; }
.btn.ghost { background:rgba(255,255,255,.06); }

/* Table */
.table-wrapper { margin-top:14px; overflow:auto; border-radius:12px; border:1px solid rgba(255,255,255,.06); }
.routers-table { width:100%; border-collapse:separate; border-spacing:0; font-size:14px; }
.routers-table thead th {
  text-align:left; padding:12px 14px; background:rgba(255,255,255,.04);
  border-bottom:1px solid rgba(255,255,255,.08); color:var(--text-secondary);
}
.routers-table tbody td { padding:14px; border-bottom:1px solid rgba(255,255,255,.05); }
.mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; letter-spacing:.2px; }
.empty { text-align:center; color:var(--text-secondary); padding:18px; }

/* Status badges */
.state { display:inline-block; padding:6px 10px; border-radius:999px; font-weight:600; font-size:12px; }
.state.green  { background:rgba(22,163,74,.15); color:#22c55e; }
.state.red    { background:rgba(239,68,68,.15); color:#f87171; }
.state.orange { background:rgba(245,158,11,.15); color:#fbbf24; }

/* Toast */
.toast { margin-top:12px; text-align:center; padding:10px 12px; border-radius:8px; font-size:14px; }
.toast.error { background:rgba(239,68,68,.12); color:var(--danger); }
</style>
