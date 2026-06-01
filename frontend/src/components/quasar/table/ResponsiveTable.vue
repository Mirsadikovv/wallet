<script setup lang="ts">
import { computed, useSlots } from "vue";

type Row = Record<string, unknown>;

type PageData = {
  data: Row[];
  totalRows: number;
  currentPage: number;
  pageSize: number;
  totalPages: number;
};

const props = withDefaults(
  defineProps<{
    models: PageData;
    hasOrder?: boolean;
    loading?: boolean;
  }>(),
  {
    hasOrder: false,
    loading: false,
  },
);

const slots = useSlots();

const fieldNames = computed(() =>
  Object.keys(slots).filter(
    (name) => !name.includes(":") && name !== "edit" && name !== "tfoot",
  ),
);
</script>

<template>
  <div class="relative">
    <q-inner-loading :showing="loading" color="primary" />

    <q-markup-table v-if="models.data.length > 0" separator="cell" flat bordered class="rounded-xl">
      <thead>
        <tr class="bg-gray-50 text-gray-600 text-sm">
          <th v-if="hasOrder" class="text-center w-12">#</th>
          <th v-for="field in fieldNames" :key="field" class="text-left">
            <slot :name="`${field}:thead`">{{ field }}</slot>
          </th>
          <th v-if="slots['edit']" class="text-center w-20">Действия</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(row, idx) in models.data" :key="idx" class="hover:bg-gray-50 transition-colors">
          <td v-if="hasOrder" class="text-center text-gray-400 text-sm">
            {{ (models.currentPage - 1) * models.pageSize + idx + 1 }}
          </td>
          <td v-for="field in fieldNames" :key="field">
            <slot :name="field" :model="row" />
          </td>
          <td v-if="slots['edit']" class="text-center">
            <slot name="edit" :model="row" />
          </td>
        </tr>
      </tbody>
    </q-markup-table>

    <div v-else-if="!loading" class="text-center text-gray-400 py-12">
      <q-icon name="inbox" size="48px" class="mb-2 opacity-40" />
      <div>Нет данных</div>
    </div>

    <slot name="tfoot" :total-pages="models.totalPages" />
  </div>
</template>
