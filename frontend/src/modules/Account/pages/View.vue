<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { AccountService, type AccountType } from "../service";
import PageLoading from "@/components/PageLoading.vue";
import LoadingSkeleton from "@/components/LoadingSkeleton.vue";
import AppFooter from "@/components/AppFooter.vue";
import { useAuthStore } from "@/store/auth-store";
import { useTelegramViewport } from "@/composables/useTelegramViewport";
import { SuccessNotify } from "@/common/Notify";

export interface Props {
  id: number | string;
}

const { id } = defineProps<Props>();

const router = useRouter();
const authStore = useAuthStore();
const { containerStyle } = useTelegramViewport();

const accountModel = ref<AccountType>({} as AccountType);

async function loadAccount() {
  const data = await AccountService.getByID(+id);
  if (data) accountModel.value = data;
}

async function toggleDeleteRestore() {
  const result = await AccountService.deleteOrRestore(+id);
  if (!result) return;
  SuccessNotify(accountModel.value.is_active ? "Счёт удалён" : "Счёт восстановлен");
  router.push({ name: "ACCOUNT_PAGE" });
}
</script>

<template>
  <PageLoading :find="loadAccount" #="{ loading }">
    <LoadingSkeleton v-if="loading" />

    <q-layout view="hHh Lpr lff" v-else>
      <q-page-container>
        <q-page :style="containerStyle" class="bg-gray-100 text-gray-900 overflow-auto p-4">
          <div class="flex gap-x-3 items-center mb-4">
            <q-btn flat round color="primary" icon="arrow_back" @click="router.back()" />
            <q-breadcrumbs>
              <q-breadcrumbs-el
                :label="$tl('account_list')"
                icon="manage_accounts"
                :to="{ name: 'ACCOUNT_PAGE' }"
              />
              <q-breadcrumbs-el :label="$tl('account_view_title')" />
            </q-breadcrumbs>
            <q-space />
            <q-btn
              flat
              round
              :icon="accountModel.is_active ? 'delete' : 'restore'"
              :color="accountModel.is_active ? 'negative' : 'positive'"
              size="sm"
              @click="toggleDeleteRestore"
            />
            <q-btn
              v-if="$canPage('ACCOUNT_EDIT')"
              flat
              round
              icon="edit"
              color="primary"
              size="sm"
              :to="{ name: 'ACCOUNT_EDIT', params: { id } }"
            />
          </div>

          <q-markup-table separator="cell" flat bordered class="rounded-xl">
            <tbody>
              <tr>
                <td class="font-bold text-left w-40">ID</td>
                <td>{{ accountModel.id }}</td>
              </tr>
              <tr>
                <td class="font-bold text-left">{{ $tl("account_name") }}</td>
                <td>{{ accountModel.name }}</td>
              </tr>
              <tr>
                <td class="font-bold text-left">{{ $tl("currency") }}</td>
                <td>
                  <q-badge color="primary">{{ accountModel.currency }}</q-badge>
                </td>
              </tr>
              <tr>
                <td class="font-bold text-left">{{ $tl("balance") }}</td>
                <td class="font-bold text-secondary">
                  {{ accountModel.balance }} {{ accountModel.currency }}
                </td>
              </tr>
              <tr>
                <td class="font-bold text-left">{{ $tl("is_active") }}</td>
                <td>
                  <q-chip
                    dense
                    :color="accountModel.is_active ? 'positive' : 'negative'"
                    text-color="white"
                    size="sm"
                  >
                    {{ accountModel.is_active ? "Активен" : "Удалён" }}
                  </q-chip>
                </td>
              </tr>
              <tr>
                <td class="font-bold text-left">{{ $tl("created_at") }}</td>
                <td>{{ new Date(accountModel.created_at).toLocaleString("ru-RU") }}</td>
              </tr>
            </tbody>
          </q-markup-table>
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
