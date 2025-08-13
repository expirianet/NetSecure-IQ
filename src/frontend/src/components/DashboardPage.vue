<template>
  <div class="dashboard-page">
    <div id="particles-js"></div>

    <div class="dashboard-wrapper">
      <div class="dashboard-container">
        <div class="dashboard-card">
          <h1 class="dashboard-title">Dashboard</h1>

          <div v-if="isAdminOrOperator" class="dashboard-actions">
            <router-link class="dashboard-button" to="/routertable">
              <i class="fas fa-network-wired"></i>
              Router Info
            </router-link>

            <router-link v-if="!needsOrganization" class="dashboard-button" to="/adduser">
              <i class="fas fa-user-plus"></i>
              Add User
            </router-link>

            <router-link v-else class="dashboard-button" to="/organization">
              <i class="fas fa-building"></i>
              Organization Info
            </router-link>

            <router-link v-if="canAddOperator" class="dashboard-button" to="/addoperator">
              <i class="fas fa-user-shield"></i>
              Add Operator
            </router-link>
          </div>

          <div v-else class="welcome-message">
            <p>Welcome to your dashboard. Select an option to get started.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, nextTick } from 'vue'

const role = computed(() => (localStorage.getItem('role') || '').toLowerCase())
const isAdmin = computed(() => role.value === 'administrator')
const isOperator = computed(() => role.value === 'operator')
const isAdminOrOperator = computed(() => isAdmin.value || isOperator.value)

const hasOrganization = computed(() => !!String(localStorage.getItem('organization_id') || '').trim())
const needsOrganization = computed(() => !hasOrganization.value)
const canAddOperator = computed(() => isAdmin.value || (isOperator.value && hasOrganization.value))

function renderParticles() {
  const dark =
    document.documentElement.getAttribute('data-theme') === 'dark' ||
    document.documentElement.classList.contains('dark')

  const old = document.querySelector('#particles-js > canvas')
  if (old) old.remove()

  if (window.particlesJS) {
    window.particlesJS('particles-js', {
      particles: {
        number: { value: 80, density: { enable: true, value_area: 800 } },
        color: { value: dark ? '#ffffff' : '#555555' },
        shape: { type: 'circle' },
        opacity: { value: 0.5 },
        size: { value: 3, random: true },
        line_linked: {
          enable: true, distance: 150,
          color: dark ? '#ffffff' : '#888888', opacity: 0.4, width: 1
        },
        move: { enable: true, speed: 6, direction: 'none', out_mode: 'bounce' }
      },
      interactivity: {
        detect_on: 'canvas',
        events: { onhover: { enable: true, mode: 'repulse' }, onclick: { enable: true, mode: 'push' }, resize: true },
        modes: { repulse: { distance: 200 }, push: { particles_nb: 4 } }
      },
      retina_detect: true
    })
  }
}

onMounted(async () => {
  document.title = 'NetSecure-IQ - Dashboard'
  if (!window.particlesJS) {
    await new Promise(resolve => {
      const s = document.createElement('script')
      s.src = '/particles/particles.min.js'
      s.onload = resolve
      document.body.appendChild(s)
    })
  }
  await nextTick()
  renderParticles()
})
</script>

<style scoped>
#particles-js {
  position: fixed;
  width: 100%;
  height: 100%;
  top: 0; left: 0;
  z-index: 0;
  pointer-events: none; /* ✅ ne bloque jamais les clics */
}
.dashboard-page { min-height: 100vh; display: flex; flex-direction: column; }
.dashboard-wrapper { flex: 1; display:flex; align-items:center; justify-content:center; padding:2rem; position:relative; z-index:1; }
.dashboard-container { width: 100%; max-width: 1000px; margin: 0 auto; }
.dashboard-card { background: var(--panel-grey); border-radius: 16px; padding: 2.5rem; border: 1px solid rgba(255,255,255,.05); box-shadow: 0 0 40px rgba(0,194,194,.05); }
.dashboard-title { text-align:center; margin-bottom:2rem; font-size:2rem; font-weight:600; }
.dashboard-actions { display:grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap:1.5rem; }

.dashboard-button {
  display:flex; flex-direction:column; align-items:center; justify-content:center;
  padding:1.5rem 1rem; min-height:110px; text-decoration:none;
  background-color: var(--button-bg, #f0f4f8); color: var(--button-text, #2c3e50);
  border: 1px solid var(--border-color, #e0e6ed); border-radius:12px; font-weight:500; transition:.2s;
}
.dashboard-button i { font-size:1.75rem; margin-bottom:.75rem; color: var(--primary-color, #3b82f6); }
.dashboard-button:hover { transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0,0,0,.1); background-color: var(--button-hover-bg, #e6edf5); }

.welcome-message { text-align:center; color: var(--text-secondary, #64748b); padding:2rem 0; }
</style>
