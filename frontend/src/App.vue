<script setup lang="ts">
import { onMounted } from "vue";
import { useAuthStore } from "@/store/auth-store";
import { RouterView } from "vue-router";

const authStore = useAuthStore();

onMounted(() => {
  window.Telegram?.WebApp?.ready();
  window.Telegram?.WebApp?.expand();
  authStore.initFromTelegram();
});
</script>

<template>
  <RouterView />
</template>
