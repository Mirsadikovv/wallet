import { ref, computed, onMounted, onUnmounted } from "vue";

export function useTelegramViewport() {
  const viewportHeight = ref(window.innerHeight);

  function onViewportChanged() {
    viewportHeight.value =
      window.Telegram?.WebApp?.viewportStableHeight || window.innerHeight;
  }

  onMounted(() => {
    window.Telegram?.WebApp?.onEvent("viewportChanged", onViewportChanged);
    onViewportChanged();
  });

  onUnmounted(() => {
    window.Telegram?.WebApp?.offEvent("viewportChanged", onViewportChanged);
  });

  const containerStyle = computed(() => ({
    minHeight: `${viewportHeight.value}px`,
    paddingBottom: "72px",
  }));

  return { viewportHeight, containerStyle };
}
