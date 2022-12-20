<script setup>
import {inject, onMounted, ref} from "vue";
import Post from "@/components/Post.vue";

const error_msg = ref(null);
const axios = inject("axios");
const router = inject("router");
const userProfile = ref(null);

async function getUserProfile() {
	const offset = 0
	const amount = 10
	const route = router.currentRoute.value
	try {
		let response = await axios.get(route.path, {
			params: {
				offset: offset,
				amount: amount
			}
		})
		userProfile.value = response.data
		console.log(userProfile.value.user_info)
	} catch (e) {
		error_msg.value = e.toString();
	}
}

onMounted(() => {
	getUserProfile()
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
			<div class="row row-cols-1 row-cols-md-3 px-5 pt-5">
				<Post v-for="photo in userProfile.photos" :userId="userProfile.user_info.id" :photo="photo"/>
			</div>
		</div>
	</div>
</template>
