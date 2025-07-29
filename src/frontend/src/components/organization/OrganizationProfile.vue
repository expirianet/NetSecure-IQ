<template>
  <div class="login-page">
    <!-- animated particles background -->
    <div id="particles-js"></div>

    <div class="login-wrapper">
      <div class="login-container">
        <div class="login-card">
          <!-- Page header -->
          <h2 class="login-title">Organization Profile</h2>
          <h3 class="login-subtitle">View and manage your organization's information</h3>

          <!-- Profile content -->
          <div class="login-form">
            <!-- Organization Info -->
            <div class="form-section">
              <h4><i class="fas fa-building"></i> Organization Information</h4>
              <input readonly :value="org.name" />
              <input readonly :value="org.vatNumber" />
              <input readonly :value="org.address" />
              <div class="form-row">
                <input readonly :value="org.city" />
                <input readonly :value="org.state" />
                <input readonly :value="org.zipCode" />
              </div>
              <input readonly :value="org.email" />
              <input readonly :value="org.phone" />
              <input readonly :value="org.pecEmail || '—'" />
              <input readonly :value="org.sdiCode || '—'" />
            </div>

            <!-- Company Manager -->
            <div class="form-section">
              <h4><i class="fas fa-user-tie"></i> Company Manager</h4>
              <input readonly :value="org.manager.name" />
              <input readonly :value="org.manager.email" />
              <input readonly :value="org.manager.phone" />
            </div>

            <!-- Data Controller -->
            <div class="form-section">
              <h4><i class="fas fa-shield-alt"></i> Data Controller</h4>
              <input readonly :value="org.controller.name" />
              <input readonly :value="org.controller.email" />
              <input readonly :value="org.controller.phone" />
            </div>

            <!-- Data Processor -->
            <div class="form-section">
              <h4><i class="fas fa-database"></i> Data Processor</h4>
              <input readonly :value="org.processor.name" />
              <input readonly :value="org.processor.email" />
              <input readonly :value="org.processor.phone" />
            </div>

            <!-- Edit button -->
            <div class="form-actions">
              <router-link to="/organization/edit" class="btn btn-primary">
                <i class="fas fa-edit"></i> Edit Profile
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'

// --- organization data (replace with real API/store) ---
const org = ref({
  name: 'Acme Corporation',
  vatNumber: 'IT12345678901',
  address: '123 Business Street',
  city: 'Milan',
  state: 'MI',
  zipCode: '20100',
  email: 'info@acmecorp.com',
  phone: '+39 02 1234567',
  pecEmail: 'acmecorp@pec.it',
  sdiCode: 'ABC12345',
  manager: {
    name: 'Mario Rossi',
    email: 'mario.rossi@acmecorp.com',
    phone: '+39 333 1234567'
  },
  controller: {
    name: 'Laura Bianchi',
    email: 'laura.bianchi@acmecorp.com',
    phone: '+39 345 9876543'
  },
  processor: {
    name: 'Giuseppe Verdi',
    email: 'giuseppe.verdi@acmecorp.com',
    phone: '+39 333 7654321'
  }
})

// --- particles setup (copied from LoginForm.vue) ---
function renderParticles() {
  const theme = document.documentElement.getAttribute('data-theme')
  const isDark = theme === 'dark'

  const oldCanvas = document.querySelector('#particles-js > canvas')
  if (oldCanvas) oldCanvas.remove()

  window.particlesJS('particles-js', {
    particles: {
      number: { value: 80, density: { enable: true, value_area: 800 } },
      color: { value: isDark ? '#ffffff' : '#555555' },
      shape: { type: 'circle' },
      opacity: { value: isDark ? 0.5 : 0.5 },
      size: { value: 3, random: true },
      line_linked: {
        enable: true,
        distance: 150,
        color: isDark ? '#ffffff' : '#888888',
        opacity: isDark ? 0.4 : 0.4,
        width: 1
      },
      move: { enable: true, speed: 6, direction: 'none', out_mode: 'bounce' }
    },
    interactivity: {
      detect_on: 'canvas',
      events: {
        onhover: { enable: true, mode: 'repulse' },
        onclick: { enable: true, mode: 'push' },
        resize: true
      },
      modes: {
        repulse: { distance: 200 },
        push: { particles_nb: 4 }
      }
    },
    retina_detect: true
  })
}

async function initializeParticles() {
  // 1) wait DOM
  await nextTick()

  // 2) load script if needed
  if (!window.particlesJS) {
    await new Promise(resolve => {
      const s = document.createElement('script')
      s.src = '/particles/particles.min.js'
      s.onload = resolve
      document.body.appendChild(s)
    })
  }

  // 3) initial render
  renderParticles()

  // 4) observe theme changes
  new MutationObserver(muts => {
    muts.forEach(m => {
      if (m.attributeName === 'data-theme' || m.attributeName === 'class') {
        renderParticles()
      }
    })
  }).observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['data-theme', 'class']
  })
}

// --- mount ---
onMounted(() => {
  document.title = 'NetSecure-IQ - Organization Profile'
  initializeParticles()
})
</script>

<style scoped>
@import url('https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css');

:root {
  --bg-dark: #0e111a;
  --panel-grey: #1a1d26;
  --divider-grey: #2a2d36;
  --text-primary: #f5f7fa;
  --text-secondary: #9ca3af;
  --primary-accent: #00c2c2;
  --primary-hover: #00a7a7;
  --border-radius: 12px;
  --transition: all 0.2s ease;
}
/* hide any theme‐toggle you might have globally */
.theme-switcher { display: none!important; }

.login-page {
  position: relative;
  min-height: 100vh;
  background-color: var(--bg-dark);
  overflow: hidden;
}

#particles-js {
  position: fixed;
  top: 0; left: 0;
  width: 100vw; height: 100vh;
  z-index: 0;
  background-color: var(--bg-dark);
  transition: background-color 0.3s ease;
  pointer-events: none;
}
[data-theme='light'] #particles-js {
  background-color: #E0E0E0;
}

.login-wrapper {
  position: relative; z-index: 10;
  display: flex; align-items: center; justify-content: center;
  padding: 32px; min-height: 100vh;
}

.login-container {
  width: 100%; max-width: 800px;
}

.login-card {
  background-color: var(--panel-grey);
  border-radius: 16px;
  padding: 32px;
  border: 1px solid rgba(255,255,255,0.05);
  box-shadow: 0 0 40px rgba(0,194,194,0.05);
  box-sizing: border-box;
}

.login-title {
  text-align: center;
  font-size: 24px;
  font-weight: 600;
  color: var(--primary-accent);
  margin-bottom: 8px;
}
.login-subtitle {
  text-align: center;
  font-size: 18px;
  color: var(--text-primary);
  margin-bottom: 32px;
  font-weight: 500;
}

.login-form {
  display: flex; flex-direction: column; gap: 16px;
}

.form-section {
  background-color: rgba(31,41,55,0.3);
  border-radius: 8px;
  padding: 16px;
  transition: background-color 0.3s ease;
}
:root:not([data-theme='dark']) .form-section {
  background-color: rgba(243,244,246,0.5);
  border: 1px solid rgba(209,213,219,0.5);
}
.form-section h4 {
  color: var(--primary-accent);
  margin: 0 0 16px;
  font-size: 16px; font-weight: 500;
  display: flex; align-items: center; gap: 8px;
}

.login-form input {
  width: 100%;
  background-color: var(--panel-grey);
  border: 1px solid var(--divider-grey);
  border-radius: 6px;
  padding: 12px 14px;
  font-size: 14px;
  color: var(--text-primary);
  transition: var(--transition);
  margin-bottom: 8px;
}

.form-row {
  display: flex; gap: 16px; margin-bottom: 8px;
}
.form-row input { margin-bottom: 0; }

.form-actions {
  display: flex; justify-content: center; margin-top: 24px; gap: 16px;
}

.btn {
  padding: 12px 20px;
  border: none;
  border-radius: 6px;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: var(--transition);
  display: inline-flex; align-items: center; gap: 8px;
}

.btn-primary {
  background-color: var(--primary-accent);
  color: #0e111a;
}
.btn-primary:hover {
  background-color: var(--primary-hover);
}

@media (max-width: 768px) {
  .login-wrapper { padding: 16px; }
  .login-card   { padding: 24px; }
  .form-row     { flex-direction: column; gap: 8px; }
  .form-actions { flex-direction: column; }
  .btn          { width: 100%; }
}
</style>
