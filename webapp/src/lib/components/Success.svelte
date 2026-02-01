<script lang="ts">
	import { Check } from '$lib/icons';
	import Card from './Card.svelte';

	let { children = null, header = '', message = '' } = $props();
</script>

<div class="success">
	<Card border kind="confirm" style="height: fit-content; width: 100%; padding: 10px;">
		<div class="header">
			<Check size={25} color="var(--color-confirm-50)" />
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
	.success {
		width: 100%;
		animation: successBounce 0.5s cubic-bezier(0.4, 0, 0.2, 1);
	}

	@keyframes successBounce {
		0% {
			transform: scale(0.95);
			opacity: 0;
		}
		50% {
			transform: scale(1.02);
		}
		100% {
			transform: scale(1);
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
		align-items: start;
		gap: 12px;
		padding-top: 12px;
		padding-left: 37px;
		font-size: 14px;
		line-height: 1.5;
	}
</style>
