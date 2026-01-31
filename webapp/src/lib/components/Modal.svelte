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
		position: fixed;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);

		max-width: min(90vw, 500px);
		max-height: 85vh;

		border: 1px solid var(--glass-border);
		border-radius: var(--border-radius-lg);

		padding: 0;

		background: var(--color-background-elevated);
		color: var(--color-text-regular);

		box-shadow: var(--shadow-lg);
		overflow: hidden;
	}

	.primary {
		border-color: var(--color-border);
	}

	.danger {
		border-color: var(--color-danger-40);
		background: var(--gradient-danger);
	}

	dialog[open] {
		animation: modalIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
	}

	@keyframes modalIn {
		0% {
			opacity: 0;
			transform: translate(-50%, -50%) scale(0.9);
		}
		100% {
			opacity: 1;
			transform: translate(-50%, -50%) scale(1);
		}
	}

	::backdrop {
		background: rgba(0, 0, 0, 0.6);
		backdrop-filter: blur(4px);
		-webkit-backdrop-filter: blur(4px);
		animation: fadeIn 0.2s ease-out;
	}

	div.modal-content {
		padding: 24px;
		overflow-y: auto;
		max-height: calc(85vh - 2px);
	}
</style>
