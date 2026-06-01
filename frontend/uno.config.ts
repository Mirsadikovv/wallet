import { defineConfig, presetUno, presetAttributify } from "unocss";

export default defineConfig({
  presets: [
    presetUno(),
    presetAttributify(),
  ],
  shortcuts: {
    "my-row": "q-col-gutter-md",
  },
  theme: {
    colors: {
      primary: "#4A90D9",
      secondary: "#26A69A",
      accent: "#9C27B0",
    },
  },
  safelist: ["text-center", "text-left", "text-right", "font-bold"],
});
