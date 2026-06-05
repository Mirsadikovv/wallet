<script setup lang="ts">
import { ref, type Ref } from "vue";
import { useRouter } from "vue-router";
import { AccountService } from "../service";
import Form from "@/components/quasar/form/Form.vue";
import Button from "@/components/quasar/btn/Button.vue";
import Title from "@/components/Title.vue";
import AppFooter from "@/components/AppFooter.vue";
import { useAuthStore } from "@/store/auth-store";
import { useTelegramViewport } from "@/composables/useTelegramViewport";
import { formRequired } from "@/common/validator";

const router = useRouter();
const authStore = useAuthStore();
const { containerStyle } = useTelegramViewport();

const accountModel = ref({ name: "", currency: "USD" });
const formModel = accountModel as unknown as Ref<Record<string, unknown>>;

const currencies = [
  { label: "USD — Доллар США", value: "USD" },
  { label: "EUR — Евро", value: "EUR" },
  { label: "RUB — Рубль", value: "RUB" },
];

async function save() {
  const userId = authStore.userId;
  if (!userId) return false;

  const response = await AccountService.create({
    name: accountModel.value.name,
    currency: accountModel.value.currency,
    user_id: userId,
  });
  if (!response) return false;

  router.push({ name: "ACCOUNT_PAGE" });
  return true;
}
</script>

<template>
  <q-layout view="hHh Lpr lff">
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
            <q-breadcrumbs-el :label="$tl('page_for_create')" />
          </q-breadcrumbs>
        </div>

        <Form v-model="formModel" :save="save">
          <template #title>
            <Title class="mb-5">{{ $tl("create_account") }}</Title>
          </template>

          <template #name>
            <q-input
              v-model="accountModel.name"
              :label="$tl('account_name')"
              outlined
              dense
              class="col-12 rounded-xl"
              :rules="[formRequired()]"
            />
          </template>

          <template #currency>
            <q-select
              v-model="accountModel.currency"
              :options="currencies"
              emit-value
              map-options
              :label="$tl('currency')"
              outlined
              dense
              class="col-12"
              :rules="[formRequired()]"
            />
          </template>

          <template #actions="{ loading }">
            <div class="row justify-center mt-2">
              <q-space />
              <div class="col-auto">
                <Button :loading="loading" type="submit" class="w-48">
                  {{ $tl("save") }}
                </Button>
              </div>
            </div>
          </template>
        </Form>
      </q-page>
    </q-page-container>

    <AppFooter :username="authStore.displayName" :show-add-button="false" />
  </q-layout>
</template>

<style scoped>
@import "@/styles/telegram-app.scss";
</style>
