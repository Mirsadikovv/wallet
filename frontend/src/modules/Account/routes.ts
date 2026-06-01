import type { RouteRecordRaw } from "vue-router";

const accountPageRoute: RouteRecordRaw = {
  path: "/accounts",
  name: "ACCOUNT_PAGE",
  component: () => import("@module/Account/pages/Page.vue"),
  meta: {
    title: "account_page_title",
    activeLinkGroup: "ACCOUNT_GROUP",
    sidebar: {
      label: "account_page_title",
      icon: "manage_accounts",
      isExpandedGroup: false,
    },
  },
};

const accountCreateRoute: RouteRecordRaw = {
  path: "/accounts/create",
  name: "ACCOUNT_CREATE",
  component: () => import("@module/Account/pages/Create.vue"),
  meta: {
    title: "account_create_title",
    activeLinkGroup: "ACCOUNT_GROUP",
  },
};

const accountViewRoute: RouteRecordRaw = {
  path: "/accounts/:id",
  name: "ACCOUNT_VIEW",
  props: true,
  component: () => import("@module/Account/pages/View.vue"),
  meta: {
    title: "account_view_title",
    activeLinkGroup: "ACCOUNT_GROUP",
  },
};

const accountEditRoute: RouteRecordRaw = {
  path: "/accounts/:id/edit",
  name: "ACCOUNT_EDIT",
  props: true,
  component: () => import("@module/Account/pages/Edit.vue"),
  meta: {
    title: "account_edit_title",
    activeLinkGroup: "ACCOUNT_GROUP",
  },
};

export function AccountRoutes(sort: number): RouteRecordRaw[] {
  return [
    accountPageRoute,
    accountCreateRoute,
    accountViewRoute,
    accountEditRoute,
  ].map((route) => {
    if (route?.meta?.sidebar) {
      return { ...route, meta: { ...route.meta, sort } };
    }
    return route;
  });
}
