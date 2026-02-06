<script lang="ts">
	import { Card } from '$lib/components';
	import type { LayoutProps } from './$types';

	let { children }: LayoutProps = $props();
</script>

<main class="page-shell auth-shell">
	<section class="auth-wrap">
		<Card border>
			<div class="form">{@render children()}</div>
		</Card>
	</section>
</main>

<style>
	.auth-shell {
		min-height: calc(100vh - 220px);
		align-items: center;
	}

	.auth-wrap {
		width: min(420px, 100%);
	}

	.form {
		display: flex;
		flex-direction: column;
		gap: 0.9rem;
		padding: 1.2rem;
	}
</style>
