<template>
  <div class="login-page">
    <!-- Canvas animé -->
    <div id="particles-js"></div>

    <!-- Formulaire de login -->
    <div class="login-wrapper">
      <div class="login-container">
        <div class="login-card">
          <h2 class="login-title">NetSecure-IQ</h2>
          <h3 class="login-subtitle">Login to your account</h3>

          <form @submit.prevent="login" class="login-form">
            <input
              v-model="email"
              type="email"
              placeholder="Email address"
              required
            />
            <input
              v-model="password"
              type="password"
              placeholder="Password"
              required
            />
            <button :disabled="loading" type="submit">
              {{ loading ? "Loading..." : "Login" }}
            </button>
          </form>

          <p class="login-footer">
            Don't have an account?
            <router-link to="/register">Register</router-link>
          </p>

          <p
            v-if="message"
            :class="['login-message', successMessage ? 'success' : 'error']"
          >
            {{ message }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'

const email = ref('')
const password = ref('')
const message = ref('')
const successMessage = ref(false)
const loading = ref(false)
const router = useRouter()

/**
 * Initialise ou recharge particles.js en fonction du thème actuel.
 */
function renderParticles() {
  const dark = document.documentElement.getAttribute('data-theme') === 'dark'
  // supprime ancien canvas
  const old = document.querySelector('#particles-js > canvas')
  if (old) old.remove()

  // (re)lance particlesJS
  window.particlesJS('particles-js', {
    particles: {
      number: { value: 80, density: { enable: true, value_area: 800 } },
      color: { value: dark ? '#ffffff' : '#555555' },
      shape: { type: 'circle' },
      opacity: { value: dark ? 0.5 : 0.5 },
      size: { value: 3, random: true },
      line_linked: {
        enable: true,
        distance: 150,
        color: dark ? '#ffffff' : '#888888',
        opacity: dark ? 0.4 : 0.4,
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

// Function to ensure particles are loaded after DOM and theme are ready
async function initializeParticles() {
  // Ensure the particles container exists
  if (!document.getElementById('particles-js')) {
    await new Promise(resolve => setTimeout(resolve, 50));
    return initializeParticles();
  }
  
  // Load particles script if not already loaded
  if (!window.particlesJS) {
    await new Promise((resolve) => {
      const script = document.createElement('script');
      script.src = '/particles/particles.min.js';
      script.onload = resolve;
      document.body.appendChild(script);
    });
  }
  
  // Ensure theme is applied
  await nextTick();
  
  renderParticles();
  
  // Set up theme change observer
  const obs = new MutationObserver((mutations) => {
    for (const m of mutations) {
      if (m.attributeName === 'data-theme') {
        renderParticles();
      }
    }
  });

  // Observe theme changes on document.documentElement
  obs.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['data-theme']
  });
}

onMounted(() => {
  initializeParticles();
})

const login = async () => {
  loading.value = true
  try {
    const res = await fetch(
      `${process.env.VUE_APP_BACKEND_URL}/api/login`,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: email.value,
          password: password.value
        })
      }
    )
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || data.message || 'Login failed')

    message.value = 'Login successful! Redirecting...'
    successMessage.value = true

    // Stockage token & infos
    localStorage.setItem('token', data.token)
    localStorage.setItem('user_id', data.user_id)
    localStorage.setItem('role', data.role?.toLowerCase() || '')
    localStorage.setItem('organization_id', data.organization_id || '')

    // Choix de la redirection
    let redirectTo = '/dashboard'
    const role = data.role?.toLowerCase()
    if (role === 'user') redirectTo = '/routertable'
    else if (role === 'operator' && !data.organization_id)
      redirectTo = '/organizationForm'

    setTimeout(() => router.push(redirectTo), 500)
  } catch (err) {
    message.value = 'Error: ' + err.message
    successMessage.value = false
  } finally {
    loading.value = false
  }
}
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
.login-page {
  position: relative;
  min-height: 100vh;
  overflow: hidden;
}

/* Canvas particles */
#particles-js {
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
[data-theme='light'] #particles-js {
  background-color: #E0E0E0;
}

/* Wrapper du formulaire */
.login-wrapper {
  position: relative;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px;
  min-height: 100vh;
}

/* Conteneur */
.login-container {
  width: 100%;
  max-width: 420px;
}

/* Carte de login */
.login-card {
  background-color: var(--panel-grey);
  border-radius: 12px;
  padding: 32px;
  box-shadow: 0 0 40px rgba(0, 194, 194, 0.05);
  box-sizing: border-box;
}

/* Titres */
.login-title {
  text-align: center;
  font-size: 20px;
  font-weight: 600;
  color: var(--primary-accent);
  margin-bottom: 8px;
}
.login-subtitle {
  text-align: center;
  font-size: 16px;
  margin-bottom: 24px;
}

/* Form */
.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* Inputs & bouton full-width */
.login-form input,
.login-form button {
  width: 100%;
  box-sizing: border-box;
}

/* Input style */
.login-form input {
  background-color: var(--panel-grey);
  border: 1px solid var(--divider-grey);
  border-radius: 6px;
  padding: 12px 14px;
  font-size: 14px;
  color: var(--text-primary);
  transition: border-color 0.2s;
}
.login-form input::placeholder {
  color: var(--text-secondary);
}
.login-form input:focus {
  outline: none;
  border-color: var(--primary-accent);
  background-color: var(--bg-dark);
}

/* Bouton */
.login-form button {
  background-color: var(--primary-accent);
  color: var(--bg-dark);
  border: none;
  border-radius: 6px;
  font-weight: 600;
  font-size: 14px;
  padding: 12px 20px;
  cursor: pointer;
  transition: background-color 0.2s;
}
.login-form button:hover {
  background-color: var(--primary-hover);
  color: var(--text-primary);
}
.login-form button:disabled {
  background-color: #2f333d;
  color: #666;
  cursor: not-allowed;
}

/* Footer */
.login-footer {
  text-align: center;
  font-size: 13px;
  margin-top: 16px;
  color: var(--text-secondary);
}
.login-footer a {
  color: var(--primary-accent);
  margin-left: 4px;
  text-decoration: none;
}
.login-footer a:hover {
  color: var(--primary-hover);
}

/* Messages */
.login-message {
  margin-top: 16px;
  font-size: 14px;
  padding: 10px 12px;
  border-radius: 6px;
  text-align: center;
}
.login-message.success {
  background-color: rgba(34, 197, 94, 0.1);
  color: var(--success);
}
.login-message.error {
  background-color: rgba(239, 68, 68, 0.1);
  color: var(--danger);
}
</style>
