<script lang="ts">
	import type { LayoutProps } from './$types';

	let { children }: LayoutProps = $props();
</script>

<main class="page-shell auth-shell">
	<section class="auth-wrap">
		<div class="form">{@render children()}</div>
	</section>
</main>

<style>
	.auth-shell {
		min-height: calc(100vh - 13.75rem);
		align-items: center;
	}

	.auth-wrap {
		width: min(26.25rem, 100%);
		border: var(--border-regular);
		background: var(--color-background-object);
	}

	.form {
		display: flex;
		flex-direction: column;
		gap: var(--space-md);
		padding: var(--space-xl);
	}
</style>
