import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 3000,
    proxy: {
      '/type': 'http://localhost:8080',
      '/function': 'http://localhost:8080',
      '/product': 'http://localhost:8080',
      '/topo': 'http://localhost:8080',
      '/analyze': 'http://localhost:8080',
      '/suggest': 'http://localhost:8080',
      '/products': 'http://localhost:8080',
      '/functions': 'http://localhost:8080'
    }
  }
})
