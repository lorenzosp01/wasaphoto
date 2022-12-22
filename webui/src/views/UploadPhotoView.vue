<script setup>
import ErrorMsg from "@/components/ErrorMsg.vue";
import {inject, ref} from "vue";

const error_msg = ref(null);
const axios = inject("axios")
const photo = ref(null)
const token = localStorage.getItem("token")
const router = inject("router")

async function uploadPhoto() {
	if (photo.value) {
		axios.post(`/profiles/${token}/photos/`, photo.value, {
			headers: {
				"Content-Type":"image/png",
			},
		}).then(() => {
			error_msg.value = null;
			router.push(`/profiles/${token}`);
		}).catch((e) => {
			error_msg.value = e.toString();
		})
	} else {
		error_msg.value = "Please select a photo";
	}
}

</script>

<template>
	<div>
		<h1 class="h1">Upload Photo</h1>
		<hr>
		<ErrorMsg v-if="error_msg" :msg="error_msg"></ErrorMsg>
		<div class="mb-3">
			<label for="photo" class="form-label">Photo</label>
			<input type="file" class="form-control" @change="(event) => photo = event.target.files[0]" aria-describedby="photoHelp">
			<div id="photoHelp" class="form-text">Upload a photo</div>
		</div>
		<div class="btn btn-primary" @click="uploadPhoto">Upload</div>
	</div>
</template>
