import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import Conversations from '../views/Conversations.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: HomeView},
		{path: '/session', component: HomeView},
		{path: '/conversations', component: Conversations},
		{path: '/settings/username', component: HomeView},
	]
})

export default router
