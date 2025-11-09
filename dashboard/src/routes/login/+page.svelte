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

<div class="flex min-h-screen items-center justify-center bg-base-100 p-4">
	<div class="w-full max-w-md">
		<!-- Logo Header -->
		<div class="mb-8 text-center">
			<div class="mb-4 inline-flex items-center gap-3">
				<div class="avatar">
					<div
						class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-primary to-secondary shadow-lg"
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
				</div>
				<div class="text-left">
					<h1 class="text-2xl font-bold text-primary">Hostamat</h1>
					<p class="text-sm text-base-content/70">Dashboard</p>
				</div>
			</div>
		</div>

		<!-- Login Card -->
		<div class="card bg-base-100 shadow-xl">
			<div class="card-body">
				<h2 class="mb-6 card-title justify-center text-xl">Sign in to your dashboard</h2>

				{#if error}
					<div class="mb-6 alert alert-error">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-5 w-5 flex-shrink-0"
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
						<span>{error}</span>
					</div>
				{/if}

				<div class="space-y-6">
					<div class="form-control">
						<label for="token-input" class="label">
							<span class="label-text">Authentication Token</span>
						</label>
						<input
							id="token-input"
							type="password"
							bind:value={token}
							placeholder="Enter your access token"
							class="input-bordered input w-full"
							required
						/>
					</div>

					<button on:click={handleLogin} class="btn w-full btn-primary"> Sign In </button>
				</div>

				<div class="divider mt-6"></div>
				<p class="text-center text-xs text-base-content/70">
					Your authentication token can be found in the dashboard configuration
				</p>
			</div>
		</div>
	</div>
</div>
