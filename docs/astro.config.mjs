import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	site: 'https://mistermoe.github.io',
	base: import.meta.env.DEV ? '/' : '/sbpfx',
	integrations: [
		starlight({
			title: 'sbpfx',
			customCss: [
				'@fontsource/lato/index.css',
				'@fontsource/lato/100.css',
				'@fontsource/lato/100-italic.css',
				'@fontsource/lato/300.css',
				'@fontsource/lato/300-italic.css',
				'@fontsource/lato/700.css',
				'@fontsource/lato/700-italic.css',
				'@fontsource/lato/900.css',
				'@fontsource/lato/900-italic.css',
				'@fontsource/lato/latin.css',
				'@fontsource/lato/latin-ext.css',
				'@fontsource/lato/latin-italic.css',
				'./src/styles/custom.css'
			],
			social: {
				github: 'https://github.com/mistermoe/sbpfx',
			},
			sidebar: [
				{ 
					label: 'Documentation',
					link: '/'
				}
			],
			lastUpdated: true
		}),
	],
});
