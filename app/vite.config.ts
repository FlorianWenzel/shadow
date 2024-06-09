import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import VueDevTools from 'vite-plugin-vue-devtools'

export default defineConfig({
  server: {
    port: 8080,
    proxy: {
      '/api': {
        target: process.env.BACKEND_URL || 'http://localhost:3000',
        changeOrigin: true
      }
    }
  },
  plugins: [vue(), VueDevTools()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
