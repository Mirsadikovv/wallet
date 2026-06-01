import type { RouteRecordRaw } from "vue-router";

const walletPageRoute: RouteRecordRaw = {
  path: "/wallets",
  name: "WALLET_PAGE",
  component: () => import("@module/Wallet/pages/Page.vue"),
  meta: {
    title: "wallet_page_title",
    activeLinkGroup: "WALLET_GROUP",
    sidebar: {
      label: "wallet_page_title",
      icon: "account_balance_wallet",
      isExpandedGroup: false,
    },
  },
};

const walletCreateRoute: RouteRecordRaw = {
  path: "/wallets/create",
  name: "WALLET_CREATE",
  component: () => import("@module/Wallet/pages/Create.vue"),
  meta: {
    title: "wallet_create_title",
    activeLinkGroup: "WALLET_GROUP",
  },
};

const walletViewRoute: RouteRecordRaw = {
  path: "/wallets/:id",
  name: "WALLET_VIEW",
  props: true,
  component: () => import("@module/Wallet/pages/View.vue"),
  meta: {
    title: "wallet_view_title",
    activeLinkGroup: "WALLET_GROUP",
  },
};

export function WalletRoutes(sort: number): RouteRecordRaw[] {
  return [walletPageRoute, walletCreateRoute, walletViewRoute].map((route) => {
    if (route?.meta?.sidebar) {
      return { ...route, meta: { ...route.meta, sort } };
    }
    return route;
  });
}
