<script setup lang="ts">
import type { RouteLocationRaw } from "vue-router";

defineProps<{
  icon: string;
  color?: string;
  to?: RouteLocationRaw;
  size?: string;
}>();
</script>

<template>
  <q-btn
    v-bind="$attrs"
    flat
    round
    :icon="icon"
    :color="color ?? 'primary'"
    :to="to"
    :size="size ?? 'sm'"
  />
</template>
