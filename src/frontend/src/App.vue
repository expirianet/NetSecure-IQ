<template>
  <div id="app">
    <nav>
      <div class="nav-links">
        <router-link to="/login">Login</router-link> |
        <router-link to="/register">Register</router-link>
      </div>
      <button class="theme-toggle" @click="toggleTheme" aria-label="Toggle theme">
        <span v-if="isDark">☀️</span>
        <span v-else>🌙</span>
      </button>
    </nav>

    <router-view />

    <!-- Canvas pour particles.js -->
    <div id="particles-js"></div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'

// -- État du thème
const theme = ref(localStorage.getItem('theme') || 'dark')
const isDark = computed(() => theme.value === 'dark')

// -- Fonction pour basculer le thème
function toggleTheme() {
  theme.value = isDark.value ? 'light' : 'dark'
  localStorage.setItem('theme', theme.value)
}

// -- Mettre à jour l'attribut data-theme sur <html>
watch(theme, (newTheme) => {
  document.documentElement.setAttribute('data-theme', newTheme)
})

// -- Fonction d'initialisation de particles.js
function initParticles() {
  const dark = theme.value === 'dark'
  // si particlesJS n'existe pas, on arrête
  if (!window.particlesJS) return

  // Destruction du canvas précédent
  const oldCanvas = document.querySelector('#particles-js > canvas')
  if (oldCanvas) oldCanvas.remove()

  // Configuration des particules
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

// -- Monter le composant
onMounted(() => {
  // Appliquer d'abord le thème
  document.documentElement.setAttribute('data-theme', theme.value)

  // Charger dynamiquement particles.min.js si nécessaire
  if (!window.particlesJS) {
    const script = document.createElement('script')
    script.src = '/particles/particles.min.js'
    script.onload = () => {
      initParticles()

      // Observer les changements de data-theme
      const observer = new MutationObserver(() => initParticles())
      observer.observe(document.documentElement, {
        attributes: true,
        attributeFilter: ['data-theme']
      })
    }
    document.body.appendChild(script)
  } else {
    // déjà chargé
    initParticles()

    const observer = new MutationObserver(() => initParticles())
    observer.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ['data-theme']
    })
  }
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
  background-color: var(--bg-dark);
  color: var(--text-primary);
  font-family: 'Inter', sans-serif;
  transition: background-color 0.3s ease, color 0.3s ease;
}

/* Navbar */
nav {
  position: relative;
  z-index: 20;
  padding: 16px;
  background-color: var(--panel-grey);
  border-bottom: 1px solid var(--divider-grey);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.nav-links a {
  margin: 0 8px;
  color: var(--primary-accent);
  text-decoration: none;
}
.nav-links a:hover {
  color: var(--primary-hover);
}

/* Bouton toggle */
.theme-toggle {
  background: none;
  border: 2px solid var(--primary-accent);
  border-radius: 8px;
  color: var(--primary-accent);
  font-size: 1.1em;
  cursor: pointer;
  padding: 4px 10px;
  transition: background-color 0.2s ease, color 0.2s ease;
}
.theme-toggle:hover {
  background-color: var(--primary-accent);
  color: var(--bg-dark);
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
