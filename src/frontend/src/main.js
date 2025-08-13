import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
if (!Array.isArray(window.pJSDom)) window.pJSDom = [];

createApp(App).use(router).mount('#app')
