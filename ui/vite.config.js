import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import { resolve } from 'path'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue({
    template: {
      compilerOptions: {
        isCustomElement: (tag) => ['altcha-widget'].includes(tag),
      }
    }
  })],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  base: '/ddos/app/',
  name: '...just wait a while',
  build: {
    rollupOptions: {
      input: {
        main: resolve(__dirname, 'index.html'),
        forbiden: resolve(__dirname, 'forbiden.html')
      }
    }
  }
})
