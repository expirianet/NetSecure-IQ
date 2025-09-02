<template>
  <div class="chart-root">
    <canvas ref="canvas"></canvas>
  </div>
</template>

<script setup>
import Chart from 'chart.js/auto' // auto-registers elements/scales/plugins
import { ref, watch, onMounted, onBeforeUnmount } from 'vue'

const props = defineProps({
  chartData: { type: Object, required: true },
  options: { type: Object, default: () => ({}) }
})

const canvas = ref(null)
let chartInstance = null

function renderChart () {
  if (!canvas.value) return
  // Détruit l’instance précédente si elle existe
  if (chartInstance) {
    chartInstance.destroy()
    chartInstance = null
  }
  const ctx = canvas.value.getContext('2d')
  chartInstance = new Chart(ctx, {
    type: 'line',
    data: props.chartData,
    options: {
      responsive: true,
      maintainAspectRatio: false,
      animation: false,
      plugins: { legend: { display: false }, tooltip: { enabled: false } },
      elements: { point: { radius: 0 }, line: { tension: 0.3, borderWidth: 2 } },
      scales: { x: { display: false }, y: { display: false } },
      // Permet de surcharger via :options="..."
      ...props.options
    }
  })
}

onMounted(renderChart)

// Re-render si les props changent
watch(() => props.chartData, renderChart, { deep: true })
watch(() => props.options, renderChart, { deep: true })

onBeforeUnmount(() => {
  if (chartInstance) chartInstance.destroy()
})
</script>

<style scoped>
.chart-root {
  position: relative;
  width: 100%;
  height: 100%;
}
canvas {
  width: 100% !important;
  height: 100% !important;
  display: block;
  /* laisse les clics passer au parent si besoin */
  pointer-events: none;
}
</style>
