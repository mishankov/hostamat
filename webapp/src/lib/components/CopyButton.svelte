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
	let copiedLeaving = $state(false);
	let copiedPulse = $state(0);
	let copyTimer: ReturnType<typeof setTimeout> | null = null;
	let hideTimer: ReturnType<typeof setTimeout> | null = null;

	const clearCopyTimers = () => {
		if (copyTimer) {
			clearTimeout(copyTimer);
			copyTimer = null;
		}

		if (hideTimer) {
			clearTimeout(hideTimer);
			hideTimer = null;
		}
	};

	const handleCopy = async () => {
		try {
			await navigator.clipboard.writeText(textToCopy);
			clearCopyTimers();
			copied = true;
			copiedLeaving = false;
			copiedPulse += 1;

			copyTimer = setTimeout(() => {
				copiedLeaving = true;
				hideTimer = setTimeout(() => {
					copied = false;
					copiedLeaving = false;
				}, 180);
			}, 1200);
		} catch {
			clearCopyTimers();
			copied = false;
			copiedLeaving = false;
		}
	};

	onDestroy(() => {
		clearCopyTimers();
	});
</script>

<div class="copy-button-wrapper">
	<IconButton {grouped} {label} {title} onclick={handleCopy}>
		<Copy {size} />
	</IconButton>
	{#if copied}
		{#key copiedPulse}
			<span class={['copy-tooltip', copiedLeaving && 'leaving']} aria-live="polite"
				>Скопировано</span
			>
		{/key}
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
		font-size: 0.75rem;
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
		transform-origin: 50% 100%;
		animation: copy-confirm-in 340ms var(--ease-out-quart) both;
	}

	.copy-tooltip.leaving {
		animation: copy-confirm-out 180ms var(--ease-out-quart) both;
	}

	@keyframes copy-confirm-in {
		0% {
			opacity: 0;
			transform: translateX(-50%) translateY(0.35rem) scale(0.94);
		}

		58% {
			opacity: 1;
			transform: translateX(-50%) translateY(-0.08rem) scale(1.02);
		}

		100% {
			opacity: 1;
			transform: translateX(-50%) translateY(0) scale(1);
		}
	}

	@keyframes copy-confirm-out {
		0% {
			opacity: 1;
			transform: translateX(-50%) translateY(0) scale(1);
		}

		100% {
			opacity: 0;
			transform: translateX(-50%) translateY(0.35rem) scale(0.94);
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.copy-tooltip {
			animation: none;
		}
	}
</style>
