<template>
  <div class="login-container">
    <h1>Connexion</h1>
    <form @submit.prevent="onLogin">
      <div class="form-group">
        <label for="email">Email</label>
        <input id="email" v-model="email" type="email" required />
      </div>
      <div class="form-group">
        <label for="password">Mot de passe</label>
        <input id="password" v-model="password" type="password" required />
      </div>
      <button type="submit" :disabled="loading">Se connecter</button>
      <div v-if="error" class="error">{{ error }}</div>
    </form>
  </div>
</template>

<script setup>
defineOptions({ name: 'LoginPage' })
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'

const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')
const router = useRouter()
const auth = useAuthStore()

async function onLogin() {
  loading.value = true
  error.value = ''
  try {
    // Remplacer par un vrai appel API plus tard
    if (email.value === 'admin@example.com' && password.value === 'admin') {
      auth.isAuthenticated = true
      auth.user = { email: email.value, role_id: 1 } // Admin
      auth.token = 'fake-token'
      router.push('/')
    } else if (email.value === 'operator@example.com' && password.value === 'operator') {
      auth.isAuthenticated = true
      auth.user = { email: email.value, role_id: 2 } // Operator
      auth.token = 'fake-token'
      router.push('/')
    } else if (email.value === 'user@example.com' && password.value === 'user') {
      auth.isAuthenticated = true
      auth.user = { email: email.value, role_id: 3 } // User
      auth.token = 'fake-token'
      router.push('/')
    } else {
      throw new Error('Identifiants invalides')
    }
  } catch (e) {
    error.value = e.message || 'Erreur de connexion'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  max-width: 350px;
  margin: 4rem auto;
  background: #181b23;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 4px 16px rgba(0,0,0,0.25);
  color: #f5f7fa;
}
h1 {
  margin-bottom: 1.5rem;
  font-size: 1.6rem;
  text-align: center;
}
.form-group {
  margin-bottom: 1rem;
}
label {
  display: block;
  margin-bottom: 0.3rem;
}
input {
  width: 100%;
  padding: 0.6rem;
  border-radius: 4px;
  border: 1px solid #2a2d36;
  background: #232635;
  color: #f5f7fa;
}
button {
  width: 100%;
  padding: 0.7rem;
  background: #00b4aa;
  color: #fff;
  border: none;
  border-radius: 4px;
  font-weight: bold;
  cursor: pointer;
  margin-top: 1rem;
  transition: background 0.2s;
}
button:disabled {
  background: #555a;
  cursor: not-allowed;
}
.error {
  color: #f44336;
  font-weight: bold;
  margin-top: 1rem;
  text-align: center;
}
</style>
