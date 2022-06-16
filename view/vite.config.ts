import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue';

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [vue()],
    build: {
        outDir: "../public",
        emptyOutDir: true,
    },
    server: {
        port: 80,
        proxy: {
            "/api/": {
                target: "http://localhost:3001",
                changeOrigin: true,
            }
        }
    }
})
