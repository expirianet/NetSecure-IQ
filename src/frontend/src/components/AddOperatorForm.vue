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

      <button type="submit">Add Operator</button>
      <p v-if="message">{{ message }}</p>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import jwtDecode from 'jwt-decode'
import { API } from '@/utils/api.js'

const email = ref(''); const firstName = ref(''); const lastName = ref(''); const message = ref('')

const submitForm = async () => {
  try {
    const token = localStorage.getItem('token') || ''
    const decoded = token ? jwtDecode(token) : {}
    const organization_id = decoded?.organization_id

    const payload = {
      email: email.value, first_name: firstName.value, last_name: lastName.value, role: 'operator',
      ...(organization_id ? { organization_id: String(organization_id) } : {})
    }

    const res = await fetch(`${API}/api/users`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', ...(token ? { Authorization: 'Bearer ' + token } : {}) },
      body: JSON.stringify(payload)
    })
    const data = await res.json().catch(() => ({}))
    message.value = res.ok ? '✅ Operator created successfully!' : '❌ Failed: ' + (data.error || data.message)
  } catch (err) {
    message.value = '❌ Internal error: ' + err.message
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
