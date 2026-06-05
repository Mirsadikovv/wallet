<script setup lang="ts">
import { useRouter, useRoute, type RouteLocationRaw } from "vue-router";

defineProps<{
  username?: string;
  showAddButton?: boolean;
  addButtonRoute?: RouteLocationRaw;
  addButtonIcon?: string;
}>();

defineEmits<{
  toggleDrawer: [];
  goToProfile: [];
  logout: [];
}>();

const router = useRouter();
const route = useRoute();

function isActive(routePrefix: string) {
  return route.path.startsWith(routePrefix);
}
</script>

<template>
  <q-footer class="bg-white border-t border-gray-200" style="height: 64px">
    <div class="flex items-center justify-around h-full px-2">
      <q-btn
        flat
        no-caps
        class="flex flex-col items-center gap-1 min-w-0"
        :class="isActive('/wallets') ? 'text-primary' : 'text-gray-400'"
        @click="router.push('/wallets')"
      >
        <q-icon name="account_balance_wallet" size="24px" />
        <span class="text-10px">Кошельки</span>
      </q-btn>

      <q-btn
        v-if="showAddButton && addButtonRoute"
        fab
        :icon="addButtonIcon || 'add'"
        color="primary"
        class="-mt-6"
        @click="router.push(addButtonRoute!)"
      />

      <q-btn
        flat
        no-caps
        class="flex flex-col items-center gap-1 min-w-0"
        :class="isActive('/accounts') ? 'text-primary' : 'text-gray-400'"
        @click="router.push('/accounts')"
      >
        <q-icon name="manage_accounts" size="24px" />
        <span class="text-10px">Счета</span>
      </q-btn>
    </div>
  </q-footer>
</template>
