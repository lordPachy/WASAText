import {createRouter, createWebHashHistory} from 'vue-router'
import Authentication from '../views/Authentication.vue'
import Conversations from '../views/Conversations.vue'
import Conversation from '../views/Conversation.vue'
import Homepage from '../views/Homepage.vue'
import Settings from '../views/Settings.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{name: 'authentication', path: '/', component: Authentication},
		{name: 'homepage', path: '/session', component: Homepage, props: true},
		{name: 'conversations', path: '/conversations', component: Conversations},
		{name: 'settings', path: '/settings', component: Settings, props: true},
		{name: 'conversation', path: '/conversations/:conversationid', component: Conversation, props: true},
	],
	id: "",
	username: ""
})

export default router
  

