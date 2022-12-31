<script setup>
import {inject, ref} from "vue";

const username = ref("");
const error_msg = ref(null);
const axios = inject("axios");
const router = inject("router");
const emit = defineEmits(["login"])

async function login() {
	if (username.value) {
		try {
			let response = await axios.post("/session", {
				username: username.value,
			})
			localStorage.clear()
			localStorage.setItem('token', response.data.identifier);
			axios.defaults.headers.common["Authorization"] = `Bearer ${localStorage.getItem('token')}`;
			emit("login")
			if (response.data.identifier !== 0) {
				router.push(`/profiles/${response.data.identifier}`);
			}
		} catch (e) {
			error_msg.value = e.toString();
		}
	} else {
		error_msg.value = "Please enter a username";
	}
}

</script>
<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Login</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
			</div>
		</div>
		<ErrorMsg v-if="error_msg" :msg="error_msg"></ErrorMsg>
		<div class="row">
			<div class="col-md-6">
				<form @submit.prevent="login">
					<div class="form-group">
						<label for="username">Username</label>
						<input type="text" class="form-control" id="username" v-model="username"
							   placeholder="Enter username">
					</div>
					<br>
					<button type="submit" class="btn btn-primary">Login</button>
				</form>
			</div>
		</div>
	</div>
</template>
