import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { configureAxiosInterceptors } from '@/service/axiosInterceptors'
import i18n from '@/plugins/i18n'

configureAxiosInterceptors()

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(i18n)

document.documentElement.lang = navigator.language.split('-')[0] ?? 'en'
app.mount('#app')
