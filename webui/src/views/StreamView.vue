<script setup>

import Post from "@/components/Post.vue";
import {inject, onBeforeUnmount, onMounted, ref} from "vue";

const axios = inject("axios");
const router = inject("router");
const token = localStorage.getItem("token");
const stream = ref([]);
const error_msg = ref(null);

let offset = 0
let amount = 10
let wantsMorePhotos = true

async function getStream() {
	axios.get(`/stream/${token}`, {
		params: {
			offset: offset,
			amount: amount
		}
	}).then((response) => {
		for (let i = 0; i < response.data.photos.length; i++) {
			stream.value.push(response.data.photos[i])
		}
	}).catch((e) => {
		if (e.response.status !== 404) {
			error_msg.value = e.response.data
		} else {
			wantsMorePhotos = false
		}
	})
}

const getMorePhotos =  (e) => {
	if (window.scrollY + window.innerHeight >= document.body.scrollHeight && wantsMorePhotos) {
		offset += amount
		getStream()
	}
}

onMounted(() => {
	window.addEventListener("scroll", getMorePhotos)
	getStream()
})

onBeforeUnmount(() => {
	window.removeEventListener("scroll", getMorePhotos)
})

</script>

<template>
	<div class="w-100">
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Stream</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
			</div>
		</div>
		<ErrorMsg v-if="error_msg" :msg="error_msg"></ErrorMsg>
		<div class="w-auto row justify-content-center ">
			<div v-if="stream.length > 0" class="flex flex-column justify-content-center pt-5" style="width: 30%;">
				<Post v-for="photo in stream" :key="photo.id" :photo="photo" :showOwner="true" :userId="photo.owner.id"></Post>
			</div>
			<div v-else class="d-flex flex-column">
				No photo to display
			</div>
		</div>

	</div>
</template>


