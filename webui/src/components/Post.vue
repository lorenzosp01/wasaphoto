<script setup>
import {inject, onMounted, ref} from "vue";

const props = defineProps({
	photo: {
		type: Object,
	},
	userId: {
		type: Number,
	},
	showOwner: {
		type: Boolean,
	},
});

const tempPhoto = ref(props.photo);
const emit = defineEmits(["delete-photo"]);
const axios = inject("axios");
const imgUrl = ref(null);
const error_msg = ref(null);
const token = localStorage.getItem("token");
const photoComments = ref([]);
const newComment = ref("");
const showComments = ref(false);

async function getPhoto() {
	try {
		console.log(props.photo)
		let response = await axios.get(`/profiles/${props.userId}/photos/${props.photo.id}`, {responseType: 'blob'})
		imgUrl.value = URL.createObjectURL(response.data)
	} catch (e) {
		error_msg.value = e.toString();
	}
}

async function deletePhoto() {
	axios.delete(`/profiles/${props.userId}/photos/${props.photo.id}`)
		.then(() => {
			emit("delete-photo")
		})
		.catch((e) => {
			error_msg.value = e.toString();
		})
}

async function getPhotoComments() {
	try {
		let response = await axios.get(`/profiles/${props.userId}/photos/${props.photo.id}/comments/`)
		photoComments.value = response.data.comments
	} catch (e) {
		error_msg.value = e.toString();
	}
}

async function likePhoto() {
	axios.put(`/profiles/${props.userId}/photos/${props.photo.id}/likes/${token}`)
		.then(() => {
			tempPhoto.value.photoInfo.likesCounter += 1
		})
		.catch((e) => {
			if (e.response.status === 409) {
				axios.delete(`/profiles/${props.userId}/photos/${props.photo.id}/likes/${token}`)
					.then(() => {
						tempPhoto.value.photoInfo.likesCounter -= 1
					})
					.catch((e) => {
						error_msg.value = e.toString();
					})
			} else {
				error_msg.value = e.toString();
			}
		})
}


async function deleteComment(commentId) {
	axios.delete(`/profiles/${props.userId}/photos/${props.photo.id}/comments/${commentId}`)
		.then(() => {
			getPhotoComments().then(() => {
				tempPhoto.value.photoInfo.commentsCounter -= 1
			}).catch((e) => {
				error_msg.value = e.toString();
			})
		})
		.catch((e) => {
			error_msg.value = e.toString();
		})
}

async function commentPhoto() {
	if (newComment.value.length > 0) {
		axios.post(`/profiles/${props.userId}/photos/${props.photo.id}/comments/`, {
			content: newComment.value
		}).then(() => {
			getPhotoComments().then(() => {
				tempPhoto.value.photoInfo.commentsCounter += 1
				newComment.value = ""
			}).catch(e => {
				error_msg.value = e.toString();
			})
		}).catch(e => {
			error_msg.value = e.toString();
		})
	} else {
		error_msg.value = "Comment cannot be empty"
	}
}

onMounted(() => {
	console.log(parseInt(token), )
	getPhoto()
	getPhotoComments()
})

</script>

<template>
	<div v-if="imgUrl" class="col mb-5">
		<ErrorMsg v-if="error_msg" :msg="error_msg"></ErrorMsg>
		<div class="card">
			<div v-if="showOwner" class="card-header">
				<RouterLink :to="`/profiles/${tempPhoto.owner.id}`">
					<div class="fw-bold">{{tempPhoto.owner.username}}</div>
				</RouterLink>
			</div>
			<img :src="imgUrl" class="card-img-top" alt="...">
			<div class="card-body">
				<div class="btn" @click="likePhoto">
					<svg class="feather">
						<use href="/feather-sprite-v4.29.0.svg#thumbs-up"/>
					</svg>
					{{ tempPhoto.photoInfo.likesCounter }}
				</div>
				<div class="btn" @click.prevent="() => showComments = !showComments">
					<svg class="feather">
						<use href="/feather-sprite-v4.29.0.svg#message-square"/>
					</svg>
					{{ tempPhoto.photoInfo.commentsCounter }}
				</div>
				<div v-if="showComments" class="overflow-auto pt-3"
					 style="max-height: 30vh">
					<div v-for='comment in photoComments' class="border border-secondary bg-white rounded-1 mb-2 mx-2 px-2 py-1">
						<div class="d-flex justify-content-between">
							<h6> {{ comment.owner.username }} </h6>
							<small class="text-muted">{{ comment.uploadedAt }}</small>
						</div>
						<p class="card-text">{{ comment.content }}</p>
						<div v-if="parseInt(token) === comment.owner.id" class="btn btn-sm btn-danger" @click="deleteComment(comment.id)">
							<svg class="feather">
								<use href="/feather-sprite-v4.29.0.svg#trash-2"/>
							</svg>
						</div>
					</div>
				</div>
				<div class="mb-2">
					<h5 class="card-title">{{ props.username }}</h5>
					<textarea class="form-control" v-model="newComment" rows="3"></textarea>
					<button class="btn btn-sm btn-primary mt-2" @click="commentPhoto">Add comment</button>
				</div>

			</div>
			<div class="card-footer d-flex justify-content-between">
				<div class="row align-items-center">
					<div>
						Uploaded time: {{ props.photo.uploadedAt }}
					</div>

				</div>
				<div v-if="parseInt(token) === props.photo.owner.id" class="btn btn-sm btn-danger" @click="deletePhoto">
					<svg class="feather">
						<use href="/feather-sprite-v4.29.0.svg#trash-2"/>
					</svg>
				</div>
			</div>
		</div>
	</div>
</template>



