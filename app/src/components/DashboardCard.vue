<template>
  <Card class="flex-grow-1 p-3">
    <template #header>
      <h2 class="text-xl font-semibold">{{ title }}</h2>
    </template>
    <template #content>
      <div class="flex justify-content-center align-items-center">
        <AppChart v-if="data" :data="data"></AppChart>
      </div>
    </template>
  </Card>
</template>
<script setup lang="ts">
import AppChart from '@/components/AppChart.vue'
import { onMounted, ref } from 'vue'
const props = defineProps({
  title: String,
  endpoint: {
    type: String,
    required: true
  }
})
const data = ref(null)

onMounted(() => {
  fetch(props.endpoint)
    .then((response) => response.json())
    .then((response) => {
      data.value = response
    })
})
</script>
