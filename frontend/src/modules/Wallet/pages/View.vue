<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import {
  WalletService,
  type WalletInfoType,
  type TransactionType,
  type SendCoinsType,
} from "../service";
import PageLoading from "@/components/PageLoading.vue";
import LoadingSkeleton from "@/components/LoadingSkeleton.vue";
import AppFooter from "@/components/AppFooter.vue";
import { useAuthStore } from "@/store/auth-store";
import { useTelegramViewport } from "@/composables/useTelegramViewport";
import { formRequired, formPositiveNumber, formTONAddress } from "@/common/validator";
import { SuccessNotify } from "@/common/Notify";

export interface Props {
  id: number | string;
}

const { id } = defineProps<Props>();

const router = useRouter();
const authStore = useAuthStore();
const { containerStyle } = useTelegramViewport();

const walletInfo = ref<WalletInfoType | null>(null);
const transactions = ref<TransactionType[]>([]);
const sendModel = ref<Partial<SendCoinsType>>({});
const sendLoading = ref(false);
const showSendDialog = ref(false);
const sendFormRef = ref<{ validate: () => Promise<boolean> } | null>(null);
const activeTab = ref("info");

async function loadWallet() {
  const [info, txData] = await Promise.all([
    WalletService.getInfo(+id),
    WalletService.getTransactions(+id, 20),
  ]);

  if (info) walletInfo.value = info;
  if (txData) transactions.value = txData.transactions;
}

async function sendTON() {
  const valid = await sendFormRef.value?.validate();
  if (!valid || !sendModel.value.recipient || !sendModel.value.amount) return;

  sendLoading.value = true;
  const result = await WalletService.send(+id, {
    recipient: sendModel.value.recipient,
    amount: sendModel.value.amount,
    comment: sendModel.value.comment,
  });
  sendLoading.value = false;

  if (!result) return;

  SuccessNotify("Транзакция отправлена!");
  showSendDialog.value = false;
  sendModel.value = {};
  await loadWallet();
}

async function removeWallet() {
  const result = await WalletService.remove(+id);
  if (!result) return;
  SuccessNotify("Кошелёк удалён");
  router.push({ name: "WALLET_PAGE" });
}

function truncateAddress(address: string) {
  if (address.length <= 20) return address;
  return `${address.slice(0, 10)}...${address.slice(-10)}`;
}

function formatDate(ts: number) {
  return new Date(ts * 1000).toLocaleString("ru-RU");
}

function copyAddress() {
  if (!walletInfo.value) return;
  navigator.clipboard.writeText(walletInfo.value.address);
  SuccessNotify("Адрес скопирован");
}
</script>

<template>
  <PageLoading :find="loadWallet" #="{ loading }">
    <LoadingSkeleton v-if="loading" />

    <q-layout view="hHh Lpr lff" v-else>
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
              <q-breadcrumbs-el :label="$tl('wallet_view_title')" />
            </q-breadcrumbs>
            <q-space />
            <q-btn
              flat
              round
              icon="delete"
              color="negative"
              size="sm"
              @click="removeWallet"
            />
          </div>

          <div v-if="walletInfo">
            <q-card flat bordered class="rounded-xl mb-4">
              <q-card-section class="p-4">
                <div class="flex items-center justify-between mb-3">
                  <div class="flex items-center gap-2">
                    <q-icon name="account_balance_wallet" color="primary" size="24px" />
                    <span class="font-bold">{{ walletInfo.wallet_type }}</span>
                  </div>
                  <q-badge
                    :class="walletInfo.network === 'mainnet' ? 'network-badge mainnet' : 'network-badge testnet'"
                    outline
                  >
                    {{ walletInfo.network }}
                  </q-badge>
                </div>

                <div
                  class="font-mono text-xs text-gray-500 mb-1 flex items-center gap-1 cursor-pointer"
                  @click="copyAddress"
                >
                  <span class="break-all">{{ truncateAddress(walletInfo.address) }}</span>
                  <q-icon name="content_copy" size="14px" class="text-gray-400" />
                  <q-tooltip>Нажмите для копирования</q-tooltip>
                </div>

                <div class="balance-large text-primary my-3">
                  {{ walletInfo.balance }}
                  <span class="text-sm font-normal text-gray-400">TON</span>
                </div>

                <div class="flex gap-2">
                  <q-btn
                    color="primary"
                    unelevated
                    no-caps
                    class="rounded-xl flex-1"
                    icon="send"
                    label="Отправить"
                    @click="showSendDialog = true"
                  />
                  <q-btn
                    outline
                    color="primary"
                    no-caps
                    class="rounded-xl"
                    icon="refresh"
                    @click="loadWallet"
                  />
                </div>
              </q-card-section>
            </q-card>

            <q-tabs
              v-model="activeTab"
              dense
              active-color="primary"
              indicator-color="primary"
              align="justify"
              class="mb-3 bg-white rounded-xl"
            >
              <q-tab name="info" label="Детали" icon="info" />
              <q-tab name="transactions" label="Транзакции" icon="receipt_long" />
            </q-tabs>

            <q-tab-panels v-model="activeTab" animated>
              <q-tab-panel name="info" class="p-0">
                <q-markup-table separator="cell" flat bordered class="rounded-xl">
                  <tbody>
                    <tr>
                      <td class="font-bold text-left w-40">{{ $tl("address") }}</td>
                      <td class="font-mono text-xs break-all">{{ walletInfo.address }}</td>
                    </tr>
                    <tr>
                      <td class="font-bold text-left">{{ $tl("balance") }}</td>
                      <td>{{ walletInfo.balance }} TON</td>
                    </tr>
                    <tr>
                      <td class="font-bold text-left">{{ $tl("wallet_type") }}</td>
                      <td>{{ walletInfo.wallet_type }}</td>
                    </tr>
                    <tr>
                      <td class="font-bold text-left">{{ $tl("network") }}</td>
                      <td>{{ walletInfo.network }}</td>
                    </tr>
                    <tr>
                      <td class="font-bold text-left">Seqno</td>
                      <td>{{ walletInfo.seqno }}</td>
                    </tr>
                    <tr>
                      <td class="font-bold text-left">{{ $tl("is_active") }}</td>
                      <td>{{ walletInfo.is_active ? "Да" : "Нет" }}</td>
                    </tr>
                    <tr>
                      <td class="font-bold text-left">{{ $tl("created_at") }}</td>
                      <td>{{ new Date(walletInfo.created_at).toLocaleString("ru-RU") }}</td>
                    </tr>
                  </tbody>
                </q-markup-table>
              </q-tab-panel>

              <q-tab-panel name="transactions" class="p-0">
                <div v-if="transactions.length > 0" class="flex flex-col gap-2">
                  <q-card
                    v-for="tx in transactions"
                    :key="tx.hash"
                    flat
                    bordered
                    class="rounded-xl"
                  >
                    <q-card-section class="p-3">
                      <div class="flex items-center justify-between mb-1">
                        <q-chip
                          dense
                          :color="tx.type === 'in' ? 'positive' : 'warning'"
                          text-color="white"
                          size="xs"
                          :icon="tx.type === 'in' ? 'arrow_downward' : 'arrow_upward'"
                        >
                          {{ tx.type === "in" ? "Входящая" : "Исходящая" }}
                        </q-chip>
                        <span :class="tx.success ? 'text-positive' : 'text-negative'" class="text-xs font-bold">
                          {{ tx.success ? "✓" : "✗" }}
                        </span>
                      </div>

                      <div class="text-lg font-bold" :class="tx.type === 'in' ? 'text-positive' : 'text-orange-500'">
                        {{ tx.type === "in" ? "+" : "-" }}{{ tx.amount }} TON
                      </div>

                      <div class="text-xs text-gray-400 mt-1">
                        <div class="tx-hash">{{ tx.hash }}</div>
                        <div>{{ formatDate(tx.timestamp) }}</div>
                        <div v-if="tx.comment" class="italic mt-1">"{{ tx.comment }}"</div>
                      </div>
                    </q-card-section>
                  </q-card>
                </div>

                <div v-else class="text-center text-gray-400 py-8">
                  <q-icon name="receipt_long" size="48px" class="opacity-30 mb-2" />
                  <p class="text-sm">Транзакций нет</p>
                </div>
              </q-tab-panel>
            </q-tab-panels>
          </div>
        </q-page>
      </q-page-container>

      <AppFooter :username="authStore.displayName" :show-add-button="false" />
    </q-layout>
  </PageLoading>

  <q-dialog v-model="showSendDialog" position="bottom">
    <q-card class="rounded-t-2xl w-full" style="max-width: 480px">
      <q-card-section class="pb-0">
        <div class="flex items-center justify-between mb-3">
          <span class="text-lg font-bold">{{ $tl("send_ton") }}</span>
          <q-btn flat round icon="close" dense @click="showSendDialog = false" />
        </div>
      </q-card-section>

      <q-card-section>
        <q-form ref="sendFormRef" class="flex flex-col gap-3">
          <q-input
            v-model="sendModel.recipient"
            :label="$tl('recipient')"
            outlined
            dense
            class="rounded-xl"
            :rules="[formRequired(), formTONAddress()]"
          />
          <q-input
            v-model="sendModel.amount"
            :label="$tl('amount')"
            type="number"
            outlined
            dense
            suffix="TON"
            class="rounded-xl"
            :rules="[formRequired(), formPositiveNumber()]"
          />
          <q-input
            v-model="sendModel.comment"
            :label="$tl('comment')"
            outlined
            dense
            class="rounded-xl"
          />
          <q-btn
            color="primary"
            unelevated
            no-caps
            class="rounded-xl w-full"
            :loading="sendLoading"
            @click="sendTON"
          >
            {{ $tl("send_ton") }}
          </q-btn>
        </q-form>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<style scoped>
@import "@/styles/telegram-app.scss";
</style>
