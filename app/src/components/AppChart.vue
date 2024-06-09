<template>
  <canvas :class="'chart-' + uuid + ' w-full'"></canvas>
</template>
<script setup lang="ts">
import { BarController, BarElement, CategoryScale, Chart, type ChartType } from 'chart.js/auto'
const uuid = Math.random().toString(36).substring(7)
import { onMounted } from 'vue'
import {
  ChoroplethController,
  ColorScale,
  GeoFeature,
  ProjectionScale,
  topojson
} from 'chartjs-chart-geo'
import { countryIsoToName } from '@/components/countryIsoToName'

Chart.registry.addControllers(BarController, ChoroplethController)
Chart.registry.addElements(BarElement, GeoFeature)
Chart.registry.addScales(CategoryScale, ColorScale, ProjectionScale)

const props = defineProps({
  data: {
    type: Object as () => any,
    required: true
  },
  type: {
    type: String as () => ChartType,
    default: 'bar'
  }
})

function renderChoropleth() {
  fetch('https://unpkg.com/world-atlas/countries-110m.json')
    .then((r) => r.json())
    .then((data) => {
      // @ts-ignore - apparently the types are wrong
      const countries = topojson.feature(data, data.objects.countries).features

      new Chart(document.querySelector('.chart-' + uuid) as HTMLCanvasElement, {
        type: 'choropleth',
        data: {
          labels: countries.map((d) => d.properties.name),
          datasets: [
            {
              label: 'Countries',
              data: countries.map((d) => {
                const value =
                  props.data.datasets.find(
                    (dataset) => countryIsoToName(dataset.label) === d.properties.name
                  )?.data[0] || 0
                return { feature: d, value }
              })
            }
          ]
        },
        options: {
          showOutline: true,
          showGraticule: true,
          plugins: {
            legend: {
              display: false
            }
          },
          scales: {
            projection: {
              axis: 'x',
              projection: 'equalEarth'
            }
          }
        }
      })
    })
}

function renderChart() {
  new Chart(document.querySelector('.chart-' + uuid) as HTMLCanvasElement, {
    type: props.type,
    data: props.data
  })
}

onMounted(() => {
  if (props.type === 'choropleth') {
    renderChoropleth()
  } else {
    renderChart()
  }
})
</script>
