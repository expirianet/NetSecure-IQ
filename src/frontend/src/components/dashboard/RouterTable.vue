<!-- /src/frontend/src/components/dashboard/RouterTable.vue -->
<template>
  <div class="routers-page">
    <BackgroundParticles />

    <div class="wrapper">
      <div class="container">
        <div class="card">
          <!-- Header -->
          <div class="header-row">
            <h2>Router Status</h2>
            <div class="actions">
              <button class="btn ghost" @click="loadRouters" :disabled="loading">
                <i class="fas fa-rotate"></i>
                <span>{{ loading ? 'Loading…' : 'Refresh' }}</span>
              </button>
            </div>
          </div>

          <!-- Table -->
          <div class="table-wrapper">
            <table v-if="routers.length" class="routers-table" aria-label="Router status">
              <thead>
                <tr>
                  <th>MAC Address</th>
                  <th>VPN IP</th>
                  <th>Status</th>
                  <th>Last Contact</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="r in routers" :key="r.mac">
                  <td class="mono">{{ r.mac }}</td>
                  <td class="mono">{{ r.ip || '—' }}</td>
                  <td>
                    <span :class="['state', stateClass(r)]" :title="statusTitle(r)">
                      {{ statusLabel(r) }}
                    </span>
                  </td>
                  <td :title="r.lastTime ? new Date(r.lastTime).toLocaleString() : ''">
                    {{ r.lastTime ? formatRelativeTime(r.lastTime) : '—' }}
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
/**
 * RouterTable.vue
 * Backend endpoints:
 *   - GET /api/mikrotik/list   -> [{ mac, ip, status, site_id }]
 *   - GET /api/data/routers    -> [{ MAC|mac, Value|value, Time|time }] (JWT required)
 * Résolution d'URL d'API compatible Webpack:
 *   1) process.env.VUE_APP_API_BASE_URL (Vue CLI)
 *   2) window.__API_BASE_URL__ (override sans rebuild)
 *   3) fallback http://<host>:8000/api
 */
import { ref, onMounted } from 'vue'
import BackgroundParticles from '@/components/common/BackgroundParticles.vue'

function resolveApiBaseUrl () {
  // Vue CLI / Webpack
  let fromVueCli
  try {
    fromVueCli = (typeof process !== 'undefined' && process && process.env && process.env.VUE_APP_API_BASE_URL)
      ? process.env.VUE_APP_API_BASE_URL
      : undefined
  } catch { fromVueCli = undefined }

  // Global override
  const fromWindow = (typeof window !== 'undefined' && window.__API_BASE_URL__)
    ? window.__API_BASE_URL__
    : undefined

  // Fallback
  const fallback = (typeof window !== 'undefined' && window.location)
    ? `${window.location.protocol}//${window.location.hostname}:8000/api`
    : 'http://localhost:8000/api'

  return String(fromVueCli || fromWindow || fallback).replace(/\/+$/, '')
}

const API_BASE_URL = resolveApiBaseUrl()

const routers = ref([]) // [{ mac, ip, assoc, online, lastTime }]
const loading = ref(false)
const error = ref('')

/* Helpers */
function normalizeMAC (mac) {
  return (mac || '').trim().toUpperCase()
}

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
  if (r.online === true) return 'Online (from Influx last ping)'
  if (r.online === false) return 'Offline (from Influx last ping)'
  return `Association: ${r.assoc || 'Unknown'}`
}

/* Data loading */
async function fetchList () {
  const res = await fetch(`${API_BASE_URL}/mikrotik/list`)
  if (!res.ok) {
    const txt = await res.text().catch(() => '')
    throw new Error(`GET /mikrotik/list failed: ${res.status} ${txt || res.statusText}`)
  }
  /** Expected: [{ mac, ip, status, site_id }] */
  const data = await res.json()
  if (!Array.isArray(data)) throw new Error('Invalid /mikrotik/list payload')

  return data.map(row => {
    const mac = normalizeMAC(row.mac || row.MAC)
    const assoc = (row.status || '').trim() || 'Unknown' // "Associated" | "Unassociated" | "Deactivated"
    const ip = row.ip || row.vpn_internal_ip || ''
    return { mac, ip, assoc, online: null, lastTime: null }
  })
}

async function fetchInfluxMap () {
  const token = (typeof localStorage !== 'undefined') ? localStorage.getItem('token') : null
  if (!token) return {} // no auth -> skip
  const res = await fetch(`${API_BASE_URL}/data/routers`, {
    headers: { Authorization: `Bearer ${token}` }
  })
  if (!res.ok) return {} // 401/403 => pas d'enrichissement

  const rows = await res.json()
  const map = {}
  for (const r of rows || []) {
    const mac = normalizeMAC(r.mac || r.MAC)
    const val = String(r.value ?? r.Value ?? '').toLowerCase()
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

    // Merge enrichment
    for (const r of base) {
      const enr = enrich[r.mac]
      if (enr) {
        r.online = enr.online
        r.lastTime = enr.time
      }
    }

    // Sort: online -> associated -> others, then lastTime desc
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

    if (!routers.value.length) {
      routers.value = [
        { mac: 'AA:BB:CC:DD:EE:01', ip: '10.100.99.10', assoc: 'Unassociated', online: null, lastTime: null }
      ]
    }
  } catch (e) {
    // eslint-disable-next-line no-console
    console.error('Failed to load routers:', e)
    error.value = 'Failed to load router data. Check the backend URL or your connection.'
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
  --danger:#ef4444; --success:#22c55e; --info:#38bdf8; --muted:#6b7280;
}

/* Layout */
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
  border-bottom:1px solid rgba(255,255,255,.08); color:var(--text-secondary); white-space:nowrap;
}
.routers-table tbody td { padding:14px; border-bottom:1px solid rgba(255,255,255,.05); color:var(--text-primary); }
.mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; letter-spacing:.2px; }
.empty { text-align:center; color:var(--text-secondary); padding:18px; }

/* Status badges */
.state { display:inline-block; padding:6px 10px; border-radius:999px; font-weight:600; font-size:12px; }
.state.green  { background:rgba(22,163,74,.15); color:#22c55e; }       /* online */
.state.red    { background:rgba(239,68,68,.15); color:#f87171; }       /* offline */
.state.blue   { background:rgba(59,130,246,.15); color:#93c5fd; }      /* associated */
.state.orange { background:rgba(245,158,11,.15); color:#fbbf24; }      /* deactivated */
.state.gray   { background:rgba(156,163,175,.15); color:#d1d5db; }     /* unknown */

/* Toast */
.toast { margin-top:12px; text-align:center; padding:10px 12px; border-radius:8px; font-size:14px; }
.toast.error { background:rgba(239,68,68,.12); color:var(--danger); }
</style>
