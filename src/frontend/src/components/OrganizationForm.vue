<template>
  <div class="organization-form">
    <h2>Register Your Organization</h2>
    <form @submit.prevent="submitForm">

      <!-- Section: Organization -->
      <fieldset>
        <legend>Organization Info</legend>
        <input v-model="form.name" placeholder="Organization Name" required />
        <input v-model="form.vat_number" placeholder="VAT Number or Fiscal Code" required />
        <input v-model="form.address" placeholder="Address" required />
        <input v-model="form.state" placeholder="State" required />
        <input v-model="form.city" placeholder="City" required />
        <input v-model="form.zip_code" placeholder="ZIP Code" required />
        <input v-model="form.email" placeholder="Email" type="email" required />
        <input v-model="form.pec" placeholder="PEC Email" />
        <input v-model="form.sdi" placeholder="SDI Code" />
        <input v-model="form.phone" placeholder="Phone Number" required />
      </fieldset>

      <!-- Section: Company Manager -->
      <fieldset>
        <legend>Company Manager</legend>
        <input v-model="form.manager_name" placeholder="Name and Surname" required />
        <input v-model="form.manager_email" placeholder="Email" type="email" required />
        <input v-model="form.manager_phone" placeholder="Phone Number" required />
      </fieldset>

      <!-- Section: Technical Manager -->
      <fieldset>
        <legend>Technical Manager</legend>
        <input v-model="form.technical_name" placeholder="Name and Surname" required />
        <input v-model="form.technical_email" placeholder="Email" type="email" required />
        <input v-model="form.technical_phone" placeholder="Phone Number" required />
      </fieldset>

      <!-- Section: GDPR / NIS2 -->
      <fieldset>
        <legend>Privacy – GDPR – NIS2</legend>

        <label>Data Controller</label>
        <input v-model="form.controller_name" placeholder="Name and Surname" required />
        <input v-model="form.controller_email" placeholder="Email" type="email" required />
        <input v-model="form.controller_phone" placeholder="Phone Number" required />

        <label>Data Processor</label>
        <input v-model="form.processor_name" placeholder="Name and Surname" required />
        <input v-model="form.processor_email" placeholder="Email" type="email" required />
        <input v-model="form.processor_phone" placeholder="Phone Number" required />
      </fieldset>

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
const message = ref('')

const form = reactive({
  // Org
  name: '', vat_number: '', address: '', state: '', city: '', zip_code: '',
  email: '', pec: '', sdi: '', phone: '',
  // Manager
  manager_name: '', manager_email: '', manager_phone: '',
  // Technical
  technical_name: '', technical_email: '', technical_phone: '',
  // GDPR
  controller_name: '', controller_email: '', controller_phone: '',
  processor_name: '', processor_email: '', processor_phone: ''
})

const goToDashboard = () => {
  router.push('/dashboard')
}

const submitForm = async () => {
  if (!form.name || !form.vat_number || !form.address || !form.email || !form.manager_name) {
    message.value = "Please fill in all required fields."
    return
  }

  try {
    const response = await fetch(`${process.env.VUE_APP_BACKEND_URL}/api/complete-organization`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        ...form,
        user_id: localStorage.getItem("user_id")
      }),
    })

    let responseText = await response.text()
    let data
    try {
      data = JSON.parse(responseText)
    } catch {
      throw new Error(responseText)
    }

    if (!response.ok) throw new Error(data.error || data.message)

    message.value = "Organization info submitted! Redirecting..."
    setTimeout(() => {
      router.push('/dashboard')
    }, 1000)

  } catch (err) {
    message.value = "Submission failed: " + err.message
  }
}
</script>

<style scoped>
.organization-form {
  max-width: 800px;
  margin: auto;
  padding: 20px;
}

fieldset {
  margin-top: 20px;
  padding: 10px 20px;
  border: 1px solid #ccc;
  border-radius: 8px;
}

legend {
  font-weight: bold;
  padding: 0 10px;
}

input {
  width: 100%;
  padding: 8px;
  margin-top: 6px;
  margin-bottom: 10px;
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
