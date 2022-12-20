<script setup>
const props = defineProps(["msg"]);

</script>

<template>
	<div class="alert alert-danger" role="alert">
		{{ props.msg }}
	</div>
</template>

<style>
</style>
