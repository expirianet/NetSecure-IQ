<template>
  <div>
    <h2 style="margin-bottom: 1rem;">Router Status Table</h2>
    <table class="router-table">
      <thead>
        <tr>
          <th>#</th>
          <th>MAC Address</th>
          <th>Status</th>
          <th>Timestamp</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(router, index) in routers" :key="index">
          <td>{{ index + 1 }}</td>
          <td>{{ router.mac }}</td>
          <td :class="router.status === 'online' ? 'online' : 'offline'">{{ router.status }}</td>
          <td>{{ formatDate(router.time) }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  name: 'RouterTable',
  data() {
    return {
      routers: []
    };
  },
  methods: {
    formatDate(isoString) {
      const date = new Date(isoString);
      return date.toLocaleString(); // adjust to your locale
    }
  },
  mounted() {
    fetch(`${process.env.VUE_APP_BACKEND_URL}/api/data/routers`)
      .then(res => res.json())
      .then(data => {
        this.routers = data.map(entry => ({
          mac: entry.mac,
          status: entry.value,
          time: entry.time
        }));
      })
      .catch(err => console.error('Error loading router status:', err));
  }
}
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
