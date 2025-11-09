<script lang="ts">
	import { goto } from '$app/navigation';

	export let data: { stats: { userCount: number; siteCount: number; uploadSessionCount: number } };

	async function handleLogout() {
		await fetch('/api/auth/logout', { method: 'POST' });
		await goto('/login');
	}
</script>

<div class="min-h-screen bg-base-100">
	<!-- Navigation -->
	<div class="navbar border-b border-base-300 bg-base-100 shadow-sm">
		<div class="flex-1">
			<div class="flex items-center gap-3">
				<div class="avatar">
					<div
						class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-primary to-secondary"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-5 w-5 text-white"
							viewBox="0 0 20 20"
							fill="currentColor"
						>
							<path
								fill-rule="evenodd"
								d="M4 4a2 2 0 012-2h8a2 2 0 012 2v12a1 1 0 110 2h-3a1 1 0 01-1-1v-2a1 1 0 00-1-1H9a1 1 0 00-1 1v2a1 1 0 01-1 1H4a1 1 0 110-2V4zm3 1h2v2H7V5zm2 4H7v2h2V9zm2-4h2v2h-2V5zm2 4h-2v2h2V9z"
								clip-rule="evenodd"
							/>
						</svg>
					</div>
				</div>
				<div class="flex items-center gap-2">
					<span class="text-xl font-bold text-primary">Hostamat</span>
					<span class="hidden font-medium text-base-content/70 sm:inline">Dashboard</span>
				</div>
			</div>
		</div>
		<div class="flex-none">
			<button on:click={handleLogout} class="btn btn-outline btn-sm">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-4 w-4"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
					/>
				</svg>
				Logout
			</button>
		</div>
	</div>

	<!-- Main Content -->
	<div class="container mx-auto px-6 py-8">
		<!-- Header -->
		<div class="mb-8">
			<h1 class="mb-2 text-4xl font-bold text-base-content">System Overview</h1>
			<p class="text-base-content/70">Monitor your Hostamat platform metrics and performance</p>
		</div>

		<!-- Metrics Grid -->
		<div class="mb-12 grid grid-cols-1 gap-6 md:grid-cols-3">
			<!-- Users Card -->
			<div class="card bg-base-100 shadow-md transition-shadow duration-300 hover:shadow-lg">
				<div class="card-body">
					<div class="mb-4 flex items-center justify-between">
						<div class="rounded-xl bg-primary/10 p-3">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-6 w-6 text-primary"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"
								/>
							</svg>
						</div>
						<div class="text-right">
						<div class="text-3xl font-bold text-base-content">
							{data.stats.userCount}
						</div>
							<div class="text-sm text-base-content/70">Users</div>
						</div>
					</div>
					<h3 class="mb-1 card-title text-base">Registered Users</h3>
					<p class="text-sm text-base-content/70">Total unique users in the system</p>
				</div>
			</div>

			<!-- Sites Card -->
			<div class="card bg-base-100 shadow-md transition-shadow duration-300 hover:shadow-lg">
				<div class="card-body">
					<div class="mb-4 flex items-center justify-between">
						<div class="rounded-xl bg-success/10 p-3">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-6 w-6 text-success"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
								/>
							</svg>
						</div>
						<div class="text-right">
						<div class="text-3xl font-bold text-base-content">
							{data.stats.siteCount}
						</div>
							<div class="text-sm text-base-content/70">Sites</div>
						</div>
					</div>
					<h3 class="mb-1 card-title text-base">Sites Created</h3>
					<p class="text-sm text-base-content/70">Total sites hosted on the platform</p>
				</div>
			</div>

			<!-- Upload Sessions Card -->
			<div class="card bg-base-100 shadow-md transition-shadow duration-300 hover:shadow-lg">
				<div class="card-body">
					<div class="mb-4 flex items-center justify-between">
						<div class="rounded-xl bg-warning/10 p-3">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-6 w-6 text-warning"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
								/>
							</svg>
						</div>
						<div class="text-right">
						<div class="text-3xl font-bold text-base-content">
							{data.stats.uploadSessionCount}
						</div>
							<div class="text-sm text-base-content/70">Sessions</div>
						</div>
					</div>
					<h3 class="mb-1 card-title text-base">Upload Sessions</h3>
					<p class="text-sm text-base-content/70">Total file upload sessions processed</p>
				</div>
			</div>
		</div>

		<!-- Additional Info Section -->
		<div class="card bg-base-100 shadow-sm">
			<div class="card-body">
				<div class="flex items-start gap-4">
					<div class="flex-shrink-0 rounded-xl bg-info/10 p-3">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-6 w-6 text-info"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
							/>
						</svg>
					</div>
					<div>
						<h3 class="mb-2 card-title text-lg">Platform Status</h3>
						<p class="mb-4 text-base-content/70">
							Your Hostamat platform is running smoothly. Additional metrics and analytics will be
							added in future updates to provide deeper insights into your system's performance.
						</p>
						<div class="flex items-center text-sm text-base-content/70">
							<div class="mr-2 badge badge-sm badge-success"></div>
							<span>All systems operational</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
