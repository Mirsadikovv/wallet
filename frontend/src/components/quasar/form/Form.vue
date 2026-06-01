<script setup lang="ts">
import { ref, useSlots } from "vue";

const props = defineProps<{
  modelValue: Record<string, unknown>;
  save: (model: Record<string, unknown>) => Promise<boolean | undefined>;
}>();

defineEmits<{
  "update:modelValue": [value: Record<string, unknown>];
}>();

const slots = useSlots();
const loading = ref(false);
const formRef = ref<{ validate: () => Promise<boolean> } | null>(null);

const fieldSlotNames = () =>
  Object.keys(slots).filter((name) => name !== "title" && name !== "actions");

async function submit() {
  const valid = await formRef.value?.validate();
  if (!valid) return;
  loading.value = true;
  await props.save(props.modelValue);
  loading.value = false;
}
</script>

<template>
  <q-form ref="formRef" @submit.prevent="submit">
    <slot name="title" />

    <q-card flat bordered class="p-4 rounded-xl mb-4">
      <div class="row q-col-gutter-md">
        <template v-for="name in fieldSlotNames()" :key="name">
          <slot :name="name" :model="modelValue" />
        </template>
      </div>
    </q-card>

    <slot name="actions" :loading="loading" />
  </q-form>
</template>
