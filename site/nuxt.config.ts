import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  ssr: true,

  app: {
    baseURL: process.env.NUXT_APP_BASE_URL || '/sleepoff/',
    head: {
      htmlAttrs: {
        lang: 'pt-BR',
        class: 'dark'
      },
      title: 'sleepoff — Sleep Timer Inteligente para Windows e Multiplataforma',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'CLI moderno, ultra-leve e seguro para agendar o desligamento do seu PC. Zero bloat, modo dry-run e proteção contra travamentos.' },
        { name: 'theme-color', content: '#0f172a' }
      ],
      link: [
        { rel: 'icon', type: 'image/svg+xml', href: (process.env.NUXT_APP_BASE_URL || '/sleepoff/').replace(/\/$/, '') + '/favicon.svg' }
      ]
    }
  },

  css: ['~/assets/css/main.css'],

  modules: [
    'shadcn-nuxt',
    'motion-v/nuxt',
    '@nuxt/icon'
  ],

  shadcn: {
    prefix: '',
    componentDir: './components/ui'
  },

  nitro: {
    preset: 'github-pages',
    prerender: {
      crawlLinks: true,
      routes: ['/']
    }
  },

  vite: {
    plugins: [
      tailwindcss()
    ]
  }
})
