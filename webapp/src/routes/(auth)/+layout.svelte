<script lang="ts">
	import { Card } from '$lib/components';
	import type { LayoutProps } from './$types';

	let { children }: LayoutProps = $props();
</script>

<main>
	<div class="auth-container">
		<Card border>
			<div class="form">
				{@render children()}
			</div>
		</Card>
	</div>
</main>

<style>
	main {
		width: 100%;
		min-height: calc(100vh - 140px);

		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;

		padding: 40px 24px;
	}

	.auth-container {
		width: 100%;
		max-width: 400px;

		animation: fadeInUp 0.4s var(--ease-out-expo) forwards;
	}

	@keyframes fadeInUp {
		from {
			opacity: 0;
			transform: translateY(20px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.form {
		display: flex;
		flex-direction: column;
		gap: 20px;

		padding: 32px;
	}

	@media (max-width: 480px) {
		main {
			padding: 24px 16px;
		}

		.form {
			padding: 24px;
		}
	}
</style>
