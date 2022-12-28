import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import ProfileView from "@/views/ProfileView.vue";
import UploadPhotoView from "@/views/UploadPhotoView.vue";
import StreamView from "@/views/StreamView.vue";
import SearchUsersView from "@/views/SearchUsersView.vue";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{name: "Profile", path: '/profiles/:id', component: ProfileView},
		{name: "UploadPhoto", path: '/upload', component: UploadPhotoView},
		{name: "Login", path: '/login', component: LoginView},
		{name: "Stream", path: '/', component: StreamView},
		{name: "SearchUsers", path: '/search', component: SearchUsersView},
	]
})

router.beforeEach((to, from) => {
	if (to.name !== 'Login' && !localStorage.getItem('token')) {
		return {name: 'Login'}
	}

	if (to.name === 'Login' && localStorage.getItem('token')) {
		return false
	}
})

export default router
