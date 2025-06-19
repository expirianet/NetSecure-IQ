<template>
  <div class="login">
    <div class="login-box">
      <h3>Register your account</h3>
      <form @submit.prevent="register">
        <input
          v-model="email"
          type="email"
          placeholder="Email"
          required
        />
        <input
          v-model="password"
          type="password"
          placeholder="Password"
          required
        />
        <button :disabled="loading" type="submit">
          {{ loading ? "Loading..." : "Register" }}
        </button>
      </form>



      <div class="link">
        <span>Have already an account ?</span>
        <router-link to="/login">Login</router-link>
      </div>

      <p v-if="message" :class="{ success: successMessage, error: !successMessage }">
        {{ message }}
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue"

const email = ref("")
const password = ref("")
const message = ref("")
const successMessage = ref(false)
const loading = ref(false)

const register = async () => {
  loading.value = true
  try {
    const res = await fetch("http://localhost:8000/api/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email: email.value, password: password.value }),
    })

    if (!res.ok) {
      const errText = await res.text()
      throw new Error(errText)
    }

    message.value = "Registration successful!"
    successMessage.value = true
  } catch (err) {
    message.value = "Erreur : " + err.message
    successMessage.value = false
  } finally {
    loading.value = false
  }
}
</script>

