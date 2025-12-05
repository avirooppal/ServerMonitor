import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': {
        target: 'http://107.150.20.37:8080',
        changeOrigin: true,
        secure: false,
      },
      '/downloads': {
        target: 'http://107.150.20.37:8080',
        changeOrigin: true,
        secure: false,
      }
    }
  }
})
