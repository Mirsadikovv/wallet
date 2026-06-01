<script setup lang="ts">
defineProps<{
  loading?: boolean;
  type?: "button" | "submit" | "reset";
  color?: string;
  outline?: boolean;
  flat?: boolean;
  disable?: boolean;
}>();

defineEmits<{
  click: [e: MouseEvent];
}>();
</script>

<template>
  <q-btn
    v-bind="$attrs"
    :loading="loading"
    :type="type ?? 'button'"
    :color="color ?? 'primary'"
    :outline="outline"
    :flat="flat"
    :disable="disable"
    no-caps
    unelevated
    class="rounded-xl"
    @click="$emit('click', $event as MouseEvent)"
  >
    <slot />
  </q-btn>
</template>
