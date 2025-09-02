
<template>
  <div class="add-operator-page">
    <h2>Add New Operator</h2>

    <form @submit.prevent="submitForm">
      <div>
        <label>Email:</label>
        <input v-model="email" type="email" required />
      </div>

      <div>
        <label>First Name:</label>
        <input v-model="firstName" type="text" required />
      </div>

      <div>
        <label>Last Name:</label>
        <input v-model="lastName" type="text" required />
      </div>

      <button type="submit" :disabled="loading">Add Operator</button>
      <p v-if="message">{{ message }}</p>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { API } from '@/appCore.js'

const email = ref('')
const firstName = ref('')
const lastName = ref('')
const message = ref('')
const loading = ref(false)

const submitForm = async () => {
  loading.value = true
  message.value = ''
  try {
    const token = localStorage.getItem('token') || ''
    // Le back ne met pas organization_id dans le JWT.
    // On le récupère de la réponse /api/login sauvegardée en localStorage.
    const orgId = localStorage.getItem('organization_id') || ''

    const payload = {
      email: email.value,
      first_name: firstName.value,
      last_name: lastName.value,
      role: 'operator',
      ...(orgId ? { organization_id: String(orgId) } : {}) // optionnel si vide
    }

    const res = await fetch(`${API}/api/users`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: 'Bearer ' + token } : {})
      },
      body: JSON.stringify(payload)
    })
    const data = await res.json().catch(() => ({}))
    if (!res.ok) {
      throw new Error(data.error || data.message || 'Request failed')
    }
    message.value = '✅ Operator created successfully!'
    email.value = ''
    firstName.value = ''
    lastName.value = ''
  } catch (err) {
    message.value = '❌ Failed: ' + (err.message || 'Unknown error')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.add-operator-page {
  max-width: 480px;
  margin: 2rem auto;
  padding: 2rem;
  background: var(--panel-grey);
  border-radius: 8px;
}
.add-operator-page h2 {
  margin-bottom: 1rem;
  color: var(--primary-accent);
}
.add-operator-page form div {
  margin-bottom: 0.75rem;
}
.add-operator-page label {
  display: block;
  margin-bottom: 0.25rem;
  color: var(--text-secondary);
}
.add-operator-page input {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid var(--divider-grey);
  border-radius: 4px;
  background: var(--bg-dark);
  color: var(--text-primary);
}
.add-operator-page button {
  margin-top: 1rem;
  background: var(--primary-accent);
  color: #fff;
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
.add-operator-page p {
  margin-top: 0.5rem;
}
</style>

