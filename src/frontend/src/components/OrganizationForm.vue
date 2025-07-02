<template>
  <div class="organization-form">
    <h2>Complete Your Organization Info</h2>
    <form @submit.prevent="submitForm">
      <div>
        <label for="name">Organization Name:</label>
        <input v-model="form.name" type="text" id="name" required />
      </div>

      <div>
        <label for="address">Address:</label>
        <input v-model="form.address" type="text" id="address" required />
      </div>

      <div>
        <label for="vat">VAT Number:</label>
        <input v-model="form.vat_number" type="text" id="vat" />
      </div>

      <div>
        <label for="contact_email">Contact Email:</label>
        <input v-model="form.contact_email" type="email" id="contact_email" />
      </div>

      <div>
        <label for="contact_phone">Contact Phone:</label>
        <input v-model="form.contact_phone" type="tel" id="contact_phone" />
      </div>

      <div class="buttons">
        <button type="submit">Submit</button>
        <button type="button" @click="goToDashboard">Go to Dashboard</button>
      </div>

      <p v-if="message" class="success">{{ message }}</p>
    </form>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const form = reactive({
  name: '',
  address: '',
  vat_number: '',
  contact_email: '',
  contact_phone: ''
})

const message = ref('')

const goToDashboard = () => {
  router.push('/dashboard')
}

const submitForm = () => {
  if (!form.name || !form.address) {
    message.value = "Please fill in all required fields."
    return
  }

  console.log("Submitted Organization Info:", { ...form })
  message.value = "Organization info submitted! Redirecting..."

  // Later replace with a POST call here
  setTimeout(() => {
    router.push('/dashboard')
  }, 1000)
}
</script>

<style scoped>
.organization-form {
  max-width: 500px;
  margin: auto;
}

.organization-form label {
  display: block;
  margin-top: 10px;
}

.organization-form input {
  width: 100%;
  padding: 8px;
  margin-top: 4px;
}

.buttons {
  display: flex;
  gap: 10px;
  margin-top: 20px;
}

button {
  padding: 10px 20px;
}

.success {
  margin-top: 10px;
  color: green;
}
</style>
