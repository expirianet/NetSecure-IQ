<template>
  <div>
    <h2>Register</h2>
    <form @submit.prevent="register">
      <input v-model="email" placeholder="Email" required />
      <input v-model="password" type="password" placeholder="Password" required />
      <button type="submit">Register</button>
    </form>
    <p v-if="message">{{ message }}</p>
  </div>
</template>

<script>
export default {
  data() {
    return {
      email: '',
      password: '',
      message: ''
    }
  },
  methods: {
    async register() {
      try {
        const res = await fetch('http://localhost:8000/api/register', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ email: this.email, password: this.password })
        })

        if (!res.ok) {
          const errText = await res.text()
          throw new Error(errText)
        }

        this.message = 'Registration successful!'
      } catch (err) {
        this.message = 'Error: ' + err.message
      }
    }
  }
}
</script>
