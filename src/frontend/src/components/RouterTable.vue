<template>
  <div>
    <h2>Router Table</h2>
    <table v-if="Array.isArray(routers) && routers.length">
      <thead>
        <tr>
          <th>MAC</th>
          <th>Status</th>
          <th>Time</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="router in routers" :key="router.mac">
          <td>{{ router.mac }}</td>
          <td>{{ router.value }}</td>
          <td>{{ router.time }}</td>
        </tr>
      </tbody>
    </table>
    <p v-else>No routers found.</p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const routers = ref([])

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
    console.log("📥 Raw router response:", text)

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

    routers.value = data
  } catch (err) {
    console.error("❌ Error loading router status:", err.message)
    routers.value = []
  }
})
</script>

<style scoped>
table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 1rem;
}

th, td {
  border: 1px solid #ccc;
  padding: 8px 12px;
  text-align: left;
}

th {
  background-color: #f4f4f4;
}

tr:hover {
  background-color: #f9f9f9;
}
</style>
