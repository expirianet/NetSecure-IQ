<template>
  <div id="app">
    <!-- Barre de navigation personnalisée -->
    <TopNavigation />
    
    <!-- Contenu principal -->
    <main class="main-content">
      <router-view />
    </main>

    <!-- Canvas pour particles.js -->
    <div id="particles-js"></div>
    <!-- Le bouton de basculement de thème a été supprimé -->
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import TopNavigation from '@/components/TopNavigation.vue'

// -- État du thème
const theme = ref(localStorage.getItem('theme') || 'dark')

// -- Mettre à jour l'attribut data-theme sur <html>
watch(theme, (newTheme) => {
  document.documentElement.setAttribute('data-theme', newTheme)
})

// -- Fonction d'initialisation de particles.js
function initParticles() {
  const dark = theme.value === 'dark'
  if (!window.particlesJS) return

  // Destruction de l'ancien canvas
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
      const observer = new MutationObserver(() => initParticles())
      observer.observe(document.documentElement, {
        attributes: true,
        attributeFilter: ['data-theme']
      })
    }
    document.body.appendChild(script)
  } else {
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

/* Reset des marges et paddings */
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
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
