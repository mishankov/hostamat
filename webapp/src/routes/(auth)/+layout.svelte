<script lang="ts">
	import { Card } from '$lib/components';
	import type { LayoutProps } from './$types';

	let { children }: LayoutProps = $props();
</script>

<main>
	<div class="auth-container">
		<Card border glass>
			<div class="form">
				{@render children()}
			</div>
		</Card>
	</div>
</main>

<style>
	main {
		min-height: calc(100vh - 140px);
		width: 100%;

		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;

		padding: 40px 20px;
	}

	.auth-container {
		animation: scaleIn 0.4s ease-out;
	}

	.form {
		width: 340px;
		max-width: 90vw;

		display: flex;
		flex-direction: column;
		gap: 20px;

		padding: 32px;
	}

	.form :global(h1) {
		font-size: 28px;
		text-align: center;
		margin-bottom: 8px;
	}
</style>
