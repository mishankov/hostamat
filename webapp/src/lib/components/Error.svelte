<script lang="ts">
	import { X } from '$lib/icons';
	import Card from './Card.svelte';

	let { children = null, header = '', message = '' } = $props();
</script>

<div class="error">
	<Card border kind="danger">
		<div class="error-content">
			<div class="error-icon">
				<X size={18} color="currentColor" />
			</div>
			<div class="error-text">
				{#if header}
					<h4 class="error-header">{header}</h4>
				{:else if message}
					<p class="error-message">{message}</p>
				{/if}
				{#if children}
					<div class="error-body">
						{@render children?.()}
					</div>
				{:else if message && header}
					<p class="error-body">{message}</p>
				{/if}
			</div>
		</div>
	</Card>
</div>

<style>
	.error {
		width: 100%;
		animation: slideUp 0.3s ease-out;
	}

	.error-content {
		display: flex;
		gap: 12px;
		padding: 16px;
	}

	.error-icon {
		flex-shrink: 0;
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--color-danger-20);
		border-radius: var(--border-radius-sm);
		color: var(--color-danger-60);
	}

	.error-text {
		display: flex;
		flex-direction: column;
		gap: 4px;
		min-width: 0;
	}

	.error-header {
		font-size: 15px;
		font-weight: 600;
		color: var(--color-danger-80);
		margin: 0;
	}

	.error-message {
		font-size: 14px;
		color: var(--color-danger-80);
		margin: 0;
	}

	.error-body {
		font-size: 14px;
		color: var(--color-danger-70);
		margin: 0;
	}
</style>
