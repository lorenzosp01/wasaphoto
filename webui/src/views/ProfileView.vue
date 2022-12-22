<script setup>
import {inject, onMounted, ref} from "vue";
import Post from "@/components/Post.vue";
import {onBeforeRouteUpdate} from "vue-router";

const error_msg = ref(null);
const axios = inject("axios");
const router = inject("router");
const userProfile = ref(null);
const userId = ref(null);
const token = localStorage.getItem("token");
const followersList = ref([]);
const bannedList = ref([]);

async function getUserProfile() {
	const offset = 0
	const amount = 10
	try {
		let response = await axios.get(`/profiles/${userId.value}`, {
			params: {
				offset: offset,
				amount: amount
			}
		})
		error_msg.value = null
		userProfile.value = response.data
	} catch (e) {
		error_msg.value = e.toString();
	}
}

async function getUserFollowers() {
	axios.get(`/profiles/${token}/following/`)
		.then((response) => {
			followersList.value = response.data.users.map((user) => {
				return user.id
			})
		})
		.catch((e) => {
			if (e.response.status !== 404) {
				error_msg.value = e.toString();
			} else {
				followersList.value = []
			}
		})
}

async function getUserBannedList() {
	axios.get(`/profiles/${token}/ban/`)
		.then((response) => {
			bannedList.value = response.data.users.map((user) => {
				return user.id
			})
		})
		.catch((e) => {
			if (e.response.status !== 404) {
				error_msg.value = e.toString();
			} else {
				bannedList.value = []
			}
		})
}

async function banUser () {
	axios.put(`/profiles/${token}/ban/${userId.value}`)
		.then(() => {
			getUserBannedList()
		})
		.catch((e) => {
			error_msg.value = e.toString();
		})
}

async function unbanUser () {
	axios.delete(`/profiles/${token}/ban/${userId.value}`)
		.then(() => {
			getUserBannedList()
		})
		.catch((e) => {
			error_msg.value = e.toString();
		})
}

async function followUser () {
	axios.put(`/profiles/${token}/following/${userId.value}`)
		.then(() => {
			getUserFollowers()
		})
		.catch((e) => {
			error_msg.value = e.toString();
		})
}

async function unfollowUser () {
	axios.delete(`/profiles/${token}/following/${userId.value}`)
		.then(() => {
			getUserFollowers()
		})
		.catch((e) => {
			error_msg.value = e.toString();
		})
}

onBeforeRouteUpdate((to, from) => {
	userId.value = to.params.id
	getUserProfile()
	getUserFollowers()
	getUserBannedList()
})

onMounted(() => {
	userId.value = router.currentRoute.value.params.id
	getUserProfile()
	getUserFollowers()
	getUserBannedList()
})

</script>
<template>
	<LoadingSpinner v-if="!userProfile"></LoadingSpinner>
	<div v-if="userProfile">
		<div class="text-center">
			<h1 class="h1">{{ userProfile.user_info.username }}'s profile</h1>
			<hr>
		</div>
		<ErrorMsg v-if="error_msg" :msg="error_msg"></ErrorMsg>
		<div class="pt-3">

			<div class="d-flex justify-content-around">
				<div class="">
					<h1 class="h4"> Photos </h1>
					<div class="fs-5 text-center">{{ userProfile.profileInfo.photosCounter }}</div>
				</div>
				<div>
					<h1 class="h4"> Following </h1>
					<div class="fs-5 text-center">{{ userProfile.profileInfo.followingCounter }}</div>
				</div>
				<div>
					<h1 class="h4"> Followers </h1>
					<div class="fs-5 text-center">{{ userProfile.profileInfo.followersCounter }}</div>
				</div>
			</div>
			<div v-if="token !== userId " class="w-100 d-flex justify-content-center">
				<div v-if="followersList.includes(parseInt(userId))" class="btn btn-danger" @click="unfollowUser">Unfollow</div>
				<div v-else class="btn btn-primary" @click="followUser">Follow</div>
				<div v-if="bannedList.includes(parseInt(userId))" class="btn btn-danger" @click="unbanUser">Unban</div>
				<div v-else class="btn btn-primary " @click="banUser">Ban</div>
			</div>
			<div class="row row-cols-1 row-cols-md-3 px-5 pt-5">
				<Post v-for="photo in userProfile.photos" @delete-photo="getUserProfile" :userId="userProfile.user_info.id" :photo="photo"/>
			</div>
		</div>
	</div>
</template>
