<script setup>
import {RouterLink, RouterView} from 'vue-router'
import {inject} from "vue";

const route = inject("route");

const logout = () => {
	localStorage.removeItem('token');
}

</script>

<template>
	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">Example App</a>
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse"
				data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false"
				aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
	</header>

	<div class="container-fluid">
		<div class="row">
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>General</span>
					</h6>
					<ul class="nav flex-column" v-if="$route.name !== 'Login'">
						<li class="nav-item">
							<RouterLink to="/login" @click='logout' class="nav-link">
								<svg class="feather">
									<use href="/feather-sprite-v4.29.0.svg#home"/>
								</svg>
								Logout
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/link1" class="nav-link">
								<svg class="feather">
									<use href="/feather-sprite-v4.29.0.svg#layout"/>
								</svg>
								Menu item 1
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/link2" class="nav-link">
								<svg class="feather">
									<use href="/feather-sprite-v4.29.0.svg#key"/>
								</svg>
								Menu item 2
							</RouterLink>
						</li>
					</ul>
				</div>
			</nav>

			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<br>
				<RouterView/>
			</main>
		</div>
	</div>
</template>

<style>
</style>
