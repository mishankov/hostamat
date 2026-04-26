<script lang="ts">
	import { onDestroy } from 'svelte';
	import IconButton from './IconButton.svelte';
	import { Copy } from '$lib/icons';

	interface Props {
		textToCopy: string;
		label?: string;
		title?: string;
		size?: number;
		grouped?: boolean;
	}

	let {
		textToCopy,
		label = 'Скопировать',
		title = 'Скопировать',
		size = 20,
		grouped = false
	}: Props = $props();

	let copied = $state(false);
	let copyTimer: ReturnType<typeof setTimeout> | null = null;

	const handleCopy = async () => {
		try {
			await navigator.clipboard.writeText(textToCopy);
			copied = true;

			if (copyTimer) {
				clearTimeout(copyTimer);
			}

			copyTimer = setTimeout(() => {
				copied = false;
			}, 1400);
		} catch {
			copied = false;
		}
	};

	onDestroy(() => {
		if (copyTimer) {
			clearTimeout(copyTimer);
		}
	});
</script>

<div class="copy-button-wrapper">
	<IconButton {grouped} {label} {title} onclick={handleCopy}>
		<Copy {size} />
	</IconButton>
	{#if copied}
		<span class="copy-tooltip" aria-live="polite">Скопировано</span>
	{/if}
</div>

<style>
	.copy-button-wrapper {
		position: relative;
		display: flex;
		align-items: center;
	}

	.copy-tooltip {
		position: absolute;
		bottom: calc(100% + 0.35rem);
		left: 50%;
		transform: translateX(-50%);
		font-size: 0.72rem;
		font-weight: 600;
		color: var(--color-success-strong);
		background: var(--color-background-object);
		border: var(--border-width) var(--border-style) var(--color-border);
		border-radius: var(--radius-sm);
		padding: 0.15rem 0.4rem;
		box-shadow: var(--shadow-sm);
		white-space: nowrap;
		pointer-events: none;
		z-index: 1;
	}
</style>
