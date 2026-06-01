import { defineStore } from "pinia";

type TelegramUser = {
  id: number;
  first_name: string;
  last_name?: string;
  username?: string;
};

type State = {
  user: TelegramUser | undefined;
};

export const useAuthStore = defineStore("auth", {
  state(): State {
    return {
      user: undefined,
    };
  },
  actions: {
    initFromTelegram() {
      const tgUser = window.Telegram?.WebApp?.initDataUnsafe?.user;
      if (tgUser) {
        this.user = {
          id: tgUser.id,
          first_name: tgUser.first_name,
          last_name: tgUser.last_name,
          username: tgUser.username,
        };
      }
    },
  },
  getters: {
    isAuthenticated(state): Readonly<boolean> {
      return !!state.user;
    },
    userId(state): Readonly<number | undefined> {
      return state.user?.id;
    },
    displayName(state): Readonly<string> {
      if (!state.user) return "User";
      return (
        state.user.username ||
        `${state.user.first_name} ${state.user.last_name ?? ""}`.trim()
      );
    },
  },
});
