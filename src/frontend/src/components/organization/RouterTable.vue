<template>
  <div class="router-table-page">
    <BackgroundParticles />
    <h2 style="margin-bottom: 1rem;">Router Status Table</h2>
    <table class="router-table" v-if="Array.isArray(routers) && routers.length">
      <thead>
        <tr>
          <th>#</th>
          <th>MAC Address</th>
          <th>Status</th>
          <th>Timestamp</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(router, index) in routers" :key="router.mac">
          <td>{{ index + 1 }}</td>
          <td>{{ router.mac }}</td>
          <td :class="router.status === 'online' ? 'online' : 'offline'">{{ router.status }}</td>
          <td>{{ formatDate(router.time) }}</td>
        </tr>
      </tbody>
    </table>
    <p v-else>No routers found.</p>
  </div>
</template>

<script setup>
import BackgroundParticles from '../common/BackgroundParticles.vue';
import { ref, onMounted } from 'vue'

const routers = ref([])

function formatDate(isoString) {
  const date = new Date(isoString)
  return date.toLocaleString()
}

onMounted(async () => {
  const token = localStorage.getItem("token")

  if (!token) {
    console.error("❌ No JWT token found in localStorage")
    return
  }

  try {
    const response = await fetch(`${process.env.VUE_APP_BACKEND_URL}/api/data/routers`, {
      method: "GET",
      headers: {
        "Authorization": `Bearer ${token}`,
        "Content-Type": "application/json"
      }
    })

    const text = await response.text()
    if (!response.ok) {
      console.error("❌ Error response:", text)
      throw new Error(text)
    }

    const data = JSON.parse(text)

    if (!Array.isArray(data)) {
      console.warn("⚠️ Invalid router data:", data)
      routers.value = []
      return
    }

    routers.value = data.map(entry => ({
      mac: entry.mac,
      status: entry.value,
      time: entry.time
    }))
  } catch (err) {
    console.error("❌ Error loading router status:", err.message)
    routers.value = []
  }
})
</script>

<style scoped>
.router-table {
  width: 100%;
  border-collapse: collapse;
  font-family: Arial, sans-serif;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

.router-table th,
.router-table td {
  padding: 12px 16px;
  border: 1px solid #ccc;
  text-align: left;
}

.router-table th {
  background-color: #f4f4f4;
}

.router-table tr:hover {
  background-color: #f9f9f9;
}

.online {
  color: green;
  font-weight: bold;
}

.offline {
  color: red;
  font-weight: bold;
}
</style>
