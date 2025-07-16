<template>
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
</template>

<script>
export default {
  data() {
    return {
      email: '',
      message: '',
      error: '',
      loading: false
    }
  },
  methods: {
    async register() {
      this.message = ''
      this.error = ''
      this.loading = true

      try {
        const res = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/register`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ email: this.email })
        })

        const data = await res.json()

        if (!res.ok) throw new Error(data.error || data.message || 'Registration failed')

        this.message = data.message || 'Registration successful. Check your email.'
        this.email = ''

        setTimeout(() => {
          this.$router.push('/login')
        }, 3000)
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
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

.register-wrapper,
.register-wrapper * {
  color: var(--text-primary);
  font-family: 'Inter', sans-serif;
}

.register-wrapper {
  min-height: 100vh;
  background-color: var(--bg-dark);
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 32px;
}

.register-container {
  width: 100%;
  max-width: 420px;
}

.register-card {
  background-color: var(--panel-grey);
  padding: 32px;
  border-radius: 12px;
  box-shadow: 0 0 40px rgba(0, 194, 194, 0.05);
  display: flex;
  flex-direction: column;
  gap: 16px;
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
  font-weight: 600;
  font-size: 14px;
  padding: 12px 20px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

button:hover {
  background-color: var(--primary-hover);
  color: var(--text-primary);
}

button:disabled {
  background-color: #2F333D;
  color: #666;
  cursor: not-allowed;
}

.register-message {
  text-align: center;
  font-size: 14px;
  padding: 10px 12px;
  border-radius: 6px;
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
