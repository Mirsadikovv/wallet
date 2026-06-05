<script setup lang="ts">
import { ref, type Ref } from "vue";
import { useRouter } from "vue-router";
import { WalletService } from "../service";
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

const walletModel = ref({ wallet_type: "V5R1Final", network: "testnet" });
const formModel = walletModel as unknown as Ref<Record<string, unknown>>;

const walletTypes = ["V5R1Final", "V4R2"];
const networks = [
  { label: "Testnet", value: "testnet" },
  { label: "Mainnet", value: "mainnet" },
];

async function save() {
  const userId = authStore.userId;
  if (!userId) return false;

  const response = await WalletService.create({
    wallet_type: walletModel.value.wallet_type,
    network: walletModel.value.network,
    user_id: userId,
  });
  if (!response) return false;

  router.push({ name: "WALLET_PAGE" });
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
              :label="$tl('wallet_list')"
              icon="account_balance_wallet"
              :to="{ name: 'WALLET_PAGE' }"
            />
            <q-breadcrumbs-el :label="$tl('page_for_create')" />
          </q-breadcrumbs>
        </div>

        <Form v-model="formModel" :save="save">
          <template #title>
            <Title class="mb-5">{{ $tl("create_wallet") }}</Title>
          </template>

          <template #wallet_type>
            <q-select
              v-model="walletModel.wallet_type"
              :options="walletTypes"
              :label="$tl('wallet_type')"
              outlined
              dense
              class="col-12"
              :rules="[formRequired()]"
            />
          </template>

          <template #network>
            <q-select
              v-model="walletModel.network"
              :options="networks"
              emit-value
              map-options
              :label="$tl('network')"
              outlined
              dense
              class="col-12"
              :rules="[formRequired()]"
            />
          </template>

          <template #actions="{ loading }">
            <div class="row q-col-gutter-md justify-center">
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
