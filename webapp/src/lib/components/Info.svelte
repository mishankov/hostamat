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
	}

	.header {
		display: flex;
		flex-direction: row;
		gap: 10px;
		align-items: center;
	}

	.content {
		display: flex;
		flex-direction: column;
		align-items: start;
		gap: 10px;

		padding-top: 10px;
	}
</style>
