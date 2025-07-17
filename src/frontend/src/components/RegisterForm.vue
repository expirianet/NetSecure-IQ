<template>
  <div>
    <!-- Fond animé derrière tout -->
    <div id="particles-js"></div>

    <!-- Contenu du register -->
    <div class="register-wrapper">
      <div class="register-container">
        <div class="register-card">
          <h2 class="register-title">NetSecure-IQ</h2>
          <h3 class="register-subtitle">Create your account</h3>

          <form @submit.prevent="register" class="register-form">
            <input
              v-model="email"
              type="email"
              placeholder="Email address"
              required
            />
            <button type="submit" :disabled="loading">
              {{ loading ? "Registering..." : "Register" }}
            </button>
          </form>

          <p v-if="message" class="register-message success">{{ message }}</p>
          <p v-if="error" class="register-message error">{{ error }}</p>

          <p class="register-footer">
            Already have an account?
            <router-link to="/login">Login</router-link>
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
const message = ref('')
const error = ref('')
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

const register = async () => {
  message.value = ''
  error.value = ''
  loading.value = true

  try {
    const res = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: email.value })
    })

    const data = await res.json()
    if (!res.ok) throw new Error(data.error || data.message || 'Registration failed')

    message.value = data.message || 'Registration successful. Check your email.'
    email.value = ''

    setTimeout(() => router.push('/login'), 3000)
  } catch (err) {
    error.value = err.message
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

/* Fond animé */
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

/* Couche par-dessus */
.register-wrapper {
  position: relative;
  z-index: 10;
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px;
}

.register-wrapper,
.register-wrapper * {
  color: var(--text-primary);
  font-family: 'Inter', sans-serif;
}

.register-container {
  width: 100%;
  max-width: 420px;
}

.register-card {
  background-color: var(--panel-grey);
  border-radius: 12px;
  padding: 32px;
  box-shadow: 0 0 40px rgba(0, 194, 194, 0.05);
  box-sizing: border-box;
}

.register-title {
  text-align: center;
  font-size: 20px;
  font-weight: 600;
  color: var(--primary-accent);
  margin-bottom: 8px;
}

.register-subtitle {
  text-align: center;
  font-size: 16px;
  margin-bottom: 24px;
}

.register-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.register-form input,
.register-form button {
  width: 100%;
  box-sizing: border-box;
}

input {
  background-color: var(--panel-grey);
  border: 1px solid var(--divider-grey);
  border-radius: 6px;
  padding: 12px 14px;
  font-size: 14px;
  color: var(--text-primary);
  transition: border-color 0.2s ease;
}

input::placeholder {
  color: var(--text-secondary);
}

input:focus {
  outline: none;
  border-color: var(--primary-accent);
  background-color: var(--bg-dark);
}

button {
  background-color: var(--primary-accent);
  color: var(--bg-dark);
  border: none;
  border-radius: 6px;
  font-weight: 600;
  font-size: 14px;
  padding: 12px 20px;
  cursor: pointer;
  transition: all 0.2s ease;
}

button:hover {
  background-color: var(--primary-hover);
  color: var(--text-primary);
}

button:disabled {
  background-color: #2f333d;
  color: #666;
  cursor: not-allowed;
}

.register-message {
  margin-top: 16px;
  font-size: 14px;
  padding: 10px 12px;
  border-radius: 6px;
  text-align: center;
}

.register-message.success {
  background-color: rgba(34, 197, 94, 0.1);
  color: var(--success);
}

.register-message.error {
  background-color: rgba(239, 68, 68, 0.1);
  color: var(--danger);
}

.register-footer {
  text-align: center;
  font-size: 13px;
  margin-top: 16px;
  color: var(--text-secondary);
}

.register-footer a {
  color: var(--primary-accent);
  text-decoration: none;
  margin-left: 4px;
}

.register-footer a:hover {
  color: var(--primary-hover);
}
</style>
