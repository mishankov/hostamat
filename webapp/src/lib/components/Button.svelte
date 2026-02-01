<script lang="ts">
	type ButtonKind = 'primary' | 'confirm' | 'danger' | 'action';

	type ButtonProps = {
		kind?: ButtonKind;
		disabled?: boolean;
		filled?: boolean;
		link?: string;
		ref?: any;
		children?: any;
		onclick?: (e: Event) => void;
	};

	let {
		kind = 'primary',
		disabled = false,
		filled = false,
		link = '',
		children = undefined,
		ref = $bindable(null),
		onclick = undefined
	}: ButtonProps = $props();
</script>

{#if link}
	<a href={link}>
		<button {disabled} bind:this={ref} class={[kind, filled && 'filled']} {onclick}>
			{@render children?.()}
		</button>
	</a>
{:else}
	<button {disabled} bind:this={ref} class={[kind, filled && 'filled']} {onclick}>
		{@render children?.()}
	</button>
{/if}

<style>
	button {
		appearance: none;
		width: fit-content;

		border: var(--border-regular);
		border-radius: var(--border-radius-regular);

		padding: 12px 24px;
		font-weight: 600;
		font-size: 14px;
		text-transform: uppercase;
		letter-spacing: 0.02em;

		text-align: center;
		user-select: none;
		cursor: pointer;

		transition: var(--transition);
		position: relative;

		color: var(--color);
		background-color: var(--color-background);
		border-color: var(--color-border);
	}

	.filled {
		color: var(--color-filled);
		background-color: var(--color-background-filled);
		border-color: var(--color-background-filled);
		box-shadow: 0 4px 0 0 var(--color-shadow);
	}

	.primary {
		--color: var(--color-primary-90);
		--color-border: var(--color-primary-40);
		--color-background: transparent;
		--color-background-filled: var(--color-accent-50);
		--color-filled: white;
		--color-shadow: var(--color-accent-70);
	}

	.confirm {
		--color: var(--color-confirm-60);
		--color-border: var(--color-confirm-50);
		--color-background: transparent;
		--color-background-filled: var(--color-confirm-50);
		--color-filled: white;
		--color-shadow: var(--color-confirm-70);
	}

	.danger {
		--color: var(--color-danger-60);
		--color-border: var(--color-danger-50);
		--color-background: transparent;
		--color-background-filled: var(--color-danger-50);
		--color-filled: white;
		--color-shadow: var(--color-danger-70);
	}

	.action {
		--color: var(--color-action-60);
		--color-border: var(--color-action-50);
		--color-background: transparent;
		--color-background-filled: var(--color-action-50);
		--color-filled: white;
		--color-shadow: var(--color-action-70);
	}

	button:hover:enabled {
		transform: translateY(-2px);
		filter: brightness(1.1);
	}

	.filled:hover:enabled {
		transform: translateY(-2px);
		box-shadow: 0 6px 0 0 var(--color-shadow);
	}

	button:active:enabled {
		transform: translateY(0);
	}

	.filled:active:enabled {
		transform: translateY(2px);
		box-shadow: 0 2px 0 0 var(--color-shadow);
	}

	button:disabled {
		background-color: var(--color-inactive-20);
		color: var(--color-inactive-50);
		border-color: var(--color-inactive-40);
		cursor: not-allowed;
		opacity: 0.6;
	}
</style>
