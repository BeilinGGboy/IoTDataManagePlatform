import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
    },
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/health': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: 'dist',
    assetsDir: 'assets',
    rollupOptions: {
      output: {
        // rolldown/vite 版本中 manualChunks 需要是函数
        // 返回 chunk 名称用于分包；返回 undefined 则使用默认分包策略
        manualChunks(id) {
          if (!id) return undefined
          if (id.includes('node_modules')) {
            if (id.includes('element-plus')) return 'element-plus'
            if (id.includes('vue-router') || id.includes('vue') || id.includes('pinia')) return 'vue-vendor'
          }
          return undefined
        },
      },
    },
  },
})
