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

		background: var(--primary-gradient);
		color: var(--color-text);

		transition: var(--transition);
	}

	.primary {
		border-color: var(--color-border);
	}

	.danger {
		border-color: var(--color-danger);
		background: var(--danger-gradient);
	}

	dialog[open] {
		animation: fadein 300ms ease-in forwards;
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
		background-color: rgba(0, 0, 0, 0.5);
	}

	div.modal-content {
		padding: 20px;
	}
</style>
