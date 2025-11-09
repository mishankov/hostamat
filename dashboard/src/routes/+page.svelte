<script lang="ts">
	// import NumberFlow from '@number-flow/svelte'
	import { goto } from '$app/navigation';

	export let data: { stats: { userCount: number; siteCount: number; uploadSessionCount: number } };

	async function handleLogout() {
		await fetch('/api/auth/logout', { method: 'POST' });
		await goto('/login');
	}
</script>

<div
	class="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50 dark:from-slate-900 dark:to-slate-800"
>
	<!-- Navigation -->
	<nav
		class="flex items-center justify-between border-b border-slate-200/50 bg-white/80 px-6 py-4 shadow-sm backdrop-blur-md"
	>
		<div class="flex items-center gap-3">
			<div
				class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-blue-500 to-purple-600"
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
			<div class="flex items-center gap-2">
				<span
					class="bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-xl font-bold text-transparent"
					>Hostamat</span
				>
				<span class="hidden font-medium text-slate-500 sm:inline">Dashboard</span>
			</div>
		</div>
		<div>
			<button
				on:click={handleLogout}
				class="flex items-center gap-2 rounded-lg border border-slate-300 bg-white px-4 py-2 text-sm font-medium text-slate-700 shadow-sm transition-colors duration-200 hover:border-slate-400 hover:bg-slate-50"
			>
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
	</nav>

	<!-- Main Content -->
	<div class="container mx-auto px-6 py-8">
		<!-- Header -->
		<div class="mb-8">
			<h1 class="mb-2 text-4xl font-bold text-slate-800 dark:text-white">System Overview</h1>
			<p class="text-slate-600 dark:text-slate-300">
				Monitor your Hostamat platform metrics and performance
			</p>
		</div>

		<!-- Metrics Grid -->
		<div class="mb-12 grid grid-cols-1 gap-6 md:grid-cols-3">
			<!-- Users Card -->
			<div class="group relative">
				<div
					class="absolute -inset-0.5 rounded-2xl bg-gradient-to-r from-blue-500 to-purple-600 opacity-10 blur transition duration-300 group-hover:opacity-20"
				></div>
				<div
					class="relative rounded-xl border border-slate-200/50 bg-white/80 shadow-sm backdrop-blur-sm transition-all duration-300 hover:shadow-md"
				>
					<div class="p-6">
						<div class="mb-4 flex items-center justify-between">
							<div class="rounded-xl bg-blue-50 p-3">
								<svg
									xmlns="http://www.w3.org/2000/svg"
									class="h-6 w-6 text-blue-600"
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
								<div class="text-3xl font-bold text-slate-800">{data.stats.userCount}</div>
								<div class="text-sm text-slate-500">Users</div>
							</div>
						</div>
						<h3 class="mb-1 font-semibold text-slate-700">Registered Users</h3>
						<p class="text-sm text-slate-500">Total unique users in the system</p>
					</div>
				</div>
			</div>

			<!-- Sites Card -->
			<div class="group relative">
				<div
					class="absolute -inset-0.5 rounded-2xl bg-gradient-to-r from-green-500 to-emerald-600 opacity-10 blur transition duration-300 group-hover:opacity-20"
				></div>
				<div
					class="relative rounded-xl border border-slate-200/50 bg-white/80 shadow-sm backdrop-blur-sm transition-all duration-300 hover:shadow-md"
				>
					<div class="p-6">
						<div class="mb-4 flex items-center justify-between">
							<div class="rounded-xl bg-green-50 p-3">
								<svg
									xmlns="http://www.w3.org/2000/svg"
									class="h-6 w-6 text-green-600"
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
								<div class="text-3xl font-bold text-slate-800">{data.stats.siteCount}</div>
								<div class="text-sm text-slate-500">Sites</div>
							</div>
						</div>
						<h3 class="mb-1 font-semibold text-slate-700">Sites Created</h3>
						<p class="text-sm text-slate-500">Total sites hosted on the platform</p>
					</div>
				</div>
			</div>

			<!-- Upload Sessions Card -->
			<div class="group relative">
				<div
					class="absolute -inset-0.5 rounded-2xl bg-gradient-to-r from-orange-500 to-red-600 opacity-10 blur transition duration-300 group-hover:opacity-20"
				></div>
				<div
					class="relative rounded-xl border border-slate-200/50 bg-white/80 shadow-sm backdrop-blur-sm transition-all duration-300 hover:shadow-md"
				>
					<div class="p-6">
						<div class="mb-4 flex items-center justify-between">
							<div class="rounded-xl bg-orange-50 p-3">
								<svg
									xmlns="http://www.w3.org/2000/svg"
									class="h-6 w-6 text-orange-600"
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
								<div class="text-3xl font-bold text-slate-800">{data.stats.uploadSessionCount}</div>
								<div class="text-sm text-slate-500">Sessions</div>
							</div>
						</div>
						<h3 class="mb-1 font-semibold text-slate-700">Upload Sessions</h3>
						<p class="text-sm text-slate-500">Total file upload sessions processed</p>
					</div>
				</div>
			</div>
		</div>

		<!-- Additional Info Section -->
		<div class="rounded-2xl border border-slate-200/50 bg-white/50 p-8 backdrop-blur-sm">
			<div class="flex items-start gap-4">
				<div class="flex-shrink-0 rounded-xl bg-blue-50 p-3">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-6 w-6 text-blue-600"
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
					<h3 class="mb-2 text-lg font-semibold text-slate-800">Platform Status</h3>
					<p class="mb-4 text-slate-600">
						Your Hostamat platform is running smoothly. Additional metrics and analytics will be
						added in future updates to provide deeper insights into your system's performance.
					</p>
					<div class="flex items-center text-sm text-slate-500">
						<div class="mr-2 h-2 w-2 rounded-full bg-green-500"></div>
						<span>All systems operational</span>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
