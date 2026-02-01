<script lang="ts">
	import { X } from '$lib/icons';
	import Card from './Card.svelte';

	let { children = null, header = '', message = '' } = $props();
</script>

<div class="error">
	<Card border kind="danger" style="height: fit-content; width: 100%; padding: 10px;">
		<div class="header">
			<X size={25} color="var(--color-danger-50)" />
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
	.error {
		width: 100%;
		animation: shake 0.4s cubic-bezier(0.4, 0, 0.2, 1);
	}

	@keyframes shake {
		0%, 100% {
			transform: translateX(0);
		}
		25% {
			transform: translateX(-8px);
		}
		75% {
			transform: translateX(8px);
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
