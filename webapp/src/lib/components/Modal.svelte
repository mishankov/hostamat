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
		border-radius: calc(var(--border-radius-regular) + 2px);

		padding: 0;

		background: var(--color-background-object);
		color: var(--color-text-regular);
		box-shadow: var(--shadow-medium);
		backdrop-filter: blur(18px);

		transition: var(--transition);
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
	}
</style>
