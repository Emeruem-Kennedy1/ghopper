import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react";
import { VitePWA } from 'vite-plugin-pwa';

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "");

  // Use environment variable with fallback
  const backendUrl =
    env.VITE_BACKEND_URL ||
    (mode === "production"
      ? "https://ghopper-backend.kemeruem.com"
      : "http://backend:9797");

  console.log(`Mode: ${mode}, Backend URL: ${backendUrl}`); // For debugging

  return {
    plugins: [
    react(),
    VitePWA({
      registerType: 'autoUpdate',
      injectRegister: 'auto',
      strategies: 'injectManifest',
      srcDir: 'public',
      filename: 'service-worker.js',
      includeAssets: ['favicon.ico', 'app-icon.svg', 'robots.txt'],
      manifest: false // We're using our own manifest.json
    })
    ],
    server: {
      host: "0.0.0.0",
      port: 3000,
      proxy: {
        "/api": {
          target: backendUrl,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api/, ""),
        },
      },
    },
    define: {
      "process.env": {},
    },
  };
});
