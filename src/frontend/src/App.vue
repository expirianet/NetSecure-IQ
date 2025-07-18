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
  margin: 0 8px;
  color: var(--primary-accent);
  text-decoration: none;
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

  /* Styles spécifiques pour le mode clair */
  .login-wrapper,
  .register-wrapper {
    background: linear-gradient(135deg, #f0f9ff 0%, #e6f3ff 100%);
    min-height: 100vh;
    overflow: hidden;
    position: relative;
  }

  /* Effet de vague animé */
  .login-wrapper::before,
  .register-wrapper::before {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    width: 100%;
    height: 120px;
    background: linear-gradient(135deg, transparent 0%, #f0f9ff 100%);
    animation: wave 10s infinite linear;
  }

  @keyframes wave {
    0% { transform: translateX(0); }
    100% { transform: translateX(-20%); }
  }

  .login-card,
  .register-card {
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.03);
    border-radius: 24px;
    backdrop-filter: blur(15px);
    background: linear-gradient(145deg, rgba(255, 255, 255, 0.98) 0%, rgba(255, 255, 255, 0.95) 100%);
    border: 1px solid rgba(0, 0, 0, 0.03);
    padding: 40px;
    transform: translateY(0);
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .login-card:hover,
  .register-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 25px 50px rgba(0, 0, 0, 0.05);
  }

  .login-form input,
  .register-form input {
    border: 2px solid transparent;
    background: linear-gradient(to right, #f0f9ff, #e6f3ff);
    background-clip: padding-box;
    border-radius: 12px;
    transition: all 0.3s ease;
    padding: 14px 20px;
  }

  .login-form input:focus,
  .register-form input:focus {
    border-color: var(--primary-accent);
    box-shadow: 0 0 0 4px rgba(0, 168, 168, 0.15);
    transform: scale(1.02);
  }

  .login-form button,
  .register-form button {
    border-radius: 12px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    padding: 14px 30px;
    position: relative;
    overflow: hidden;
    background: linear-gradient(135deg, var(--primary-accent) 0%, #00c2c2 100%);
    box-shadow: 0 8px 20px rgba(0, 168, 168, 0.2);
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .login-form button::after,
  .register-form button::after {
    content: '';
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
    transition: all 0.6s ease;
  }

  .login-form button:hover::after,
  .register-form button:hover::after {
    left: 100%;
  }

  .login-form button:hover,
  .register-form button:hover {
    transform: translateY(-2px);
    box-shadow: 0 12px 30px rgba(0, 168, 168, 0.3);
  }

  .login-message,
  .register-message {
    border-radius: 16px;
    padding: 16px 24px;
    font-weight: 500;
    position: relative;
    overflow: hidden;
  }

  .login-message::before,
  .register-message::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: linear-gradient(45deg, transparent, rgba(255, 255, 255, 0.2), transparent);
    animation: shine 2s infinite;
  }

  @keyframes shine {
    0% { transform: translateX(-100%); }
    100% { transform: translateX(100%); }
  }

  .login-message.success,
  .register-message.success {
    background: linear-gradient(135deg, #defce8 0%, #c1f5d0 100%);
    color: #166534;
  }

  .login-message.error,
  .register-message.error {
    background: linear-gradient(135deg, #fee2e2 0%, #fecdd3 100%);
    color: #991b1b;
  }

  /* Transitions douces */
  * {
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }
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

#particles-js {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: -1;
  background-color: #0E111A;
}
</style>
