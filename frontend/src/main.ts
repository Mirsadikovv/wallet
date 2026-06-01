import "@quasar/extras/material-icons/material-icons.css";
import "quasar/dist/quasar.css";
import "./styles/utilities.css";
import { createApp } from "vue";
import App from "./App.vue";
import { install } from "./plugins/install";

const app = createApp(App);
install(app);
app.mount("#app");
