<template> 
    <div id="app">
        <h1>NetSecure IQ – Device Status</h1>

      <form @submit.prevent="fetchStatus">
        <input v-model="mac" placeholder="Enter MAC address" />
        <button>Check Status</button>
    </form>

    <div v-if="status !== null">
        <p><strong>Status:</strong> {{ status }}</p>
    </div>

    <div v-if="error" class="error">
        <p>Error: {{ error }}</p>
    </div>
</div>
</template >

<script>
export default {
  data() {
    return {
      mac: '',
      status: null,
      error: null
    }
  },
  methods: {
    async fetchStatus() {
      try {
        this.error = null;
        const response = await fetch('http://localhost:8080/api/ping', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ mac: this.mac, status: 'online' }) // simulate status
        });

        if (!response.ok) {
          const text = await response.text();
          throw new Error(text);
        }

        this.status = 'Online (simulated)';
      } catch (err) {
        this.status = null;
        this.error = err.message;
      }
    }
  }
}
</script>

<style>
#error {
  color: red;
}
</style>
