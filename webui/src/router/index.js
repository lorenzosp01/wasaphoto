import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/ProfileView.vue'
import LoginView from '../views/LoginView.vue'
import ProfileView from "@/views/ProfileView.vue";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/profile/:id', component: ProfileView},
		{path: '/login', component: LoginView},
		{path: '/link1', component: HomeView},
		{path: '/link2', component: HomeView},
		{path: '/some/:id/link', component: HomeView},
	]
})

export default router
