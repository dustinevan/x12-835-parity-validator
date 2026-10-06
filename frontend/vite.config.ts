import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// The validator app (index.html → src/app), built into the directory the Go
// binary embeds. The local docs site has its own config: vite.docs.config.ts.
export default defineConfig({
  plugins: [react(), tailwindcss()],
  build: {
    outDir: '../backend/app/site/dist',
    emptyOutDir: true,
  },
  server: {
    port: 5179,
    strictPort: true,
    proxy: {
      '/api': {
        target: process.env.VALIDATOR_API_URL ?? 'http://127.0.0.1:8089',
      },
    },
  },
})
