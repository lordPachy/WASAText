import {createApp} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'
<<<<<<< HEAD
=======
import { createPinia } from 'pinia';
import piniaPersist from "pinia-plugin-persist"
>>>>>>> working

import './assets/dashboard.css'
import './assets/main.css'

const pinia = createPinia();
pinia.use(piniaPersist);
const app = createApp(App);
app.config.globalProperties.$axios = axios;
app.component("ErrorMsg", ErrorMsg);
<<<<<<< HEAD
app.use(router)
app.mount('#app')
=======
app.use(router);
app.use(pinia);
app.mount('#app');
>>>>>>> working
