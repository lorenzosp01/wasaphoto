<script setup>
import {inject, onMounted, ref} from "vue";

const axios = inject("axios");
const router = inject("router");
const token = localStorage.getItem("token");
const users = ref([]);
const error_msg = ref(null);
const pattern = ref("")

async function searchUsers() {
	if (pattern.value.length > 0) {
		axios.get(`/search`, {
			params: {
				pattern: pattern.value
			}
		}).then((response) => {
			users.value = response.data.users
		}).catch((e) => {
			if (e.response.status !== 404) {
				error_msg.value = e.toString();
			} else {
				users.value = []
			}
		})
	} else {
		users.value = []
	}
}
</script>

<template>
	<div class="w-100">
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Search users</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
			</div>
		</div>
		<ErrorMsg v-if="error_msg" :msg="error_msg"></ErrorMsg>
		<div class="d-flex justify-content-center w-100">
			<div class="form-outline mx-2" style="width: 26%; ">
				<input type="text" v-model="pattern" @input="searchUsers" class="form-control"  placeholder="Search users"/>
			</div>
			<div class="btn btn-primary" @click="searchUsers">
				<svg class="feather">
					<use href="/feather-sprite-v4.29.0.svg#search"/>
				</svg>
			</div>
		</div>
		<div class="w-auto row justify-content-center">
			<div v-if="users" class="flex flex-column justify-content-center pt-5" style="width: 30%;">
				<div v-for="user in users" :key="user.id" class="col mb-5">
					<ErrorMsg v-if="error_msg" :msg="error_msg"></ErrorMsg>
					<div class="card">
						<div class="card-body">
							<div class="mb-2">
								<h5 class="card-title">{{ user.username }}</h5>
							</div>
							<div>
								<RouterLink :to="{name: 'Profile', params: {id: user.id}}" >
									<svg class="feather">
										<use href="/feather-sprite-v4.29.0.svg#user"/>
									</svg>
									Profile
								</RouterLink>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>
