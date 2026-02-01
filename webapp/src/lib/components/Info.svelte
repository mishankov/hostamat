<script lang="ts">
	import { I } from '$lib/icons';
	import Card from './Card.svelte';

	let { children = null, header = '', message = '' } = $props();
</script>

<div class="info">
	<Card border kind="action" style="height: fit-content; width: 100%; padding: 10px;">
		<div class="header">
			<I size={25} color="var(--color-action-50)" />
			{#if header}
				<h4>{header}</h4>
			{:else if message}
				<p>{message}</p>
			{/if}
		</div>
		{#if children}
			<div class="content">
				{@render children?.()}
			</div>
		{:else if message && header}
			<div class="content">
				{message}
			</div>
		{/if}
	</Card>
</div>

<style>
	.info {
		width: 100%;
		animation: slideIn 0.4s cubic-bezier(0.4, 0, 0.2, 1);
	}

	@keyframes slideIn {
		from {
			transform: translateX(-20px);
			opacity: 0;
		}
		to {
			transform: translateX(0);
			opacity: 1;
		}
	}

	.header {
		display: flex;
		flex-direction: row;
		gap: 12px;
		align-items: center;
	}

	.header h4 {
		margin: 0;
		font-size: 15px;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.02em;
	}

	.header p {
		margin: 0;
		font-size: 14px;
	}

	.content {
		display: flex;
		flex-direction: column;
		align-items: start;
		gap: 12px;
		padding-top: 12px;
		padding-left: 37px;
		font-size: 14px;
		line-height: 1.5;
	}
</style>
