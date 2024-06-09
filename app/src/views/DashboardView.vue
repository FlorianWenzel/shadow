<template>
  <div class="dashboard">
    <DashboardCard
      v-for="card of cards"
      :key="card.id"
      :title="card.title"
      :endpoint="card.endpoint"
      :type="card.type"
    ></DashboardCard>
  </div>
</template>

<script setup lang="ts">
import DashboardCard from '@/components/DashboardCard.vue'
import { ref } from 'vue'

const cards = ref([
  { id: 1, title: 'Card 1', endpoint: '/api/requests-per-page' },
  { id: 2, title: 'Card 2', endpoint: '/api/requests-per-country', type: 'choropleth' },
  { id: 3, title: 'Card 2', endpoint: '/api/requests-per-ip', type: 'bar' }
])
</script>
<style>
.dashboard {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(800px, 1fr));
  grid-auto-rows: 1fr;
  gap: 1rem;
}
</style>
<style src="gridstack/dist/gridstack.min.css"></style>
