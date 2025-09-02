<!-- src/frontend/src/components/DashboardPage.vue -->
<template>
  <div class="dash-page">
    <BackgroundParticles />

    <div class="wrapper">
      <div class="container">
        <div class="card">
          <header class="head">
            <div class="title">
              <h1>Dashboard</h1>
              <p class="muted">Overview & quick access</p>
            </div>

            <div v-if="hasOrganization" class="status-pill">
              <span class="dot"></span> org : <code>{{ orgId }}</code>
            </div>
          </header>

          <section class="grid">
            <router-link class="tile" to="/organization">
              <div class="tile-head">
                <div class="icon-wrap"><i class="fas fa-building"></i></div>
                <div class="labels">
                  <h3>Organization</h3>
                  <p class="muted">Profile & compliance</p>
                </div>
              </div>
              <ul class="info">
                <li><span class="hint">Name</span><span class="val">{{ orgName || '—' }}</span></li>
                <li><span class="hint">City</span><span class="val">{{ orgCity || '—' }}</span></li>
                <li><span class="hint">Email</span><span class="val">{{ orgEmail || '—' }}</span></li>
              </ul>
              <div class="cta">Open profile</div>
            </router-link>

            <router-link class="tile" to="/routertable">
              <div class="tile-head">
                <div class="icon-wrap"><i class="fas fa-network-wired"></i></div>
                <div class="labels">
                  <h3>RouterTable</h3>
                  <p class="muted">Router status</p>
                </div>
              </div>
              <p class="muted">Check statuses at a glance.</p>
              <div class="cta">View table</div>
            </router-link>

            <router-link class="tile" to="/adduser">
              <div class="tile-head">
                <div class="icon-wrap"><i class="fas fa-user-plus"></i></div>
                <div class="labels">
                  <h3>Add User</h3>
                  <p class="muted">Invite a user</p>
                </div>
              </div>
              <p v-if="!hasOrganization" class="warn">Attach an organization first</p>
              <div v-else class="cta">Create user</div>
            </router-link>

            <router-link class="tile" to="/addoperator">
              <div class="tile-head">
                <div class="icon-wrap"><i class="fas fa-user-gear"></i></div>
                <div class="labels">
                  <h3>Add Operator</h3>
                  <p class="muted">Invite an operator</p>
                </div>
              </div>
              <div class="cta">Create operator</div>
            </router-link>

            <router-link class="tile" to="/agents">
              <div class="tile-head">
                <div class="icon-wrap"><i class="fas fa-microchip"></i></div>
                <div class="labels">
                  <h3>Agents</h3>
                  <p class="muted">Deployed MikroTik</p>
                </div>
              </div>
              <div class="cta">Manage agents</div>
            </router-link>

            <router-link class="tile" to="/agents/register">
              <div class="tile-head">
                <div class="icon-wrap"><i class="fas fa-keyboard"></i></div>
                <div class="labels">
                  <h3>Pre-Registration</h3>
                  <p class="muted">MikroTik .rsc script</p>
                </div>
              </div>
              <div class="cta">Generate a script</div>
            </router-link>
          </section>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import BackgroundParticles from '@/components/BackgroundParticles.vue'
import { computed, ref } from 'vue'

const orgId = ref(localStorage.getItem('organization_id') || '')
const hasOrganization = computed(() => !!String(orgId.value).trim())

const orgProfile = computed(() => {
  try { return JSON.parse(localStorage.getItem('organization_profile') || 'null') }
  catch { return null }
})
const orgName  = computed(() => orgProfile.value?.name)
const orgCity  = computed(() => orgProfile.value?.city)
const orgEmail = computed(() => orgProfile.value?.contact_email)
</script>

<style scoped>
@import url('https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css');

:root {
  --bg-dark:#0e111a; --panel-grey:#1a1d26; --divider-grey:#2a2d36;
  --text-primary:#f5f7fa; --text-secondary:#9ca3af;
  --primary-accent:#00c2c2; --primary-hover:#00a7a7;
  --danger:#ef4444; --success:#22c55e;
}

.dash-page { position:relative; min-height:100vh; }
.wrapper { position:relative; z-index:1; display:flex; justify-content:center; padding:28px; }
.container { width:100%; max-width:1200px; }
.card {
  background: var(--panel-grey);
  border: 1px solid rgba(255,255,255,.06);
  border-radius: 18px;
  box-shadow: 0 20px 60px rgba(0,0,0,.35), inset 0 1px 0 rgba(255,255,255,.03);
  padding: 22px;
}

/* Header */
.head { display:flex; align-items:center; justify-content:space-between; gap:12px; margin-bottom:14px; }
.title h1 { margin:0; font-size:20px; }
.muted { color: var(--text-secondary); font-size:13px; }
.status-pill {
  display:inline-flex; align-items:center; gap:8px;
  padding:8px 10px; border-radius:999px;
  background: rgba(0,194,194,.10); border: 1px solid rgba(0,194,194,.25);
  font-size: 12px;
}
.status-pill .dot { width: 8px; height: 8px; border-radius: 50%; background: var(--primary-accent); }

/* Grid tiles */
.grid {
  display:grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 14px;
}
.tile {
  display:flex; flex-direction:column; gap:10px;
  padding:14px; border-radius:14px;
  background: rgba(31,41,55,.55);
  border:1px solid rgba(255,255,255,.06);
  text-decoration:none; color:inherit;
  transition: transform .2s, box-shadow .2s, background .2s;
}
.tile:hover { transform: translateY(-2px); box-shadow: 0 10px 24px rgba(0,0,0,.2); background: rgba(41,53,72,.55); }

.tile-head { display:flex; align-items:center; gap:12px; }
.icon-wrap {
  width:38px; height:38px; border-radius:12px; display:grid; place-items:center;
  background: rgba(0,194,194,.12); border:1px solid rgba(0,194,194,.22);
  color: var(--primary-accent); flex:none;
}
.icon-wrap i { font-size:18px; }
.labels h3 { margin:0; font-size:15px; }

.info { list-style:none; display:grid; gap:6px; padding:0; margin:0; }
.info .hint { color: var(--text-secondary); font-size:12px; }
.info .val  { font-size:13px; }

.cta { margin-top:auto; display:inline-flex; align-items:center; gap:8px; color:var(--primary-accent); font-weight:700; font-size:13px; }
.warn { color:#f59e0b; font-weight:700; font-size:13px; }
</style>
