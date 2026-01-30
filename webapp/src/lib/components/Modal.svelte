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

		max-width: min(90vw, 600px);
		max-height: 85vh;
		overflow: auto;

		border: 1px solid var(--color-border);
		border-radius: var(--border-radius-xl);

		padding: 0;

		background: var(--gradient-primary-7-10);
		color: var(--color-text-regular);

		box-shadow: var(--shadow-lg);
	}

	/* Inner highlight effect */
	dialog::before {
		content: '';
		position: absolute;
		inset: 0;
		border-radius: inherit;
		padding: 1px;
		background: linear-gradient(
			135deg,
			var(--color-primary-30) 0%,
			transparent 40%,
			transparent 100%
		);
		mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
		mask-composite: exclude;
		-webkit-mask-composite: xor;
		pointer-events: none;
		opacity: 0.3;
	}

	.primary {
		border-color: var(--color-border);
	}

	.danger {
		border-color: var(--color-danger-40);
		background: var(--gradient-danger);
	}

	.danger::before {
		background: linear-gradient(
			135deg,
			var(--color-danger-40) 0%,
			transparent 40%,
			transparent 100%
		);
	}

	dialog[open] {
		animation: modalIn 0.3s var(--ease-out-expo) forwards;
	}

	@keyframes modalIn {
		from {
			opacity: 0;
			transform: translate(-50%, -50%) scale(0.95);
		}
		to {
			opacity: 1;
			transform: translate(-50%, -50%) scale(1);
		}
	}

	::backdrop {
		background: linear-gradient(
			135deg,
			hsla(220, 40%, 4%, 0.9),
			hsla(220, 40%, 8%, 0.85)
		);
		backdrop-filter: blur(4px);
		-webkit-backdrop-filter: blur(4px);
		animation: backdropIn 0.2s ease-out forwards;
	}

	@keyframes backdropIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	div.modal-content {
		position: relative;
		padding: 28px;
	}
</style>
