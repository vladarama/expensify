import path from "path";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  optimizeDeps: {
    include: ["react-datepicker"],
  },

  server: {
    port: 3000,
  },

  build: {
    rollupOptions: {
      output: {
        manualChunks: {},
      },
    },
  },
});
