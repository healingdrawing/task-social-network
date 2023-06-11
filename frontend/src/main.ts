import { createApp } from 'vue'
import App from './App.vue'
import './registerServiceWorker'
import router from './router'
import { pinia } from './store/pinia';

createApp(App).use(pinia).use(router).mount('#app')
