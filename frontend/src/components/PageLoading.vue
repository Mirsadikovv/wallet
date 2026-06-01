<script setup lang="ts">
import { ref, onMounted } from "vue";

const props = defineProps<{
  find: (query?: string) => Promise<void> | void;
}>();

const loading = ref(true);

async function fetch(query = "") {
  loading.value = true;
  await props.find(query);
  loading.value = false;
}

onMounted(() => fetch());
</script>

<template>
  <slot :fetch="fetch" :loading="loading" />
</template>
