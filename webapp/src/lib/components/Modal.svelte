<script lang="ts">
	import type { Snippet } from 'svelte';
	import { X } from '$lib/icons';

	let {
		children,
		kind = 'primary',
		openCallback = () => {},
		closeCallback = () => {},
		titleId = '',
		descriptionId = '',
		initialFocusSelector = '',
		closeOnBackdrop = true,
		showCloseButton = true,
		closeLabel = 'Закрыть окно'
	}: {
		children?: Snippet;
		kind?: 'primary' | 'danger';
		openCallback?: () => void;
		closeCallback?: () => void;
		titleId?: string;
		descriptionId?: string;
		initialFocusSelector?: string;
		closeOnBackdrop?: boolean;
		showCloseButton?: boolean;
		closeLabel?: string;
	} = $props();

	let dialog = $state<HTMLDialogElement>();
	let previousFocusedElement: HTMLElement | null = null;
	const focusInitialElement = () => {
		if (!dialog) {
			return;
		}

		const explicitTarget = initialFocusSelector
			? dialog.querySelector<HTMLElement>(initialFocusSelector)
			: null;

		const defaultTarget = dialog.querySelector<HTMLElement>(
			'[data-autofocus], button:not([disabled]), [href], input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])'
		);

		(explicitTarget || defaultTarget)?.focus();
	};

	const restoreFocus = () => {
		if (previousFocusedElement && document.contains(previousFocusedElement)) {
			previousFocusedElement.focus();
		}

		previousFocusedElement = null;
	};

	export const open = () => {
		openCallback();
		previousFocusedElement = document.activeElement instanceof HTMLElement ? document.activeElement : null;

		if (dialog?.open) {
			return;
		}

		dialog?.showModal();
		queueMicrotask(() => {
			focusInitialElement();
		});
	};

	export const close = () => {
		dialog?.close();
	};

	const onDialogClose = () => {
		closeCallback();
		restoreFocus();
	};
</script>

<dialog
	bind:this={dialog}
	class={kind}
	aria-labelledby={titleId || undefined}
	aria-describedby={descriptionId || undefined}
	aria-modal="true"
	onclose={onDialogClose}
	onclick={(e: MouseEvent) => {
		if (closeOnBackdrop && e.target === dialog) {
			close();
		}
	}}
>
		<div class={['modal-content', showCloseButton && 'has-close']}>
			{#if showCloseButton}
				<button class="close-modal" type="button" aria-label={closeLabel} title={closeLabel} onclick={close}>
					<X size={18} color="currentColor" />
				</button>
			{/if}
		{@render children?.()}
	</div>
</dialog>

<style>
	dialog {
		top: 50%;
		left: 50%;
		transform: translateX(-50%) translateY(-50%);

		border: var(--border-regular);
		border-radius: calc(var(--border-radius-regular) + 2px);

		padding: 0;

		background: var(--color-background-object);
		color: var(--color-text-regular);
		box-shadow: var(--shadow-lg);
		backdrop-filter: blur(18px);

		transition: var(--transition);
		max-width: min(680px, calc(100vw - 1rem));
		max-height: calc(100vh - 1rem);
		overflow: auto;
	}

	.primary {
		border-color: var(--color-border);
	}

	.danger {
		border-color: var(--color-danger-60);
		background: color-mix(in hsl, var(--color-danger-10) 92%, var(--color-background-object));
	}

	dialog[open] {
		animation: fadein 180ms cubic-bezier(0.22, 1, 0.36, 1) forwards;
	}

	@keyframes fadein {
		0% {
			opacity: 0;
			transform: translateX(-50%) translateY(-50%) scale(0.95);
		}
		100% {
			opacity: 1;
			transform: translateX(-50%) translateY(-50%) scale(1);
		}
	}

	::backdrop {
		background-color: hsl(210 28% 10% / 0.35);
		backdrop-filter: blur(2px);
	}

	div.modal-content {
		padding: 20px;
		position: relative;
	}

	div.modal-content.has-close {
		padding-right: 56px;
	}

	.close-modal {
		all: unset;
		position: absolute;
		top: 14px;
		right: 14px;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 28px;
		height: 28px;
		border-radius: var(--border-radius-regular);
		border: var(--border-regular);
		background: var(--color-background-object);
		color: var(--color-primary-90);
		cursor: pointer;
		transition: var(--transition);
	}

	.close-modal:hover {
		background: color-mix(in hsl, var(--color-background-object) 94%, var(--color-accent-30) 6%);
	}

	dialog.danger .close-modal {
		border-color: var(--color-danger-60);
		background: transparent;
		color: var(--color-danger-60);
	}

	dialog.danger .close-modal:hover {
		background: color-mix(in hsl, var(--color-danger-60) 90%, var(--color-danger-90) 10%);
		color: var(--color-background-object);
	}
</style>
