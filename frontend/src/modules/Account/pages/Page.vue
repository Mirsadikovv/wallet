<script setup lang="ts">
import { ref } from "vue";
import { AccountService, type AccountType } from "../service";
import PageLoading from "@/components/PageLoading.vue";
import LoadingSkeleton from "@/components/LoadingSkeleton.vue";
import AppFooter from "@/components/AppFooter.vue";
import IconBtn from "@/components/quasar/btn/IconBtn.vue";
import { useAuthStore } from "@/store/auth-store";
import { useTelegramViewport } from "@/composables/useTelegramViewport";
import { SuccessNotify } from "@/common/Notify";

const authStore = useAuthStore();
const { containerStyle } = useTelegramViewport();

const accounts = ref<AccountType[]>([]);

async function loadAccounts() {
  const userId = authStore.userId;
  if (!userId) return;
  const response = await AccountService.list(userId);
  if (!response) return;
  accounts.value = response.accounts;
}

const currencyColors: Record<string, string> = {
  USD: "positive",
  EUR: "primary",
  RUB: "warning",
};

function currencyColor(currency: string) {
  return currencyColors[currency] ?? "grey";
}
</script>

<template>
  <PageLoading :find="loadAccounts" #="{ loading, fetch }">
    <LoadingSkeleton v-if="loading" />

    <q-layout view="hHh Lpr lff" v-else>
      <q-page-container>
        <q-page :style="containerStyle" class="bg-gray-100 text-gray-900 overflow-auto p-4">
          <div class="flex items-center justify-between mb-4">
            <h1 class="text-xl font-bold">{{ $tl("account_page_title") }}</h1>
            <q-btn flat round icon="refresh" color="primary" size="sm" @click="fetch()" />
          </div>

          <div v-if="accounts.length > 0" class="flex flex-col gap-3">
            <q-card
              v-for="account in accounts"
              :key="account.id"
              flat
              bordered
              class="rounded-xl"
            >
              <q-card-section class="p-4">
                <div class="flex items-center justify-between mb-2">
                  <div class="flex items-center gap-2">
                    <q-icon name="account_balance" color="secondary" size="20px" />
                    <span class="font-semibold">{{ account.name }}</span>
                  </div>
                  <q-badge :color="currencyColor(account.currency)" class="text-xs">
                    {{ account.currency }}
                  </q-badge>
                </div>

                <div class="text-2xl font-bold text-secondary mb-2">
                  {{ account.balance }}
                  <span class="text-sm font-normal text-gray-400">{{ account.currency }}</span>
                </div>

                <div class="flex items-center justify-between">
                  <q-chip
                    dense
                    :color="account.is_active ? 'positive' : 'negative'"
                    text-color="white"
                    size="xs"
                  >
                    {{ account.is_active ? "Активен" : "Неактивен" }}
                  </q-chip>

                  <div class="flex gap-1">
                    <IconBtn
                      v-if="$canPage('ACCOUNT_EDIT')"
                      :to="{ name: 'ACCOUNT_EDIT', params: { id: account.id } }"
                      icon="edit"
                    />
                    <IconBtn
                      v-if="$canPage('ACCOUNT_VIEW')"
                      :to="{ name: 'ACCOUNT_VIEW', params: { id: account.id } }"
                      icon="visibility"
                    />
                  </div>
                </div>
              </q-card-section>
            </q-card>
          </div>

          <div v-else class="flex flex-col items-center justify-center py-16 text-gray-400">
            <q-icon name="account_balance" size="64px" class="opacity-30 mb-4" />
            <p class="text-center text-sm">У вас нет счетов</p>
            <p class="text-center text-xs mt-1">Нажмите «+» чтобы создать первый</p>
          </div>
        </q-page>
      </q-page-container>

      <AppFooter
        :username="authStore.displayName"
        :show-add-button="true"
        :add-button-route="{ name: 'ACCOUNT_CREATE' }"
        add-button-icon="add"
      />
    </q-layout>
  </PageLoading>
</template>

<style scoped>
@import "@/styles/telegram-app.scss";
</style>
