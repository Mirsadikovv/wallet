<script setup lang="ts">
defineProps<{
  modelValue?: string | number | null;
  label?: string;
  placeholder?: string;
  type?: "text" | "password" | "email" | "search" | "tel" | "file" | "number" | "url" | "time" | "date" | "datetime-local" | "textarea";
  rules?: ((val: unknown) => boolean | string)[];
  disable?: boolean;
  readonly?: boolean;
  hint?: string;
  prefix?: string;
  suffix?: string;
  autofocus?: boolean;
}>();

defineEmits<{
  "update:modelValue": [value: string];
}>();
</script>

<template>
  <q-input
    v-bind="$attrs"
    :model-value="modelValue"
    :label="label"
    :placeholder="placeholder"
    :type="type ?? 'text'"
    :rules="rules"
    :disable="disable"
    :readonly="readonly"
    :hint="hint"
    :prefix="prefix"
    :suffix="suffix"
    :autofocus="autofocus"
    outlined
    dense
    class="rounded-xl"
    @update:model-value="$emit('update:modelValue', $event as string)"
  />
</template>
