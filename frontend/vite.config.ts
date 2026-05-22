import { sveltekit } from "@sveltejs/kit/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [sveltekit(), tailwindcss()],
  server: {
    host: true,
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
        cookieDomainRewrite: "localhost",
      },
      "/me": {
        target: "http://localhost:8080",
        changeOrigin: true,
        cookieDomainRewrite: "localhost",
      },
    },
  },
});
