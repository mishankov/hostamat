<script lang="ts">
	import type { Snippet } from 'svelte';
	import { X } from '$lib/icons';

	type ModalTone = 'neutral' | 'danger';
	type LegacyModalKind = 'primary' | 'danger';

	let {
		children,
		tone = undefined,
		kind = undefined,
		openCallback = () => {},
		closeCallback = () => {},
		title = '',
		panelWidth = '',
		titleId = '',
		descriptionId = '',
		initialFocusSelector = '',
		closeOnBackdrop = true,
		showCloseButton = true,
		closeLabel = 'Закрыть окно'
	}: {
		children?: Snippet;
		tone?: ModalTone;
		kind?: LegacyModalKind;
		openCallback?: () => void;
		closeCallback?: () => void;
		title?: string;
		panelWidth?: string;
		titleId?: string;
		descriptionId?: string;
		initialFocusSelector?: string;
		closeOnBackdrop?: boolean;
		showCloseButton?: boolean;
		closeLabel?: string;
	} = $props();

	const resolvedTone = $derived<ModalTone>(tone || (kind === 'danger' ? 'danger' : 'neutral'));

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
		previousFocusedElement =
			document.activeElement instanceof HTMLElement ? document.activeElement : null;

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
	class={resolvedTone}
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
	<div class="modal-panel" style:width={panelWidth || undefined}>
		{#if title || showCloseButton}
			<div class="modal-header">
				{#if title}
					<h2 id={titleId || undefined} class={['modal-title', showCloseButton && 'with-close']}>
						{title}
					</h2>
				{/if}

				{#if showCloseButton}
					<button
						class="close-modal"
						type="button"
						aria-label={closeLabel}
						title={closeLabel}
						onclick={close}
					>
						<X size={18} color="currentColor" />
					</button>
				{/if}
			</div>
		{/if}

		<div class="modal-body">
			{@render children?.()}
		</div>
	</div>
</dialog>

<style>
	dialog {
		top: 50%;
		left: 50%;
		transform: translateX(-50%) translateY(-50%);

		border: var(--border-regular);
		border-radius: var(--radius-lg);

		padding: 0;

		background: var(--color-background-object);
		color: var(--color-text-regular);
		box-shadow: var(--shadow-xl);
		backdrop-filter: blur(18px);

		transition: var(--transition);
		max-width: min(42.5rem, calc(100vw - 1rem));
		max-height: calc(100vh - 1rem);
		overflow: auto;
	}

	.neutral {
		border-color: var(--color-border);
	}

	.danger {
		border-color: var(--color-danger-border);
		background: var(--color-danger-surface);
	}

	dialog.danger .modal-title {
		color: var(--color-danger-text);
	}

	dialog[open] {
		animation: modal-panel-in 220ms var(--ease-out-quart) forwards;
	}

	@keyframes modal-panel-in {
		0% {
			opacity: 0;
			transform: translateX(-50%) translateY(-50%) scale(0.96);
		}
		100% {
			opacity: 1;
			transform: translateX(-50%) translateY(-50%) scale(1);
		}
	}

	::backdrop {
		background-color: hsl(210 28% 10% / 0.35);
		backdrop-filter: blur(2px);
		animation: backdrop-in 200ms var(--ease-out-quart) forwards;
	}

	div.modal-panel {
		box-sizing: border-box;
		max-width: calc(100vw - 1rem);
		padding: var(--space-xl);
	}

	div.modal-header {
		position: relative;
		margin-bottom: var(--space-xl);
		min-height: 1.75rem;
	}

	.modal-title {
		margin: 0;
	}

	.modal-title.with-close {
		padding-right: 3rem;
	}

	.close-modal {
		all: unset;
		position: absolute;
		top: 0;
		right: 0;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 1.75rem;
		height: 1.75rem;
		border-radius: var(--radius-md);
		border: var(--border-regular);
		background: var(--color-background-object);
		color: var(--color-primary-90);
		cursor: pointer;
		transition: var(--transition);
	}

	.close-modal:hover {
		background: var(--color-accent-surface);
	}

	dialog.danger .close-modal {
		border-color: var(--color-danger-border);
		background: transparent;
		color: var(--color-danger-strong);
	}

	dialog.danger .close-modal:hover {
		background: var(--color-danger-strong);
		color: var(--color-background-object);
	}

	@media (max-width: 520px) {
		div.modal-panel {
			padding: var(--space-lg);
		}

		.modal-title {
			padding-right: 0;
		}

		.modal-title.with-close {
			padding-right: 2.75rem;
		}
	}
</style>
