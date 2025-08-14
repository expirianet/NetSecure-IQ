<template>
  <div id="app">
    <!-- Barre de navigation (switch auto selon le rôle) -->
    <component :is="navComponent" />

    <!-- Contenu principal -->
    <main class="main-content">
      <router-view />
    </main>
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

/* ---------- Lifecycle ---------- */
onMounted(async () => {
  // Appliquer le thème au démarrage
  document.documentElement.setAttribute('data-theme', theme.value)

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
/* Override pour le mode clair (v2) */
[data-theme='light'] {
  /* surfaces */
  --bg-dark: #f6f8fb;            /* fond page */
  --panel-grey: #ffffff;         /* cartes/panels */
  --divider-grey: #e5e7eb;       /* traits/borders */

  /* textes */
  --text-primary: #0f172a;       /* slate-900 */
  --text-secondary: #475569;     /* slate-600 */

  /* accent */
  --primary-accent: #0ea5a5;     /* teal-600 */
  --primary-hover: #0b8f8f;      /* teal-700 */

  --danger: #dc2626;
  --success: #16a34a;

  /* nouveaux tokens */
  --panel-border: rgba(2, 6, 23, 0.08);
  --panel-shadow: 0 12px 32px rgba(2, 6, 23, 0.08), 0 1px 0 rgba(2, 6, 23, 0.04);
  --muted-surface: rgba(2, 6, 23, 0.04);
}

/* Améliorations globales light sans changer la forme */
[data-theme='light'] .login-card,
[data-theme='light'] .register-card,
[data-theme='light'] .adduser-card,
[data-theme='light'] .opdash-shell,
[data-theme='light'] .card,
[data-theme='light'] .tile,
[data-theme='light'] .form-section,
[data-theme='light'] .table-wrapper,
[data-theme='light'] .org-card {
  border-color: var(--panel-border) !important;
  box-shadow: var(--panel-shadow) !important;
  background: var(--panel-grey);
}

/* Surfaces “muted” (tiles, sections…) */
[data-theme='light'] .tile {
  background: var(--muted-surface);
}
[data-theme='light'] .form-section {
  background: rgba(2, 6, 23, 0.03);
}

/* Focus ring propre et accessible */
:focus-visible {
  outline: 2px solid var(--primary-accent);
  outline-offset: 2px;
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
</style>
