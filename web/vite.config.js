import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import vueJsx from "@vitejs/plugin-vue-jsx";
import UnoCSS from "unocss/vite";

import AutoImport from "unplugin-auto-import/vite";
import Components from "unplugin-vue-components/vite";
import { NaiveUiResolver } from "unplugin-vue-components/resolvers";
import { VantResolver } from "unplugin-vue-components/resolvers";
import { fileURLToPath } from "url";

// const debugMode = process.env.APP_ENV !== 'prod'

export default defineConfig({
  base: "./",
  build: {
    outDir: "../",
  },
  plugins: [
    vue(),
    vueJsx(),
    AutoImport({
      imports: [
        "vue",
        {
          "naive-ui": [
            "useDialog",
            "useMessage",
            "useNotification",
            "useLoadingBar",
          ],
        },
      ],
    }),
    Components({
      resolvers: [NaiveUiResolver(), VantResolver()],
    }),
    UnoCSS(),
  ],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  server: {
    proxy: {
      "/api": {
        // target: "http://127.0.0.1:8080",
        target: "http://175.155.64.58:8080",
        changeOrigin: true,
        // rewrite: (path) => path.replace(/^\/api/, ""),
      },
    },
  },
});
