<template>
  <div class="add-user-page">
    <h2>Add New User</h2>

    <form @submit.prevent="submitForm">
      <div>
        <label>First Name:</label>
        <input v-model="firstName" type="text" required />
      </div>

      <div>
        <label>Last Name:</label>
        <input v-model="lastName" type="text" required />
      </div>

      <div>
        <label>Email:</label>
        <input v-model="email" type="email" required />
      </div>

      <div v-if="isAdmin">
        <label>Select Organization:</label>
        <select v-model="selectedOrg" required>
          <option v-for="org in organizations" :key="org.id" :value="org.id">
            {{ org.name }}
          </option>
        </select>
      </div>

      <button type="submit">Add User</button>
      <p v-if="message">{{ message }}</p>
    </form>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const email = ref('')
const firstName = ref('')
const lastName = ref('')
const organizations = ref([])
const selectedOrg = ref(null)
const message = ref('')

const role = localStorage.getItem('role')
const userOrgId = localStorage.getItem('organization_id')
const isAdmin = role === 'administrator'

onMounted(async () => {
  if (isAdmin) {
    const res = await fetch(
      `${process.env.VUE_APP_BACKEND_URL}/api/organizations`,
      {
        headers: {
          Authorization: 'Bearer ' + localStorage.getItem('token'),
        },
      }
    )
    const data = await res.json().catch(() => ({ organizations: [] }))
    organizations.value = data.organizations || []
  }
})

const submitForm = async () => {
  const payload = {
    email: email.value,
    first_name: firstName.value,
    last_name: lastName.value,
    organization_id: isAdmin ? selectedOrg.value : userOrgId,
  }

  const res = await fetch(
    `${process.env.VUE_APP_BACKEND_URL}/api/users`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: 'Bearer ' + localStorage.getItem('token'),
      },
      body: JSON.stringify(payload),
    }
  )
  const data = await res.json().catch(() => ({ error: 'Invalid response' }))
  message.value = res.ok
    ? '✅ User created and password sent via email!'
    : '❌ Failed: ' + (data.error || data.message)
}
</script>

<style scoped>
.add-user-page {
  max-width: 480px;
  margin: 2rem auto;
  padding: 2rem;
  background: var(--panel-grey);
  border-radius: 8px;
}
.add-user-page h2 {
  margin-bottom: 1rem;
  color: var(--primary-accent);
}
.add-user-page form div {
  margin-bottom: 0.75rem;
}
.add-user-page label {
  display: block;
  margin-bottom: 0.25rem;
  color: var(--text-secondary);
}
.add-user-page input,
.add-user-page select {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid var(--divider-grey);
  border-radius: 4px;
  background: var(--bg-dark);
  color: var(--text-primary);
}
.add-user-page button {
  margin-top: 1rem;
  background: var(--primary-accent);
  color: #fff;
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
.add-user-page p {
  margin-top: 0.5rem;
}
</style>
