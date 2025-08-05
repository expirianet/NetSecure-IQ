<template>
  <div class="chart-container">
    <canvas ref="chartCanvas"></canvas>
  </div>
</template>

<script>
import { ref, onMounted, watch } from 'vue';
import Chart from 'chart.js/auto';

export default {
  name: 'LineChart',
  props: {
    chartData: {
      type: Object,
      required: true,
      default: () => ({
        labels: [],
        datasets: []
      })
    },
    options: {
      type: Object,
      default: () => ({
        responsive: true,
        maintainAspectRatio: false
      })
    }
  },
  setup(props) {
    const chartCanvas = ref(null);
    let chart = null;

    const initChart = () => {
      if (!chartCanvas.value) return;
      
      if (chart) {
        chart.destroy();
      }
      
      const ctx = chartCanvas.value.getContext('2d');
      if (!ctx) return;
      
      chart = new Chart(ctx, {
        type: 'line',
        data: props.chartData,
        options: {
          ...props.options,
          responsive: true,
          maintainAspectRatio: false,
        },
      });
    };

    onMounted(() => {
      initChart();
    });

    watch(() => props.chartData, () => {
      initChart();
    }, { deep: true });

    return {
      chartCanvas
    };
  }
};
</script>

<style scoped>
.chart-container {
  position: relative;
  height: 100%;
  min-height: 250px;
  width: 100%;
}
</style>
