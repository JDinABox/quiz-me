import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss()],
	build: {
		// generate .vite/manifest.json in outDir
		manifest: true,
		rollupOptions: {
			// overwrite default .html entry
			input: ['./web/main.ts', './web/style.css', './web/templates/showquiz/quiz.ts'],
			output: {
				dir: './web/dist'
			}
		},
		modulePreload: {
			polyfill: false
		}
	}
});
