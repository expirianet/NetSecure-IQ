<template>
  <div class="home-page">
    <!-- Canvas animé -->
    <div id="home-particles"></div>

    <!-- Contenu principal -->
    <div class="home-wrapper">
      <div class="home-container">
        <div class="home-card">
          <h1 class="welcome-title">Welcome to NetSecure-IQ</h1>
          <p class="welcome-subtitle">Your network security management solution</p>
          <router-link to="/login" class="login-button">Get Started</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onBeforeUnmount } from 'vue';
import TopNavigation from '@/components/navigation/TopNavigation.vue';
import {
  ensurePJSDom, loadParticlesScript, defaultConfig,
  safeRender, observeTheme, destroyForId, themeIsDark
} from '@/appCore.js'

const ID = 'home-particles'
let stopObs = () => {}

function render() {
  return safeRender(ID, defaultConfig(themeIsDark()))
}

onMounted(async () => {
  try { await loadParticlesScript() } catch (e) {
    console.debug('[particles] load failed (non-blocking)', e)
  }
  ensurePJSDom()
  render()
  stopObs = observeTheme(ID, render)
})

onBeforeUnmount(() => {
  stopObs?.()
  ensurePJSDom()
  destroyForId(ID)
})
</script>

<style scoped>
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

/* Page entière */
.home-page {
  position: relative;
  min-height: 100vh;
  overflow: hidden;
}

/* Canvas particles */
#home-particles {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 0;
  background-color: var(--bg-dark);
  transition: background-color 0.3s ease;
}

/* override light mode */
[data-theme='light'] #home-particles { background-color: #f6f8fb; }


/* Wrapper du contenu */
.home-wrapper {
  position: relative;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 32px;
  text-align: center;
}

/* Conteneur */
.home-container {
  width: 100%;
  max-width: 800px;
}

/* Carte de bienvenue */
.home-card {
  background-color: var(--panel-grey);
  border-radius: 12px;
  padding: 48px 32px;
  box-shadow: 0 0 40px rgba(0, 194, 194, 0.05);
  box-sizing: border-box;
}

/* Titres */
.welcome-title {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--primary-accent);
  margin-bottom: 16px;
}

.welcome-subtitle {
  font-size: 1.25rem;
  color: var(--text-secondary);
  margin-bottom: 32px;
}

/* Bouton */
.login-button {
  display: inline-block;
  background-color: var(--primary-accent);
  color: var(--bg-dark);
  text-decoration: none;
  font-weight: 600;
  font-size: 1rem;
  padding: 12px 32px;
  border-radius: 6px;
  transition: background-color 0.2s, color 0.2s;
}

.login-button:hover {
  background-color: var(--primary-hover);
  color: var(--text-primary);
}

/* Responsive */
@media (max-width: 768px) {
  .welcome-title {
    font-size: 2rem;
  }
  
  .welcome-subtitle {
    font-size: 1.1rem;
  }
  
  .home-card {
    padding: 32px 24px;
  }
}
</style>
