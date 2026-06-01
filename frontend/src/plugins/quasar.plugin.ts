import type { App } from "vue";
import { Quasar, Notify, Dialog, Loading } from "quasar";
import iconSet from "quasar/icon-set/material-icons";

export function installQuasar(app: App) {
  app.use(Quasar, {
    plugins: { Notify, Dialog, Loading },
    iconSet,
    config: {
      notify: {
        position: "top",
        timeout: 3000,
      },
    },
  });
}
