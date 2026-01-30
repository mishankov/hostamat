<script lang="ts">
	import { I } from '$lib/icons';
	import Card from './Card.svelte';

	let { children = null, header = '', message = '' } = $props();
</script>

<div class="info">
	<Card border kind="action">
		<div class="info-inner">
			<div class="icon-wrapper">
				<I size={18} color="currentColor" />
			</div>
			<div class="content">
				{#if header}
					<h4 class="header">{header}</h4>
				{:else if message}
					<p class="message">{message}</p>
				{/if}
				{#if children}
					<div class="body">
						{@render children?.()}
					</div>
				{:else if message && header}
					<div class="body">
						{message}
					</div>
				{/if}
			</div>
		</div>
	</Card>
</div>

<style>
	.info {
		width: 100%;
	}

	.info-inner {
		display: flex;
		gap: 14px;
		padding: 16px;
	}

	.icon-wrapper {
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;

		width: 28px;
		height: 28px;

		background: var(--color-action-20);
		border-radius: var(--border-radius-sm);

		color: var(--color-action-60);
	}

	.content {
		display: flex;
		flex-direction: column;
		gap: 6px;
		flex: 1;
		min-width: 0;
	}

	.header {
		font-family: var(--font-display);
		font-size: 0.9375rem;
		font-weight: 600;
		color: var(--color-action-80);
		margin: 0;
	}

	.message {
		font-size: 0.875rem;
		color: var(--color-action-70);
		margin: 0;
		line-height: 1.5;
	}

	.body {
		font-size: 0.875rem;
		color: var(--color-action-60);
		line-height: 1.5;
	}
</style>
