<script lang="ts">
	import { Check } from '$lib/icons';
	import Card from './Card.svelte';

	let { children = null, header = '', message = '' } = $props();
</script>

<div class="success">
	<Card border kind="confirm">
		<div class="success-content">
			<div class="success-icon">
				<Check size={18} color="currentColor" />
			</div>
			<div class="success-text">
				{#if header}
					<h4 class="success-header">{header}</h4>
				{:else if message}
					<p class="success-message">{message}</p>
				{/if}
				{#if children}
					<div class="success-body">
						{@render children?.()}
					</div>
				{:else if message && header}
					<p class="success-body">{message}</p>
				{/if}
			</div>
		</div>
	</Card>
</div>

<style>
	.success {
		width: 100%;
		animation: slideUp 0.3s ease-out;
	}

	.success-content {
		display: flex;
		gap: 12px;
		padding: 16px;
	}

	.success-icon {
		flex-shrink: 0;
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--color-confirm-20);
		border-radius: var(--border-radius-sm);
		color: var(--color-confirm-60);
	}

	.success-text {
		display: flex;
		flex-direction: column;
		gap: 4px;
		min-width: 0;
	}

	.success-header {
		font-size: 15px;
		font-weight: 600;
		color: var(--color-confirm-80);
		margin: 0;
	}

	.success-message {
		font-size: 14px;
		color: var(--color-confirm-80);
		margin: 0;
	}

	.success-body {
		font-size: 14px;
		color: var(--color-confirm-70);
		margin: 0;
	}
</style>
