import { createApp } from 'vue'
import { createPinia } from 'pinia'

import { createI18n } from 'vue-i18n'

import ForbidenApp from './ForbidenApp.vue'

import './assets/index.css'
import en from './locales/en.json'
import tr from './locales/tr.json'
import nl from './locales/nl.json'

const i18n = createI18n({
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: en,
    tr: tr,
    nl: nl
  }
})

const app = createApp(ForbidenApp)

app.use(createPinia())
app.use(i18n)
app.mount('#app')
