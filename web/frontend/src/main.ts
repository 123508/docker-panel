import { createApp } from 'vue'
import App from './App.vue'

import '@/styles/base.css'
import '@/styles/theme.css'
import '@/store/ui'
import router from './router'


createApp(App).use(router).mount('#app')
