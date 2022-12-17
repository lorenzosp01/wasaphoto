<script setup>
import {inject, ref} from "vue";
import LoginView from "@/views/LoginView.vue";

const error_msg = ref(null);

async function getMyProfile() {
}

</script>
<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Profile</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
			</div>
		</div>
		<ErrorMsg v-if="error_msg" :msg="error_msg"></ErrorMsg>
		<div class="row">
			<div class="col-md-6">
			</div>
		</div>
	</div>
</template>
