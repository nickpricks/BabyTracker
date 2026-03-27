import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { VitePWA } from "vite-plugin-pwa";

export default defineConfig({
  plugins: [
    react(),
    VitePWA({
      registerType: "autoUpdate",
      manifest: {
        short_name: "BabyTracker",
        name: "Baby Tracker",
        description:
          "Track your baby's feeds, sleep, growth, and diaper changes",
        icons: [
          { src: "icon-192.png", type: "image/png", sizes: "192x192" },
          { src: "icon-512.png", type: "image/png", sizes: "512x512" },
        ],
        start_url: ".",
        display: "standalone",
        theme_color: "#4a90d9",
        background_color: "#ffffff",
      },
    }),
  ],
  server: {
    port: 3000,
  },
  build: {
    outDir: "build",
  },
  test: {
    environment: "jsdom",
    setupFiles: "./src/setupTests.js",
    globals: true,
  },
});
