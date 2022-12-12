<template>
	   <div>
			<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			  <h1 class="h2">Login</h1>
			  <div class="btn-toolbar mb-2 mb-md-0">
			  </div>
			</div>
		<div class="row">
		  <div class="col-md-6">
			<form @submit.prevent="login">
			  <div class="form-group">
				<label for="username">Username</label>
				<input type="text" class="form-control" id="username" v-model="username" placeholder="Enter username">
			  	<button type="submit" class="btn btn-primary">Login</button>
			  </div>
			</form>
		  </div>
		</div>
	  </div>
</template>
<script>
export default {
	name: "Login"
}
</script>

<style scoped>

</style>
