<script lang="ts">
	import type { Snippet } from 'svelte';

	let {
		children,
		kind = 'primary',
		openCallback = () => {},
		closeCallback = () => {}
	}: {
		children?: Snippet;
		kind?: 'primary' | 'danger';
		openCallback?: () => void;
		closeCallback?: () => void;
	} = $props();

	let dialog = $state<HTMLDialogElement>();
	export const open = () => {
		openCallback();
		dialog?.showModal();
	};

	export const close = () => {
		closeCallback();
		dialog?.close();
	};
</script>

<dialog
	bind:this={dialog}
	class={kind}
	onclick={(e: MouseEvent) => {
		if (e.target === dialog) {
			close();
		}
	}}
>
	<div class="modal-content">
		{@render children?.()}
	</div>
</dialog>

<style>
	dialog {
		top: 50%;
		left: 50%;
		transform: translateX(-50%) translateY(-50%);

		border: var(--border-regular);
		border-radius: var(--border-radius-regular);

		padding: 0;
		max-width: 90vw;
		max-height: 90vh;

		background: var(--color-primary-7);
		color: var(--color-text-regular);

		transition: var(--transition);
		box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
	}

	.primary {
		border-color: var(--color-border);
	}

	.danger {
		border-color: var(--color-danger-50);
		background: var(--gradient-danger);
	}

	dialog[open] {
		animation: modalFadeIn 0.3s cubic-bezier(0.4, 0, 0.2, 1) forwards;
	}

	@keyframes modalFadeIn {
		0% {
			opacity: 0;
			transform: translateX(-50%) translateY(-48%) scale(0.96);
		}
		100% {
			opacity: 1;
			transform: translateX(-50%) translateY(-50%) scale(1);
		}
	}

	::backdrop {
		background-color: rgba(0, 0, 0, 0.6);
		backdrop-filter: blur(4px);
		animation: backdropFadeIn 0.3s ease-out;
	}

	@keyframes backdropFadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	div.modal-content {
		padding: 28px;
	}
</style>
