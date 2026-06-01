import { createRouter, createWebHistory } from "vue-router";
import { WalletRoutes } from "@module/Wallet/routes";
import { AccountRoutes } from "@module/Account/routes";

export const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      redirect: "/wallets",
    },
    ...WalletRoutes(1),
    ...AccountRoutes(2),
    {
      path: "/:pathMatch(.*)*",
      redirect: "/wallets",
    },
  ],
});
