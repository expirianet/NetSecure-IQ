<template>
  <div id="app">
    <!-- Barre de navigation (switch auto selon le rôle) -->
    <component :is="navComponent" />

    <!-- Contenu principal -->
    <main class="main-content">
      <router-view />
    </main>

    <!-- Canvas pour particles.js -->
    <div id="particles-js"></div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted, computed, nextTick } from 'vue'
import TopNavigation from '@/components/TopNavigation.vue'
import TopNavigationUser from '@/components/TopNavigationUser.vue'
import TopNavigationOperator from '@/components/TopNavigationOperator.vue'

/* ---------- Rôle & navigation ---------- */
const role = ref(localStorage.getItem('role') || '')

const navComponent = computed(() => {
  if (role.value === 'user') return TopNavigationUser
  if (role.value === 'operator') return TopNavigationOperator
  return TopNavigation // admin / défaut
})

function updateRoleFromStorage() {
  role.value = localStorage.getItem('role') || ''
}

/* ---------- Thème ---------- */
const theme = ref(localStorage.getItem('theme') || 'dark')

watch(theme, (newTheme) => {
  document.documentElement.setAttribute('data-theme', newTheme)
  localStorage.setItem('theme', newTheme)
})

/* ---------- particles.js ---------- */
function initParticles() {
  const dark = theme.value === 'dark'
  if (!window.particlesJS) return

  // Détruire l'ancien canvas
  const oldCanvas = document.querySelector('#particles-js > canvas')
  if (oldCanvas) oldCanvas.remove()

  window.particlesJS('particles-js', {
    particles: {
      number: { value: 80, density: { enable: true, value_area: 800 } },
      color: { value: dark ? '#ffffff' : '#888888' },
      shape: { type: 'circle' },
      opacity: { value: dark ? 0.5 : 0.3 },
      size: { value: 3, random: true },
      line_linked: {
        enable: true,
        distance: 150,
        color: dark ? '#ffffff' : '#E0E0E0',
        opacity: dark ? 0.4 : 0.3,
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

function attachThemeObserver() {
  const observer = new MutationObserver(() => initParticles())
  observer.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['data-theme']
  })
  return observer
}

/* ---------- Lifecycle ---------- */
let themeObserver

onMounted(async () => {
  // Appliquer le thème
  document.documentElement.setAttribute('data-theme', theme.value)

  // Charger particles si nécessaire
  if (!window.particlesJS) {
    const script = document.createElement('script')
    script.src = '/particles/particles.min.js'
    script.onload = () => {
      initParticles()
      themeObserver = attachThemeObserver()
    }
    document.body.appendChild(script)
  } else {
    initParticles()
    themeObserver = attachThemeObserver()
  }

  // Écoutes pour mettre à jour la nav quand l’auth change
  window.addEventListener('storage', updateRoleFromStorage)
  window.addEventListener('auth-changed', updateRoleFromStorage)

  // S’assurer que la nav se rend correctement après le mount
  await nextTick()
  updateRoleFromStorage()
})

onUnmounted(() => {
  window.removeEventListener('storage', updateRoleFromStorage)
  window.removeEventListener('auth-changed', updateRoleFromStorage)
  if (themeObserver?.disconnect) themeObserver.disconnect()
})
</script>

<style>
/* Variables par défaut */
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
}

/* Override pour le mode clair */
[data-theme='light'] {
  --bg-dark: #ffffff;
  --panel-grey: #f5f5f5;
  --divider-grey: #e0e0e0;
  --text-primary: #1f2937;
  --text-secondary: #6b7280;
  --primary-accent: #00a8a8;
  --primary-hover: #008a8a;
  --danger: #dc2626;
  --success: #22c55e;
}

/* Style global */
body {
  margin: 0;
  padding: 0;
  background-color: var(--bg-dark);
  color: var(--text-primary);
  font-family: 'Inter', sans-serif;
  transition: background-color 0.3s ease, color 0.3s ease;
  min-height: 100vh;
}

/* Contenu principal */
.main-content {
  padding-top: 80px; /* Hauteur de la barre de navigation */
  min-height: calc(100vh - 80px);
  box-sizing: border-box;
}

/* Reset */
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

/* Canvas particles */
#particles-js {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: -1;
  background-color: var(--bg-dark);
  transition: background-color 0.3s ease;
}
</style>
