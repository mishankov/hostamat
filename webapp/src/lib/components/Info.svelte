<script lang="ts">
	import { I } from '$lib/icons';
	import Card from './Card.svelte';

	let { children = null, header = '', message = '' } = $props();
</script>

<div class="info">
	<Card border kind="action">
		<div class="info-content">
			<div class="info-icon">
				<I size={18} color="currentColor" />
			</div>
			<div class="info-text">
				{#if header}
					<h4 class="info-header">{header}</h4>
				{:else if message}
					<p class="info-message">{message}</p>
				{/if}
				{#if children}
					<div class="info-body">
						{@render children?.()}
					</div>
				{:else if message && header}
					<p class="info-body">{message}</p>
				{/if}
			</div>
		</div>
	</Card>
</div>

<style>
	.info {
		width: 100%;
		animation: slideUp 0.3s ease-out;
	}

	.info-content {
		display: flex;
		gap: 12px;
		padding: 16px;
	}

	.info-icon {
		flex-shrink: 0;
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--color-action-20);
		border-radius: var(--border-radius-sm);
		color: var(--color-action-60);
	}

	.info-text {
		display: flex;
		flex-direction: column;
		gap: 4px;
		min-width: 0;
	}

	.info-header {
		font-size: 15px;
		font-weight: 600;
		color: var(--color-action-80);
		margin: 0;
	}

	.info-message {
		font-size: 14px;
		color: var(--color-action-80);
		margin: 0;
	}

	.info-body {
		font-size: 14px;
		color: var(--color-action-70);
		margin: 0;
	}
</style>
