<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { AccountService, type AccountType, type AccountUpdateType } from "../service";
import PageLoading from "@/components/PageLoading.vue";
import LoadingSkeleton from "@/components/LoadingSkeleton.vue";
import Form from "@/components/quasar/form/Form.vue";
import Button from "@/components/quasar/btn/Button.vue";
import Title from "@/components/Title.vue";
import AppFooter from "@/components/AppFooter.vue";
import { useAuthStore } from "@/store/auth-store";
import { useTelegramViewport } from "@/composables/useTelegramViewport";
import { formRequired } from "@/common/validator";

export interface Props {
  id: number | string;
}

const { id } = defineProps<Props>();

const router = useRouter();
const authStore = useAuthStore();
const { containerStyle } = useTelegramViewport();

const accountModel = ref<AccountType>({
  id: 0,
  user_id: 0,
  name: "",
  currency: "",
  balance: "0",
  is_active: true,
  created_at: "",
});

async function loadAccount() {
  const data = await AccountService.getByID(+id);
  if (data) accountModel.value = data;
}

async function save(model: AccountUpdateType) {
  const response = await AccountService.update(+id, { name: model.name });
  if (!response) return false;

  router.push({ name: "ACCOUNT_PAGE" });
  return true;
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
              <q-breadcrumbs-el :label="$tl('page_for_edit')" />
            </q-breadcrumbs>
          </div>

          <Form
            v-model="accountModel as unknown as Record<string, unknown>"
            :save="save as (model: Record<string, unknown>) => Promise<boolean>"
          >
            <template #title>
              <Title class="mb-5">{{ $tl("account_edit_title") }}</Title>
            </template>

            <template #name="{ model }">
              <q-input
                v-model="(model as AccountType).name"
                :label="$tl('account_name')"
                outlined
                dense
                class="col-12 rounded-xl"
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
  </PageLoading>
</template>

<style scoped>
@import "@/styles/telegram-app.scss";
</style>
