import { defineConfig, type Plugin } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import { fileURLToPath } from 'node:url'

// The local docs site: viewers, the compare page, compared fields, file formats, and the reference.
// It shares src/components, src/data, and index.css with the validator app,
// runs on its own port, and never enters the app build or the Go binary.
//
// Usage: npm run docs (or make docs from the repository root).

// Route navigations (extensionless paths such as /835-viewer) get the docs
// shell, docs/index.html, instead of the app's index.html. Files, Vite
// internals, and /api pass through.
function docsHtmlPlugin(): Plugin {
  return {
    name: 'docs-html',
    configureServer(server) {
      server.middlewares.use((req, _res, next) => {
        const url = req.url ?? ''
        const path = url.split('?')[0]
        const isFile = path.includes('.')
        const isViteInternal = url.startsWith('/@') || url.startsWith('/__')
        const isApi = url.startsWith('/api')
        if (!isFile && !isViteInternal && !isApi) req.url = '/docs/index.html'
        next()
      })
    },
  }
}

export default defineConfig({
  plugins: [docsHtmlPlugin(), react(), tailwindcss()],
  server: {
    port: 5180,
    strictPort: true,
    proxy: {
      '/api': {
        target: process.env.VALIDATOR_API_URL ?? 'http://127.0.0.1:8089',
      },
    },
  },
  // A local build for checking only; it is never embedded in the Go binary.
  build: {
    outDir: 'dist-docs',
    emptyOutDir: true,
    rolldownOptions: {
      input: fileURLToPath(new URL('docs/index.html', import.meta.url)),
    },
  },
})
