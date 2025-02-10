import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import GoogleSignInPlugin from 'vue3-google-signin'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.use(GoogleSignInPlugin, {
  // TODO secret? pass as env variable
  clientId: '281221302650-5kinlm2k9b8rllcb4blh72lcj3qe0jbj'
})

app.mount('#app')
