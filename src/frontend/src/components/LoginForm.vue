<template>
  <div>
    <!-- Fond animé derrière tout -->
    <div id="particles-js"></div>

    <!-- Contenu du login -->
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
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const email = ref('')
const password = ref('')
const message = ref('')
const successMessage = ref(false)
const loading = ref(false)
const router = useRouter()

onMounted(() => {
  const script = document.createElement('script')
  script.src = '/particles/particles.min.js'
  script.onload = () => {
    if (window.particlesJS) {
      window.particlesJS('particles-js', {
        particles: {
          number: { value: 80, density: { enable: true, value_area: 800 } },
          color: { value: '#ffffff' },
          shape: { type: 'circle' },
          opacity: { value: 0.5 },
          size: { value: 3, random: true },
          line_linked: {
            enable: true,
            distance: 150,
            color: '#ffffff',
            opacity: 0.4,
            width: 1
          },
          move: {
            enable: true,
            speed: 6,
            direction: 'none',
            out_mode: 'bounce'
          }
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
  document.body.appendChild(script)
})

const login = async () => {
  loading.value = true
  try {
    const res = await fetch(
      `${import.meta.env.VITE_BACKEND_URL}/api/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: email.value, password: password.value })
      }
    )
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || data.message || 'Login failed')

    message.value = 'Login successful! Redirecting...'
    successMessage.value = true

    localStorage.setItem('token', data.token)
    localStorage.setItem('user_id', data.user_id)
    localStorage.setItem('role', data.role?.toLowerCase() || '')
    localStorage.setItem('organization_id', data.organization_id || '')

    let redirectTo = '/dashboard'
    if (data.role?.toLowerCase() === 'user') redirectTo = '/routertable'
    else if (
      data.role?.toLowerCase() === 'operator' &&
      !data.organization_id
    ) {
      redirectTo = '/organizationForm'
    }

    setTimeout(() => {
      router.push(redirectTo)
    }, 500)
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

/* Fond animé derrière tout */
#particles-js {
  position: fixed;
  top: 0;
  left: 0;
  z-index: 0;
  width: 100vw;
  height: 100vh;
  background-color: var(--bg-dark);
  pointer-events: none;
}

/* Contenu du formulaire devant */
.login-wrapper {
  position: relative;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
}

.login-container {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 400px;
  padding: 24px;
}

.login-card {
  background-color: var(--panel-grey);
  border-radius: 12px;
  padding: 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.login-title {
  font-size: 2.5em;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 16px;
  text-align: center;
}

.login-subtitle {
  font-size: 1.1em;
  color: var(--text-secondary);
  margin-bottom: 32px;
  text-align: center;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.login-form input {
  width: 100%;
  padding: 12px;
  border: none;
  border-radius: 8px;
  background-color: var(--panel-grey);
  color: var(--text-primary);
  font-size: 1em;
  box-sizing: border-box;
}

.login-form input:focus {
  outline: none;
  border: 1px solid var(--primary-accent);
  box-shadow: 0 0 0 2px var(--primary-accent);
}

.login-form button {
  width: 100%;
  padding: 12px;
  border: none;
  border-radius: 8px;
  background-color: var(--primary-accent);
  color: white;
  font-size: 1em;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
  box-sizing: border-box;
}

.login-form button:hover {
  background-color: var(--primary-hover);
}

.login-form button:disabled {
  background-color: var(--divider-grey);
  cursor: not-allowed;
}

.login-footer {
  text-align: center;
  margin-top: 24px;
  color: var(--text-secondary);
}

.login-footer a {
  color: var(--primary-accent);
  text-decoration: none;
}

.login-footer a:hover {
  color: var(--primary-hover);
}

.login-message {
  margin-top: 16px;
  padding: 12px;
  border-radius: 8px;
  text-align: center;
}

.login-message.success {
  background-color: var(--success);
  color: white;
}

.login-message.error {
  background-color: var(--danger);
  color: white;
}

</style>

