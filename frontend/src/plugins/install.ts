import type { App } from "vue";
import { installPinia } from "./pinia.plugin";
import { installRouter } from "./router.plugin";
import { installQuasar } from "./quasar.plugin";
import { installGlobals } from "./globals.plugin";

export function install(app: App) {
  installPinia(app);
  installRouter(app);
  installQuasar(app);
  installGlobals(app);
}
