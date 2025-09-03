<template>
  <div class="opdash-page">
    <BackgroundParticles />

    <div class="opdash-wrapper">
      <div class="opdash-shell">
        <!-- Header -->
        <header class="opdash-header">
          <div class="title">
            <h1>Operator Control</h1>
            <p class="muted">Quick overview and direct access to tools</p>
          </div>

          <div class="status-pill" v-if="hasOrganization">
            <span class="dot"></span> org&nbsp;: <code>{{ orgId }}</code>
          </div>
        </header>

        <!-- Mini-views grid -->
        <section class="grid">
          <!-- RouterTable -->
          <button class="tile" @click="go('/routertable')" aria-label="Open RouterTable">
            <div class="tile-head">
              <div class="icon-wrap"><i class="fas fa-network-wired"></i></div>
              <div class="labels">
                <h3>RouterTable</h3>
                <p class="muted">Status summary</p>
              </div>
            </div>

            <div class="row kpis">
              <div class="kpi">
                <span class="kpi-val green">{{ routers.online }}</span>
                <span class="kpi-label">Online</span>
              </div>
              <div class="kpi">
                <span class="kpi-val red">{{ routers.offline }}</span>
                <span class="kpi-label">Offline</span>
              </div>
              <div class="kpi">
                <span class="kpi-val amber">{{ routers.unknown }}</span>
                <span class="kpi-label">Unknown</span>
              </div>
            </div>

            <div class="spark">
              <LineChart :chartData="sparks.router" :options="sparkOpts" />
            </div>

            <ul class="mini-list">
              <li v-for="r in miniRouters" :key="r.mac">
                <span class="mono">{{ r.mac }}</span>
                <span :class="['badge', stateClass(r.status)]">{{ label(r.status) }}</span>
              </li>
            </ul>
          </button>

          <!-- Organization -->
          <button class="tile" @click="go('/organization')" aria-label="Open Organization">
            <div class="tile-head">
              <div class="icon-wrap"><i class="fas fa-building"></i></div>
              <div class="labels">
                <h3>Organization</h3>
                <p class="muted">Profile & compliance</p>
              </div>
            </div>

            <div class="org-card">
              <div class="org-row">
                <span class="hint">Name</span>
                <span class="val">{{ orgName || '—' }}</span>
              </div>
              <div class="org-row">
                <span class="hint">City</span>
                <span class="val">{{ orgCity || '—' }}</span>
              </div>
              <div class="org-row">
                <span class="hint">Email</span>
                <span class="val">{{ orgEmail || '—' }}</span>
              </div>
            </div>

            <div class="cta">Open profile</div>
          </button>

          <!-- Add User -->
          <button
            class="tile"
            :disabled="!hasOrganization"
            @click="go('/adduser')"
            aria-label="Open Add User"
          >
            <div class="tile-head">
              <div class="icon-wrap"><i class="fas fa-user-plus"></i></div>
              <div class="labels">
                <h3>Add User</h3>
                <p class="muted">Invite a user</p>
              </div>
            </div>

            <div class="giant-number">
              <span class="num">{{ recentUsers }}</span>
              <span class="unit">recent adds</span>
            </div>

            <div class="spark">
              <LineChart :chartData="sparks.users" :options="sparkOpts" />
            </div>

            <div v-if="!hasOrganization" class="warn">Attach an organization first</div>
            <div v-else class="cta">Create user</div>
          </button>

          <!-- Agents -->
          <button class="tile" @click="go('/agents')" aria-label="Open Agents">
            <div class="tile-head">
              <div class="icon-wrap"><i class="fas fa-microchip"></i></div>
              <div class="labels">
                <h3>Agents</h3>
                <p class="muted">Deployed MikroTik</p>
              </div>
            </div>

            <div class="row kpis">
              <div class="kpi">
                <span class="kpi-val green">{{ agents.associated }}</span>
                <span class="kpi-label">Associated</span>
              </div>
              <div class="kpi">
                <span class="kpi-val amber">{{ agents.unassociated }}</span>
                <span class="kpi-label">Unassociated</span>
              </div>
              <div class="kpi">
                <span class="kpi-val red">{{ agents.deactivated }}</span>
                <span class="kpi-label">Deactivated</span>
              </div>
            </div>

            <div class="spark"><LineChart :chartData="sparks.agents" :options="sparkOpts" /></div>
            <div class="cta">Manage agents</div>
          </button>

          <!-- Agent Pre-Registration -->
          <button class="tile" @click="go('/agents/register')" aria-label="Open Pre-Registration">
            <div class="tile-head">
              <div class="icon-wrap"><i class="fas fa-keyboard"></i></div>
              <div class="labels">
                <h3>Agent Pre-Registration</h3>
                <p class="muted">MikroTik .rsc script</p>
              </div>
            </div>

            <div class="quick-form" @click.stop>
              <input
                v-model="quickMac"
                class="mac-input mono"
                placeholder="AA:BB:CC:DD:EE:FF"
                @input="formatMac"
              />
              <button class="btn" :disabled="!macValid" @click="go('/agents/register')">
                Generate
              </button>
            </div>

            <div class="spark"><LineChart :chartData="sparks.wg" :options="sparkOpts" /></div>
            <div class="cta">Open tool</div>
          </button>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import BackgroundParticles from '@/components/common/BackgroundParticles.vue'
import LineChart from '@/components/common/LineChart.vue'

const router = useRouter()
const go = (p) => router.push(p)

/* --------- Context (org / role) --------- */
const orgId = ref(localStorage.getItem('organization_id') || '')
const hasOrganization = computed(() => !!String(orgId.value).trim())

// Read org profile from storage if present (populated by OrganizationForm)
const orgProfile = computed(() => {
  try { return JSON.parse(localStorage.getItem('organization_profile') || 'null') } catch { return null }
})
const orgName  = computed(() => orgProfile.value?.name)
const orgCity  = computed(() => orgProfile.value?.city)
const orgEmail = computed(() => orgProfile.value?.contact_email)

/* --------- Demo data (local) --------- */
const routers = ref({ online: 7, offline: 2, unknown: 1 })
const miniRouters = ref([
  { mac: 'E4:8D:8C:AA:01:11', status: 'online'  },
  { mac: '58:EF:68:02:7C:22', status: 'offline' },
  { mac: 'C0:56:27:9A:33:44', status: 'unknown' }
])
const agents = ref({ associated: 12, unassociated: 4, deactivated: 1 })
const recentUsers = ref(3)

/* --------- Micro-charts (sparklines) --------- */
const labels = Array.from({ length: 16 }, (_, i) => i + 1)
function series(seed, jitter = 3) {
  let v = seed
  return labels.map(() => (v = Math.max(0, v + (Math.random() * jitter - jitter / 2))))
}
const sparkOpts = {
  plugins: { legend: { display: false }, tooltip: { enabled: false } },
  elements: { point: { radius: 0 }, line: { tension: 0.3, borderWidth: 2 } },
  scales: { x: { display: false }, y: { display: false } }
}
const sparks = ref({
  router: { labels, datasets: [{ data: series(6),      borderColor: '#22c55e' }] },
  users:  { labels, datasets: [{ data: series(2, 2),   borderColor: '#60a5fa' }] },
  agents: { labels, datasets: [{ data: series(5, 3),   borderColor: '#f59e0b' }] },
  wg:     { labels, datasets: [{ data: series(4, 2.5), borderColor: '#00c2c2' }] }
})

/* --------- Mini helpers --------- */
const macRe = /^[0-9A-F]{2}(:[0-9A-F]{2}){5}$/i
const quickMac = ref('')
const macValid = computed(() => macRe.test(quickMac.value))
function formatMac() {
  let v = quickMac.value.replace(/[^0-9a-f]/gi, '').toUpperCase().slice(0, 12)
  quickMac.value = v.match(/.{1,2}/g)?.join(':') ?? ''
}
function label(s) {
  return s === 'online' ? 'Online' : s === 'offline' ? 'Offline' : 'Unknown'
}
function stateClass(s) {
  return s === 'online' ? 'green' : s === 'offline' ? 'red' : 'amber'
}

/* --------- Light animation of figures --------- */
let iv
onMounted(() => {
  iv = setInterval(() => {
    routers.value.online  = Math.max(0, routers.value.online  + (Math.random() > .5 ? 1 : -1))
    routers.value.offline = Math.max(0, routers.value.offline + (Math.random() > .7 ? 1 : -1))
    sparks.value.router.datasets[0].data = series(6)
    sparks.value.agents.datasets[0].data = series(5)
    sparks.value.users.datasets[0].data  = series(2, 2)
    sparks.value.wg.datasets[0].data     = series(4, 2.5)
  }, 4000)
})
onUnmounted(() => { if (iv) clearInterval(iv) })
</script>

<style scoped>
@import url('https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css');

.opdash-page { position: relative; min-height: 100vh; }
.opdash-wrapper { position: relative; z-index: 1; display: flex; justify-content: center; padding: 28px; }
.opdash-shell {
  width: 100%; max-width: 1200px;
  background: var(--panel-grey);
  border: 1px solid rgba(255,255,255,.06);
  border-radius: 18px;
  box-shadow: 0 20px 60px rgba(0,0,0,.35), inset 0 1px 0 rgba(255,255,255,.03);
  padding: 22px;
}

/* Header */
.opdash-header { display:flex; align-items:center; justify-content:space-between; gap: 12px; margin-bottom: 14px; }
.title h1 { margin: 0; font-size: 20px; }
.muted { color: var(--text-secondary); font-size: 13px; }
.status-pill {
  display:inline-flex; align-items:center; gap:8px;
  padding:8px 10px; border-radius:999px;
  background: rgba(0,194,194,.10); border: 1px solid rgba(0,194,194,.25);
  font-size: 12px;
}
.status-pill .dot { width: 8px; height: 8px; border-radius: 50%; background: var(--primary-accent); }

/* Grid */
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 14px;
}

/* Tile */
.tile {
  display: flex; flex-direction: column; gap: 10px;
  padding: 14px; border-radius: 14px;
  background: rgba(31,41,55,.55);
  border: 1px solid rgba(255,255,255,.06);
  cursor: pointer; text-align: left;
  transition: transform .2s, box-shadow .2s, background .2s;
}
.tile:hover { transform: translateY(-2px); box-shadow: 0 10px 24px rgba(0,0,0,.2); background: rgba(41,53,72,.55); }
.tile:disabled { opacity: .6; cursor: not-allowed; }

.tile-head { display:flex; align-items:center; gap: 12px; }
.icon-wrap {
  width: 38px; height: 38px; border-radius: 12px; display:grid; place-items:center;
  background: rgba(0,194,194,.12); border:1px solid rgba(0,194,194,.22);
  color: var(--primary-accent); flex: none;
}
.icon-wrap i { font-size: 18px; }
.labels h3 { margin: 0; font-size: 15px; }
.labels p { margin: 0; }

/* KPIs */
.row.kpis { display:flex; gap: 16px; }
.kpi { display:flex; flex-direction:column; gap:2px; }
.kpi-val { font-weight: 800; font-size: 18px; }
.kpi-val.green { color:#22c55e; }
.kpi-val.red { color:#ef4444; }
.kpi-val.amber { color:#f59e0b; }
.kpi-label { font-size: 11px; color: var(--text-secondary); }

.spark { height: 64px; }

/* Router mini-list */
.mini-list { list-style: none; padding: 0; margin: 0; display: grid; gap: 6px; }
.mini-list li { display:flex; align-items:center; justify-content:space-between; }
.mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 12px; }
.badge { padding: 4px 8px; border-radius: 999px; font-size: 11px; font-weight: 700; }
.badge.green { background: rgba(34,197,94,.15); color: #22c55e; }
.badge.red { background: rgba(239,68,68,.15); color: #ef4444; }
.badge.amber { background: rgba(245,158,11,.15); color: #f59e0b; }

/* Organization card */
.org-card { display:grid; gap:6px; border:1px solid rgba(255,255,255,.06); border-radius:10px; padding:10px; background: rgba(12,16,24,.35); }
.org-row { display:flex; align-items:center; justify-content:space-between; gap: 12px; }
.org-row .hint { color: var(--text-secondary); font-size: 12px; }
.org-row .val { font-size: 13px; }

/* Add User */
.giant-number { display:flex; align-items:baseline; gap:8px; margin-top: 2px; }
.giant-number .num { font-size: 28px; font-weight: 800; }
.giant-number .unit { font-size: 12px; color: var(--text-secondary); }

/* Pre-registration form */
.quick-form { display: grid; grid-template-columns: 1fr auto; gap: 8px; }
.mac-input {
  border: 1px solid var(--divider-grey); background: #0b0e16; color: var(--text-primary);
  border-radius: 8px; padding: 10px 12px; font-size: 13px;
}
.btn {
  border: 1px solid rgba(0,194,194,.35);
  background: rgba(0,194,194,.12); color: var(--primary-accent);
  border-radius: 8px; font-weight: 700; padding: 10px 12px;
}
.btn:disabled { opacity: .5; }

/* CTA */
.cta {
  margin-top: auto;
  display: inline-flex; align-items:center; gap: 8px;
  color: var(--primary-accent); font-weight: 700; font-size: 13px;
}
.warn { margin-top: auto; color: #f59e0b; font-weight: 700; font-size: 13px; }

/* Responsive */
@media (max-width: 780px) { .opdash-wrapper { padding: 16px; } .opdash-shell { padding: 16px; } }
</style>

<style>
/* --- LIGHT THEME OVERRIDES (does not affect dark) --- */

[data-theme='light'] .opdash-shell {
  background: #fff !important;
  border-color: var(--panel-border) !important;
  box-shadow: var(--panel-shadow) !important;
}

/* lighter cards/tiles */
[data-theme='light'] .tile {
  background: #ffffff !important;
  border-color: var(--panel-border) !important;
  box-shadow: 0 4px 12px rgba(2,6,23,0.06) !important;
}
[data-theme='light'] .tile:hover {
  background: #f9fafb !important;
  box-shadow: 0 8px 20px rgba(2,6,23,0.08) !important;
}

/* icon chip (soft teal) */
[data-theme='light'] .icon-wrap {
  background: rgba(14,165,165,0.10) !important;
  border-color: rgba(14,165,165,0.28) !important;
  color: var(--primary-accent) !important;
}

/* “Organization” card on a neutral surface */
[data-theme='light'] .org-card {
  background: var(--muted-surface) !important;
  border-color: var(--panel-border) !important;
}

/* pills/cta readable on light background */
[data-theme='light'] .status-pill {
  background: rgba(14,165,165,0.10) !important;
  border-color: rgba(14,165,165,0.28) !important;
  color: var(--text-primary) !important;
}
[data-theme='light'] .cta { color: var(--primary-accent) !important; }

/* MAC mini form */
[data-theme='light'] .mac-input {
  background: #ffffff !important;
  color: #0f172a !important;               /* slate-900 */
  border-color: #e5e7eb !important;        /* gray-200 */
}

/* state badges: same colors with opacity adapted to white background */
[data-theme='light'] .badge.green { background: rgba(34,197,94,.12) !important; color:#16a34a !important; }
[data-theme='light'] .badge.red   { background: rgba(239,68,68,.12) !important; color:#dc2626 !important; }
[data-theme='light'] .badge.amber { background: rgba(245,158,11,.12) !important; color:#d97706 !important; }
</style>
