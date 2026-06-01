<script setup lang="ts">
import { ref } from "vue";
import { WalletService, type WalletType } from "../service";
import PageLoading from "@/components/PageLoading.vue";
import LoadingSkeleton from "@/components/LoadingSkeleton.vue";
import AppFooter from "@/components/AppFooter.vue";
import { useAuthStore } from "@/store/auth-store";
import { useTelegramViewport } from "@/composables/useTelegramViewport";

const authStore = useAuthStore();
const { containerStyle } = useTelegramViewport();

const wallets = ref<WalletType[]>([]);

async function loadWallets() {
  const userId = authStore.userId;
  if (!userId) return;
  const response = await WalletService.list(userId);
  if (!response) return;
  wallets.value = response.wallets;
}

function truncateAddress(address: string) {
  if (address.length <= 16) return address;
  return `${address.slice(0, 8)}...${address.slice(-8)}`;
}
</script>

<template>
  <PageLoading :find="loadWallets" #="{ loading }">
    <LoadingSkeleton v-if="loading" />

    <q-layout view="hHh Lpr lff" v-else>
      <q-page-container>
        <q-page :style="containerStyle" class="bg-gray-100 text-gray-900 overflow-auto p-4">
          <div class="flex items-center justify-between mb-4">
            <h1 class="text-xl font-bold">{{ $tl("wallet_page_title") }}</h1>
            <span class="text-sm text-gray-400">{{ authStore.displayName }}</span>
          </div>

          <div v-if="wallets.length > 0" class="flex flex-col gap-3">
            <q-card
              v-for="wallet in wallets"
              :key="wallet.id"
              flat
              bordered
              class="rounded-xl cursor-pointer hover:shadow-md transition-shadow"
              @click="$router.push({ name: 'WALLET_VIEW', params: { id: wallet.id } })"
            >
              <q-card-section class="p-4">
                <div class="flex items-center justify-between mb-2">
                  <div class="flex items-center gap-2">
                    <q-icon name="account_balance_wallet" color="primary" size="20px" />
                    <span class="font-semibold text-sm">{{ wallet.wallet_type }}</span>
                  </div>
                  <q-badge
                    :class="wallet.network === 'mainnet' ? 'network-badge mainnet' : 'network-badge testnet'"
                    outline
                  >
                    {{ wallet.network }}
                  </q-badge>
                </div>

                <div class="font-mono text-xs text-gray-500 mb-2 address-truncated">
                  {{ truncateAddress(wallet.address) }}
                </div>

                <div class="flex items-center justify-between text-xs text-gray-400">
                  <span>ID: {{ wallet.id }}</span>
                  <q-chip
                    dense
                    :color="wallet.is_active ? 'positive' : 'negative'"
                    text-color="white"
                    size="xs"
                  >
                    {{ wallet.is_active ? "Активен" : "Неактивен" }}
                  </q-chip>
                </div>
              </q-card-section>
            </q-card>
          </div>

          <div v-else class="flex flex-col items-center justify-center py-16 text-gray-400">
            <q-icon name="account_balance_wallet" size="64px" class="opacity-30 mb-4" />
            <p class="text-center text-sm">У вас нет кошельков</p>
            <p class="text-center text-xs mt-1">Нажмите «+» чтобы создать первый</p>
          </div>
        </q-page>
      </q-page-container>

      <AppFooter
        :username="authStore.displayName"
        :show-add-button="true"
        :add-button-route="{ name: 'WALLET_CREATE' }"
        add-button-icon="add"
      />
    </q-layout>
  </PageLoading>
</template>

<style scoped>
@import "@/styles/telegram-app.scss";
</style>
