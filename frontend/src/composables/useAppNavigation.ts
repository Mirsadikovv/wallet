import { useRouter } from "vue-router";

export function useAppNavigation() {
  const router = useRouter();

  function toggleLeftDrawer() {
    window.Telegram?.WebApp?.HapticFeedback?.impactOccurred("light");
  }

  function setLang(_id: number) {
    // TODO: добавить поддержку языков
  }

  function logout() {
    window.Telegram?.WebApp?.close();
  }

  function goBack() {
    router.back();
  }

  return { toggleLeftDrawer, setLang, logout, goBack };
}
