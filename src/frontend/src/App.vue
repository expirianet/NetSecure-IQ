<template>
  <div id="app">
    <nav>
      <router-link to="/login">Login</router-link> |
      <router-link to="/register">Register</router-link>
      <button class="theme-toggle" @click="toggleTheme" aria-label="Toggle theme">
        <span v-if="isDark">☀️</span>
        <span v-else>🌙</span>
      </button>
    </nav>
    <router-view />
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'

// État du thème
const theme = ref(localStorage.getItem('theme') || 'dark')
const isDark = computed(() => theme.value === 'dark')

// Fonction pour basculer le thème
const toggleTheme = () => {
  theme.value = theme.value === 'dark' ? 'light' : 'dark'
  localStorage.setItem('theme', theme.value)
}

// Observer les changements de thème
watch(theme, (newTheme) => {
  document.documentElement.setAttribute('data-theme', newTheme)
})

// Initialiser le thème
onMounted(() => {
  document.documentElement.setAttribute('data-theme', theme.value)
})
</script>

<style>
:root {
  --bg-dark: #0E111A;
  --panel-grey: #1A1D26;
  --divider-grey: #2A2D36;
  --text-primary: #F5F7FA;
  --text-secondary: #9CA3AF;
  --primary-accent: #00C2C2;
  --primary-hover: #00A7A7;
  --danger: #EF4444;
  --success: #22C55E;
}

/* **Cette règle s’appliquera enfin** */
body {
  margin: 0;
  background-color: var(--bg-dark);
  color: var(--text-primary);
  font-family: 'Inter', sans-serif;
}

/* Le nav reste au-dessus */
nav {
  position: relative;
  z-index: 20;
  padding: 16px;
  background-color: var(--panel-grey);
  border-bottom: 1px solid var(--divider-grey);
}

nav a {
  color: var(--primary-accent);
}

nav a:hover {
  color: var(--primary-hover);
}

/* Et par défaut, tous tes titres, paragraphes, liens... */
h1,h2,h3,h4,h5,h6,p,a,span,label {
  color: var(--text-primary);
}

/* Thème clair */
[data-theme='light'] {
  --bg-dark: #ffffff;
  --panel-grey: #f9fafb;
  --divider-grey: #e5e7eb;
  --text-primary: #1f2937;
  --text-secondary: #6b7280;
  --primary-accent: #00a8a8;
  --primary-hover: #008a8a;
  --danger: #dc2626;
  --success: #22c55e;
}

/* Bouton de toggle du thème */
.theme-toggle {
  margin-left: 16px;
  background: none;
  border: none;
  color: var(--text-primary);
  font-size: 1.2em;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.theme-toggle:hover {
  background-color: var(--divider-grey);
}

.theme-toggle:focus {
  outline: none;
  box-shadow: 0 0 0 2px var(--primary-accent);
}
</style>



nav {
  padding: 16px;
  background-color: var(--panel-grey);
  border-bottom: 1px solid var(--divider-grey);
}

nav a {
  margin: 0 8px;
  color: var(--primary-accent);
  text-decoration: none;
}

nav a:hover {
  color: var(--primary-hover);
}

#particles-js {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: -1;
  background-color: #0E111A;
}

