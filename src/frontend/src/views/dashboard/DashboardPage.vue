<template>
  <div class="dashboard-page">
    <!-- Particles Background -->
    <div id="particles-js"></div>

    <!-- Main Content -->
    <div class="dashboard-wrapper">
      <div class="dashboard-container">
        <div class="dashboard-card">
          <h1 class="dashboard-title">Dashboard</h1>
          
          <div v-if="isAdminOrOperator" class="dashboard-actions">
            <button class="dashboard-button" @click="goToRouterInfo">
              <i class="fas fa-network-wired"></i> Router Info
            </button>
            <button v-if="!needsOrganization" class="dashboard-button" @click="addUser">
              <i class="fas fa-user-plus"></i> Add User
            </button>
            <button v-if="needsOrganization" class="dashboard-button" @click="goToOrganizationInfo">
              <i class="fas fa-building"></i> Organization Info
            </button>
            <button v-if="canAddOperator" class="dashboard-button" @click="addOperator">
              <i class="fas fa-user-shield"></i> Add Operator
            </button>
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
import { useRouter } from 'vue-router'

const router = useRouter()

// Computed properties (alignées avec le localStorage)
const role = computed(() => (localStorage.getItem("role") || '').toLowerCase())
const isAdmin = computed(() => role.value === "administrator")
const isOperator = computed(() => role.value === "operator")
const hasOrganization = computed(() => !!localStorage.getItem("organization_id"))
const isAdminOrOperator = computed(() => isAdmin.value || isOperator.value)
const needsOrganization = computed(() => isOperator.value && !hasOrganization.value)
const canAddOperator = computed(() => isAdmin.value || (isOperator.value && hasOrganization.value))

// Navigation methods (alignées avec tes routes)
const goToRouterInfo = () => {
  router.push('/routertable')
}
const addUser = () => {
  router.push('/adduser')
}
const goToOrganizationInfo = () => {
  router.push('/organization/edit')
}
const addOperator = () => {
  router.push('/addoperator')
}

// Particles initialization
const renderParticles = () => {
  const dark = document.documentElement.getAttribute('data-theme') === 'dark' || 
              document.documentElement.classList.contains('dark')
  
  const old = document.querySelector('#particles-js > canvas')
  if (old) old.remove()

  const savedTheme = localStorage.getItem('theme')
  const isDark = savedTheme ? savedTheme === 'dark' : dark

  if (window.particlesJS) {
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
}

const initializeParticles = async () => {
  if (!document.getElementById('particles-js')) {
    await new Promise(resolve => setTimeout(resolve, 50));
    return initializeParticles();
  }
  
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme) {
    document.documentElement.setAttribute('data-theme', savedTheme)
    if (savedTheme === 'dark') {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }
  
  if (!window.particlesJS) {
    await new Promise((resolve) => {
      const script = document.createElement('script');
      script.src = '/particles/particles.min.js';
      script.onload = resolve;
      document.body.appendChild(script);
    });
  }
  
  await nextTick();
  renderParticles();
  
  const obs = new MutationObserver((mutations) => {
    for (const m of mutations) {
      if (m.attributeName === 'data-theme' || m.attributeName === 'class') {
        const theme = document.documentElement.getAttribute('data-theme') || 
                     (document.documentElement.classList.contains('dark') ? 'dark' : 'light')
        localStorage.setItem('theme', theme)
        
        const old = document.querySelector('#particles-js > canvas')
        if (old) old.remove()
        renderParticles()
      }
    }
  });

  obs.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['data-theme', 'class']
  });
}

// Lifecycle hooks
onMounted(() => {
  document.title = 'NetSecure-IQ - Dashboard'
  initializeParticles()
})
</script>

<style scoped>
.dashboard-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: var(--background-color, #f5f7fa);
  color: var(--text-color, #2c3e50);
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
}

#particles-js {
  position: fixed;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  z-index: 0;
}

.dashboard-wrapper {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  position: relative;
  z-index: 1;
}

.dashboard-container {
  width: 100%;
  max-width: 1000px;
  margin: 0 auto;
}

.dashboard-card {
  background-color: var(--panel-grey);
  border-radius: 16px;
  padding: 2.5rem;
  box-shadow: 0 0 40px rgba(0, 194, 194, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.dashboard-title {
  text-align: center;
  color: var(--text-color, #2c3e50);
  margin-bottom: 2rem;
  font-size: 2rem;
  font-weight: 600;
}

.dashboard-actions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
  margin-top: 1.5rem;
}

.dashboard-button {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 1.5rem 1rem;
  background-color: var(--button-bg, #f0f4f8);
  color: var(--button-text, #2c3e50);
  border: 1px solid var(--border-color, #e0e6ed);
  border-radius: 12px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  min-height: 100px;
}

.dashboard-button i {
  font-size: 1.75rem;
  margin-bottom: 0.75rem;
  color: var(--primary-color, #3b82f6);
}

.dashboard-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  background-color: var(--button-hover-bg, #e6edf5);
}

.welcome-message {
  text-align: center;
  padding: 2rem 0;
  color: var(--text-secondary, #64748b);
}

/* Dark mode styles */
:global(.dark) .dashboard-card {
  background-color: rgba(30, 41, 59, 0.7);
  border-color: rgba(255, 255, 255, 0.1);
}

:global(.dark) .dashboard-button {
  background-color: rgba(30, 41, 59, 0.7);
  border-color: rgba(255, 255, 255, 0.1);
  color: #e2e8f0;
}

:global(.dark) .dashboard-button:hover {
  background-color: rgba(41, 55, 78, 0.7);
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .dashboard-actions {
    grid-template-columns: 1fr;
  }
  
  .dashboard-card {
    padding: 1.5rem;
  }
}
</style>
