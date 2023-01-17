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
const username = ref("")
const isEditingName = ref(false)
const emit = defineEmits(["login"]);
const photos = ref([])

let amount = 10
let offset = 0
let wantsMorePhotos = true

async function getUserProfile() {
	axios.get(`/profiles/${userId.value}`, {
		params: {
			offset: offset,
			amount: amount
		}
	}).then((response) => {
		userProfile.value = response.data
		wantsMorePhotos = (userProfile.value.photos !== null)
		if (wantsMorePhotos) {
			let photosId = photos.value.map(photo => photo.id)
			userProfile.value.photos.forEach(photo => {
				if (!photosId.includes(photo.id)) {
					photos.value.push(photo)
				}
			})
		}
		username.value = userProfile.value.user_info.username
		error_msg.value = null
	}).catch((e) => {
		switch (e.response.status) {
			case 404:
				error_msg.value = "User not found"
				break
			case 403:
				error_msg.value = "You are not allowed to see this profile"
				break
			default:
				error_msg.value = e.toString()
		}
	})

}

async function getUserFollowed() {
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

async function banUser() {
	axios.put(`/profiles/${token}/ban/${userId.value}`)
		.then(() => {
			bannedList.value.push(parseInt(userId.value))
		})
		.catch((e) => {
			error_msg.value = e.toString();
		})
}

async function unbanUser() {
	axios.delete(`/profiles/${token}/ban/${userId.value}`)
		.then(() => {
			bannedList.value = bannedList.value.filter((id) => {
				return id !== parseInt(userId.value)
			})
		})
		.catch((e) => {
			error_msg.value = e.toString();
		})
}

async function followUser() {
	axios.put(`/profiles/${token}/following/${userId.value}`)
		.then(() => {
			userProfile.value.profileInfo.followersCounter++
			followersList.value.push(parseInt(userId.value))
		})
		.catch((e) => {
			error_msg.value = e.toString();
		})
}

async function unfollowUser() {
	axios.delete(`/profiles/${token}/following/${userId.value}`)
		.then(() => {
			userProfile.value.profileInfo.followersCounter--
			followersList.value = followersList.value.filter((id) => {
				return id !== parseInt(userId.value)
			})
		})
		.catch((e) => {
			error_msg.value = e.toString();
		})
}


async function editName() {
	if (username.value !== userProfile.value.user_info.username) {
		axios.put(`/profiles/${token}/name`, {
			username: username.value
		}).then(() => {
			getUserProfile()
			isEditingName.value = false
		}).catch((e) => {
			if (e.response.status === 409) {
				error_msg.value = "Someone else is using this username"
			} else {
				error_msg.value = e.toString();
			}
		})
	}
}

const deletePhoto = (id) => {
	userProfile.value.profileInfo.photosCounter--
	photos.value = photos.value.filter((photo) => {
		return photo.id !== id
	})
}

const getMorePhotos =  (e) => {
	if (window.scrollY + window.innerHeight >= document.body.scrollHeight && wantsMorePhotos) {
		offset += amount
		getUserProfile()
	}
}


onBeforeRouteUpdate((to, from) => {
	userId.value = to.params.id
	photos.value = []
	offset = 0
	amount = 10
	getUserProfile()
	getUserFollowed()
	getUserBannedList()
})

onMounted(() => {
	userId.value = router.currentRoute.value.params.id
	window.addEventListener("scroll", getMorePhotos)
	getUserProfile()
	getUserFollowed()
	getUserBannedList()
})

</script>
<template>
	<ErrorMsg v-if="error_msg" :msg="error_msg"></ErrorMsg>
	<div v-if="userProfile && !error_msg" class="h-100">
		<div class="d-flex justify-content-center align-items-center">
			<div class="h1 align-items-center" v-if="!isEditingName">{{ username }}</div>
			<input v-else type="text" class="form-control w-25" v-model="username">
			<div v-if="isEditingName" @click="editName" class="btn btn-sm btn-primary mx-2 align-text-bottom">Conferma</div>
			<div v-if="token === userId" class="mx-2" @click="isEditingName=!isEditingName">
				<svg class="feather">
					<use href="/feather-sprite-v4.29.0.svg#edit"/>
				</svg>
			</div>
			<div v-if="token !== userId">
				<div v-if="followersList.includes(parseInt(userId))" class="btn btn-danger mx-2" @click="unfollowUser">
					Unfollow
				</div>
				<div v-else class="btn btn-primary mx-2" @click="followUser">Follow</div>
				<div v-if="bannedList.includes(parseInt(userId))" class="btn btn-danger ms-2" @click="unbanUser">Unban</div>
				<div v-else class="btn btn-primary ms-2" @click="banUser">Ban</div>
			</div>
		</div>
		<hr>
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

			<div class="row row-cols-1 row-cols-md-3 px-5 pt-5">
				<Post v-for="photo in photos" :key="photo.id" @delete-photo="deletePhoto(photo.id)"
					  :userId="userProfile.user_info.id" :photo="photo"/>
			</div>
		</div>
	</div>
</template>
