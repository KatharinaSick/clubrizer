import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import GoogleSignInPlugin from 'vue3-google-signin'
import { configureAxiosInterceptors } from '@/service/axiosInterceptors'
import i18n from '@/plugins/i18n'

configureAxiosInterceptors()

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(i18n)

app.use(GoogleSignInPlugin, {
  clientId: import.meta.env.VITE_OAUTH_GOOGLE_CLIENT_ID
})

app.mount('#app')
