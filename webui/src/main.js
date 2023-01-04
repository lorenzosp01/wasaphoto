import {createApp} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'

import './assets/dashboard.css'
import './assets/main.css'

const app = createApp(App)
app.config.globalProperties.$axios = axios;
app.provide('axios', app.config.globalProperties.$axios);
app.component("ErrorMsg", ErrorMsg);
app.use(router)
app.provide('router', app.config.globalProperties.$router)
app.mount('#app')

export default app;
