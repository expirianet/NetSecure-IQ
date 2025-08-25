<!-- src/frontend/src/components/TopNavigationOperator.vue -->
<template>
  <nav class="navbar">
    <!-- Left: brand + links -->
    <div class="nav-left">
      <router-link to="/" class="brand">NetSecure-IQ</router-link>
      <router-link to="/" class="nav-link" :class="{ active: $route.path === '/' }">Home</router-link>

      <!-- Not authenticated -->
      <template v-if="!isAuthenticated">
        <router-link to="/login" class="nav-link" :class="{ active: $route.path === '/login' }">Login</router-link>
        <router-link to="/register" class="nav-link" :class="{ active: $route.path === '/register' }">Register</router-link>
      </template>

      <!-- Authenticated (OPERATOR) -->
      <template v-else>
        <router-link
          to="/dashboard-operator"
          class="nav-link"
          :class="{ active: $route.path === '/dashboard-operator' }"
        >Dashboard</router-link>

        <router-link
          to="/routertable"
          class="nav-link"
          :class="{ active: $route.path === '/routertable' }"
        >RouterTable</router-link>

        <router-link
          to="/organization"
          class="nav-link"
          :class="{ active: $route.path === '/organization' }"
        >Organization</router-link>

        <!-- Add User only if attached to an organization -->
        <router-link
          v-if="hasOrganization"
          to="/adduser"
          class="nav-link"
          :class="{ active: $route.path === '/adduser' }"
        >Add User</router-link>

        <!-- Agents & Pre-registration visible only if org -->
        <router-link
          v-if="hasOrganization"
          to="/agents"
          class="nav-link"
          :class="{ active: $route.path === '/agents' }"
        >Agents</router-link>

        <router-link
          v-if="hasOrganization"
          to="/agents/register"
          class="nav-link"
          :class="{ active: $route.path === '/agents/register' }"
        >Agent Pre-Registration</router-link>

        <button @click="logout" class="nav-link">Logout</button>
      </template>
    </div>

    <!-- Right: status + theme toggle -->
    <div class="nav-right">
      <div v-if="isAuthenticated" class="role-badge" title="You are signed in as operator">
        <span class="dot online" aria-hidden="true"></span>
        <span class="role-text">Operator signed in</span>
        <span v-if="hasOrganization" class="org-hint">org: {{ organizationId }}</span>
      </div>

      <button
        class="theme-toggle"
        @click="toggleTheme"
        :aria-label="isDark ? 'Switch to light mode' : 'Switch to dark mode'"
      >
        <span v-if="isDark">☀️</span>
        <span v-else>🌙</span>
      </button>
    </div>
  </nav>
</template>

<script setup>
import { useAuth } from '@/composables/useAuth.js'
import { ref, computed, onMounted, onUnmounted } from 'vue'

const { isAuthenticated, logout } = useAuth()

/** Organization: reactive via localStorage + events */
const organizationId = ref(localStorage.getItem('organization_id') || '')
const hasOrganization = computed(() => !!(organizationId.value && String(organizationId.value).trim()))
function syncOrgId() {
  organizationId.value = localStorage.getItem('organization_id') || ''
}

/** Theme */
const theme = ref(localStorage.getItem('theme') || 'dark')
const isDark = computed(() => theme.value === 'dark')
function toggleTheme() {
  theme.value = isDark.value ? 'light' : 'dark'
  document.documentElement.setAttribute('data-theme', theme.value)
  localStorage.setItem('theme', theme.value)
}

onMounted(() => {
  document.documentElement.setAttribute('data-theme', theme.value)
  window.addEventListener('storage', syncOrgId)
  window.addEventListener('auth-changed', syncOrgId)
})
onUnmounted(() => {
  window.removeEventListener('storage', syncOrgId)
  window.removeEventListener('auth-changed', syncOrgId)
})
</script>

<style scoped>
:root {
  --bg-dark: #0e111a;
  --panel-grey: #1a1d26;
  --divider-grey: #2a2d36;
  --text-primary: #f5f7fa;
  --text-secondary: #9ca3af;
  --primary-accent: #00b4aa;
  --primary-hover: #008a8a;
  --success: #22c55e;
}

/* Navbar */
.navbar {
  position: fixed; top: 0; left: 0; right: 0;
  z-index: 1000;
  display: flex; align-items: center; justify-content: space-between;
  padding: 0 32px; height: 64px;
  background-color: var(--bg-dark);
  border-bottom: 1px solid var(--divider-grey);
  transition: background-color 0.3s ease;
}
[data-theme='light'] .navbar { background-color: #ffffff; box-shadow: 0 2px 10px rgba(0,0,0,0.05); }

/* Left */
.nav-left { display: flex; align-items: center; gap: 24px; }
.brand { font-weight: bold; font-size: 18px; color: var(--primary-accent); text-decoration: none; }
.brand:hover { color: var(--primary-hover); }

.nav-link {
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 14px;
  padding: 8px 12px;
  border-radius: 4px;
  transition: all 0.2s ease;
  background: transparent;
  border: none;
  cursor: pointer;
}
.nav-link:hover,
.nav-link.active {
  color: var(--primary-accent);
  background-color: rgba(0,194,194,0.1);
}
.nav-link.active { font-weight: 600; }

/* Right */
.nav-right { margin-left: auto; display:flex; align-items:center; gap:12px; }

.role-badge {
  display: inline-flex; align-items: center; gap: 8px;
  padding: 6px 10px; border-radius: 999px;
  background: rgba(34, 197, 94, 0.12);
  border: 1px solid rgba(34, 197, 94, 0.25);
  color: var(--text-primary); font-size: 12px;
}
.dot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; }
.dot.online { background: var(--success); }
.role-text { font-weight: 600; }
.org-hint { opacity: .75; font-family: ui-monospace, SFMono-Regular, Menlo, monospace; }

/* Theme toggle */
.theme-toggle {
  background: none; border: none;
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 1.25rem; padding: 8px;
  border-radius: 50%;
  transition: background-color 0.2s ease;
}
.theme-toggle:hover {
  background-color: var(--divider-grey);
  color: var(--text-primary);
}
[data-theme='light'] .theme-toggle:hover { background-color: rgba(0,0,0,0.05); }
</style>
