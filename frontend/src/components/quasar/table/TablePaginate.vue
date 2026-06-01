<script setup lang="ts">
import { ref, watch } from "vue";

type Pikers = {
  page?: number;
  perpage?: number;
};

const props = withDefaults(
  defineProps<{
    pikers: Pikers;
    total: number;
    perpageOptions?: number[];
  }>(),
  {
    perpageOptions: () => [10, 20, 50],
  },
);

const emit = defineEmits<{
  "update:pikers": [val: Pikers];
  page: [query: string];
}>();

const currentPage = ref(props.pikers.page ?? 1);
const perpage = ref(props.pikers.perpage ?? 10);

function buildQuery() {
  const params = new URLSearchParams();
  params.append("page", String(currentPage.value));
  params.append("perpage", String(perpage.value));
  return params.toString();
}

function onPageChange(page: number) {
  currentPage.value = page;
  emit("update:pikers", { page: currentPage.value, perpage: perpage.value });
  emit("page", buildQuery());
}

function onPerpageChange(val: number) {
  perpage.value = val;
  currentPage.value = 1;
  emit("update:pikers", { page: 1, perpage: perpage.value });
  emit("page", buildQuery());
}

watch(
  () => props.pikers,
  (val) => {
    if (val.page) currentPage.value = val.page;
    if (val.perpage) perpage.value = val.perpage;
  },
);
</script>

<template>
  <div class="flex items-center justify-between p-3 border-t border-gray-200 flex-wrap gap-2">
    <div class="flex items-center gap-2 text-sm text-gray-500">
      <span>По</span>
      <q-select
        :model-value="perpage"
        :options="perpageOptions"
        dense
        outlined
        class="w-16"
        @update:model-value="onPerpageChange"
      />
      <span>строк</span>
    </div>

    <q-pagination
      v-if="total > 1"
      :model-value="currentPage"
      :max="total"
      :max-pages="5"
      boundary-numbers
      direction-links
      color="primary"
      active-color="primary"
      @update:model-value="onPageChange"
    />
  </div>
</template>
