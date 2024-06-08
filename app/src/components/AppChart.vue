<template>
  <canvas :class="'chart-' + uuid + ' w-full'"></canvas>
</template>
<script setup lang="ts">
import { Chart } from 'chart.js/auto'

const uuid = Math.random().toString(36).substring(7)

import { onMounted } from 'vue'

const props = defineProps({
  data: {
    type: Object as () => any,
    required: true
  }
})

onMounted(() => {
  ;(async function () {
    const chart = new Chart(document.querySelector('.chart-' + uuid) as HTMLCanvasElement, {
      type: 'bar',
      data: props.data
    })

    window.addEventListener('beforeprint', () => {
      console.log('?')
      chart.resize(600, 600)
    })
    window.addEventListener('afterprint', () => {
      chart.resize()
    })
  })()
})
</script>
