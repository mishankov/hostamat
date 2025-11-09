<script lang="ts">
	import { goto } from '$app/navigation';

	let token = '';
	let error = '';

	async function handleLogin() {
		error = '';

		console.log('Login attempt with token:', token);

		// Simple test - if token is empty, show error
		if (!token.trim()) {
			error = 'Please enter a token';
			return;
		}

		try {
			const response = await fetch('/api/auth/login', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ token })
			});

			console.log('Login response status:', response.status);

			if (response.ok) {
				console.log('Login successful, redirecting...');
				await goto('/');
			} else {
				const errorData = await response.text();
				console.log('Login error:', errorData);
				error = 'Invalid token';
			}
		} catch (err) {
			console.error('Login request failed:', err);
			error = 'Login failed. Please try again.';
		}
	}
</script>

<div
	class="flex min-h-screen items-center justify-center bg-gradient-to-br from-slate-50 to-blue-50 p-4 dark:from-slate-900 dark:to-slate-800"
>
	<div class="w-full max-w-md">
		<!-- Logo Header -->
		<div class="mb-8 text-center">
			<div class="mb-4 inline-flex items-center gap-3">
				<div
					class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-blue-500 to-purple-600 shadow-lg"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-6 w-6 text-white"
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
				<div class="text-left">
					<h1
						class="bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-2xl font-bold text-transparent"
					>
						Hostamat
					</h1>
					<p class="text-sm text-slate-500">Dashboard</p>
				</div>
			</div>
		</div>

		<!-- Login Card -->
		<div class="rounded-2xl border border-slate-200/50 bg-white/80 p-8 shadow-xl backdrop-blur-sm">
			<h2 class="mb-6 text-center text-xl font-semibold text-slate-800">
				Sign in to your dashboard
			</h2>

			{#if error}
				<div class="mb-6 flex items-start gap-3 rounded-xl border border-red-200 bg-red-50 p-4">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="mt-0.5 h-5 w-5 flex-shrink-0 text-red-500"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
						/>
					</svg>
					<div class="text-sm text-red-700">{error}</div>
				</div>
			{/if}

			<div class="space-y-6">
				<div>
					<label for="token-input" class="mb-2 block text-sm font-medium text-slate-700">
						Authentication Token
					</label>
					<input
						id="token-input"
						type="password"
						bind:value={token}
						placeholder="Enter your access token"
						class="w-full rounded-xl border border-slate-300 bg-white/50 px-4 py-3 backdrop-blur-sm transition-colors duration-200 focus:border-blue-500 focus:ring-2 focus:ring-blue-500"
						required
					/>
				</div>

				<button
					on:click={handleLogin}
					class="w-full rounded-xl bg-gradient-to-r from-blue-500 to-purple-600 px-4 py-3 font-medium text-white shadow-lg transition-all duration-200 hover:from-blue-600 hover:to-purple-700 hover:shadow-xl active:scale-95"
				>
					Sign In
				</button>
			</div>

			<div class="mt-6 border-t border-slate-200 pt-6">
				<p class="text-center text-xs text-slate-500">
					Your authentication token can be found in the dashboard configuration
				</p>
			</div>
		</div>
	</div>
</div>
