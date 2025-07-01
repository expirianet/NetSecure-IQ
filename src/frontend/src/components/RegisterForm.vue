<template>
  <div class="register-form">
    <h2>Register</h2>
    <form @submit.prevent="register">
      <div>
        <label for="email">Email:</label>
        <input v-model="email" type="email" required />
      </div>

      <button type="submit">Register</button>
    </form>

    <p v-if="message">{{ message }}</p>
    <p v-if="error" class="error">{{ error }}</p>
  </div>
</template>

<script>
export default {
  data() {
    return {
      email: '',
      message: '',
      error: ''
    }
  },
  methods: {
    async register() {
      this.error = ''
      this.message = ''

      try {
        const res = await fetch(`${process.env.VUE_APP_BACKEND_URL}/api/register`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ email: this.email })
        })

        const data = await res.json()

        if (!res.ok) {
          throw new Error(data.error || data.message || 'Registration failed')
        }

        this.message = data.message || 'Registration successful. Check your email.'
        this.email = ''

        // optional: redirect to login after delay
        setTimeout(() => {
          this.$router.push('/login') // or '/verify'
        }, 3000)

      } catch (err) {
        this.error = err.message
      }
    }
  }
}
</script>

<style scoped>
.register-form {
  max-width: 400px;
  margin: auto;
}

.error {
  color: red;
}
</style>
