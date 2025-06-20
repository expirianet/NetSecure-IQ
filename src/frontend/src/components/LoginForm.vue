<template>
  <div>
    <h2>Login</h2>
    <form @submit.prevent="login">
      <input v-model="email" placeholder="Email" required />
      <input v-model="password" type="password" placeholder="Password" required />
      <button type="submit">Login</button>
    </form>
    <p v-if="message">{{ message }}</p>
  </div>
</template>

<script>
export default {
  name: 'LoginForm',
  data() {
    return {
      email: '',
      password: '',
      message: ''
    }
  },
  methods: {
    async login() {
      try {
        const res = await fetch('http://localhost:8000/api/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ email: this.email, password: this.password })
        })

        const text = await res.text()
        if (!res.ok) throw new Error(text)
        
        this.message = 'Login successful! Redirecting...'
        setTimeout(() => {
          this.$router.push('/dashboard')
        }, 500) // or even 100ms works
      } catch (err) {
        this.message = 'Error: ' + err.message
      }
    }
  }
}
</script>
