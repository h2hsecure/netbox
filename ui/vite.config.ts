import { defineConfig } from 'vite'
import { resolve } from 'path'
import vue from '@vitejs/plugin-vue'
import viteCompress from 'vite-plugin-compression'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue({
    template: {
      compilerOptions: {
        isCustomElement: (tag) => ['altcha-widget'].includes(tag),
      }
    }
  }), viteCompress(),
  {
    name: 'build-html',
    apply: 'build',
    transformIndexHtml: (html) => {
      return {
        html,
        tags: [
          {
            tag: 'script',
            attrs: {
              src: '/ddos/ntb_dds',
            },
            injectTo: 'head',
          },
          {
            tag: 'title',
            var: '...Just wait a while',
            injectTo: 'head',
          }
        ],
      };
    },
  },],
  base: '/ddos/app/',
  build: {
    minify: true,
    rollupOptions: {
      input: {
        index: resolve(__dirname, 'index.html'),
        forbiden: resolve(__dirname, 'forbiden.html'),
        queue: resolve(__dirname, 'queue.html'),
            }
    }
  },
  esbuild: {
    legalComments: 'none'
  }
})
