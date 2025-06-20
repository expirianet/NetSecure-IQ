<template>
  <div>
    <h1>Router Dashboard</h1>
    <RouterTable />
  </div>
</template>

<script>
import RouterTable from './RouterTable.vue'

export default {
  name: 'DashboardPage',
  components: {
    RouterTable
  }
}
</script>
